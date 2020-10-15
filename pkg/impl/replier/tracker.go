package replier

import (
	"sync"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"
)

// Tracker is a plugin.Replier that tracks the messages that were sent.
//
// Usage
//
// The tracker should be used in a middleware.
// Make sure that no middleware before the one using tracker sends message.
//
// In the middleware replace the Replier of the Context call next and then
// handle the results.
//
// 		func(next bot.CommandFunc) bot.CommandFunc {
//			return func(s *state.State, ctx *plugin.Context) error {
//				t := NewTracker(s, ctx.Author.ID, ctx.ChannelID)
//				ctx.Replier = t
//
//				err := next(s, ctx)
//				if err != nil {
//					return err
//				}
//
//				// do something with t.DMs() and t.GuildMessages()
//			}
//		}
type Tracker struct {
	s *state.State

	dms     []discord.Message
	dmMutex sync.RWMutex
	dmID    discord.ChannelID
	userID  discord.UserID

	guildMessages      []discord.Message
	guildMessagesMutex sync.RWMutex
	guildChannelID     discord.ChannelID
}

// NewTracker creates a new tracker using the passed state, with the passed
// invoking user and the passed guild channel.
func NewTracker(s *state.State, invokingUserID discord.UserID, guildChannelID discord.ChannelID) *Tracker {
	return &Tracker{
		s:              s,
		userID:         invokingUserID,
		guildChannelID: guildChannelID,
	}
}

// GuildMessages returns the guild messages that were sent.
func (t *Tracker) GuildMessages() (cp []discord.Message) {
	cp = make([]discord.Message, len(t.guildMessages))
	copy(cp, t.guildMessages)

	return
}

// DMs returns the direct messages that were sent.
func (t *Tracker) DMs() (cp []discord.Message) {
	cp = make([]discord.Message, len(t.dms))
	copy(cp, t.dms)

	return
}

func (t *Tracker) ReplyMessage(data api.SendMessageData) (*discord.Message, error) {
	t.guildMessagesMutex.Lock()
	defer t.guildMessagesMutex.Unlock()

	msg, err := t.s.SendMessageComplex(t.guildChannelID, data)
	if err != nil {
		return nil, err
	}

	t.guildMessages = append(t.guildMessages, *msg)

	return msg, nil
}

func (t *Tracker) ReplyDM(data api.SendMessageData) (*discord.Message, error) {
	if !t.dmID.IsValid() {
		c, err := t.s.CreatePrivateChannel(t.userID)
		if err != nil {
			return nil, err
		}

		t.dmID = c.ID
	}

	t.dmMutex.Lock()
	defer t.dmMutex.Unlock()

	msg, err := t.s.SendMessageComplex(t.dmID, data)
	if err != nil {
		return nil, err
	}

	t.dms = append(t.dms, *msg)

	return msg, nil
}
