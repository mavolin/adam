package mock

import (
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type Module struct {
	plugin.ModuleMeta
	CommandsReturn []plugin.Command
	ModulesReturn  []plugin.Module
}

var _ plugin.Module = Module{}

func (c Module) Commands() []plugin.Command { return c.CommandsReturn }
func (c Module) Modules() []plugin.Module   { return c.ModulesReturn }

// GenerateRegisteredModule generates a RegisteredModule from the passed source
// module using the passed provider name.
func GenerateRegisteredModule(providerName string, smod plugin.Module) *plugin.RegisteredModule {
	rmod := plugin.GenerateRegisteredModules([]plugin.Repository{
		{
			ProviderName: providerName,
			Modules:      []plugin.Module{smod},
		},
	})

	return rmod[0]
}

type ModuleMeta struct {
	Name             string
	ShortDescription string
	LongDescription  string

	Hidden                bool
	DefaultChannelTypes   plugin.ChannelTypes
	DefaultBotPermissions *discord.Permissions
	DefaultRestrictions   plugin.RestrictionFunc
	DefaultThrottler      plugin.Throttler
}

var _ plugin.ModuleMeta = ModuleMeta{}

func (c ModuleMeta) GetName() string                                   { return c.Name }
func (c ModuleMeta) GetShortDescription(*i18n.Localizer) string        { return c.ShortDescription }
func (c ModuleMeta) GetLongDescription(*i18n.Localizer) string         { return c.LongDescription }
func (c ModuleMeta) IsHidden() bool                                    { return c.Hidden }
func (c ModuleMeta) GetDefaultChannelTypes() plugin.ChannelTypes       { return c.DefaultChannelTypes }
func (c ModuleMeta) GetDefaultBotPermissions() *discord.Permissions    { return c.DefaultBotPermissions }
func (c ModuleMeta) GetDefaultRestrictionFunc() plugin.RestrictionFunc { return c.DefaultRestrictions }
func (c ModuleMeta) GetDefaultThrottler() plugin.Throttler             { return c.DefaultThrottler }
