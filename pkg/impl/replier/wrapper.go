package replier

import (
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

type wrappedReplier struct {
	s    *state.State
	dmID discord.ChannelID
}

var _ plugin.Replier = new(wrappedReplier)

// WrapState wraps the passed state and id of the invoking user into a
// plugin.Replier.
func WrapState(s *state.State) plugin.Replier {
	return &wrappedReplier{s: s}
}

func (r *wrappedReplier) Reply(ctx *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, errors.NewInsufficientPermissionsError(discord.PermissionSendMessages)
	}

	return r.s.SendMessageComplex(ctx.ChannelID, data)
}

func (r *wrappedReplier) ReplyDM(ctx *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, errors.NewInsufficientPermissionsError(discord.PermissionSendMessages)
	}

	if !r.dmID.IsValid() { // lazily load dm id
		c, err := r.s.CreatePrivateChannel(ctx.Author.ID)
		if err != nil {
			return nil, err
		}

		r.dmID = c.ID
	}

	return r.s.SendMessageComplex(r.dmID, data)
}
