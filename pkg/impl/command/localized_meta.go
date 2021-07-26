package command

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// LocalizedMeta is the localized implementation of the plugin.CommandMeta
// interface.
type LocalizedMeta struct {
	// Name is the name of the command.
	// It may not contain whitespace or dots.
	Name string
	// Aliases are the optional aliases of the command.
	// They may not contain whitespace or dots.
	Aliases []string
	// ShortDescription is an optional short description of the command.
	ShortDescription *i18n.Config
	// LongDescription is an optional long description of the command.
	LongDescription *i18n.Config

	// Args is the argument configuration of the command.
	// If this is left empty, the command won't accept any arguments.
	Args plugin.ArgConfig
	// ArgParser is the optional custom ArgParser of the command.
	ArgParser plugin.ArgParser
	// ExampleArgs contains the optional example arguments of the command.
	ExampleArgs exampleArgsGetter

	// Hidden specifies whether this command should be hidden from the help
	// message.
	Hidden bool

	// ChannelTypes are the plugin.ChannelTypes the command may be executed in.
	//
	// If this is not set, AllChannels.
	ChannelTypes plugin.ChannelTypes
	// BotPermissions are the permissions the bot needs to execute this
	// command.
	BotPermissions discord.Permissions
	// Restrictions contains the optional restrictions of the command.
	Restrictions plugin.RestrictionFunc
	// Throttler is the optional plugin.Throttler of the command.
	Throttler plugin.Throttler
}

var _ plugin.CommandMeta = LocalizedMeta{}

func (m LocalizedMeta) GetName() string      { return m.Name }
func (m LocalizedMeta) GetAliases() []string { return m.Aliases }

func (m LocalizedMeta) GetShortDescription(l *i18n.Localizer) string {
	desc, err := l.Localize(m.ShortDescription)
	if err != nil {
		return ""
	}

	return desc
}

func (m LocalizedMeta) GetLongDescription(l *i18n.Localizer) string {
	desc, err := l.Localize(m.LongDescription)
	if err != nil {
		return ""
	}

	return desc
}

func (m LocalizedMeta) GetExampleArgs(l *i18n.Localizer) plugin.ExampleArgs {
	if m.ExampleArgs == nil {
		return nil
	}

	return m.ExampleArgs.BaseType(l)
}

func (m LocalizedMeta) GetArgs() plugin.ArgConfig              { return m.Args }
func (m LocalizedMeta) GetArgParser() plugin.ArgParser         { return m.ArgParser }
func (m LocalizedMeta) IsHidden() bool                         { return m.Hidden }
func (m LocalizedMeta) GetChannelTypes() plugin.ChannelTypes   { return m.ChannelTypes }
func (m LocalizedMeta) GetBotPermissions() discord.Permissions { return m.BotPermissions }

func (m LocalizedMeta) IsRestricted(s *state.State, ctx *plugin.Context) error {
	if m.Restrictions == nil {
		return nil
	}

	return m.Restrictions(s, ctx)
}

func (m LocalizedMeta) GetThrottler() plugin.Throttler { return m.Throttler }

// =============================================================================
// ExampleArgs
// =====================================================================================

type exampleArgsGetter interface {
	BaseType(*i18n.Localizer) plugin.ExampleArgs
}

var _ exampleArgsGetter = plugin.ExampleArgs{}

type LocalizedExampleArgs []struct {
	// Flags is a map of example flags.
	Flags map[string]*i18n.Config
	// Args contains the example arguments.
	Args []*i18n.Config
}

var _ exampleArgsGetter = LocalizedExampleArgs{}

func (lexamples LocalizedExampleArgs) BaseType(l *i18n.Localizer) plugin.ExampleArgs {
	base := make(plugin.ExampleArgs, len(lexamples))
	var i int

Examples:
	for _, lexample := range lexamples {
		if len(lexample.Flags) > 0 {
			base[i].Flags = make(map[string]string, len(lexample.Flags))
			for name, contentConfig := range lexample.Flags {
				content, err := l.Localize(contentConfig)
				if err != nil {
					continue Examples
				}

				base[i].Flags[name] = content
			}
		}

		if len(lexample.Args) > 0 {
			base[i].Args = make([]string, len(lexample.Args))
			for j, argConfig := range lexample.Args {
				arg, err := l.Localize(argConfig)
				if err != nil {
					continue Examples
				}

				base[i].Args[j] = arg
			}
		}

		i++ // kinda hacky, but it gets the job done
	}

	// don't use i+1, since i is incremented at the end
	return base[:i]
}
