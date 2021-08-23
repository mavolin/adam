package resolve

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/internal/resolved"
	"github.com/mavolin/adam/pkg/plugin"
)

// Provider mocks the plugin.Provider interface.
type Provider struct {
	// Sources are the plugin.Sources the plugin provider provides.
	//
	// Modifying this after using one of the Providers methods will have no
	// effect.
	Sources []plugin.Source
	// UnavailableSources are the plugin.UnavailableSources the plugin provider
	// could not provide.
	//
	// Modifying this after using one of the Providers methods will have no
	// effect.
	UnavailableSources []plugin.UnavailableSource

	mut sync.Mutex
	p   plugin.Provider
}

var _ plugin.Provider = new(Provider)

func (p *Provider) lazyInit() {
	p.mut.Lock()
	defer p.mut.Unlock()

	if p.p != nil {
		return
	}

	r := resolved.NewPluginResolver(nil)

	for _, src := range p.Sources {
		source := src

		if src.Name == plugin.BuiltInSource {
			for _, cmd := range src.Commands {
				r.AddBuiltInCommand(cmd)
			}

			for _, mod := range src.Modules {
				r.AddBuiltInModule(mod)
			}
		} else {
			r.AddSource(source.Name,
				func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
					return source.Commands, source.Modules, nil
				})
		}
	}

	for _, usrc := range p.UnavailableSources {
		unavailableSource := usrc
		r.AddSource(unavailableSource.Name,
			func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				return nil, nil, unavailableSource.Error
			})
	}

	provider := r.NewProvider(nil, nil)
	provider.Resolve()
	p.p = provider
}

func (p *Provider) PluginSources() []plugin.Source {
	p.lazyInit()
	return p.p.PluginSources()
}

func (p *Provider) Commands() []plugin.ResolvedCommand {
	p.lazyInit()
	return p.p.Commands()
}

func (p *Provider) Modules() []plugin.ResolvedModule {
	p.lazyInit()
	return p.p.Modules()
}

func (p *Provider) Command(id plugin.ID) plugin.ResolvedCommand {
	p.lazyInit()
	return p.p.Command(id)
}

func (p *Provider) Module(id plugin.ID) plugin.ResolvedModule {
	p.lazyInit()
	return p.p.Module(id)
}

func (p *Provider) FindCommand(invoke string) plugin.ResolvedCommand {
	p.lazyInit()
	return p.p.FindCommand(invoke)
}

func (p *Provider) FindCommandWithArgs(invoke string) (cmd plugin.ResolvedCommand, args string) {
	p.lazyInit()
	return p.p.FindCommandWithArgs(invoke)
}

func (p *Provider) FindModule(invoke string) plugin.ResolvedModule {
	p.lazyInit()
	return p.p.FindModule(invoke)
}

func (p *Provider) UnavailablePluginSources() []plugin.UnavailableSource {
	p.lazyInit()
	return p.p.UnavailablePluginSources()
}
