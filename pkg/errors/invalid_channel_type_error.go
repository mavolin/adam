package errors

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

// InvalidChannelTypeError is the error returned if a command is invoked in
// an channel that is not supported by that command.
type InvalidChannelTypeError struct {
	// AllowedChannelTypes are the plugin.ChannelTypes that the command supports.
	AllowedChannelTypes plugin.ChannelTypes
}

// NewInvalidChannelTypeError creates a new InvalidChannelTypeError with the
// passed allowed plugin.ChannelTypes.
func NewInvalidChannelTypeError(allowed plugin.ChannelTypes) *InvalidChannelTypeError {
	return &InvalidChannelTypeError{
		AllowedChannelTypes: allowed,
	}
}

// Description returns the description containing the types of channels this
// command may be used in.
func (e *InvalidChannelTypeError) Description(l *localization.Localizer) (desc string) {
	switch {
	// ----- singles -----
	case e.AllowedChannelTypes == plugin.GuildTextChannels:
		desc, _ = l.Localize(channelTypeErrorGuildText)
	case e.AllowedChannelTypes == plugin.GuildNewsChannels:
		desc, _ = l.Localize(channelTypeErrorGuildNews)
	case e.AllowedChannelTypes == plugin.DirectMessages:
		desc, _ = l.Localize(channelTypeErrorDirectMessage)
	// ----- combos -----
	case e.AllowedChannelTypes == plugin.GuildChannels:
		desc, _ = l.Localize(channelTypeErrorGuild)
	case e.AllowedChannelTypes == (plugin.DirectMessages | plugin.GuildTextChannels):
		desc, _ = l.Localize(channelTypeErrorDirectMessageAndGuildText)
	case e.AllowedChannelTypes == (plugin.DirectMessages | plugin.GuildNewsChannels):
		desc, _ = l.Localize(channelTypeErrorDirectMessageAndGuildNews)
	default:
		desc, _ = l.Localize(channelTypeErrorFallback)
	}

	return
}

func (e *InvalidChannelTypeError) Error() string {
	return "invalid channel type error"
}

func (e *InvalidChannelTypeError) Is(target error) bool {
	casted, ok := target.(*InvalidChannelTypeError)
	if !ok {
		return false
	}

	return e.AllowedChannelTypes == casted.AllowedChannelTypes
}

// Handle sends an error message stating the allowed channel types permissions.
func (e *InvalidChannelTypeError) Handle(_ *state.State, ctx *plugin.Context) (err error) {
	embed := newErrorEmbedBuilder(ctx.Localizer).
		WithDescription(e.Description(ctx.Localizer))

	_, err = ctx.ReplyEmbedBuilder(embed)
	return
}
