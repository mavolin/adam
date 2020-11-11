package messageutil

import (
	"context"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
)

var (
	// DefaultReplyWaiter is the waiter used for NewReplyWaiterFromDefault.
	// DefaultReplyWaiter must not be used directly to handleMessages reply.
	DefaultReplyWaiter = &ReplyWaiter{
		cancelKeywords: []*i18nutil.Text{i18nutil.NewTextl(defaultCancelKeyword)},
		maxTimeout:     MaxTimeout,
	}

	// MaxTimeout is the default maximum amount of time a ReplyWaiter will wait,
	// even if a user is still typing.
	MaxTimeout = 30 * time.Minute

	// Canceled is the error that gets returned, if a user signals the bot
	// should not continue waiting for a reply.
	Canceled = errors.NewInformationalError("canceled")
)

// typingInterval is the interval in which the client of the user sends the
// typing event, if the user is continuously typing.
//
// It has been observed that the first follow-up event is received after about
// 9.5 seconds, all successive events are received in an intervall of approx.
// 8.25 seconds.
// Additionally, there is a 1.5 second margin for network delays.
var typingInterval = 11 * time.Second

type (
	// A ReplyWaiter is used to await messages.
	// Wait can be cancelled either by the user by using a cancel keyword or
	// using a cancel reaction.
	// Furthermore, wait may be cancelled by the ReplyWaiter if the initial timeout
	// expires and the user is not typing or the user stopped typing and the
	// typing timeout expired.
	ReplyWaiter struct {
		state *state.State
		ctx   *plugin.Context

		userID    discord.UserID
		channelID discord.ChannelID

		caseSensitive bool
		noAutoReact   bool

		cancelKeywords  []*i18nutil.Text
		cancelReactions []cancelReaction

		maxTimeout time.Duration

		middlewares []interface{}
	}

	cancelReaction struct {
		messageID discord.MessageID
		reaction  api.Emoji
	}
)

// NewReplyWaiter creates a new reply waiter using the passed state and context.
// It will wait for a message from the message author in the channel the
// command was invoked in.
// Additionally, the ReplyMiddlewares stored in the Context will be added to
// the waiter.
// ctx.Author will be assumed as the user allowed to make the reply and
// ctx.ChannelID will be assumed as the channel the reply will be made in.
func NewReplyWaiter(s *state.State, ctx *plugin.Context) (w *ReplyWaiter) {
	w = &ReplyWaiter{
		state:      s,
		ctx:        ctx,
		userID:     ctx.Author.ID,
		channelID:  ctx.ChannelID,
		maxTimeout: MaxTimeout,
	}

	w.WithMiddlewares(ctx.ReplyMiddlewares...)

	return w
}

// NewReplyWaiterFromDefault creates a new default waiter using the DefaultReplyWaiter
// variable as a template.
func NewReplyWaiterFromDefault(s *state.State, ctx *plugin.Context) (w *ReplyWaiter) {
	w = DefaultReplyWaiter.Copy()
	w.state = s
	w.ctx = ctx
	w.userID = ctx.Author.ID
	w.channelID = ctx.ChannelID

	w.WithMiddlewares(ctx.ReplyMiddlewares...)

	return
}

// WithUser changes the user that is expected to reply to the user with the
// passed id.
func (w *ReplyWaiter) WithUser(id discord.UserID) *ReplyWaiter {
	w.userID = id
	return w
}

// InChannel changes the channel to listen for the reply to the channel with
// the passed id.
func (w *ReplyWaiter) InChannel(id discord.ChannelID) *ReplyWaiter {
	w.channelID = id
	return w
}

// CaseSensitive makes the cancel keywords check case sensitive.
func (w *ReplyWaiter) CaseSensitive() *ReplyWaiter {
	w.caseSensitive = true
	return w
}

