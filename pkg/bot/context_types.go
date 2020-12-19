package bot

import (
	"sort"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// =============================================================================
// plugin.Provider
// =====================================================================================

type ctxPluginProvider struct {
	base *state.Base
	msg  *discord.Message

	// repos contains the already collected repositories
	repos []plugin.Repository
	// remProviders are the remaining, i.e. uncalled, plugin providers.
	remProviders []pluginProvider

	commands []*plugin.RegisteredCommand
	modules  []*plugin.RegisteredModule

	unavailableProviders []plugin.UnavailablePluginProvider
}

func (p *ctxPluginProvider) PluginRepositories() []plugin.Repository {
	p.lazyRepos()
	return p.repos
}

func (p *ctxPluginProvider) lazyRepos() {
	for _, remp := range p.remProviders {
		cmds, mods, err := remp.provider(p.base, p.msg)
		if err != nil {
			p.unavailableProviders = append(p.unavailableProviders, plugin.UnavailablePluginProvider{
				Name:  remp.name,
				Error: err,
			})
		} else {
			p.repos = append(p.repos, plugin.Repository{
				ProviderName: remp.name,
				Modules:      mods,
				Commands:     cmds,
				Defaults:     remp.defaults,
			})
		}
	}

	p.remProviders = nil
}

func (p *ctxPluginProvider) Commands() []*plugin.RegisteredCommand {
	p.lazyCommands()
	return p.commands
}

func (p *ctxPluginProvider) lazyCommands() {
	if p.commands == nil {
		p.lazyRepos()
		p.commands = plugin.GenerateRegisteredCommands(p.repos)
	}
}

func (p *ctxPluginProvider) Modules() []*plugin.RegisteredModule {
	p.lazyModules()
	return p.modules
}

func (p *ctxPluginProvider) lazyModules() {
	if p.modules == nil {
		p.lazyRepos()
		p.modules = plugin.GenerateRegisteredModules(p.repos)
	}
}

func (p *ctxPluginProvider) Command(id plugin.Identifier) *plugin.RegisteredCommand {
	if id.IsRoot() {
		return nil
	}

	if id.Parent().IsRoot() { // top-lvl command
		p.lazyCommands()

		name := id.Name()

		i := sort.Search(len(p.commands), func(i int) bool {
			return p.commands[i].Name >= name
		})

		if i == len(p.commands) || p.commands[i].Name != name { // nothing found
			return nil
		}

		return p.commands[i]
	}

	mod := p.Module(id.Parent())
	if mod == nil {
		return nil
	}

	return mod.FindCommand(id.Name())
}

func (p *ctxPluginProvider) Module(id plugin.Identifier) *plugin.RegisteredModule {
	p.lazyModules()

	all := id.All()
	if len(all) <= 1 { // invalid or just root
		return nil
	}

	all = all[1:]

	name := all[0].Name()

	i := sort.Search(len(p.modules), func(i int) bool {
		return p.modules[i].Name >= name
	})

	if i == len(p.modules) { // nothing found
		return nil
	}

	mod := p.modules[i]

	for _, id := range all[1:] {
		mod = mod.FindModule(id.Name())
		if mod == nil {
			return nil
		}
	}

	return mod
}

func (p *ctxPluginProvider) FindCommand(invoke string) *plugin.RegisteredCommand {
	id := plugin.NewIdentifierFromInvoke(invoke)
	return p.Command(id)
}

func (p *ctxPluginProvider) FindModule(invoke string) *plugin.RegisteredModule {
	id := plugin.NewIdentifierFromInvoke(invoke)
	return p.Module(id)
}

func (p *ctxPluginProvider) UnavailablePluginProviders() []plugin.UnavailablePluginProvider {
	p.lazyRepos()
	return p.unavailableProviders
}
