package reply

import (
	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/emoji"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
)

var (
	// DefaultWaiter is the waiter used for NewDefaultWaiter.
	// DefaultWaiter must not be used directly to handleMessages reply.
	DefaultWaiter = &Waiter{
		cancelKeywords: []*i18nutil.Text{i18nutil.NewTextl(defaultCancelKeyword)},
	}

	// TimeExtensionReaction is the reaction used to prolong the wait for a
	// reply, if a time extension is possible.
	TimeExtensionReaction = emoji.CheckMarkButton
)

// Canceled is the error that gets returned, if a user signals the bot
// should not continue waiting for a reply.
var Canceled = errors.NewInformationalError("canceled")

type (
	// Waiter is used to handleMessages a reply from the user.
	// It filters responses by user and channel.
	// Additionally, Waiter provides several ways for an user to abort waiting
	// for a reply.
	Waiter struct {
		state *state.State
		ctx   *plugin.Context

		caseSensitive bool
		noAutoReact   bool
		// timeExtensions is the number of extension the user will be given.
		// the following values are accepted
		//
		//		• timeExtensions < 0:  unlimited
		//		• timeExtensions == 0: none
		// 		• timeExtensions > 0: timeExtension times
		timeExtensions int

		cancelKeywords  []*i18nutil.Text
		cancelReactions []cancelReaction

		middlewares []interface{}
	}

	cancelReaction struct {
		messageID discord.MessageID
		reaction  api.Emoji
	}
)

// NewWaiter creates a new reply waiter using the passed state and context.
// It will wait for a message from the message author in the channel the
// command was invoked in.
// Additionally, the ReplyMiddlewares stored in the Context will be added to
// the waiter.
func NewWaiter(s *state.State, ctx *plugin.Context) (w *Waiter) {
	w = &Waiter{
		state: s,
		ctx:   ctx,
	}

	w.WithMiddlewares(ctx.ReplyMiddlewares...)

	return w
}

// NewDefaultWaiter creates a new default waiter using the DefaultWaiter
// variable as a template.
func NewDefaultWaiter(s *state.State, ctx *plugin.Context) (w *Waiter) {
	w = DefaultWaiter.copy()
	w.state = s
	w.ctx = ctx

	w.WithMiddlewares(ctx.ReplyMiddlewares...)

	return
}

// CaseSensitive makes the cancel keywords check case sensitive.
func (w *Waiter) CaseSensitive() *Waiter {
	w.caseSensitive = true
	return w
}

// NoAutoReact disables the automatic reaction of the cancel reaction messages
// with their respective emojis.
func (w *Waiter) NoAutoReact() *Waiter {
	w.noAutoReact = true
	return w
}

// WithTimeExtensions sets the amount of time extensions the user will be
// given.
// 0 equals to unlimited.
func (w *Waiter) WithTimeExtensions(qty uint) *Waiter {
	if qty == 0 {
		w.timeExtensions = -1 // internally, unlimited is stored as -1
	}

	w.timeExtensions = int(qty)
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
func (w *Waiter) WithMiddlewares(middlewares ...interface{}) *Waiter {
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
// If the user filtered for writes this keyword Await will return Canceled.
func (w *Waiter) WithCancelKeyword(keyword string) *Waiter {
	w.cancelKeywords = append(w.cancelKeywords, i18nutil.NewText(keyword))
	return w
}

// WithCancelKeywordk adds the passed keyword to the cancel keywords.
// If the user filtered for writes this keyword Await will return Canceled.
func (w *Waiter) WithCancelKeywordl(keyword *i18n.Config) *Waiter {
	w.cancelKeywords = append(w.cancelKeywords, i18nutil.NewTextl(keyword))
	return w
}

// WithCancelKeywordlt adds the passed keyword to the cancel keywords.
// If the user filtered for writes this keyword Await will return Canceled.
func (w *Waiter) WithCancelKeywordlt(keyword i18n.Term) *Waiter {
	return w.WithCancelKeywordl(keyword.AsConfig())
}

// WithCancelReaction adds the passed reaction to the cancel reactions.
// The passed message must be in the same channel, that a message is being
// waited for.
// If the user filtered for reacts with this reaction Await will return
// Canceled.
func (w *Waiter) WithCancelReaction(messageID discord.MessageID, reaction api.Emoji) *Waiter {
	w.cancelReactions = append(w.cancelReactions, cancelReaction{
		messageID: messageID,
		reaction:  reaction,
	})

	return w
}

// copy creates a copy of the Waiter.
func (w *Waiter) copy() (cp *Waiter) {
	cp = &Waiter{
		state:          w.state,
		ctx:            w.ctx,
		caseSensitive:  w.caseSensitive,
		noAutoReact:    w.noAutoReact,
		timeExtensions: w.timeExtensions,
	}

	cp.cancelKeywords = make([]*i18nutil.Text, len(w.cancelKeywords))
	copy(cp.cancelKeywords, w.cancelKeywords)

	cp.cancelReactions = make([]cancelReaction, len(w.cancelReactions))
	copy(cp.cancelReactions, w.cancelReactions)

	cp.middlewares = make([]interface{}, len(w.middlewares))
	copy(cp.middlewares, w.middlewares)

	return
}

// Reset resets Waiter.
// Await may never be called on waiters that were resetted.
// This is meant for use of DefaultWaiter only.
func (w *Waiter) Reset() {
	w.state = nil
	w.ctx = nil

	w.caseSensitive = false
	w.noAutoReact = false
	w.timeExtensions = 0
	w.cancelKeywords = nil
	w.cancelReactions = nil
}
