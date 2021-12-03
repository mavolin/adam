package bot

import (
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/diamondburned/arikawa/v3/api"

	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/discorderr"
	"github.com/mavolin/adam/pkg/utils/msgbuilder"
	"github.com/mavolin/adam/pkg/utils/permutil"
)

// SendTyping sends a typing event every 6 seconds until the command finishes
// executing.
func SendTyping(next CommandFunc) CommandFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if !ctx.InvokedCommand.BotPermissions().Has(discord.PermissionSendMessages) {
			return next(s, ctx)
		}

		stop := make(chan struct{})
		defer close(stop)

		go func() {
			t := time.NewTicker(6 * time.Second)

			err := s.Typing(ctx.ChannelID)
			if err != nil {
				ctx.HandleErrorSilently(err)
			}

			for {
				select {
				case <-stop:
					t.Stop()
					return
				case <-t.C:
					err := s.Typing(ctx.ChannelID)
					if err != nil {
						ctx.HandleErrorSilently(err)
					}
				}
			}
		}()

		return next(s, ctx)
	}
}

// CheckMessageType checks if the invoking message is of type
// discord.DefaultMessage.
// If so it calls the next middleware, otherwise it aborts with an
// *errors.InformationalError.
func CheckMessageType(next CommandFunc) CommandFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if ctx.Type != discord.DefaultMessage {
			return errors.NewInformationalError("bot: router: message is not of type default message")
		}

		return next(s, ctx)
	}
}

// CheckHuman checks if the invoking message was written by a human.
// If so it calls the next middleware, otherwise it aborts with an
// *errors.InformationalError.
func CheckHuman(next CommandFunc) CommandFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if ctx.Author.Bot {
			return errors.NewInformationalError("bot: router: message was written by a bot")
		}

		return next(s, ctx)
	}
}

// NewSettingsRetriever creates a new settings retriever middleware that
// retrieves the settings from the passed SettingsProvider.
// If the settings provider returns !ok, the returned middleware will abort
// by returning an *errors.InformationalError.
//
// The returned middleware will set the Prefixes and Localizer context fields.
func NewSettingsRetriever(settingsProvider SettingsProvider) Middleware {
	return func(next CommandFunc) CommandFunc {
		return func(s *state.State, ctx *plugin.Context) error {
			var ok bool

			ctx.Prefixes, ctx.Localizer, ok = settingsProvider(ctx.Base, &ctx.Message)
			if !ok {
				return errors.NewInformationalError("bot: router: settings provider returned not ok")
			}

			if ctx.Localizer == nil {
				ctx.Localizer = i18n.NewFallbackLocalizer()
			}

			return next(s, ctx)
		}
	}
}

// CheckPrefix checks if the message starts with the prefix.
// The prefix must either be the mention of the bot, or one of the prefixes
// found in the context.
//
// Direct messages don't require prefixes, however, if a message starts with a
// prefix, it will still be stripped from the invoke.
//
// If the prefix doesn't match, an *errors.InformationalError will be returned.
//
// The middleware sets the ctx.InvokeIndex context field.
func CheckPrefix(next CommandFunc) CommandFunc {
	var selfMentionRegexp *regexp.Regexp
	var once sync.Once

	return func(s *state.State, ctx *plugin.Context) (err error) {
		once.Do(func() {
			var self *discord.User

			self, err = s.Me()
			if err != nil {
				err = errors.WithStack(err)
				return
			}

			selfMentionRegexp = regexp.MustCompile("^<@!?" + self.ID.String() + ">")
		})
		if err != nil {
			return err
		}

		indexes := selfMentionRegexp.FindStringIndex(ctx.Content)
		if indexes != nil { // invoked by mention
			ctx.InvokeIndex = len(ctx.Content) - len(strings.TrimLeft(ctx.Content[indexes[1]:], shared.Whitespace))
			return next(s, ctx)
		}

		for _, p := range ctx.Prefixes {
			if strings.HasPrefix(ctx.Content, p) {
				ctx.InvokeIndex = len(ctx.Content) - len(strings.TrimLeft(ctx.Content[len(p):], shared.Whitespace))
				return next(s, ctx)
			}
		}

		// prefixes aren't required in direct messages, so DMs always "match"
		if ctx.GuildID == 0 {
			return next(s, ctx)
		}

		return errors.NewInformationalError("bot: router: prefix does not match")
	}
}

// FindCommand attempts to find the command being invoked by the message.
// If no matching command is found, the middleware returns ErrUnknownCommand.
//
// The middleware sets the InvokedCommand and ArgsIndex context fields.
func FindCommand(next CommandFunc) CommandFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		cmd, rawArgs := ctx.FindCommandWithArgs(ctx.Content[ctx.InvokeIndex:])
		if cmd == nil {
			return ErrUnknownCommand
		}

		ctx.InvokedCommand = cmd
		ctx.ArgsIndex = len(ctx.Content) - len(rawArgs)

		return next(s, ctx)
	}
}

