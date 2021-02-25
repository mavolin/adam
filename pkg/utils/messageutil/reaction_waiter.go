package messageutil

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/discorderr"
)

type (
	// A ReactionWaiter is used to await reactions.
	// Wait can either be cancelled by the user through a cancel reaction, or
	// by the ReactionWaiter if the timeout expires.
	ReactionWaiter struct {
		state *state.State
		ctx   *plugin.Context

		userID    discord.UserID
		messageID discord.MessageID
		channelID discord.ChannelID

		reactions, cancelReactions []discord.APIEmoji
		noAutoReact                bool
		noAutoDelete               bool

		middlewares []interface{}
	}
)

// NewReactionWaiter creates a new ReactionWaiter using the passed state.State
// and plugin.Context.
// ctx.Author will be assumed as the user to make the reaction in
// ctx.ChannelID.
func NewReactionWaiter(s *state.State, ctx *plugin.Context, messageID discord.MessageID) *ReactionWaiter {
	return &ReactionWaiter{
		state:     s,
		ctx:       ctx,
		userID:    ctx.Author.ID,
		messageID: messageID,
		channelID: ctx.ChannelID,
	}
}

// WithUser changes the user that is expected to react to the user with the
// passed id.
func (w *ReactionWaiter) WithUser(id discord.UserID) *ReactionWaiter {
	w.userID = id
	return w
}

// InChannel changes the channel id to the passed discord.ChannelID.
func (w *ReactionWaiter) InChannel(id discord.ChannelID) *ReactionWaiter {
	w.channelID = id
	return w
}

// WithReactions adds the passed reaction to the wait list.
func (w *ReactionWaiter) WithReactions(reactions ...discord.APIEmoji) *ReactionWaiter {
	w.reactions = append(w.reactions, reactions...)
	return w
}

// WithCancelReactions adds the passed cancel reactions.
// If the user reacts with one of the passed emojis, AwaitReply will return
// errors.Abort.
func (w *ReactionWaiter) WithCancelReactions(reactions ...discord.APIEmoji) *ReactionWaiter {
	w.cancelReactions = append(w.cancelReactions, reactions...)
	return w
}

// NoAutoReact disables automatic reaction and deletion of the reactions.
func (w *ReactionWaiter) NoAutoReact() *ReactionWaiter {
	w.noAutoReact = true
	w.noAutoDelete = true
	return w
}

// NoAutoDelete disables the automatic deletion of the reactions.
func (w *ReactionWaiter) NoAutoDelete() *ReactionWaiter {
	w.noAutoDelete = true
	return w
}

// WithMiddlewares adds the passed middlewares to the waiter.
// All middlewares of invalid type will be discarded.
//
// The following types are permitted:
//		• func(*state.State, interface{})
//		• func(*state.State, interface{}) error
//		• func(*state.State, *state.Base)
//		• func(*state.State, *state.Base) error
//		• func(*state.State, *state.MessageReactionAddEvent)
//		• func(*state.State, *state.MessageReactionAddEvent) error
func (w *ReactionWaiter) WithMiddlewares(middlewares ...interface{}) *ReactionWaiter {
	if len(w.middlewares) == 0 {
		w.middlewares = make([]interface{}, 0, len(middlewares))
	}

	for _, m := range middlewares {
		switch m.(type) { // check that the middleware is of a valid type
		case func(*state.State, interface{}):
		case func(*state.State, interface{}) error:
		case func(*state.State, *state.Base):
		case func(*state.State, *state.Base) error:
		case func(*state.State, *state.MessageReactionAddEvent):
		case func(*state.State, *state.MessageReactionAddEvent) error:
		default:
			continue
		}

		w.middlewares = append(w.middlewares, m)
	}

	return w
}

// Clone creates a deep copy of the ReactionWaiter.
func (w *ReactionWaiter) Clone() (cp *ReactionWaiter) {
	cp = &ReactionWaiter{
		noAutoReact: w.noAutoReact,
	}

	cp.reactions = make([]discord.APIEmoji, len(w.reactions))
	copy(cp.reactions, w.reactions)

	cp.cancelReactions = make([]discord.APIEmoji, len(w.cancelReactions))
	copy(cp.cancelReactions, w.cancelReactions)

	cp.middlewares = make([]interface{}, len(w.middlewares))
	copy(cp.middlewares, w.middlewares)

	return
}

