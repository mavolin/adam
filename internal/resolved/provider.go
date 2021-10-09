package resolved

import (
	"sync"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/event"

	"github.com/mavolin/adam/pkg/plugin"
)

type PluginProvider struct {
	resolver *PluginResolver

	// mut is the sync.RWMutex used to secure sources, commands, modules,
	// unavailableSources and userNames before all plugins are fully resolved.
	// After fully resolving, no locks are needed, as the data won't be
	// modified anymore.
	mut sync.RWMutex

	base *event.Base
	msg  *discord.Message

	sources []plugin.Source

	commands []plugin.ResolvedCommand
	modules  []plugin.ResolvedModule

	unavailableSources []plugin.UnavailableSource
	usedNames          map[string]struct{}
}

var _ plugin.Provider = new(PluginProvider)

func (r *PluginResolver) NewProvider(base *event.Base, msg *discord.Message) *PluginProvider {
	p := &PluginProvider{
		resolver: r,
		base:     base,
		msg:      msg,
	}

	p.commands = replaceCommandProvider(nil, r.builtinProvider.commands, p)
	p.modules = replaceModuleProvider(nil, r.builtinProvider.modules, p)

	p.usedNames = make(map[string]struct{}, len(r.builtinProvider.usedNames))
	for name := range r.builtinProvider.usedNames {
		p.usedNames[name] = struct{}{}
	}

	p.sources = make([]plugin.Source, 1, len(r.CustomSources)+1)
	p.sources[0] = plugin.Source{
		Name:     plugin.BuiltInSource,
		Commands: r.Commands,
		Modules:  r.Modules,
	}

	return p
}

// ================================ Getters ================================

func (p *PluginProvider) PluginSources() []plugin.Source {
	p.Resolve()

	// p.sources won't be modified again after resolving.
	// No mutex or copy required.
	return p.sources
}

func (p *PluginProvider) Commands() []plugin.ResolvedCommand {
	p.Resolve()

	// p.commands won't be modified again after resolving.
	// No mutex or copy required.
	return p.commands
}

func (p *PluginProvider) Modules() []plugin.ResolvedModule {
	p.Resolve()

	// p.modules won't be modified again after resolving.
	// No mutex or copy required.
	return p.modules
}

func (p *PluginProvider) Command(id plugin.ID) plugin.ResolvedCommand {
	// do we have a match among the built-in commands?
	p.mut.RLock()
	rcmd := p.command(id)
	p.mut.RUnlock()

	if rcmd != nil {
		return rcmd
	}

	// if there are any unresolved commands, resolve all and search among them
	if p.Resolve() {
		return p.command(id)
	}

	// already searched anything, found nothing
	return nil
}

// command returns the plugin.ResolvedCommand with the passed id, as found
// among the already resolved commands.
//
// Callers must ensure that a read lock exist, if needed.
func (p *PluginProvider) command(id plugin.ID) plugin.ResolvedCommand {
	if id.Parent().IsRoot() {
		return findCommand(p.commands, id.Name(), false)
	}

	rmod := p.module(id.Parent())
	if rmod == nil {
		return nil
	}

	return findCommand(rmod.Commands(), id.Name(), false)
}

func (p *PluginProvider) Module(id plugin.ID) plugin.ResolvedModule {
	p.Resolve()
	return p.module(id)
}

// module returns the plugin.ResolvedModule with the passed id, as found
// among the already resolved modules.
//
// Callers must ensure that a read lock exist, if needed.
func (p *PluginProvider) module(id plugin.ID) plugin.ResolvedModule {
	all := id.All()
	if len(all) < 1 {
		return nil
	}

	all = all[1:]

	rmod := findModule(p.modules, all[0].Name())
	if rmod == nil {
		return nil
	}

	for _, subID := range all[1:] {
		rmod = rmod.FindModule(subID.Name())
		if rmod == nil {
			return nil
		}
	}

	return rmod
}

func (p *PluginProvider) FindCommand(invoke string) plugin.ResolvedCommand {
	id := plugin.NewIDFromInvoke(invoke)

	// try built-in sources first
	p.mut.RLock()
	rcmd := p.findCommand(id)
	p.mut.RUnlock()

	if rcmd != nil {
		return rcmd
	}

	// if didn't search all sources already, resolve and search again
	if p.Resolve() {
		return p.findCommand(id)
	}

	return nil
}

