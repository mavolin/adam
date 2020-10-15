package errors

import (
	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// wrappedReplier is a copy of replier.wrappedReplier, used to prevent import
// cycles.
type wrappedReplier struct {
	s      *state.State
	userID discord.UserID
	dmID   discord.ChannelID
}

func replierFromState(s *state.State, userID discord.UserID) plugin.Replier {
	return &wrappedReplier{
		s:      s,
		userID: userID,
	}
}

func (r *wrappedReplier) SendMessageComplex(
	channelID discord.ChannelID, data api.SendMessageData,
) (*discord.Message, error) {
	return r.s.SendMessageComplex(channelID, data)
}

func (r *wrappedReplier) PrivateChannelID() (discord.ChannelID, error) {
	if r.dmID.IsValid() {
		return r.dmID, nil
	}

	c, err := r.s.CreatePrivateChannel(r.userID)
	if err != nil {
		return 0, err
	}

	return c.ID, nil
}
