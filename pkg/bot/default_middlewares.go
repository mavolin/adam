package bot

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/diamondburned/arikawa/v2/api"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/discorderr"
	"github.com/mavolin/adam/pkg/utils/embedutil"
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

// NewThrottlerChecker creates a new bot.MiddlewareFunc that checks if a
// command is being throttled.
// Additionally, it signals cancellation to the throttler
func NewThrottlerChecker(cancelChecker func(err error) bool) MiddlewareFunc {
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
			// losing stack
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

		if err := sendReply(reply, s, ctx); err != nil {
			return err
		}

		return next(s, ctx)
	}
}

func sendReply(reply interface{}, s *state.State, ctx *plugin.Context) (err error) { //nolint:funlen,gocyclo
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
			err = reply.SendReply(s, ctx)
		}
	default:
		err = &ReplyTypeError{Reply: reply}
	}

	return err
}
