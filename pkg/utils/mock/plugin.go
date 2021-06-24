package mock

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/internal/resolved"
	"github.com/mavolin/adam/pkg/plugin"
)

// =============================================================================
// Command
// =====================================================================================

type Command struct {
	plugin.CommandMeta
	InvokeFunc func(*state.State, *plugin.Context) (interface{}, error)
}

var _ plugin.Command = Command{}

func (c Command) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	return c.InvokeFunc(s, ctx)
}

// =============================================================================
// Module
// =====================================================================================

type Module struct {
	plugin.ModuleMeta
	CommandsReturn []plugin.Command
	ModulesReturn  []plugin.Module
}

var _ plugin.Module = Module{}

func (m Module) Commands() []plugin.Command { return m.CommandsReturn }
func (m Module) Modules() []plugin.Module   { return m.ModulesReturn }

// =============================================================================
// Resolver
// =====================================================================================

// NewPluginProvider creates a new plugin.Provider using the passed
// []plugin.Source and []plugin.UnavailablePluginSource.
func NewPluginProvider(sources []plugin.Source, unavailableSources []plugin.UnavailablePluginSource) plugin.Provider {
	r := resolved.NewPluginResolver(nil)

	for _, source := range sources {
		source := source
		r.AddSource(source.Name,
			func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				return source.Commands, source.Modules, nil
			})
	}

	for _, unavailableSource := range unavailableSources {
		unavailableSource := unavailableSource
		r.AddSource(unavailableSource.Name,
			func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
				return nil, nil, unavailableSource.Error
			})
	}

	p := r.NewProvider(nil, nil)
	p.Resolve()
	return p
}

// ResolveCommand creates a plugin.ResolvedCommand from the passed
// plugin.Command using the passed source name.
func ResolveCommand(sourceName string, cmd plugin.Command) plugin.ResolvedCommand {
	r := resolved.NewPluginResolver(nil)
	r.AddSource(sourceName,
		func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
			return []plugin.Command{cmd}, nil, nil
		})

	return r.NewProvider(nil, nil).Commands()[0]
}

// ResolveModule creates a plugin.ResolvedModule from the passed
// plugin.Module using the passed source name.
func ResolveModule(sourceName string, mod plugin.Module) plugin.ResolvedModule {
	r := resolved.NewPluginResolver(nil)
	r.AddSource(sourceName,
		func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
			return nil, []plugin.Module{mod}, nil
		})

	return r.NewProvider(nil, nil).Modules()[0]
}
