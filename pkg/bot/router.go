package bot

import (
	"strings"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/replier"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

var ErrUnknownCommand = errors.NewUserErrorl(unknownCommandErrorDescription)

// Route attempts to route the passed message.
// It aborts, if the message is not a valid invoke.
func (b *Bot) Route(base *state.Base, msg *discord.Message, member *discord.Member) {
	// Only accept regular text messages.
	// Also check if a bot wrote the message, if !b.AllowBot.
	if msg.Type != discord.DefaultMessage || (!b.AllowBot && msg.Author.Bot) {
		return
	}

	prefixes, lang := b.SettingsProvider(base, msg)

	ok, invoke := b.hasPrefix(msg.Content, prefixes)
	if !ok { // not an invoke
		return
	}

	if member != nil {
		member.User = msg.Author
	}

	ctx := &plugin.Context{
		Message:          *msg,
		Member:           member,
		Base:             base,
		Localizer:        b.LocalizationManager.Localizer(lang),
		Prefixes:         prefixes,
		BotOwnerIDs:      b.Owners,
		ReplyMiddlewares: b.ReplyMiddlewares,
		Replier:          replier.WrapState(b.State),
		DiscordDataProvider: &discordDataProvider{
			s:         b.State,
			guildID:   msg.GuildID,
			channelID: msg.ChannelID,
			selfID:    b.selfID,
		},
	}
	ctx.ErrorHandler = newCtxErrorHandler(b.State, ctx, b.ErrorHandler)

	var args string

	if b.AsyncPluginProviders {
		ctx.InvokedCommand, ctx.Provider, args = b.routeCommandAsync(invoke, base, msg)
		if ctx.InvokedCommand == nil {
			b.ErrorHandler(ErrUnknownCommand, b.State, ctx)
		}
	} else {
		ctx.InvokedCommand, ctx.Provider, args = b.routeCommand(invoke, base, msg)
		if ctx.InvokedCommand == nil {
			b.ErrorHandler(ErrUnknownCommand, b.State, ctx)
		}
	}

	err := b.invoke(ctx, args)
	if err != nil {
		b.ErrorHandler(err, b.State, ctx)
	}
}

// hasPrefix checks if the passed invoke starts with one of the passed
// prefixes or a mention of the bot.
// If so it returns true and the invoke stripped of the prefix.
func (b *Bot) hasPrefix(invoke string, prefixes []string) (bool, string) {
	indexes := b.selfMentionRegexp.FindStringIndex(invoke)
	if indexes != nil {
		return true, strings.TrimLeft(invoke[indexes[1]:], whitespace)
	}

	for _, p := range prefixes {
		if strings.HasPrefix(invoke, p) {
			return true, strings.TrimLeft(invoke[len(p):], whitespace)
		}
	}

	return false, ""
}

func (b *Bot) routeCommand(
	invoke string, base *state.Base, msg *discord.Message,
) (*plugin.RegisteredCommand, plugin.Provider, string) {
	ctxprovider := &ctxPluginProvider{
		base: base,
		msg:  msg,
		repos: []plugin.Repository{
			{
				ProviderName: plugin.BuiltInProvider,
				Commands:     b.commands,
				Modules:      b.modules,
				Defaults:     b.PluginDefaults,
			},
		},
	}

	cmd, args := findCommand(invoke, ctxprovider, ctxprovider.repos[0])
	if cmd != nil {
		ctxprovider.remProviders = b.pluginProviders
		return cmd, ctxprovider, args
	}

	for i, p := range b.pluginProviders {
		cmds, mods, err := p.provider(base, msg)
		if err != nil {
			ctxprovider.unavailableProviders = append(ctxprovider.unavailableProviders, plugin.UnavailablePluginProvider{
				Name:  p.name,
				Error: err,
			})
		} else {
			repo := plugin.Repository{
				ProviderName: p.name,
				Modules:      mods,
				Commands:     cmds,
				Defaults:     p.defaults,
			}

			ctxprovider.repos = append(ctxprovider.repos, repo)

			cmd, args := findCommand(invoke, ctxprovider, repo)
			if cmd != nil {
				ctxprovider.remProviders = b.pluginProviders[i+1:]
				return cmd, ctxprovider, args
			}
		}
	}

	return nil, nil, ""
}

func (b *Bot) routeCommandAsync(
	invoke string, base *state.Base, msg *discord.Message,
) (*plugin.RegisteredCommand, plugin.Provider, string) {
	ctxprovider := &ctxPluginProvider{
		base: base,
		msg:  msg,
		repos: []plugin.Repository{
			{
				ProviderName: plugin.BuiltInProvider,
				Commands:     b.commands,
				Modules:      b.modules,
				Defaults:     b.PluginDefaults,
			},
		},
		async: true,
	}

	cmd, args := findCommand(invoke, ctxprovider, ctxprovider.repos[0])
	if cmd != nil {
		ctxprovider.remProviders = b.pluginProviders
		return cmd, ctxprovider, args
	}

	repos, up := pluginProvidersAsync(base, msg, b.pluginProviders)
	ctxprovider.repos = append(ctxprovider.repos, repos...)
	ctxprovider.unavailableProviders = append(ctxprovider.unavailableProviders, up...)

	for _, r := range repos {
		cmd, args := findCommand(invoke, ctxprovider, r)
		if cmd != nil {
			return cmd, ctxprovider, args
		}
	}

	return nil, nil, ""
}

func (b *Bot) invoke(ctx *plugin.Context, args string) error {
	middlewares := b.Middlewares()

	for _, mod := range ctx.InvokedCommand.SourceParents {
		if m, ok := mod.(Middlewarer); ok {
			middlewares = append(middlewares, m.Middlewares()...)
		}
	}

	if m, ok := ctx.InvokedCommand.Source.(Middlewarer); ok {
		middlewares = append(middlewares, m.Middlewares()...)
	}

	inv := func(_ *state.State, ctx *plugin.Context) error { return b.invokeCommand(ctx, args) }

	for i := len(middlewares) - 1; i >= 0; i-- {
		inv = middlewares[i](inv)
	}

	return inv(b.State, ctx)
}

func (b *Bot) invokeCommand(ctx *plugin.Context, args string) error {
	err := ctx.InvokedCommand.IsRestricted(b.State, ctx)
	if err != nil {
		return err
	}

	if ctx.InvokedCommand.Args != nil {
		ctx.Args, ctx.Flags, err = ctx.InvokedCommand.Args.Parse(args, b.State, ctx)
		if err != nil {
			return err
		}
	}

	reply, err := ctx.InvokedCommand.Invoke(b.State, ctx)
	rerr := b.handleReply(reply, ctx)
	if err != nil {
		ctx.HandleErrorSilent(rerr)

		return err
	}

	return rerr
}

func (b *Bot) handleReply(reply interface{}, ctx *plugin.Context) (err error) {
	if reply == nil {
		return nil
	}

	switch reply := reply.(type) {
	case uint:
		_, err = ctx.Reply(reply)
	case uint8:
		_, err = ctx.Reply(reply)
	case uint16:
		_, err = ctx.Reply(reply)
	case uint32:
		_, err = ctx.Reply(reply)
	case uint64:
		_, err = ctx.Reply(reply)
	case int:
		_, err = ctx.Reply(reply)
	case int8:
		_, err = ctx.Reply(reply)
	case int16:
		_, err = ctx.Reply(reply)
	case int32:
		_, err = ctx.Reply(reply)
	case int64:
		_, err = ctx.Reply(reply)
	case float32:
		_, err = ctx.Reply(reply)
	case float64:
		_, err = ctx.Reply(reply)
	case string:
		_, err = ctx.Reply(reply)
	case discord.Embed:
		_, err = ctx.ReplyEmbed(reply)
	case *discord.Embed:
		_, err = ctx.ReplyEmbed(*reply)
	case *embedutil.Builder:
		_, err = ctx.ReplyEmbedBuilder(reply)
	case api.SendMessageData:
		_, err = ctx.ReplyMessage(reply)
	case i18n.Term:
		_, err = ctx.Replylt(reply)
	case *i18n.Config:
		_, err = ctx.Replyl(reply)
	case plugin.Reply:
		err = reply.SendReply(b.State, ctx)
	default:
		err = &ReplyTypeError{Reply: reply}
	}

	return err
}
