package replier

import (
	"sync"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/discorderr"
)

type wrappedReplier struct {
	s           *state.State
	inlineReply bool

	dmID   discord.ChannelID
	dmErr  error
	dmOnce sync.Once
}

var _ plugin.Replier = new(wrappedReplier)

// WrapState wraps the passed state and id of the invoking user into a
// plugin.Replier.
// If inlineReply is set to true, messages will reference the invoke, unless
// MessageReference is non-nil.
//
// The returned plugin.Replier is safe for concurrent use.
func WrapState(s *state.State, inlineReply bool) plugin.Replier {
	return &wrappedReplier{s: s, inlineReply: inlineReply}
}

func (r *wrappedReplier) Reply(ctx *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	perms, err := ctx.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, plugin.NewBotPermissionsError(discord.PermissionSendMessages)
	}

	if data.Reference == nil && r.inlineReply {
		data.Reference = &discord.MessageReference{MessageID: ctx.Message.ID}
	}

	msg, err := r.s.SendMessageComplex(ctx.ChannelID, data)
	if discorderr.Is(discorderr.As(err), discorderr.UnknownChannel) {
		// user deleted channel
		return nil, errors.Abort
	}

	return msg, errors.WithStack(err)
}

func (r *wrappedReplier) ReplyDM(ctx *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	dmID, err := r.lazyDMID(ctx)
	if err != nil {
		return nil, err
	}

	msg, err := r.s.SendMessageComplex(dmID, data)
	return msg, errors.WithStack(err)
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

	msg, err := r.s.EditMessageComplex(ctx.ChannelID, messageID, data)
	if discorderr.Is(discorderr.As(err), discorderr.UnknownChannel) {
		// user deleted channel
		return nil, errors.Abort
	}

	return msg, errors.WithStack(err)
}

func (r *wrappedReplier) EditDM(
	ctx *plugin.Context, messageID discord.MessageID, data api.EditMessageData,
) (*discord.Message, error) {
	dmID, err := r.lazyDMID(ctx)
	if err != nil {
		return nil, err
	}

	msg, err := r.s.EditMessageComplex(dmID, messageID, data)
	return msg, errors.WithStack(err)
}

// lazyDMID lazily gets the id of the direct message channel with the invoking
// user.
func (r *wrappedReplier) lazyDMID(ctx *plugin.Context) (discord.ChannelID, error) {
	r.dmOnce.Do(func() {
		c, err := r.s.CreatePrivateChannel(ctx.Author.ID)
		r.dmErr = err
		if err == nil {
			r.dmID = c.ID
		}
	})

	return r.dmID, r.dmErr
}
