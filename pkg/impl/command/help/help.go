// Package help provides the default help command.
// It is used to create help messages that can list all commands and modules,
// commands and modules in a specific module, and details about a command.
package help

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

var BaseEmbed = embedutil.NewBuilder().
	WithColor(0x6eb7b1)

// Help is the default help command.
// It should fully suffice for most bots, however, it is restricted to the same
// restrictions any embed has.
// Particularly, the following rules apply:
//
// You can have at most 23 top-level modules and any amount of top-level
// commands, or 24 top-level modules and no top-level commands.
//
// The total amount characters used for listing all top-level commands, as well
// as the commands of individual modules, may not exceed 1024.
// Otherwise, another embed field will be used, which means the number of
// allowed top-level modules shrinks by one.
//
// The total amount of characters of the embed cannot exceed 6000.
// A few of those are already taken, to display title and prefixes.
//
// Note also that some commands may be cut off, in order to stay within that
// character limit.
type Help struct {
	command.LocalizedMeta
	Options
}

var _ plugin.Command = New(Options{})

type Options struct {
	// HideFuncs are the functions used to determine the HiddenLevel of a
	// command.
	// The highest HiddenLevel returned by any the functions will be used.
	// Invalid HiddenLevels will default to Hide.
	//
	// Modules may also be hidden as a result of this, if all their commands
	// are hidden.
	// If that is the case, the lowest HiddenLevel any subcommand of the module
	// has, will be used for that module.
	//
	// Defaults to:
	// 	[]HideFunc{
	// 		CheckHidden(HideList), CheckChannelTypes(HideList),
	//		CheckRestrictions(HideList),
	//	}
	//
	// Use an empty slice to always show commands.
	HideFuncs []HideFunc
	// NoPrefix toggles whether in a guild the all embed should list the
	// available prefixes.
	NoPrefix bool
}

// New creates a new help command using the passed Options.
func New(o Options) *Help {
	if o.HideFuncs == nil {
		o.HideFuncs = []HideFunc{
			CheckHidden(HideList), CheckChannelTypes(HideList),
			CheckRestrictions(HideList),
		}
	}

	return &Help{
		LocalizedMeta: command.LocalizedMeta{
			Name:             "help",
			Aliases:          []string{"how"},
			ShortDescription: shortDescription,
			LongDescription:  longDescription,
			Examples:         examples,
			Args: &arg.LocalizedCommaConfig{
				Optional: []arg.LocalizedOptionalArg{
					{
						Name:        argsPluginName,
						Type:        arg.Plugin,
						Description: argsPluginDescription,
					},
				},
			},
			ChannelTypes:   plugin.AllChannels,
			BotPermissions: discord.PermissionSendMessages,
		},
		Options: o,
	}
}

func (h *Help) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	if ctx.Args[0] == nil {
		return h.all(s, ctx)
	}

	switch p := ctx.Args[0].(type) {
	case *plugin.RegisteredCommand:
		// do sth
		return nil, nil
	case *plugin.RegisteredModule:
		return h.module(s, ctx, p)
	default:
		panic(fmt.Sprintf("got illegal argument type %T from arg.Plugin, but expected only (interface{})(nil), "+
			"*plugin.RegisteredCommand, or *plugin.RegisteredModule", ctx.Args[0]))
	}
}

func (h *Help) commands(
	b *cappedBuilder, s *state.State, ctx *plugin.Context, cmds []*plugin.RegisteredCommand,
) (f discord.EmbedField) {
	b.reset(1024)

	h.formatCommands(b, cmds, s, ctx, Show)
	if b.b.Len() == 0 {
		return
	}

	f.Name = ctx.MustLocalize(commandsFieldName)
	b.use(len(f.Name))

	f.Value = b.string()
	return
}
