package reactawait

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/discorderr"
)

// A Waiter is used to await reactions.
// Wait can either be cancelled by the user through a cancel reaction, or
// by the Waiter if the timeout expires.
type Waiter struct {
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

// New creates a new Waiter using the passed state.State and
// plugin.Context.
// ctx.Author will be assumed as the user to make the reaction in
// ctx.ChannelID.
func New(s *state.State, ctx *plugin.Context, messageID discord.MessageID) *Waiter {
	return &Waiter{
		state:     s,
		ctx:       ctx,
		userID:    ctx.Author.ID,
		messageID: messageID,
		channelID: ctx.ChannelID,
	}
}

// WithUser changes the user that is expected to react to the user with the
// passed id.
func (w *Waiter) WithUser(id discord.UserID) *Waiter {
	w.userID = id
	return w
}

// InChannel changes the channel id to the passed discord.ChannelID.
func (w *Waiter) InChannel(id discord.ChannelID) *Waiter {
	w.channelID = id
	return w
}

// WithReactions adds the passed reaction to the wait list.
func (w *Waiter) WithReactions(reactions ...discord.APIEmoji) *Waiter {
	w.reactions = append(w.reactions, reactions...)
	return w
}

// WithCancelReactions adds the passed cancel reactions.
// If the user reacts with one of the passed emojis, AwaitReply will return
// errors.Abort.
func (w *Waiter) WithCancelReactions(reactions ...discord.APIEmoji) *Waiter {
	w.cancelReactions = append(w.cancelReactions, reactions...)
	return w
}

// NoAutoReact disables automatic reaction and deletion of the reactions.
func (w *Waiter) NoAutoReact() *Waiter {
	w.noAutoReact = true
	w.noAutoDelete = true
	return w
}

// NoAutoDelete disables the automatic deletion of the reactions.
func (w *Waiter) NoAutoDelete() *Waiter {
	w.noAutoDelete = true
	return w
}

// WithMiddlewares adds the passed middlewares to the waiter.
// All middlewares of invalid type will be discarded.
//
// The following types are permitted:
// 	• func(*state.State, interface{})
//	• func(*state.State, interface{}) error
//	• func(*state.State, *event.Base)
//	• func(*state.State, *event.Base) error
//	• func(*state.State, *state.MessageReactionAddEvent)
//	• func(*state.State, *state.MessageReactionAddEvent) error
func (w *Waiter) WithMiddlewares(middlewares ...interface{}) *Waiter {
	if len(w.middlewares) == 0 {
		w.middlewares = make([]interface{}, 0, len(middlewares))
	}

	for _, m := range middlewares {
		switch m.(type) { // check that the middleware is of a valid type
		case func(*state.State, interface{}):
		case func(*state.State, interface{}) error:
		case func(*state.State, *event.Base):
		case func(*state.State, *event.Base) error:
		case func(*state.State, *event.MessageReactionAdd):
		case func(*state.State, *event.MessageReactionAdd) error:
		default:
			continue
		}

		w.middlewares = append(w.middlewares, m)
	}

	return w
}

// Clone creates a deep copy of the Waiter.
func (w *Waiter) Clone() (cp *Waiter) {
	cp = &Waiter{
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

// Await creates a context.Context with the given timeout and calls
// AwaitContext with it.
func (w *Waiter) Await(timeout time.Duration) (discord.APIEmoji, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return w.AwaitContext(ctx)
}

// AwaitContext awaits a reaction of the user until they signal cancellation, the
// timeout expires or the context expires.
//
// If the user cancels the wait or deletes the message, errors.Abort will be
// returned.
// Furthermore, if the guild, channel or message becomes unavailable while
// adding reactions, errors.Abort will be returned as well.
// If the context expires, a *TimeoutError with Cause set to ctx.Err() will be
// returned.
//
// Besides that, the wait can also be canceled through a middleware.
func (w *Waiter) AwaitContext(ctx context.Context) (discord.APIEmoji, error) {
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
	case r := <-result:
		switch r := r.(type) {
		case discord.APIEmoji:
			return r, nil
		case error:
			return "", r
		default: // this should never happen
			return "", errors.NewWithStack("reactawait: unexpected return value of result channel")
		}
	}
}

//nolint:gocognit,funlen
func (w *Waiter) handleReactions(ctx context.Context, result chan<- interface{}) (func(), error) {
	rmReact := w.state.AddHandler(func(s *state.State, e *event.MessageReactionAdd) {
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

	rmMsgDel := w.state.AddHandler(func(s *state.State, e *event.MessageDelete) {
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
