package mock

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

type Command struct {
	MetaReturn plugin.CommandMeta
	InvokeFunc func(*state.State, *plugin.Context) (interface{}, error)
}

func (c Command) Meta() plugin.CommandMeta { return c.MetaReturn }

func (c Command) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	return c.InvokeFunc(s, ctx)
}

type CommandMeta struct {
	Name              string
	Aliases           []string
	Args              ArgConfig
	ShortDescription  string
	LongDescription   string
	Examples          []string
	Hidden            bool
	ChannelTypes      plugin.ChannelTypes
	BotPermissions    *discord.Permissions
	Restrictions      plugin.RestrictionFunc
	ThrottlingOptions plugin.ThrottlingOptions
}

func (c CommandMeta) GetName() string                                    { return c.Name }
func (c CommandMeta) GetAliases() []string                               { return c.Aliases }
func (c CommandMeta) GetArgs() plugin.ArgConfig                          { return c.Args }
func (c CommandMeta) GetShortDescription(*localization.Localizer) string { return c.ShortDescription }
func (c CommandMeta) GetLongDescription(*localization.Localizer) string  { return c.LongDescription }
func (c CommandMeta) GetExamples(*localization.Localizer) []string       { return c.Examples }
func (c CommandMeta) IsHidden() bool                                     { return c.Hidden }
func (c CommandMeta) GetChannelTypes() plugin.ChannelTypes               { return c.ChannelTypes }
func (c CommandMeta) GetBotPermissions() *discord.Permissions            { return c.BotPermissions }
func (c CommandMeta) GetRestrictionFunc() plugin.RestrictionFunc         { return c.Restrictions }
func (c CommandMeta) GetThrottlingOptions() plugin.ThrottlingOptions     { return c.ThrottlingOptions }

type ArgConfig struct {
	Expect string

	ArgsReturn  plugin.Args
	FlagsReturn plugin.Flags
	ErrorReturn error
}

func (a ArgConfig) Parse(_ string, _ *state.State, _ *plugin.Context) (plugin.Args, plugin.Flags, error) {
	return a.ArgsReturn, a.FlagsReturn, a.ErrorReturn
}
