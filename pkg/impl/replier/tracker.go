package replier

import (
	"sync"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/pkg/plugin"
)

// Tracker is a plugin.Replier that tracks the messages that were sent.
type Tracker struct {
	r plugin.Replier

	guildMessages       []discord.Message
	editedGuildMessages []discord.Message
	guildMessagesMutex  sync.Mutex

	dms       []discord.Message
	editedDMs []discord.Message
	dmMutex   sync.Mutex
}

var _ plugin.Replier = new(Tracker)

// NewTracker creates a new tracker that tracks the replies sent by the passed
// baseReplier.
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
// 			t := NewTracker(ctx.Replier)
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
func NewTracker(baseReplier plugin.Replier) *Tracker {
	return &Tracker{r: baseReplier}
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
	msg, err := t.r.Reply(ctx, data)
	if err != nil {
		return nil, err
	}

	t.guildMessagesMutex.Lock()
	t.guildMessages = append(t.guildMessages, *msg)
	t.guildMessagesMutex.Unlock()

	return msg, nil
}

func (t *Tracker) ReplyDM(ctx *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	msg, err := t.r.ReplyDM(ctx, data)
	if err != nil {
		return nil, err
	}

	t.dmMutex.Lock()
	t.dms = append(t.dms, *msg)
	t.dmMutex.Unlock()

	return msg, nil
}

func (t *Tracker) Edit(
	ctx *plugin.Context, messageID discord.MessageID, data api.EditMessageData,
) (*discord.Message, error) {
	msg, err := t.r.Edit(ctx, messageID, data)
	if err != nil {
		return nil, err
	}

	t.storeEdit(*msg)
	return msg, nil
}

func (t *Tracker) EditDM(
	ctx *plugin.Context, messageID discord.MessageID, data api.EditMessageData,
) (*discord.Message, error) {
	msg, err := t.r.EditDM(ctx, messageID, data)
	if err != nil {
		return nil, err
	}

	t.storeEditDM(*msg)
	return msg, nil
}

func (t *Tracker) storeEdit(m discord.Message) {
	t.guildMessagesMutex.Lock()
	defer t.guildMessagesMutex.Unlock()

	for i, old := range t.guildMessages {
		if old.ID == m.ID {
			t.guildMessages[i] = m
			return
		}
	}

	t.editedGuildMessages = append(t.editedGuildMessages, m)
}

func (t *Tracker) storeEditDM(m discord.Message) {
	t.dmMutex.Lock()
	defer t.dmMutex.Unlock()

	for i, old := range t.dms {
		if old.ID == m.ID {
			t.dms[i] = m
			return
		}
	}

	t.editedDMs = append(t.editedDMs, m)
}
