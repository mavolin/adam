package module

import (
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// Meta is the static, unlocalized, implementation of the plugin.ModuleMeta
// interface.
type Meta struct {
	// Name is the name of the module.
	// It may not contain whitespace or dots.
	Name string
	// ShortDescription is an optional short description of the module.
	ShortDescription string
	// LongDescription is an optional long description of the module.
	LongDescription string

	// Hidden specifies whether this module should be hidden from the help
	// message.
	//
	// All subcommands and submodules will be hidden as well.
	Hidden bool

	// DefaultChannelTypes are the plugin.ChannelTypes the used as default, for
	// all submodule and -commands that don't specify some.
	//
	// If this is not set, the channel types of the parent will be used.
	DefaultChannelTypes plugin.ChannelTypes
	// DefaultBotPermissions are the discord.Permissions used as default, for
	// all submodule and -commands that don't specify some.
	//
	// If this is not set, the bot permissions of the parent will be used.
	DefaultBotPermissions *discord.Permissions
	// DefaultRestrictions are the restrictions used as default, for all
	// submodule and -commands that don't specify some.
	//
	// If this is not set, the restrictions of the parent will be used.
	DefaultRestrictions plugin.RestrictionFunc
	// DefaultThrottler is the plugin.Throttler used as default, for all
	// submodule and -commands that don't specify some.
	//
	// If this is not set, the throttler of the parent will be used.
	DefaultThrottler plugin.Throttler
}

var _ plugin.ModuleMeta = Meta{}

func (m Meta) GetName() string                                   { return m.Name }
func (m Meta) GetShortDescription(*i18n.Localizer) string        { return m.ShortDescription }
func (m Meta) GetLongDescription(*i18n.Localizer) string         { return m.LongDescription }
func (m Meta) IsHidden() bool                                    { return m.Hidden }
func (m Meta) GetDefaultChannelTypes() plugin.ChannelTypes       { return m.DefaultChannelTypes }
func (m Meta) GetDefaultBotPermissions() *discord.Permissions    { return m.DefaultBotPermissions }
func (m Meta) GetDefaultRestrictionFunc() plugin.RestrictionFunc { return m.DefaultRestrictions }
func (m Meta) GetDefaultThrottler() plugin.Throttler             { return m.DefaultThrottler }
