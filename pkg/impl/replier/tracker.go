package replier

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// Tracker is a plugin.Replier that tracks the messages that were sent.
type Tracker struct {
	s *state.State

	guildMessages      []discord.Message
	guildMessagesMutex sync.Mutex

	editedGuildMessages      []discord.Message
	editedGuildMessagesMutex sync.Mutex

	dms     []discord.Message
	dmMutex sync.Mutex
	dmID    discord.ChannelID

	editedDMs   []discord.Message
	editDMMutex sync.Mutex
}

var _ plugin.Replier = new(Tracker)

// NewTracker creates a new tracker using the passed state, with the passed
// invoking user and the passed guild channel.
//
// Example Usage
//
//  b, _ := bot.New(bot.Options{Token: "abc"})
//
//	// A tracker is typically added to a Context through a middleware.
//	// Make sure that the middleware replacing the default replier is executed
//	// before any middlewares that could send replies.
//
//	b.MustAddMiddleware(func(next bot.CommandFunc) bot.CommandFunc {
//		return func(s *state.State, ctx *plugin.Context) error {
//			t := NewTracker(s)
//			ctx.Replier = t // replace the default replier
//
//			err := next(s, ctx)
//			if err != nil {
//				return err
//			}
//
//			// do something with t.DMs() and t.GuildMessages()
//
//			return nil
//		}
//	})
func NewTracker(s *state.State) *Tracker {
	return &Tracker{s: s}
}

// GuildMessages returns the guild messages that were sent.
// This may include edited messages, if they were sent earlier by the command.
func (t *Tracker) GuildMessages() (cp []discord.Message) {
	cp = make([]discord.Message, len(t.guildMessages))
	copy(cp, t.guildMessages)

	return
}

// DMs returns the direct messages that were sent.
// This may include edited messages, if they were sent earlier by the command.
func (t *Tracker) DMs() (cp []discord.Message) {
	cp = make([]discord.Message, len(t.dms))
	copy(cp, t.dms)

	return
}

// EditedGuildMessages returns the guild messages that were edited, but not
// previously sent by the command.
func (t *Tracker) EditedGuildMessages() (cp []discord.Message) {
	cp = make([]discord.Message, len(t.editedGuildMessages))
	copy(cp, t.editedGuildMessages)

	return
}

// EditedDMs returns the guild messages that were edited, but not previously
// sent by the command.
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

	msg, err := t.s.SendMessageComplex(ctx.ChannelID, data)
	if err != nil {
		return nil, err
	}

	t.guildMessagesMutex.Lock()
	defer t.guildMessagesMutex.Unlock()

	t.guildMessages = append(t.guildMessages, *msg)

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

	err = t.lazyDM(ctx)
	if err != nil {
		return nil, err
	}

	msg, err := t.s.SendMessageComplex(t.dmID, data)
	if err != nil {
		return nil, err
	}

	t.dmMutex.Lock()
	defer t.dmMutex.Unlock()

	t.dms = append(t.dms, *msg)

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
		return nil, err
	}

	t.editedGuildMessagesMutex.Lock()
	defer t.editedGuildMessagesMutex.Unlock()

	t.editedGuildMessages = append(t.editedGuildMessages, *msg)

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

	err = t.lazyDM(ctx)
	if err != nil {
		return nil, err
	}

	msg, err := t.s.EditMessageComplex(t.dmID, messageID, data)
	if err != nil {
		return nil, err
	}

	t.editDMMutex.Lock()
	defer t.editDMMutex.Unlock()

	t.editedDMs = append(t.editedDMs, *msg)

	return msg, nil
}

// lazyDM lazily gets the id of the direct message channel with the invoking
// user.
func (t *Tracker) lazyDM(ctx *plugin.Context) error {
	if !t.dmID.IsValid() {
		c, err := t.s.CreatePrivateChannel(ctx.Author.ID)
		if err != nil {
			return err
		}

		t.dmID = c.ID
	}

	return nil
}
