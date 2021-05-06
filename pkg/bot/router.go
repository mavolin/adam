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
	"github.com/mavolin/adam/pkg/utils/discorderr"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

var ErrUnknownCommand = errors.NewUserErrorl(unknownCommandErrorDescription)

// Route attempts to route the passed message.
// It aborts, if the message is not a valid invoke.
func (b *Bot) Route(base *state.Base, msg *discord.Message, member *discord.Member) { //nolint:funlen
	// Only accept regular text messages.
	// Also check if a bot wrote the message, if !b.AllowBot.
	// Lastly, also discard if this bot wrote this message, even if b.AllowBot.
	if msg.Type != discord.DefaultMessage || (!b.AllowBot && msg.Author.Bot) || msg.Author.ID == b.selfID {
		return
	}

	prefixes, localizer, ok := b.SettingsProvider(base, msg)
	if !ok {
		return
	}

	if localizer == nil {
		localizer = i18n.NewFallbackLocalizer()
	}

	invoke := b.hasPrefix(msg.Content, prefixes, msg.GuildID.IsValid())
	if len(invoke) == 0 {
		return
	}

	if member != nil {
		member.User = msg.Author
	}

	ctx := &plugin.Context{
		Message:     *msg,
		Member:      member,
		Base:        base,
		Localizer:   localizer,
		Prefixes:    prefixes,
		BotOwnerIDs: b.Owners,
		Replier:     replier.WrapState(b.State, false),
		DiscordDataProvider: &discordDataProvider{
			s:         b.State,
			guildID:   msg.GuildID,
			channelID: msg.ChannelID,
			selfID:    b.selfID,
		},
	}
	ctx.ErrorHandler = newCtxErrorHandler(b.State, ctx, b.ErrorHandler)

	if b.AsyncPluginProviders {
		ctx.InvokedCommand, ctx.Provider, ctx.RawArgs = b.findCommandAsync(invoke, base, msg)
	} else {
		ctx.InvokedCommand, ctx.Provider, ctx.RawArgs = b.findCommand(invoke, base, msg)
	}

	if ctx.InvokedCommand == nil {
		ctx.HandleError(ErrUnknownCommand)
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			b.PanicHandler(rec, b.State, ctx)
		}
	}()

	inv := b.applyMiddlewares(ctx)
	if err := inv(b.State, ctx); err != nil {
		b.ErrorHandler(err, b.State, ctx)
	}
}

// hasPrefix checks if the passed invoke starts with one of the passed
// prefixes or a mention of the bot.
// If so it the invoke stripped of the prefix.
// Otherwise it returns an empty string.
func (b *Bot) hasPrefix(invoke string, prefixes []string, guild bool) string {
	indexes := b.selfMentionRegexp.FindStringIndex(invoke)
	if indexes != nil {
		return strings.TrimLeft(invoke[indexes[1]:], whitespace)
	}

	for _, p := range prefixes {
		if strings.HasPrefix(invoke, p) {
			return strings.TrimLeft(invoke[len(p):], whitespace)
		}
	}

	if guild {
		return ""
	}

	return invoke // prefix isn't required in direct messages, so this is valid
}

func (b *Bot) findCommand(
	invoke string, base *state.Base, msg *discord.Message,
) (*plugin.ResolvedCommand, plugin.Provider, string) {
	ctxprovider := &ctxPluginProvider{
		base: base,
		msg:  msg,
		repos: []plugin.Repository{
			{
				ProviderName: plugin.BuiltInProvider,
				Commands:     b.commands,
				Modules:      b.modules,
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
			ctxprovider.unavailableProviders = append(ctxprovider.unavailableProviders,
				plugin.UnavailablePluginProvider{
					Name:  p.name,
					Error: err,
				})
		} else {
			repo := plugin.Repository{
				ProviderName: p.name,
				Modules:      mods,
				Commands:     cmds,
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

func (b *Bot) findCommandAsync(
	invoke string, base *state.Base, msg *discord.Message,
) (*plugin.ResolvedCommand, plugin.Provider, string) {
	ctxprovider := &ctxPluginProvider{
		base: base,
		msg:  msg,
		repos: []plugin.Repository{
			{
				ProviderName: plugin.BuiltInProvider,
				Commands:     b.commands,
				Modules:      b.modules,
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

func (b *Bot) applyMiddlewares(ctx *plugin.Context) CommandFunc {
	middlewares := b.Middlewares()

	for _, mod := range ctx.InvokedCommand.SourceParents {
		if m, ok := mod.(Middlewarer); ok && m != nil {
			middlewares = append(middlewares, m.Middlewares()...)
		}
	}

	if m, ok := ctx.InvokedCommand.Source.(Middlewarer); ok && m != nil {
		middlewares = append(middlewares, m.Middlewares()...)
	}

	middlewares = append(middlewares, b.postMiddlewares.Middlewares()...)

	inv := b.invoke

	for i := len(middlewares) - 1; i >= 0; i-- {
		inv = middlewares[i](inv)
	}

	return inv
}

func (b *Bot) invoke(_ *state.State, ctx *plugin.Context) error {
	reply, err := ctx.InvokedCommand.Invoke(b.State, ctx)
	if err != nil {
		// special case, prevent this from going through as an *InternalError
		if discorderr.Is(discorderr.As(err), discorderr.InsufficientPermissions) {
			err = plugin.DefaultBotPermissionsError
		}

		return err
	}

	return b.sendReply(reply, ctx)
}

func (b *Bot) sendReply(reply interface{}, ctx *plugin.Context) (err error) { //nolint:funlen,gocyclo
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
		if len(reply) > 0 {
			_, err = ctx.Reply(reply)
		}
	case discord.Embed:
		_, err = ctx.ReplyEmbed(reply)
		if discorderr.Is(discorderr.As(err), discorderr.CannotSendEmptyMessage) {
			err = nil
		}
	case *discord.Embed:
		if reply != nil {
			_, err = ctx.ReplyEmbed(*reply)
			if discorderr.Is(discorderr.As(err), discorderr.CannotSendEmptyMessage) {
				err = nil
			}
		}
	case *embedutil.Builder:
		if reply != nil {
			_, err = ctx.ReplyEmbedBuilder(reply)
			if discorderr.Is(discorderr.As(err), discorderr.CannotSendEmptyMessage) {
				err = nil
			}
		}
	case api.SendMessageData:
		_, err = ctx.ReplyMessage(reply)
		if discorderr.Is(discorderr.As(err), discorderr.CannotSendEmptyMessage) {
			err = nil
		}
	case i18n.Term:
		if len(reply) > 0 {
			_, err = ctx.Replylt(reply)
		}
	case *i18n.Config:
		if reply != nil {
			_, err = ctx.Replyl(reply)
		}
	case plugin.Reply:
		if reply != nil {
			err = reply.SendReply(b.State, ctx)
		}
	default:
		err = &ReplyTypeError{Reply: reply}
	}

	return err
}
