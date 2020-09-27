package mock

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

type Command struct {
	plugin.CommandMeta
	InvokeFunc func(*state.State, *plugin.Context) (interface{}, error)
}

func (c Command) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	return c.InvokeFunc(s, ctx)
}

type RegisteredCommand struct {
	ParentReturn plugin.RegisteredModule
	ParentError  error

	IdentifierReturn       plugin.Identifier
	NameReturn             string
	AliasesReturn          []string
	ArgsReturn             plugin.ArgConfig
	ShortDescriptionReturn string
	LongDescriptionReturn  string
	ExamplesReturn         []string
	IsHiddenReturn         bool
	ChannelTypesReturn     plugin.ChannelTypes
	BotPermissionsReturn   discord.Permissions
	IsRestrictedReturn     error
	ThrottlerReturn        plugin.Throttler
	InvokeFunc             func(s *state.State, ctx *plugin.Context) (interface{}, error)
}

func (r RegisteredCommand) Parent() (plugin.RegisteredModule, error) {
	return r.ParentReturn, r.ParentError
}

func (r RegisteredCommand) Identifier() plugin.Identifier { return r.IdentifierReturn }
func (r RegisteredCommand) Name() string                  { return r.NameReturn }
func (r RegisteredCommand) Aliases() []string             { return r.AliasesReturn }
func (r RegisteredCommand) Args() plugin.ArgConfig        { return r.ArgsReturn }

func (r RegisteredCommand) ShortDescription(*localization.Localizer) string {
	return r.ShortDescriptionReturn
}
func (r RegisteredCommand) LongDescription(*localization.Localizer) string {
	return r.LongDescriptionReturn
}

func (r RegisteredCommand) Examples(*localization.Localizer) []string { return r.ExamplesReturn }
func (r RegisteredCommand) IsHidden() bool                            { return r.IsHiddenReturn }
func (r RegisteredCommand) ChannelTypes() plugin.ChannelTypes         { return r.ChannelTypesReturn }
func (r RegisteredCommand) BotPermissions() discord.Permissions       { return r.BotPermissionsReturn }

func (r RegisteredCommand) IsRestricted(*state.State, *plugin.Context) error {
	return r.IsRestrictedReturn
}

func (r RegisteredCommand) Throttler() plugin.Throttler { return r.ThrottlerReturn }

func (r RegisteredCommand) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	return r.InvokeFunc(s, ctx)
}

type CommandMeta struct {
	Name             string
	Aliases          []string
	Args             ArgConfig
	ShortDescription string
	LongDescription  string
	Examples         []string
	Hidden           bool
	ChannelTypes     plugin.ChannelTypes
	BotPermissions   *discord.Permissions
	Restrictions     plugin.RestrictionFunc
	Throttler        plugin.Throttler
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
func (c CommandMeta) GetThrottler() plugin.Throttler                     { return c.Throttler }

type ArgConfig struct {
	Expect string

	ArgsReturn  plugin.Args
	FlagsReturn plugin.Flags
	ErrorReturn error
}

func (a ArgConfig) Parse(_ string, _ *state.State, _ *plugin.Context) (plugin.Args, plugin.Flags, error) {
	return a.ArgsReturn, a.FlagsReturn, a.ErrorReturn
}