// findCommand attempts to find the command whose parent id matches id.Parent()
// and whose name or alias matches id.Name(), among the already resolved
// commands.
//
// Callers must ensure that a read lock exist, if needed.
func (p *PluginProvider) findCommand(id plugin.ID) plugin.ResolvedCommand {
	if id.Parent().IsRoot() {
		return findCommand(p.commands, id.Name(), true)
	}

	rmod := p.module(id.Parent())
	if rmod == nil {
		return nil
	}

	return rmod.FindCommand(id.Name())
}

func (p *PluginProvider) FindCommandWithArgs(invoke string) (plugin.ResolvedCommand, string) {
	p.mut.RLock()
	rcmd, args := p.findCommandWithArgs(invoke)
	p.mut.RUnlock()

	if rcmd != nil {
		return rcmd, args
	}

	p.Resolve()
	return p.findCommandWithArgs(invoke)
}

// findCommandWithArgs attempts to find a command with the given invoke among
// the already resolved commands.
//
// Callers must ensure that a read lock exist, if needed.
func (p *PluginProvider) findCommandWithArgs(invoke string) (rcmd plugin.ResolvedCommand, args string) {
	var word string
	word, invoke = firstWord(invoke)
	if len(word) == 0 {
		return nil, ""
	}

	id := plugin.NewIDFromInvoke(word)

	rcmd = p.findCommand(id) // top-level command?
	if rcmd != nil {
		return rcmd, invoke
	}

	rmod := p.module(id)
	if rmod == nil {
		return nil, ""
	}

	for {
		word, invoke = firstWord(invoke)
		if word == "" {
			return nil, ""
		}

		rcmd = rmod.FindCommand(word)
		if rcmd != nil {
			return rcmd, invoke
		}

		rmod = rmod.FindModule(word)
		if rmod == nil {
			return nil, ""
		}
	}
}

func (p *PluginProvider) FindModule(invoke string) plugin.ResolvedModule {
	return p.Module(plugin.NewIDFromInvoke(invoke))
}

func (p *PluginProvider) UnavailablePluginSources() []plugin.UnavailableSource {
	p.Resolve()
	return p.unavailableSources
}

func (p *PluginProvider) Resolve() bool {
	p.mut.Lock()
	defer p.mut.Unlock()

	if len(p.sources) > 1 || len(p.resolver.CustomSources) == 0 {
		return false
	}

	type result struct {
		scmds []plugin.Command
		smods []plugin.Module
		err   error
	}

	results := make([]result, len(p.resolver.CustomSources))

	var wg sync.WaitGroup
	var mut sync.Mutex

	wg.Add(len(p.resolver.CustomSources))
	for i, src := range p.resolver.CustomSources {
		go func(i int, src UnqueriedPluginSource) {
			scmds, smods, err := src.Func(p.base, p.msg)

			mut.Lock()
			results[i] = result{
				scmds: scmds,
				smods: smods,
				err:   err,
			}
			mut.Unlock()

			wg.Done()
		}(i, src)
	}

	wg.Wait()

	for i, r := range results {
		sourceName := p.resolver.CustomSources[i].Name

		if r.err != nil {
			p.unavailableSources = append(p.unavailableSources, plugin.UnavailableSource{
				Name:  sourceName,
				Error: r.err,
			})
		} else {
			p.sources = append(p.sources, plugin.Source{
				Name:     sourceName,
				Commands: r.scmds,
				Modules:  r.smods,
			})

			p.addCommands(sourceName, r.scmds)
			p.addModules(sourceName, r.smods)
		}
	}

	return true
}

func (p *PluginProvider) addCommands(sourceName string, scmds []plugin.Command) {
	cp := make([]plugin.ResolvedCommand, len(p.commands), len(p.commands)+len(scmds))
	copy(cp, p.commands)
	p.commands = cp

	for _, scmd := range scmds {
		rcmd := newCommand(nil, p, sourceName, nil, scmd)
		if rcmd != nil {
			p.commands = insertCommand(p.commands, rcmd, -1)
		}
	}
}

func (p *PluginProvider) addModules(sourceName string, smods []plugin.Module) {
	cp := make([]plugin.ResolvedCommand, len(p.commands), len(p.commands)+len(smods))
	copy(cp, p.commands)
	p.commands = cp

	for _, smod := range smods {
		i := searchModule(p.modules, smod.GetName())
		if i < len(p.modules) && p.modules[i].Name() == smod.GetName() {
			p.modules[i].(*Module).update(p, sourceName, smod)
		} else {
			rmod := newModule(nil, p, sourceName, smod)
			if rmod != nil {
				p.modules = insertModule(p.modules, rmod, i)
			}
		}
	}
}
