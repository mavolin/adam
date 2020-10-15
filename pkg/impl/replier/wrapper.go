package replier

import (
	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

type wrappedReplier struct {
	s      *state.State
	userID discord.UserID
	dmID   discord.ChannelID
}

// WrapState wraps the passed state and id of the invoking user into a
// plugin.Replier.
func WrapState(s *state.State, invokingUserID discord.UserID) plugin.Replier {
	return &wrappedReplier{
		s:      s,
		userID: invokingUserID,
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
		return 0, errors.WithStack(err)
	}

	return c.ID, nil
}
