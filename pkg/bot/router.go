package bot

import (
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/replier"
	"github.com/mavolin/adam/pkg/plugin"
)

var ErrUnknownCommand = errors.NewUserErrorl(unknownCommandErrorDescription)

// Route attempts to route the passed message.
// It aborts if the message is not a valid invoke.
func (b *Bot) Route(base *state.Base, msg *discord.Message, member *discord.Member) { //nolint:funlen
	// Only accept regular text messages.
	// Also check if a bot wrote the message, if !b.AllowBot.
	// Lastly, also discard if this bot wrote this message, even if b.AllowBot.
	if msg.Type != discord.DefaultMessage || (!b.AllowBot && msg.Author.Bot) || msg.Author.ID == b.selfID {
		return
	}

	ctx := &plugin.Context{
		Message: *msg,
		Member:  member,
		Base:    base,
	}

	if !b.checkPrefix(ctx) {
		return
	}

	ctx.BotOwnerIDs = b.Owners
	ctx.Replier = replier.WrapState(b.State, false)
	ctx.Provider = b.pluginResolver.NewProvider(base, msg)
	ctx.DiscordDataProvider = &discordDataProvider{
		s:         b.State,
		guildID:   msg.GuildID,
		channelID: msg.ChannelID,
		selfID:    b.selfID,
	}

	ctx.ErrorHandler = newCtxErrorHandler(b.State, ctx, b.ErrorHandler)

	var ok bool

	ctx.Prefixes, ctx.Localizer, ok = b.SettingsProvider(base, msg)
	if !ok {
		return
	}

	if ctx.Localizer == nil {
		ctx.Localizer = i18n.NewFallbackLocalizer()
	}

	if member != nil {
		member.User = msg.Author
	}

	cmd, rawArgs := ctx.FindCommandWithArgs(ctx.Content[ctx.InvokeIndex:])
	if cmd == nil {
		ctx.HandleError(ErrUnknownCommand)
		return
	}

	ctx.InvokedCommand = cmd
	ctx.ArgsIndex = len(ctx.Content) - len(rawArgs)

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

func (b *Bot) checkPrefix(ctx *plugin.Context) bool {
	indexes := b.selfMentionRegexp.FindStringIndex(ctx.Content)
	if indexes != nil { // invoked by mention
		ctx.InvokeIndex = len(ctx.Content) - len(strings.TrimLeft(ctx.Content[indexes[1]:], shared.Whitespace))
		return true
	}

	for _, p := range ctx.Prefixes {
		if strings.HasPrefix(ctx.Content, p) {
			ctx.InvokeIndex = len(ctx.Content) - len(strings.TrimLeft(ctx.Content[len(p):], shared.Whitespace))
			return true
		}
	}

	// prefix isn't required in direct messages, so DM's always "match"
	return !ctx.GuildID.IsValid()
}

func (b *Bot) applyMiddlewares(ctx *plugin.Context) CommandFunc {
	middlewares := b.Middlewares()

	for _, mod := range ctx.InvokedCommand.SourceParents() {
		if m, ok := mod.(Middlewarer); ok && m != nil {
			middlewares = append(middlewares, m.Middlewares()...)
		}
	}

	if m, ok := ctx.InvokedCommand.Source().(Middlewarer); ok && m != nil {
		middlewares = append(middlewares, m.Middlewares()...)
	}

	middlewares = append(middlewares, b.postMiddlewares.Middlewares()...)

	f := func(*state.State, *plugin.Context) error { return nil }

	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}

	return f
}
