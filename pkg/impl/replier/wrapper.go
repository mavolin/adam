package replier

import (
	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

type wrappedReplier struct {
	s         *state.State
	channelID discord.ChannelID

	userID discord.UserID
	dmID   discord.ChannelID
}

// WrapState wraps the passed state and id of the invoking user into a
// plugin.Replier.
func WrapState(s *state.State, invokingUserID discord.UserID, channelID discord.ChannelID) plugin.Replier {
	return &wrappedReplier{
		s:         s,
		channelID: channelID,
		userID:    invokingUserID,
	}
}

func (r *wrappedReplier) ReplyMessage(data api.SendMessageData) (*discord.Message, error) {
	return r.s.SendMessageComplex(r.channelID, data)
}

func (r *wrappedReplier) ReplyDM(data api.SendMessageData) (*discord.Message, error) {
	if !r.dmID.IsValid() {
		c, err := r.s.CreatePrivateChannel(r.userID)
		if err != nil {
			return nil, err
		}

		r.dmID = c.ID
	}

	return r.s.SendMessageComplex(r.dmID, data)
}
