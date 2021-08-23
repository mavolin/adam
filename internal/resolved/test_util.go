package resolved

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

func newProviderFromSources(sources []plugin.Source) *PluginProvider {
	r := NewPluginResolver(nil)

	for _, src := range sources {
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

	p := r.NewProvider(nil, nil)
	p.Resolve()
	return p
}
