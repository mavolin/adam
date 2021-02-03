package mock

import (
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type Module struct {
	plugin.ModuleMeta
	CommandsReturn []plugin.Command
	ModulesReturn  []plugin.Module
}

var _ plugin.Module = Module{}

func (m Module) Commands() []plugin.Command { return m.CommandsReturn }
func (m Module) Modules() []plugin.Module   { return m.ModulesReturn }

// GenerateRegisteredModule generates a RegisteredModule from the passed source
// module using the passed provider name.
func GenerateRegisteredModule(providerName string, mod plugin.Module) *plugin.RegisteredModule {
	rmod := plugin.GenerateRegisteredModules([]plugin.Repository{
		{
			ProviderName: providerName,
			Modules:      []plugin.Module{mod},
		},
	})

	return rmod[0]
}

type ModuleMeta struct {
	Name             string
	ShortDescription string
	LongDescription  string
}

var _ plugin.ModuleMeta = ModuleMeta{}

func (m ModuleMeta) GetName() string                            { return m.Name }
func (m ModuleMeta) GetShortDescription(*i18n.Localizer) string { return m.ShortDescription }
func (m ModuleMeta) GetLongDescription(*i18n.Localizer) string  { return m.LongDescription }
