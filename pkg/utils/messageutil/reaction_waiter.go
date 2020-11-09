package messageutil

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

type (
	// A ReactionWaiter is used to await reactions.
	// Wait can either be cancelled by the user through a cancel reaction, or
	// by the ReactionWaiter if the timeout expires.
	ReactionWaiter struct {
		state *state.State
		ctx   *plugin.Context

		userID    discord.UserID
		channelID discord.ChannelID

		reactions, cancelReactions []reaction
		noAutoReact                bool

		middlewares []interface{}
	}

	reaction struct {
		messageID discord.MessageID
		reaction  api.Emoji
	}
)

// NewReactionWaiter creates a new ReactionWaiter using the passed state.State
// and plugin.Context.
// ctx.Author will be assumed as the user to make the reaction in
// ctx.ChannelID.
func NewReactionWaiter(s *state.State, ctx *plugin.Context) *ReactionWaiter {
	return &ReactionWaiter{
		state:     s,
		ctx:       ctx,
		userID:    ctx.Author.ID,
		channelID: ctx.ChannelID,
	}
}

// WithUser changes the user that is expected to react to the user with the
// passed id.
func (w *ReactionWaiter) WithUser(id discord.UserID) *ReactionWaiter {
	w.userID = id
	return w
}

func (w *ReactionWaiter) InChannel(id discord.ChannelID) *ReactionWaiter {
	w.channelID = id
	return w
}

// WithReaction adds the passed reaction to the wait list.
func (w *ReactionWaiter) WithReaction(messageID discord.MessageID, react api.Emoji) *ReactionWaiter {
	w.reactions = append(w.reactions, reaction{
		messageID: messageID,
		reaction:  react,
	})

	return w
}

// NoAutoReact disables automatic reaction.
func (w *ReactionWaiter) NoAutoReact() *ReactionWaiter {
	w.noAutoReact = true
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
func (w *ReactionWaiter) WithMiddleware(middlewares ...interface{}) *ReactionWaiter {
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

// WithCancelReaction adds the passed cancel reaction.
// If the user reacts with the passed emoji, AwaitReply will return with error
// Canceled.
func (w *ReactionWaiter) WithCancelReaction(messageID discord.MessageID, react api.Emoji) *ReactionWaiter {
	w.cancelReactions = append(w.cancelReactions, reaction{
		messageID: messageID,
		reaction:  react,
	})

	return w
}

// Copy creates a copy of the ReactionWaiter.
func (w *ReactionWaiter) Copy() (cp *ReactionWaiter) {
	cp = &ReactionWaiter{
		noAutoReact: w.noAutoReact,
	}

	cp.reactions = make([]reaction, len(w.reactions))
	copy(cp.reactions, w.reactions)

	cp.cancelReactions = make([]reaction, len(w.cancelReactions))
	copy(cp.cancelReactions, w.cancelReactions)

	cp.middlewares = make([]interface{}, len(w.middlewares))
	copy(cp.middlewares, w.middlewares)

	return
}

// AwaitReply awaits a reaction of the user until they signal cancellation or the
// timeout expires.
//
// If the timeout is reached, a *TimeoutError will be returned.
// If the user cancels the wait, Canceled will be returned.
//
// Besides that, the Wait can also be canceled through a middleware.
// If one the middlewares returns state.Filtered, errors.Abort will be
// returned.
func (w *ReactionWaiter) Await(timeout time.Duration) (api.Emoji, error) {
	return w.AwaitWithContext(context.Background(), timeout)
}

// AwaitWithContext awaits a reaction of the user until they signal
// cancellation, the timeout expires or the context expires.
//
// If the timeout is reached, a *TimeoutError will be returned.
// If the user cancels the wait, Canceled will be returned.
//
// Besides that, the Wait can also be canceled through a middleware.
// If one the middlewares returns state.Filtered, errors.Abort will be
// returned.
func (w *ReactionWaiter) AwaitWithContext(ctx context.Context, timeout time.Duration) (api.Emoji, error) {
	perms, err := w.ctx.SelfPermissions()
	if err != nil {
		return "", err
	}

	if !perms.Has(discord.PermissionAddReactions) {
		return "", errors.NewInsufficientBotPermissionsError(discord.PermissionAddReactions)
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
		return "", ctx.Err()
	case <-time.After(timeout):
		return "", &TimeoutError{UserID: w.ctx.Author.ID}
	case r := <-result:
		switch r := r.(type) {
		case api.Emoji:
			return r, nil
		case error:
			return "", r
		default: // this should never happen
			return "", errors.NewWithStack("messageutil: unexpected return value of result channel")
		}
	}
}

func (w *ReactionWaiter) handleReactions(ctx context.Context, result chan<- interface{}) (func(), error) {
	if !w.noAutoReact {
		for _, r := range w.reactions {
			if err := w.state.React(w.ctx.ChannelID, r.messageID, r.reaction); err != nil {
				w.ctx.HandleErrorSilent(err)
			}
		}

		for _, r := range w.cancelReactions {
			if err := w.state.React(w.ctx.ChannelID, r.messageID, r.reaction); err != nil {
				w.ctx.HandleErrorSilent(err)
			}
		}
	}

	rm, err := w.state.AddHandler(func(s *state.State, e *state.MessageReactionAddEvent) {
		if e.UserID != w.userID {
			return
		}

		if err := invokeReactionAddMiddlewares(s, e, w.middlewares); err != nil {
			sendResult(ctx, result, err)
			return
		}

		for _, r := range w.reactions {
			if e.MessageID == r.messageID && e.Emoji.APIString() == r.reaction {
				sendResult(ctx, result, r.reaction)
				return
			}
		}

		for _, r := range w.cancelReactions {
			if e.MessageID == r.messageID && e.Emoji.APIString() == r.reaction {
				sendResult(ctx, result, Canceled)
				return
			}
		}
	})
	if err != nil { // this should never happen
		return nil, errors.WithStack(err)
	}

	return func() {
		rm()

		if !w.noAutoReact {
			go func() {
				for _, r := range w.reactions {
					err := w.state.DeleteReactions(w.ctx.ChannelID, r.messageID, r.reaction)
					if err != nil {
						w.ctx.HandleErrorSilent(err)
					}
				}

				for _, r := range w.cancelReactions {
					err := w.state.DeleteReactions(w.ctx.ChannelID, r.messageID, r.reaction)
					if err != nil {
						w.ctx.HandleErrorSilent(err)
					}
				}
			}()
		}
	}, nil
}
