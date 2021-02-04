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

type mockCommand struct {
	name            string
	aliases         []string
	args            ArgConfig
	shortDesc       string
	longDesc        string
	exampleArgs     []string
	hidden          bool
	channelTypes    ChannelTypes
	botPermissions  discord.Permissions
	restrictionFunc RestrictionFunc
	throttler       Throttler
	invokeFunc      func(*state.State, *Context) (interface{}, error)
}

func (c mockCommand) GetName() string                            { return c.name }
func (c mockCommand) GetAliases() []string                       { return c.aliases }
func (c mockCommand) GetArgs() ArgConfig                         { return c.args }
func (c mockCommand) GetShortDescription(*i18n.Localizer) string { return c.shortDesc }
func (c mockCommand) GetLongDescription(*i18n.Localizer) string  { return c.longDesc }
func (c mockCommand) GetExampleArgs(*i18n.Localizer) []string    { return c.exampleArgs }
func (c mockCommand) IsHidden() bool                             { return c.hidden }
func (c mockCommand) GetChannelTypes() ChannelTypes              { return c.channelTypes }
func (c mockCommand) GetBotPermissions() discord.Permissions     { return c.botPermissions }

func (c mockCommand) IsRestricted(s *state.State, ctx *Context) error {
	if c.restrictionFunc == nil {
		return nil
	}

	return c.restrictionFunc(s, ctx)
}

func (c mockCommand) GetThrottler() Throttler {
	return c.throttler
}

func (c mockCommand) Invoke(s *state.State, ctx *Context) (interface{}, error) {
	return c.invokeFunc(s, ctx)
}

type mockModule struct {
	name      string
	shortDesc string
	longDesc  string
	commands  []Command
	modules   []Module
}

func (c mockModule) GetName() string                            { return c.name }
func (c mockModule) GetShortDescription(*i18n.Localizer) string { return c.shortDesc }
func (c mockModule) GetLongDescription(*i18n.Localizer) string  { return c.longDesc }
func (c mockModule) Commands() []Command                        { return c.commands }
func (c mockModule) Modules() []Module                          { return c.modules }

type mockThrottler struct {
	cmp string // used to make throttlers unique
}

func (m mockThrottler) Check(*state.State, *Context) (func(), error) { return func() {}, nil }

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
