package bot

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/replier"
	"github.com/mavolin/adam/pkg/plugin"
)

// ErrUnknownCommand is the error used if a message with a matching prefix
// does not contain a valid command invoke.
var ErrUnknownCommand = errors.NewUserErrorl(unknownCommandErrorDescription)

// Route attempts to route the passed message.
// It aborts if the message is not a valid invoke.
//
// When calling the bot's middlewares, it guarantees that Message, Member,
// Base, BotOwnerIDs, Replier, Provider, DiscordDataProvider, and ErrorHandler
// are set.
// Further, Localizer will be set to a fallback localizer.
//nolint:funlen
func (b *Bot) Route(base *event.Base, msg *discord.Message, member *discord.Member) {
	// discard the message if THIS bot wrote it, even if b.AllowBot
	if msg.Author.ID == b.selfID {
		return
	}

	if member != nil {
		member.User = msg.Author
	}

	ctx := &plugin.Context{
		Message:     *msg,
		Member:      member,
		Base:        base,
		Localizer:   i18n.NewFallbackLocalizer(),
		BotOwnerIDs: b.Owners,
		Replier:     replier.WrapState(b.State, false),
		Provider:    b.pluginResolver.NewProvider(base, msg),
		DiscordDataProvider: &discordDataProvider{
			s:         b.State,
			guildID:   msg.GuildID,
			channelID: msg.ChannelID,
			selfID:    b.selfID,
		},
	}
	ctx.ErrorHandler = newCtxErrorHandler(b.State, ctx, b.ErrorHandler)

	defer func() {
		if rec := recover(); rec != nil {
			b.PanicHandler(rec, b.State, ctx)
		}
	}()

	inv := b.applyMiddlewares()
	if err := inv(b.State, ctx); err != nil {
		b.ErrorHandler(err, b.State, ctx)
	}
}

func (b *Bot) applyMiddlewares() CommandFunc {
	middlewares := b.Middlewares()

	middlewares = append(middlewares, func(next CommandFunc) CommandFunc {
		return func(s *state.State, ctx *plugin.Context) error {
			var middlewares []Middleware

			for _, mod := range ctx.InvokedCommand.SourceParents() {
				if m, ok := mod.(Middlewarer); ok && m != nil {
					middlewares = append(middlewares, m.Middlewares()...)
				}
			}

			if m, ok := ctx.InvokedCommand.Source().(Middlewarer); ok && m != nil {
				middlewares = append(middlewares, m.Middlewares()...)
			}

			middlewares = append(middlewares, b.postMiddlewares.Middlewares()...)

			for i := len(middlewares) - 1; i >= 0; i-- {
				next = middlewares[i](next)
			}

			return next(s, ctx)
		}
	})

	f := func(*state.State, *plugin.Context) error { return nil }

	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}

	return f
}