// Await is the same as AwaitWithContext, but uses a context.Background() as
// context.Context.
func (w *ReactionWaiter) Await(timeout time.Duration) (discord.APIEmoji, error) {
	return w.AwaitWithContext(context.Background(), timeout)
}

// AwaitWithContext awaits a reaction of the user until they signal
// cancellation, the timeout expires or the context expires.
//
// If the timeout is reached, a *TimeoutError will be returned.
// If the user cancels the wait or deletes the message, errors.Abort will be
// returned.
// Furthermore, if the guild, channel or message becomes unavailable while
// adding reactions, errors.Abort will be returned as well.
// If the context expires, a *TimeoutError with Cause set to ctx.Err() will be
// returned.
// This error is also available through .Unwrap(), so errors.Is is safe to use.
//
// Besides that, the wait can also be canceled through a middleware.
func (w *ReactionWaiter) AwaitWithContext(ctx context.Context, timeout time.Duration) (discord.APIEmoji, error) {
	perms, err := w.ctx.SelfPermissions()
	if err != nil {
		return "", err
	}

	if !perms.Has(discord.PermissionAddReactions) {
		return "", plugin.NewBotPermissionsError(discord.PermissionAddReactions)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	result := make(chan interface{})

	reactCleanup, err := w.handleReactions(ctx, result)
	if err != nil {
		return "", err
	}

	defer reactCleanup()

	select {
	case <-ctx.Done():
		return "", &TimeoutError{UserID: w.userID, Cause: ctx.Err()}
	case <-time.After(timeout):
		return "", &TimeoutError{UserID: w.userID}
	case r := <-result:
		switch r := r.(type) {
		case discord.APIEmoji:
			return r, nil
		case error:
			return "", r
		default: // this should never happen
			return "", errors.NewWithStack("messageutil: unexpected return value of result channel")
		}
	}
}

func (w *ReactionWaiter) handleReactions(ctx context.Context, result chan<- interface{}) (func(), error) { //nolint:gocognit,funlen
	rmReact := w.state.MustAddHandler(func(s *state.State, e *state.MessageReactionAddEvent) {
		if e.UserID != w.userID || e.MessageID != w.messageID {
			return
		}

		if err := invokeReactionAddMiddlewares(s, e, w.middlewares); err != nil {
			sendResult(ctx, result, err)
			return
		}

		for _, r := range w.reactions {
			if e.Emoji.APIString() == r {
				sendResult(ctx, result, r)
				return
			}
		}

		for _, r := range w.cancelReactions {
			if e.Emoji.APIString() == r {
				sendResult(ctx, result, errors.Abort)
				return
			}
		}
	})

	rmMsgDel := w.state.MustAddHandler(func(s *state.State, e *state.MessageDeleteEvent) {
		if e.ID == w.messageID {
			sendResult(ctx, result, errors.Abort)
		}
	})

	if !w.noAutoReact {
		for _, r := range w.reactions {
			if err := w.state.React(w.channelID, w.messageID, r); err != nil {
				// someone deleted the channel or message
				if discorderr.Is(discorderr.As(err), discorderr.UnknownResource...) {
					rmReact()
					rmMsgDel()
					return nil, errors.Abort
				}

				w.ctx.HandleErrorSilently(err)
			}
		}

		for _, r := range w.cancelReactions {
			if err := w.state.React(w.channelID, w.messageID, r); err != nil {
				// someone deleted the channel or message
				if discorderr.Is(discorderr.As(err), discorderr.UnknownResource...) {
					rmReact()
					rmMsgDel()
					return nil, errors.Abort
				}

				w.ctx.HandleErrorSilently(err)
			}
		}
	}

	return func() {
		rmReact()
		rmMsgDel()

		if !w.noAutoDelete {
			go func() {
				for _, r := range w.reactions {
					err := w.state.DeleteReactions(w.channelID, w.messageID, r)
					if err != nil {
						// someone else deleted the resource we are accessing
						if discorderr.Is(discorderr.As(err), discorderr.UnknownResource...) {
							return
						}

						w.ctx.HandleErrorSilently(err)
					}
				}

				for _, r := range w.cancelReactions {
					err := w.state.DeleteReactions(w.channelID, w.messageID, r)
					if err != nil {
						// someone else deleted the resource we are accessing
						if discorderr.Is(discorderr.As(err), discorderr.UnknownResource...) {
							return
						}

						w.ctx.HandleErrorSilently(err)
					}
				}
			}()
		}
	}, nil
}
