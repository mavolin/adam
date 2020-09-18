package mock

import (
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

type Module struct {
	MetaReturn     plugin.ModuleMeta
	CommandsReturn []plugin.Command
	ModulesReturn  []plugin.Module
}

func (c Module) Meta() plugin.ModuleMeta    { return c.MetaReturn }
func (c Module) Commands() []plugin.Command { return c.CommandsReturn }
func (c Module) Modules() []plugin.Module   { return c.ModulesReturn }

type ModuleMeta struct {
	Name              string
	ShortDescription  string
	LongDescription   string
	Hidden            bool
	ChannelTypes      plugin.ChannelTypes
	BotPermissions    *discord.Permissions
	Restrictions      plugin.RestrictionFunc
	ThrottlingOptions plugin.ThrottlingOptions
}

func (c ModuleMeta) GetName() string                                    { return c.Name }
func (c ModuleMeta) GetShortDescription(*localization.Localizer) string { return c.ShortDescription }
func (c ModuleMeta) GetLongDescription(*localization.Localizer) string  { return c.LongDescription }
func (c ModuleMeta) IsHidden() bool                                     { return c.Hidden }
func (c ModuleMeta) GetDefaultChannelTypes() plugin.ChannelTypes        { return c.ChannelTypes }
func (c ModuleMeta) GetDefaultBotPermissions() *discord.Permissions     { return c.BotPermissions }
func (c ModuleMeta) GetDefaultRestrictionFunc() plugin.RestrictionFunc  { return c.Restrictions }
func (c ModuleMeta) GetThrottlingOptions() plugin.ThrottlingOptions     { return c.ThrottlingOptions }
