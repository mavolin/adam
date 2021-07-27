package command

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// Meta is the static, unlocalized, implementation of the plugin.CommandMeta
// interface.
type Meta struct {
	// Name is the name of the command.
	// It may not contain whitespace or dots.
	Name string
	// Aliases are the optional aliases of the command.
	// They may not contain whitespace or dots.
	Aliases []string
	// ShortDescription is an optional short description of the command.
	ShortDescription string
	// LongDescription is an optional long description of the command.
	LongDescription string

	// Args is the argument configuration of the command.
	// If this is left empty, the command won't accept any arguments.
	Args plugin.ArgConfig
	// ArgParser is the optional custom ArgParser of the command.
	ArgParser plugin.ArgParser
	// ExampleArgs contains the optional example arguments of the command.
	ExampleArgs plugin.ExampleArgs

	// Hidden specifies whether this command should be hidden from the help
	// message.
	Hidden bool

	// ChannelTypes are the plugin.ChannelTypes the command may be executed in.
	//
	// If this is not set, AllChannels will be used.
	ChannelTypes plugin.ChannelTypes
	// BotPermissions are the permissions the bot needs to execute this
	// command.
	BotPermissions discord.Permissions
	// Restrictions contains the optional restrictions of the command.
	Restrictions plugin.RestrictionFunc
	// Throttler is the optional plugin.Throttler of the command.
	Throttler plugin.Throttler
}

var _ plugin.CommandMeta = Meta{}

func (m Meta) GetName() string                                   { return m.Name }
func (m Meta) GetAliases() []string                              { return m.Aliases }
func (m Meta) GetShortDescription(*i18n.Localizer) string        { return m.ShortDescription }
func (m Meta) GetLongDescription(*i18n.Localizer) string         { return m.LongDescription }
func (m Meta) GetExampleArgs(*i18n.Localizer) plugin.ExampleArgs { return m.ExampleArgs }
func (m Meta) GetArgs() plugin.ArgConfig                         { return m.Args }
func (m Meta) GetArgParser() plugin.ArgParser                    { return m.ArgParser }
func (m Meta) IsHidden() bool                                    { return m.Hidden }
func (m Meta) GetChannelTypes() plugin.ChannelTypes              { return m.ChannelTypes }
func (m Meta) GetBotPermissions() discord.Permissions            { return m.BotPermissions }

func (m Meta) IsRestricted(s *state.State, ctx *plugin.Context) error {
	if m.Restrictions == nil {
		return nil
	}

	return m.Restrictions(s, ctx)
}

func (m Meta) GetThrottler() plugin.Throttler { return m.Throttler }

type ExampleArgs = plugin.ExampleArgs
