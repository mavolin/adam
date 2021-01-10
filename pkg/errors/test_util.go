package errors

import (
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// wrappedReplier is a copy of replier.wrappedReplier, used to prevent import
// cycles.
type wrappedReplier struct {
	s         *state.State
	channelID discord.ChannelID

	userID discord.UserID
	dmID   discord.ChannelID
}

func replierFromState(s *state.State, channelID discord.ChannelID, userID discord.UserID) plugin.Replier {
	return &wrappedReplier{
		s:         s,
		userID:    userID,
		channelID: channelID,
	}
}

func (r *wrappedReplier) Reply(_ *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	return r.s.SendMessageComplex(r.channelID, data)
}

func (r *wrappedReplier) ReplyDM(_ *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	if !r.dmID.IsValid() {
		c, err := r.s.CreatePrivateChannel(r.userID)
		if err != nil {
			return nil, err
		}

		r.dmID = c.ID
	}

	return r.s.SendMessageComplex(r.dmID, data)
}
