package reply

import (
	"context"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
)

// Await awaits a reply of the user until the user signals cancellation, the
// initial timeout expires and the user is not typing or the user stops typing
// and the typing timeout is reached. Note you need the typing intent to
// monitor typing.
//
// If one of the timeouts is reached, a *errors.UserInfo containing a timeout
// info message will be returned.
// If the user cancels the reply, Canceled will be returned.
//
// Besides that, a reply can also be canceled through a middleware.
// If one the middlewares returns state.Filtered, errors.Abort will be
// returned.
func (w *Waiter) Await(initialTimeout, typingTimeout time.Duration) (*discord.Message, error) {
	return w.AwaitWithContext(context.Background(), initialTimeout, typingTimeout)
}

// Await awaits a reply of the user until the user signals cancellation, the
// initial timeout expires and the user is not typing or the user stops typing
// and the typing timeout is reached. Note you need the typing intent to
// monitor typing.
//
// If one of the timeouts is reached, a *errors.UserInfo containing a timeout
// info message will be returned.
// If the user cancels the reply, Canceled will be returned.
// If the context expires, context.Canceled will be returned.
//
// Besides that, a reply can also be canceled through a middleware.
// If one the middlewares returns state.Filtered, errors.Abort will be
// returned.
func (w *Waiter) AwaitWithContext(
	ctx context.Context, initialTimeout, typingTimeout time.Duration,
) (*discord.Message, error) {
	perms, err := w.ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	// make sure we have permission to send messages and create reactions, if
	// time extensions are enabled or we have cancel reactions.
	if !perms.Has(discord.PermissionSendMessages) {
		return nil, errors.NewInsufficientBotPermissionsError(discord.PermissionSendMessages)
	} else if !w.noAutoReact && len(w.cancelReactions) > 0 && !perms.Has(discord.PermissionAddReactions) {
		return nil, errors.NewInsufficientBotPermissionsError(discord.PermissionAddReactions)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	result := make(chan interface{})

	awaitCleanup, err := w.handleMessages(ctx, result)
	if err != nil {
		return nil, err
	}

	defer awaitCleanup()

	if !w.noAutoReact && len(w.cancelReactions) > 0 {
		reactCleanup, err := w.handleCancelReactions(ctx, result)
		if err != nil {
			return nil, err
		}

		defer reactCleanup()
	}

	timeoutCleanup, err := w.watchTimeout(ctx, initialTimeout, typingTimeout, result)
	if err != nil {
		return nil, err
	}

	defer timeoutCleanup()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-result:
		switch r := r.(type) {
		case *discord.Message:
			return r, nil
		case error:
			return nil, r
		default: // this should never happen
			return nil, errors.NewWithStack("reply: unexpected return value of result channel")
		}
	}
}

func (w *Waiter) handleMessages(ctx context.Context, result chan<- interface{}) (func(), error) {
	rm, err := w.state.AddHandler(func(s *state.State, e *state.MessageCreateEvent) {
		if e.ChannelID != w.ctx.ChannelID || e.Author.ID != w.ctx.Author.ID { // not the message we are waiting for
			return
		}

		if err := invokeMiddlewares(s, e, w.middlewares); err != nil {
			sendResult(ctx, result, err)
			return
		}

		// check if the message is a cancel keyword
		for _, kt := range w.cancelKeywords {
			k, err := kt.Get(w.ctx.Localizer)
			if err != nil {
				w.ctx.HandleErrorSilent(err)
				continue
			}

			if (w.caseSensitive && k == e.Content) || (!w.caseSensitive && strings.EqualFold(k, e.Content)) {
				sendResult(ctx, result, Canceled)
				return
			}
		}

		sendResult(ctx, result, &e.Message)
	})

	return rm, errors.WithStack(err)
}

func (w *Waiter) handleCancelReactions(ctx context.Context, result chan<- interface{}) (func(), error) {
	for _, r := range w.cancelReactions {
		if err := w.state.React(w.ctx.ChannelID, r.messageID, r.reaction); err != nil {
			w.ctx.HandleErrorSilent(err)
		}
	}

	rm, err := w.state.AddHandler(func(s *state.State, e *state.MessageReactionAddEvent) {
		for _, r := range w.cancelReactions {
			if e.MessageID == r.messageID && e.Emoji.APIString() == r.reaction && e.UserID == w.ctx.Author.ID {
				select {
				case result <- Canceled:
				case <-ctx.Done():
				}
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

func (w *Waiter) watchTimeout(
	ctx context.Context, initialTimeout, typingTimeout time.Duration, result chan<- interface{},
) (rm func(), err error) {
	maxTimer := time.NewTimer(w.maxTimeout)
	typing := make(chan struct{})

	rm, err = w.state.AddHandler(func(s *state.State, e *state.TypingStartEvent) {
		if e.ChannelID != w.ctx.ChannelID || e.UserID != w.ctx.Author.ID {
			return
		}

		select {
		case typing <- struct{}{}:
		case <-ctx.Done():
		}
	})
	if err != nil {
		return
	}

	t := time.NewTimer(initialTimeout)

	go func() {
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-typing:
				if !t.Stop() {
					<-t.C
				}

				t.Reset(typingTimeout)
			case <-t.C:
				maxTimer.Stop()

				result <- errors.NewUserInfol(timeoutInfo.
					WithPlaceholders(timeoutInfoPlaceholders{
						ResponseUserMention: w.ctx.Author.Mention(),
					}))
				return
			case <-maxTimer.C:
				t.Stop()

				result <- errors.NewUserInfol(timeoutInfo.
					WithPlaceholders(timeoutInfoPlaceholders{
						ResponseUserMention: w.ctx.Author.Mention(),
					}))
				return
			}
		}
	}()

	return
}