// CheckChannelTypes checks if the plugin.ChannelTypes of the command are
// satisfied.
func CheckChannelTypes(next CommandFunc) CommandFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		ok, err := ctx.InvokedCommand.ChannelTypes().Check(ctx)
		if err != nil {
			return err
		} else if !ok {
			return plugin.NewChannelTypeError(ctx.InvokedCommand.ChannelTypes())
		}

		return next(s, ctx)
	}
}

// CheckBotPermissions checks if the discord.Permissions the bot requires for
// the command are satisfied.
func CheckBotPermissions(next CommandFunc) CommandFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if ctx.InvokedCommand.BotPermissions() == 0 {
			return next(s, ctx)
		}

		if ctx.GuildID == 0 && !permutil.DMPermissions.Has(ctx.InvokedCommand.BotPermissions()) {
			return plugin.NewChannelTypeError(plugin.DirectMessages & ctx.InvokedCommand.ChannelTypes())
		} else if ctx.GuildID != 0 {
			p, err := ctx.SelfPermissions()
			if err != nil {
				return err
			}

			if !p.Has(ctx.InvokedCommand.BotPermissions()) {
				missing := (p & ctx.InvokedCommand.BotPermissions()) ^ ctx.InvokedCommand.BotPermissions()
				return plugin.NewBotPermissionsError(missing)
			}
		}

		return next(s, ctx)
	}
}

// NewThrottlerChecker creates a new bot.Middleware that checks if a
// command is being throttled.
// Additionally, it signals cancellation to the throttler.
func NewThrottlerChecker(cancelChecker func(err error) bool) Middleware {
	return func(next CommandFunc) CommandFunc {
		return func(s *state.State, ctx *plugin.Context) error {
			if ctx.InvokedCommand.Throttler() == nil {
				return next(s, ctx)
			}

			rm, err := ctx.InvokedCommand.Throttler().Check(s, ctx)
			if err != nil {
				return err
			}

			panicked := true

			// hacky way to check if we panicked, without repanicking and
			// losing part of the stack trace
			defer func() {
				if panicked {
					rm()
				}
			}()

			err = next(s, ctx)
			panicked = false

			if err != nil && cancelChecker(err) {
				rm()
			}

			return err
		}
	}
}

// CheckRestrictions checks if the command is restricted.
func CheckRestrictions(next CommandFunc) CommandFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		err := ctx.InvokedCommand.IsRestricted(s, ctx)
		if err != nil {
			return err
		}

		return next(s, ctx)
	}
}

// ParseArgs parses the ctx.RawArgs using the commands plugin.ArgConfig.
func ParseArgs(next CommandFunc) CommandFunc {
	return func(s *state.State, ctx *plugin.Context) (err error) {
		if ctx.InvokedCommand.Args() != nil {
			err = ctx.InvokedCommand.ArgParser().
				Parse(ctx.RawArgs(), ctx.InvokedCommand.Args(), s, ctx)
			if err != nil {
				return err
			}
		}

		return next(s, ctx)
	}
}

// InvokeCommand invokes the command and sends a reply, if the command returned
// one.
func InvokeCommand(next CommandFunc) CommandFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		reply, err := ctx.InvokedCommand.Invoke(s, ctx)
		if err != nil {
			// special case, prevent this from going through as an
			// *InternalError
			if discorderr.Is(discorderr.As(err), discorderr.InsufficientPermissions) {
				err = plugin.DefaultBotPermissionsError
			}

			return err
		}

		if err := sendReply(reply, ctx); err != nil {
			return err
		}

		return next(s, ctx)
	}
}

func sendReply(reply interface{}, ctx *plugin.Context) (err error) {
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
		if reply != "" {
			_, err = ctx.Reply(reply)
		}
	case discord.Embed:
		_, err = ctx.ReplyEmbeds(reply)
	case *discord.Embed:
		if reply != nil {
			_, err = ctx.ReplyEmbeds(*reply)
		}
	case *msgbuilder.Builder:
		if reply != nil {
			_, err = reply.Reply()
		}
	case *msgbuilder.EmbedBuilder:
		if reply != nil {
			_, err = msgbuilder.ReplyEmbedBuilders(ctx, reply)
		}
	case api.SendMessageData:
		_, err = ctx.ReplyMessage(reply)
	case i18n.Term:
		if len(reply) > 0 {
			_, err = ctx.Replyl(reply.AsConfig())
		}
	case *i18n.Config:
		if reply != nil {
			_, err = ctx.Replyl(reply)
		}
	default:
		err = &ReplyTypeError{Reply: reply}
	}

	if discorderr.Is(discorderr.As(err), discorderr.CannotSendEmptyMessage) {
		return nil
	}

	return err
}
