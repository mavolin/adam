package resolved

import (
	"sync"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/event"

	"github.com/mavolin/adam/pkg/plugin"
)

type (
	PluginResolver struct {
		Sources []UnqueriedPluginSource

		Commands []plugin.Command
		Modules  []plugin.Module

		provider  *PluginProvider
		argParser plugin.ArgParser
	}

	UnqueriedPluginSource struct {
		Name string
		Func PluginSourceFunc
	}

	PluginSourceFunc func(*event.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error)
)

func NewPluginResolver(defaultArgParser plugin.ArgParser) *PluginResolver {
	return &PluginResolver{
		provider:  &PluginProvider{usedNames: make(map[string]struct{})},
		argParser: defaultArgParser,
	}
}

func (r *PluginResolver) AddSource(name string, f PluginSourceFunc) {
	for i, rp := range r.Sources {
		if rp.Name == name {
			r.Sources = append(r.Sources[:i], r.Sources[i+1:]...)
		}
	}

	r.Sources = append(r.Sources, UnqueriedPluginSource{
		Name: name,
		Func: f,
	})
}

func (r *PluginResolver) AddBuiltInCommand(scmd plugin.Command) {
	r.Commands = append(r.Commands, scmd)
	r.provider.commands = insertCommand(r.provider.commands,
		newCommand(nil, r.provider, plugin.BuiltInSource, scmd), -1)
}

func (r *PluginResolver) AddBuiltInModule(smod plugin.Module) {
	r.Modules = append(r.Modules, smod)
	r.provider.modules = insertModule(r.provider.modules,
		newModule(nil, r.provider, plugin.BuiltInSource, smod), -1)
}

type PluginProvider struct {
	resolver *PluginResolver

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
		resolver:  r,
		base:      base,
		msg:       msg,
		sources:   make([]plugin.Source, 0, len(r.Sources)),
		usedNames: make(map[string]struct{}, len(r.provider.usedNames)),
	}

	return p
}

// ================================ Getters ================================

func (p *PluginProvider) PluginSources() []plugin.Source {
	p.Resolve()
	return p.sources
}

func (p *PluginProvider) Commands() []plugin.ResolvedCommand {
	p.Resolve()
	return p.commands
}

func (p *PluginProvider) Modules() []plugin.ResolvedModule {
	p.Resolve()
	return p.modules
}

func (p *PluginProvider) Command(id plugin.ID) plugin.ResolvedCommand {
	rcmd := p.FindCommand(id.AsInvoke())
	if rcmd != nil && rcmd.ID() != id {
		return nil
	}

	return rcmd
}

func (p *PluginProvider) Module(id plugin.ID) plugin.ResolvedModule {
	p.Resolve()

	all := id.All()
	if len(all) < 1 {
		return nil
	}

	all = all[1:]

	rmod := findModule(p.modules, all[0].Name())
	if rmod == nil {
		return nil
	}

	for _, id := range all[1:] {
		rmod = findModule(rmod.Modules(), id.Name())
		if rmod == nil {
			return nil
		}
	}

	return rmod
}

func (p *PluginProvider) FindCommand(invoke string) plugin.ResolvedCommand {
	if invoke == "" {
		return nil
	}

	id := plugin.NewIDFromInvoke(invoke)
	if id.Parent().IsRoot() {
		// try built-in sources first
		if rcmd := findCommand(p.commands, id.Name(), true); rcmd != nil {
			return rcmd
		}

		// didn't work, maybe a source other than built-in has it?
		p.Resolve()
		return findCommand(p.commands, id.Name(), true)
	}

	// try built-in sources first
	rmod := p.Module(id.Parent())
	if rmod == nil {
		return nil
	}

	return rmod.FindCommand(id.Name())
}

func (p *PluginProvider) FindCommandWithArgs(invoke string) (rcmd plugin.ResolvedCommand, args string) {
	var word string
	word, invoke = firstWord(invoke)
	if len(word) == 0 {
		return nil, ""
	}

	rcmd = p.FindCommand(word)
	if rcmd != nil {
		return rcmd, invoke
	}

	rmod := p.FindModule(word)
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

func (p *PluginProvider) Resolve() {
	p.resolveBuiltIn()

	if len(p.sources) > 1 || len(p.resolver.Sources) == 0 {
		return
	}

	type result struct {
		scmds []plugin.Command
		smods []plugin.Module
		err   error
	}

	results := make([]result, len(p.resolver.Sources))

	var wg sync.WaitGroup
	var mut sync.Mutex

	wg.Add(len(p.resolver.Sources))
	for i, src := range p.resolver.Sources {
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
		sourceName := p.resolver.Sources[i].Name

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
}

func (p *PluginProvider) resolveBuiltIn() {
	if len(p.sources) > 0 {
		return
	}

	p.commands = replaceCommandProvider(nil, p.resolver.provider.commands, p)

	p.modules = replaceModuleProvider(nil, p.resolver.provider.modules, p)

	for name := range p.resolver.provider.usedNames {
		p.usedNames[name] = struct{}{}
	}

	p.sources = append(p.sources, plugin.Source{
		Name:     plugin.BuiltInSource,
		Commands: p.resolver.Commands,
		Modules:  p.resolver.Modules,
	})
}

func (p *PluginProvider) addCommands(sourceName string, scmds []plugin.Command) {
	cp := make([]plugin.ResolvedCommand, len(p.commands), len(p.commands)+len(scmds))
	copy(cp, p.commands)
	p.commands = cp

	for _, scmd := range scmds {
		rcmd := newCommand(nil, p, sourceName, scmd)
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
			updateModule(p.modules[i].(*Module), p, sourceName, smod)
		} else {
			rmod := newModule(nil, p, sourceName, smod)
			if rmod != nil {
				p.modules = insertModule(p.modules, rmod, i)
			}
		}
	}
}
