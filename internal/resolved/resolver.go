package resolved

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/event"

	"github.com/mavolin/adam/pkg/plugin"
)

type (
	PluginResolver struct {
		CustomSources []UnqueriedPluginSource

		Commands []plugin.Command
		Modules  []plugin.Module

		builtinProvider *PluginProvider
		argParser       plugin.ArgParser
	}

	UnqueriedPluginSource struct {
		Name string
		Func PluginSourceFunc
	}

	PluginSourceFunc func(*event.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error)
)

func NewPluginResolver(defaultArgParser plugin.ArgParser) *PluginResolver {
	return &PluginResolver{
		builtinProvider: &PluginProvider{usedNames: make(map[string]struct{})},
		argParser:       defaultArgParser,
	}
}

func (r *PluginResolver) AddSource(name string, f PluginSourceFunc) {
	for i, rp := range r.CustomSources {
		if rp.Name == name {
			r.CustomSources = append(r.CustomSources[:i], r.CustomSources[i+1:]...)
		}
	}

	r.CustomSources = append(r.CustomSources, UnqueriedPluginSource{
		Name: name,
		Func: f,
	})
}

func (r *PluginResolver) AddBuiltInCommand(scmd plugin.Command) {
	r.Commands = append(r.Commands, scmd)
	r.builtinProvider.commands = insertCommand(r.builtinProvider.commands,
		newCommand(nil, r.builtinProvider, plugin.BuiltInSource, nil, scmd), -1)
}

func (r *PluginResolver) AddBuiltInModule(smod plugin.Module) {
	r.Modules = append(r.Modules, smod)
	r.builtinProvider.modules = insertModule(r.builtinProvider.modules,
		newModule(nil, r.builtinProvider, plugin.BuiltInSource, smod), -1)
}
