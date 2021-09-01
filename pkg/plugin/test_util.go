package plugin

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"
)

// mockDiscordDataProvider is a copy of mock.DiscordDataProvider to prevent
// import cycles.
type mockDiscordDataProvider struct {
	ChannelReturn *discord.Channel
	ChannelError  error

	ParentChannelReturn *discord.Channel
	ParentChannelError  error

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

func (d mockDiscordDataProvider) ParentChannelAsync() func() (*discord.Channel, error) {
	return func() (*discord.Channel, error) {
		return d.ParentChannelReturn, d.ParentChannelError
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

func newMockedWrappedReplier(s *state.State, channelID discord.ChannelID, userID discord.UserID) *wrappedReplier {
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
