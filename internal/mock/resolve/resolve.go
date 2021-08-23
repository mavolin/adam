package resolve

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/internal/resolved"
	"github.com/mavolin/adam/pkg/plugin"
)

// Command creates a plugin.ResolvedCommand from the passed
// plugin.Command using the passed source name.
func Command(sourceName string, cmd plugin.Command) plugin.ResolvedCommand {
	r := resolved.NewPluginResolver(nil)
	r.AddSource(sourceName,
		func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
			return []plugin.Command{cmd}, nil, nil
		})

	return r.NewProvider(nil, nil).Commands()[0]
}

// Module creates a plugin.ResolvedModule from the passed
// plugin.Module using the passed source name.
func Module(sourceName string, mod plugin.Module) plugin.ResolvedModule {
	r := resolved.NewPluginResolver(nil)
	r.AddSource(sourceName,
		func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
			return nil, []plugin.Module{mod}, nil
		})

	return r.NewProvider(nil, nil).Modules()[0]
}
