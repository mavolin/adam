package reply

import (
	"context"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
)

// Await awaits a reply of the user until the passed timout is reached.
// If they responds, the reply is returned.
//
// If the timeout passes and time extension are enabled, the timeout will be
// reset until the user responds, or the limit of of time extensions is
// reached, if set.
// Otherwise, a UserInfo containing a timeout info will be returned.
//
// If the user cancels the reply, Canceled will be returned.
//
// Besides that, a reply can also be canceled through a middleware.
// If one the middlewares returns state.Filtered, errors.Abort will be
// returned.
func (w *Waiter) Await(timeout time.Duration) (*discord.Message, error) {
	perms, err := w.ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	// make sure we have permission to send messages and create reactions, if
	// time extensions are enabled or we have cancel reactions.
	if !perms.Has(discord.PermissionSendMessages) {
		return nil, errors.NewInsufficientBotPermissionsError(discord.PermissionSendMessages)
	} else if (w.timeExtensions != 0 || (!w.noAutoReact && len(w.cancelReactions) > 0)) &&
		!perms.Has(discord.PermissionAddReactions) {
		return nil, errors.NewInsufficientBotPermissionsError(discord.PermissionAddReactions)
	}

	ctx, cancel := context.WithCancel(context.Background())
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

	w.watchTimeout(ctx, timeout, w.timeExtensions, result)

	r := <-result

	switch r := r.(type) {
	case *discord.Message:
		return r, nil
	case error:
		return nil, r
	default: // this should never happen
		return nil, errors.NewWithStack("reply: unexpected return value of result channel")
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

		go func() {
			for _, r := range w.cancelReactions {
				err := w.state.DeleteReactions(w.ctx.ChannelID, r.messageID, r.reaction)
				if err != nil {
					w.ctx.HandleErrorSilent(err)
				}
			}
		}()
	}, nil
}

func (w *Waiter) watchTimeout(
	ctx context.Context, timeout time.Duration, timeExtensions int, result chan<- interface{},
) {
	t := time.NewTimer(timeout)

	go func() {
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
			}

			if timeExtensions == 0 {
				err := errors.NewUserInfol(timeoutInfo.
					WithPlaceholders(timeoutInfoPlaceholders{
						ResponseUserMention: w.ctx.Author.Mention(),
					}))
				sendResult(ctx, result, err)
				return
			}

			if timeExtensions > 0 {
				timeExtensions--
			}

			if err := w.askForTimeExtension(ctx); err != nil {
				sendResult(ctx, result, err)
				return
			}
			// else if err == nil: the user passed time extension check or the context was canceled

			t.Reset(timeout)
		}
	}()
}

func (w *Waiter) askForTimeExtension(ctx context.Context) error {
	msg, err := w.sendTimeExtensionMessage()
	if err != nil {
		return err
	}

	defer func() {
		err := w.state.DeleteMessage(msg.ChannelID, msg.ID)
		if err != nil {
			w.ctx.HandleErrorSilent(err)
		}
	}()

	react := make(chan struct{})

	rm, err := w.state.AddHandler(func(_ *state.State, e *state.MessageReactionAddEvent) {
		if e.MessageID == msg.ID && e.UserID == w.ctx.Author.ID && e.Emoji.APIString() == TimeExtensionReaction {
			select {
			case <-ctx.Done():
			case react <- struct{}{}:
			}
		}
	})
	if err != nil { // this should never happen
		return errors.WithStack(err)
	}

	defer rm()

	select {
	case <-time.After(8 * time.Second):
		return errors.NewUserInfol(timeoutInfo.
			WithPlaceholders(timeoutInfoPlaceholders{
				ResponseUserMention: w.ctx.Author.Mention(),
			}))
	case <-react:
		return nil
	case <-ctx.Done():
		return nil
	}
}

func (w *Waiter) sendTimeExtensionMessage() (*discord.Message, error) {
	embed, err := errors.NewCustomUserInfo().
		WithSimpleTitlel(timeExtensionTitle).
		WithDescriptionl(timeExtensionDescription.
			WithPlaceholders(timeExtensionDescriptionPlaceholders{
				ResponseUserMention:   "<@" + w.ctx.Author.ID.String() + ">",
				TimeExtensionReaction: TimeExtensionReaction,
			})).
		Embed(w.ctx.Localizer)
	if err != nil {
		return nil, err
	}

	msg, err := w.ctx.ReplyEmbed(embed)
	if err != nil {
		return nil, err
	}

	if err = w.state.React(msg.ChannelID, msg.ID, TimeExtensionReaction); err != nil {
		// attempt to delete if something failed.
		err = w.state.DeleteMessage(msg.ChannelID, msg.ID)
	}

	return msg, errors.WithStack(err)
}
