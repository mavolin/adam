package replier

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

// Tracker is a plugin.Replier that tracks the messages that were sent.
type Tracker struct {
	s *state.State

	dms     []discord.Message
	dmMutex sync.RWMutex
	dmID    discord.ChannelID

	guildMessages      []discord.Message
	guildMessagesMutex sync.RWMutex
}

var _ plugin.Replier = new(Tracker)

// NewTracker creates a new tracker using the passed state, with the passed
// invoking user and the passed guild channel.
func NewTracker(s *state.State) *Tracker {
	return &Tracker{s: s}
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

func (t *Tracker) ReplyMessage(ctx *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, errors.NewInsufficientPermissionsError(discord.PermissionSendMessages)
	}

	t.guildMessagesMutex.Lock()
	defer t.guildMessagesMutex.Unlock()

	msg, err := t.s.SendMessageComplex(ctx.ChannelID, data)
	if err != nil {
		return nil, err
	}

	t.guildMessages = append(t.guildMessages, *msg)

	return msg, nil
}

func (t *Tracker) ReplyDM(ctx *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, errors.NewInsufficientPermissionsError(discord.PermissionSendMessages)
	}

	if !t.dmID.IsValid() { // lazily load dm id
		c, err := t.s.CreatePrivateChannel(ctx.Author.ID)
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
