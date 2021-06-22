package plugin

import (
	"errors"
	"testing"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/i18n"
)

// mockLocalizer is a copy of mock.Localizer, used to prevent import cycles.
type mockLocalizer struct {
	t *testing.T

	def      string
	onReturn map[i18n.Term]string
	errOn    map[i18n.Term]struct{}
}

func newMockedLocalizer(t *testing.T) *mockLocalizer {
	return &mockLocalizer{
		t:        t,
		onReturn: make(map[i18n.Term]string),
		errOn:    make(map[i18n.Term]struct{}),
	}
}

func (l *mockLocalizer) on(term i18n.Term, response string) *mockLocalizer {
	l.onReturn[term] = response
	return l
}

func (l *mockLocalizer) build() *i18n.Localizer {
	return i18n.NewLocalizer("dev", func(term i18n.Term, _ map[string]interface{}, _ interface{}) (string, error) {
		r, ok := l.onReturn[term]
		if ok {
			return r, nil
		}

		_, ok = l.errOn[term]
		if ok {
			return r, errors.New("error")
		}

		if l.def == "" {
			assert.Failf(l.t, "unexpected call to Localize", "unknown term %s", term)

			return string(term), errors.New("unknown term")
		}

		return l.def, nil
	})
}

// mockDiscordDataProvider is a copy of mock.DiscordDataProvider to prevent
// import cycles.
type mockDiscordDataProvider struct {
	ChannelReturn *discord.Channel
	ChannelError  error

	GuildReturn *discord.Guild
	GuildError  error

	SelfReturn *discord.Member
	SelfError  error
}

func (d mockDiscordDataProvider) ChannelAsync() func() (*discord.Channel, error) {
	return func() (*discord.Channel, error) {
		return d.ChannelReturn, d.ChannelError
	}
}

func (d mockDiscordDataProvider) GuildAsync() func() (*discord.Guild, error) {
	return func() (*discord.Guild, error) {
		return d.GuildReturn, d.GuildError
	}
}

func (d mockDiscordDataProvider) SelfAsync() func() (*discord.Member, error) {
	return func() (*discord.Member, error) {
		return d.SelfReturn, d.SelfError
	}
}

// wrappedReplier is a copy of replier.wrappedReplier, used to prevent import
// cycles.
type wrappedReplier struct {
	s         *state.State
	channelID discord.ChannelID

	userID discord.UserID
	dmID   discord.ChannelID
}

func replierFromState(s *state.State, channelID discord.ChannelID, userID discord.UserID) *wrappedReplier {
	return &wrappedReplier{
		s:         s,
		userID:    userID,
		channelID: channelID,
	}
}

func (r *wrappedReplier) Reply(_ *Context, data api.SendMessageData) (*discord.Message, error) {
	return r.s.SendMessageComplex(r.channelID, data)
}

func (r *wrappedReplier) ReplyDM(_ *Context, data api.SendMessageData) (*discord.Message, error) {
	if !r.dmID.IsValid() {
		c, err := r.s.CreatePrivateChannel(r.userID)
		if err != nil {
			return nil, err
		}

		r.dmID = c.ID
	}

	return r.s.SendMessageComplex(r.dmID, data)
}

func (r *wrappedReplier) Edit(
	_ *Context, messageID discord.MessageID, data api.EditMessageData,
) (*discord.Message, error) {
	return r.s.EditMessageComplex(r.channelID, messageID, data)
}

func (r *wrappedReplier) EditDM(
	_ *Context, messageID discord.MessageID, data api.EditMessageData,
) (*discord.Message, error) {
	if !r.dmID.IsValid() {
		c, err := r.s.CreatePrivateChannel(r.userID)
		if err != nil {
			return nil, err
		}

		r.dmID = c.ID
	}

	return r.s.EditMessageComplex(r.dmID, messageID, data)
}
