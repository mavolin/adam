// Package help provides the default help command.
// It is used to create help messages that can list all commands and modules,
// commands and modules in a specific module, and details about a command.
package help

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/internal/capbuilder"
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

	// ArgFormatter is the plugin.ArgFormatter used to generate command usages.
	//
	// Defaults to DefaultArgFormatter
	ArgFormatter plugin.ArgFormatter
}

// New creates a new help command using the passed Options.
func New(o Options) *Help {
	if o.HideFuncs == nil {
		o.HideFuncs = []HideFunc{
			CheckHidden(HideList), CheckChannelTypes(HideList),
			CheckRestrictions(HideList),
		}
	}

	if o.ArgFormatter == nil {
		o.ArgFormatter = DefaultArgFormatter
	}

	return &Help{
		LocalizedMeta: command.LocalizedMeta{
			Name:             "help",
			Aliases:          []string{"how"},
			ShortDescription: shortDescription,
			LongDescription:  longDescription,
			ExampleArgs:      exampleArgs,
			Args: &arg.LocalizedCommaConfig{
				Optional: []arg.LocalizedOptionalArg{
					{
						Name:        argPluginName,
						Type:        arg.Plugin,
						Description: argPluginDescription,
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
	case *plugin.RegisteredModule:
		return h.module(s, ctx, p)
	case *plugin.RegisteredCommand:
		return h.command(s, ctx, p)
	default:
		panic(fmt.Sprintf("got illegal argument type %T from arg.Plugin, but expected only (interface{})(nil), "+
			"*plugin.RegisteredCommand, or *plugin.RegisteredModule", ctx.Args[0]))
	}
}

func (h *Help) all(s *state.State, ctx *plugin.Context) (discord.Embed, error) {
	eb := BaseEmbed.Clone().
		WithSimpleTitlel(allTitle)

	if ctx.GuildID == 0 {
		eb.WithDescriptionl(allDescriptionDM)
	} else {
		eb.WithDescriptionl(allDescriptionGuild)
	}

	e, err := eb.Build(ctx.Localizer)
	if err != nil {
		return discord.Embed{}, err
	}

	b := capbuilder.New(embedutil.MaxChars-embedutil.CountChars(e), 1024)

	maxMods := 25 - len(e.Fields)

	if ctx.GuildID > 0 && !h.NoPrefix {
		prefixes, err := h.genPrefixesField(b, ctx)
		if err != nil {
			return discord.Embed{}, err
		}

		e.Fields = append([]discord.EmbedField{prefixes}, e.Fields...)
		maxMods--
	}

	if f := h.genCommandsField(b, s, ctx, ctx.Commands()); len(f.Name) > 0 {
		e.Fields = append(e.Fields, f)
		maxMods--
	}

	e.Fields = append(e.Fields, h.genModuleFields(b, s, ctx, ctx.Modules(), maxMods)...)
	return e, nil
}

func (h *Help) module(s *state.State, ctx *plugin.Context, mod *plugin.RegisteredModule) (discord.Embed, error) {
	eb := BaseEmbed.Clone().
		WithSimpleTitlel(moduleTitle.
			WithPlaceholders(moduleTitlePlaceholders{
				Module: mod.Identifier.AsInvoke(),
			}))

	if desc := mod.LongDescription(ctx.Localizer); len(desc) > 0 {
		eb.WithDescription(desc)
	}

	e, err := eb.Build(ctx.Localizer)
	if err != nil {
		return discord.Embed{}, nil
	}

	maxMods := 25 - len(e.Fields)

	b := capbuilder.New(embedutil.MaxChars-embedutil.CountChars(e), 1024)

	if f := h.genCommandsField(b, s, ctx, mod.Commands); len(f.Name) > 0 {
		e.Fields = append(e.Fields, f)
		maxMods--
	}

	e.Fields = append(e.Fields, h.genModuleFields(b, s, ctx, mod.Modules, maxMods)...)

	if len(e.Fields) == 0 {
		return discord.Embed{}, plugin.NewArgumentErrorl(pluginNotFoundError)
	}

	return e, nil
}

func (h *Help) command(s *state.State, ctx *plugin.Context, cmd *plugin.RegisteredCommand) (discord.Embed, error) {
	if len(filterCommands([]*plugin.RegisteredCommand{cmd}, s, ctx, Show, h.Options.HideFuncs...)) == 0 {
		return discord.Embed{}, plugin.NewArgumentErrorl(pluginNotFoundError)
	}

	eb := BaseEmbed.Clone().
		WithSimpleTitlel(commandTitle.
			WithPlaceholders(commandTitlePlaceholders{
				Command: cmd.Name,
			}))

	if desc := cmd.LongDescription(ctx.Localizer); len(desc) > 0 {
		eb.WithDescription(desc)
	}

	e, err := eb.Build(ctx.Localizer)
	if err != nil {
		return discord.Embed{}, err
	}

	var b strings.Builder
	b.Grow(1024)

	if aliases := h.genAliases(&b, cmd, ctx.Localizer); aliases != nil {
		e.Fields = append(e.Fields, *aliases)
	}

	e.Fields = append(e.Fields, h.genUsages(&b, ctx, cmd)...)

	if ex := h.genExamples(&b, cmd, ctx.Localizer); ex != nil {
		e.Fields = append(e.Fields, *ex)
	}

	return e, nil
}

func DefaultArgFormatter(i plugin.ArgInfo, optional, variadic bool) string {
	if optional {
		if variadic {
			return "[" + i.Name + "+]"
		}

		return "[" + i.Name + "]"
	}

	if variadic {
		return "<" + i.Name + "+>"
	}

	return "<" + i.Name + ">"
}
