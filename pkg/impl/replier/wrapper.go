package replier

import (
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

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
		return nil, plugin.NewBotPermissionsError(discord.PermissionSendMessages)
	}

	return r.s.SendMessageComplex(ctx.ChannelID, data)
}

func (r *wrappedReplier) ReplyDM(ctx *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, plugin.NewBotPermissionsError(discord.PermissionSendMessages)
	}

	err = r.lazyDM(ctx)
	if err != nil {
		return nil, err
	}

	return r.s.SendMessageComplex(r.dmID, data)
}

func (r *wrappedReplier) Edit(
	ctx *plugin.Context, messageID discord.MessageID, data api.EditMessageData,
) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, plugin.NewBotPermissionsError(discord.PermissionSendMessages)
	}

	return r.s.EditMessageComplex(ctx.ChannelID, messageID, data)
}

func (r *wrappedReplier) EditDM(
	ctx *plugin.Context, messageID discord.MessageID, data api.EditMessageData,
) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, plugin.NewBotPermissionsError(discord.PermissionSendMessages)
	}

	err = r.lazyDM(ctx)
	if err != nil {
		return nil, err
	}

	return r.s.EditMessageComplex(r.dmID, messageID, data)
}

// lazyDM lazily gets the id of the direct message channel with the invoking
// user.
func (r *wrappedReplier) lazyDM(ctx *plugin.Context) error {
	if !r.dmID.IsValid() {
		c, err := r.s.CreatePrivateChannel(ctx.Author.ID)
		if err != nil {
			return err
		}

		r.dmID = c.ID
	}

	return nil
}
