package replier

import (
	"sync"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/discorderr"
)

// Tracker is a plugin.Replier that tracks the messages that were sent.
type Tracker struct {
	s           *state.State
	inlineReply bool

	guildMessages       []discord.Message
	editedGuildMessages []discord.Message
	guildMessagesMutex  sync.Mutex

	dmID  discord.ChannelID
	dmErr error

	dms       []discord.Message
	editedDMs []discord.Message
	dmMutex   sync.Mutex
}

var _ plugin.Replier = new(Tracker)

// NewTracker creates a new tracker using the passed state, with the passed
// invoking user and the passed guild channel.
//
// If inlineReply is set to true, messages will reference the invoke, unless
// MessageReference is non-nil.
//
// Example Usage
//
// 	b, _ := bot.New(bot.Options{Token: "abc"})
//
// 	// A tracker is typically added to a Context through a middleware.
// 	// Make sure that the middleware replacing the default replier is executed
// 	// before any middlewares that could send replies.
//
// 	b.AddMiddleware(func(next bot.CommandFunc) bot.CommandFunc {
// 		return func(s *state.State, ctx *plugin.Context) error {
// 			t := NewTracker(s, false)
// 			ctx.Replier = t // replace the default replier
//
// 			if err := next(s, ctx); err != nil {
// 				return err
// 			}
//
// 			// do something with t.DMs() and t.GuildMessages()
//
// 			return nil
// 		}
// 	})
//
// Creating replies with the Tracker is concurrent-safe.
// However, retrieving the message sent is not.
func NewTracker(s *state.State, inlineReply bool) *Tracker {
	return &Tracker{s: s, inlineReply: inlineReply}
}

// GuildMessages returns the guild messages that were sent.
// This may include edited messages, if they were sent earlier by this replier.
func (t *Tracker) GuildMessages() (cp []discord.Message) {
	cp = make([]discord.Message, len(t.guildMessages))
	copy(cp, t.guildMessages)

	return
}

// DMs returns the direct messages that were sent.
// This may include edited messages, if they were sent earlier by this replier.
func (t *Tracker) DMs() (cp []discord.Message) {
	cp = make([]discord.Message, len(t.dms))
	copy(cp, t.dms)

	return
}

// EditedGuildMessages returns the guild messages that were edited, but not
// previously sent by this replier.
func (t *Tracker) EditedGuildMessages() (cp []discord.Message) {
	cp = make([]discord.Message, len(t.editedGuildMessages))
	copy(cp, t.editedGuildMessages)

	return
}

// EditedDMs returns the guild messages that were edited, but not previously
// sent by this replier.
func (t *Tracker) EditedDMs() (cp []discord.Message) {
	cp = make([]discord.Message, len(t.editedDMs))
	copy(cp, t.editedDMs)

	return
}

func (t *Tracker) Reply(ctx *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, plugin.NewBotPermissionsError(discord.PermissionSendMessages)
	}

	if data.Reference == nil && t.inlineReply {
		data.Reference = &discord.MessageReference{MessageID: ctx.Message.ID}
	}

	msg, err := t.s.SendMessageComplex(ctx.ChannelID, data)
	if err != nil {
		// user deleted channel
		if discorderr.Is(discorderr.As(err), discorderr.UnknownChannel) {
			return nil, errors.Abort
		}

		return nil, errors.WithStack(err)
	}

	t.guildMessagesMutex.Lock()
	t.guildMessages = append(t.guildMessages, *msg)
	t.guildMessagesMutex.Unlock()

	return msg, nil
}

func (t *Tracker) ReplyDM(ctx *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, plugin.NewBotPermissionsError(discord.PermissionSendMessages)
	}

	dmID, err := t.lazyDM(ctx)
	if err != nil {
		return nil, err
	}

	msg, err := t.s.SendMessageComplex(dmID, data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	t.dmMutex.Lock()
	t.dms = append(t.dms, *msg)
	t.dmMutex.Unlock()

	return msg, nil
}

func (t *Tracker) Edit(
	ctx *plugin.Context, messageID discord.MessageID, data api.EditMessageData,
) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, plugin.NewBotPermissionsError(discord.PermissionSendMessages)
	}

	msg, err := t.s.EditMessageComplex(ctx.ChannelID, messageID, data)
	if err != nil {
		// user deleted channel
		if discorderr.Is(discorderr.As(err), discorderr.UnknownChannel) {
			return nil, errors.Abort
		}

		t.guildMessagesMutex.Lock()
		defer t.guildMessagesMutex.Unlock()

		// we sent that message before, so it was deleted by someone else
		if t.hasID(messageID) {
			return nil, errors.Abort
		}

		return nil, errors.WithStack(err)
	}

	t.guildMessagesMutex.Lock()
	defer t.guildMessagesMutex.Unlock()

	if !t.hasID(messageID) {
		t.editedGuildMessages = append(t.editedGuildMessages, *msg)
	}

	return msg, nil
}

func (t *Tracker) EditDM(
	ctx *plugin.Context, messageID discord.MessageID, data api.EditMessageData,
) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, plugin.NewBotPermissionsError(discord.PermissionSendMessages)
	}

	dmID, err := t.lazyDM(ctx)
	if err != nil {
		return nil, err
	}

	msg, err := t.s.EditMessageComplex(dmID, messageID, data)
	if err != nil {
		t.dmMutex.Lock()
		defer t.dmMutex.Unlock()

		// we sent that message before, so it was deleted by someone else
		if t.hasDMID(messageID) {
			return nil, errors.Abort
		}

		return nil, errors.WithStack(err)
	}

	t.dmMutex.Lock()
	defer t.dmMutex.Unlock()

	if !t.hasDMID(messageID) {
		t.editedDMs = append(t.editedDMs, *msg)
	}

	return msg, nil
}

// lazyDM lazily gets the id of the direct message channel with the invoking
// user.
func (t *Tracker) lazyDM(ctx *plugin.Context) (discord.ChannelID, error) {
	t.dmMutex.Lock()
	defer t.dmMutex.Unlock()

	if t.dmID != 0 || t.dmErr != nil {
		return t.dmID, t.dmErr
	}

	c, err := t.s.CreatePrivateChannel(ctx.Author.ID)
	t.dmErr = err
	if err == nil {
		t.dmID = c.ID
	}

	return t.dmID, t.dmErr
}

func (t *Tracker) hasID(id discord.MessageID) bool {
	for _, msg := range t.guildMessages {
		if msg.ID == id {
			return true
		}
	}

	return false
}

func (t *Tracker) hasDMID(id discord.MessageID) bool {
	for _, msg := range t.dms {
		if msg.ID == id {
			return true
		}
	}

	return false
}