// NoAutoReact disables the automatic reaction of the cancel reaction messages
// with their respective emojis.
func (w *ReplyWaiter) NoAutoReact() *ReplyWaiter {
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
//		• func(*state.State, *state.MessageCreateEvent)
//		• func(*state.State, *state.MessageCreateEvent) error
func (w *ReplyWaiter) WithMiddlewares(middlewares ...interface{}) *ReplyWaiter {
	if len(w.middlewares) == 0 {
		w.middlewares = make([]interface{}, 0, len(middlewares))
	}

	for _, m := range middlewares {
		switch m.(type) { // check that the middleware is of a valid type
		case func(*state.State, interface{}):
		case func(*state.State, interface{}) error:
		case func(*state.State, *state.Base):
		case func(*state.State, *state.Base) error:
		case func(*state.State, *state.MessageCreateEvent):
		case func(*state.State, *state.MessageCreateEvent) error:
		default:
			continue
		}

		w.middlewares = append(w.middlewares, m)
	}

	return w
}

// WithCancelKeyword adds the passed keyword to the cancel keywords.
// If the user filtered for writes this keyword AwaitReply will return Canceled.
func (w *ReplyWaiter) WithCancelKeyword(keyword string) *ReplyWaiter {
	w.cancelKeywords = append(w.cancelKeywords, i18nutil.NewText(keyword))
	return w
}

// WithCancelKeywordl adds the passed keyword to the cancel keywords.
// If the user filtered for writes this keyword AwaitReply will return Canceled.
func (w *ReplyWaiter) WithCancelKeywordl(keyword *i18n.Config) *ReplyWaiter {
	w.cancelKeywords = append(w.cancelKeywords, i18nutil.NewTextl(keyword))
	return w
}

// WithCancelKeywordlt adds the passed keyword to the cancel keywords.
// If the user filtered for writes this keyword AwaitReply will return Canceled.
func (w *ReplyWaiter) WithCancelKeywordlt(keyword i18n.Term) *ReplyWaiter {
	return w.WithCancelKeywordl(keyword.AsConfig())
}

// WithMaxTimeout changes the maximum timeout of the waiter to max.
// The maximum timeout is the timeout after which the ReplyWaiter will exit, even if
// the user is still typing.
func (w *ReplyWaiter) WithMaxTimeout(max time.Duration) *ReplyWaiter {
	if w.maxTimeout > 0 {
		w.maxTimeout = max
	}

	return w
}

// WithCancelReactions adds the passed cancel reaction.
// If the user reacts with the passed emoji, AwaitReply will return with error
// Canceled.
func (w *ReplyWaiter) WithCancelReactions(messageID discord.MessageID, reactions ...api.Emoji) *ReplyWaiter {
	for _, r := range reactions {
		w.cancelReactions = append(w.cancelReactions, cancelReaction{
			messageID: messageID,
			reaction:  r,
		})
	}

	return w
}

// Copy creates a copy of the ReplyWaiter.
func (w *ReplyWaiter) Copy() (cp *ReplyWaiter) {
	cp = &ReplyWaiter{
		caseSensitive: w.caseSensitive,
		noAutoReact:   w.noAutoReact,
		maxTimeout:    w.maxTimeout,
	}

	cp.cancelKeywords = make([]*i18nutil.Text, len(w.cancelKeywords))
	copy(cp.cancelKeywords, w.cancelKeywords)

	cp.cancelReactions = make([]cancelReaction, len(w.cancelReactions))
	copy(cp.cancelReactions, w.cancelReactions)

	cp.middlewares = make([]interface{}, len(w.middlewares))
	copy(cp.middlewares, w.middlewares)

	return
}

// AwaitReply awaits a reply of the user until the user signals cancellation, the
// initial timeout expires and the user is not typing or the user stops typing
// and the typing timeout is reached. Note you need the typing intent to
// monitor typing.
//
// If one of the timeouts is reached, a *TimeoutError will be returned.
// If the user cancels the reply, Canceled will be returned.
//
// The typing timeout will start after the user stops typing.
// Because Discord sends the typing event in an interval of about 10 seconds,
// the user might have stopped typing before the waiter notices that the typing
// status was not updated.
//
// Besides that, a reply can also be canceled through a middleware.
// If one the middlewares returns state.Filtered, errors.Abort will be
// returned.
func (w *ReplyWaiter) Await(initialTimeout, typingTimeout time.Duration) (*discord.Message, error) {
	return w.AwaitWithContext(context.Background(), initialTimeout, typingTimeout)
}

// AwaitReply awaits a reply of the user until the user signals cancellation, the
// initial timeout expires and the user is not typing or the user stops typing
// and the typing timeout is reached. Note you need the typing intent to
// monitor typing.
//
// If one of the timeouts is reached, a *TimeoutError will be returned.
// If the user cancels the reply, Canceled will be returned.
// If the context expires, context.Canceled will be returned.
//
// The typing timeout will start after the user stops typing.
// Because Discord sends the typing event in an interval of about 10 seconds,
// the user might have stopped typing before the waiter notices that the typing
// status was not updated.
//
// Besides that, a reply can also be canceled through a middleware.
// If one the middlewares returns state.Filtered, errors.Abort will be
// returned.
func (w *ReplyWaiter) AwaitWithContext(
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

	msgCleanup, err := w.handleMessages(ctx, result)
	if err != nil {
		return nil, err
	}

	defer msgCleanup()

	if len(w.cancelReactions) > 0 {
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
			return nil, errors.NewWithStack("messageutil: unexpected return value of result channel")
		}
	}
}

func (w *ReplyWaiter) handleMessages(ctx context.Context, result chan<- interface{}) (func(), error) {
	rm, err := w.state.AddHandler(func(s *state.State, e *state.MessageCreateEvent) {
		if e.ChannelID != w.channelID || e.Author.ID != w.userID { // not the message we are waiting for
			return
		}

		if err := invokeMessageMiddlewares(s, e, w.middlewares); err != nil {
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

func (w *ReplyWaiter) handleCancelReactions(ctx context.Context, result chan<- interface{}) (func(), error) {
	if !w.noAutoReact {
		for _, r := range w.cancelReactions {
			if err := w.state.React(w.channelID, r.messageID, r.reaction); err != nil {
				w.ctx.HandleErrorSilent(err)
			}
		}
	}

	rm, err := w.state.AddHandler(func(s *state.State, e *state.MessageReactionAddEvent) {
		if e.UserID != w.userID {
			return
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
				for _, r := range w.cancelReactions {
					err := w.state.DeleteReactions(w.channelID, r.messageID, r.reaction)
					if err != nil {
						w.ctx.HandleErrorSilent(err)
					}
				}
			}()
		}
	}, nil
}

func (w *ReplyWaiter) watchTimeout(
	ctx context.Context, initialTimeout, typingTimeout time.Duration, result chan<- interface{},
) (rm func(), err error) {
	maxTimer := time.NewTimer(w.maxTimeout)
	t := time.NewTimer(initialTimeout)

	if typingTimeout > 0 {
		rm, err = w.state.AddHandler(func(s *state.State, e *state.TypingStartEvent) {
			if e.ChannelID != w.channelID || e.UserID != w.userID {
				return
			}

			// this should always return true, except if timer expired after
			// the typing event was received and this handler called
			if t.Stop() {
				t.Reset(typingTimeout + typingInterval)
			}
		})
		if err != nil {
			return
		}
	} else {
		rm = func() {}
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				maxTimer.Stop()

				return
			case <-t.C:
				maxTimer.Stop()

				result <- &TimeoutError{UserID: w.ctx.Author.ID}
				return
			case <-maxTimer.C:
				t.Stop()

				result <- &TimeoutError{UserID: w.ctx.Author.ID}
				return
			}
		}
	}()

	return
}
