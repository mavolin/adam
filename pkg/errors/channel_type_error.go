package errors

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// ChannelTypeError is the error returned if a command is invoked in a channel
// that is not supported by that command.
type ChannelTypeError struct {
	// Allowed are the plugin.ChannelTypes that the command supports.
	Allowed plugin.ChannelTypes
}

var _ Error = new(ChannelTypeError)

// NewChannelTypeError creates a new ChannelTypeError with the passed allowed
// plugin.ChannelTypes.
func NewChannelTypeError(allowed plugin.ChannelTypes) *ChannelTypeError {
	return &ChannelTypeError{Allowed: allowed}
}

// Description returns the description containing the types of channels this
// command may be used in.
func (e *ChannelTypeError) Description(l *i18n.Localizer) (desc string) {
	switch {
	// ----- singles -----
	case e.Allowed == plugin.GuildTextChannels:
		desc, _ = l.Localize(channelTypeErrorGuildText)
	case e.Allowed == plugin.GuildNewsChannels:
		desc, _ = l.Localize(channelTypeErrorGuildNews)
	case e.Allowed == plugin.DirectMessages:
		desc, _ = l.Localize(channelTypeErrorDirectMessage)
	// ----- combos -----
	case e.Allowed == plugin.GuildChannels:
		desc, _ = l.Localize(channelTypeErrorGuild)
	case e.Allowed == (plugin.DirectMessages | plugin.GuildTextChannels):
		desc, _ = l.Localize(channelTypeErrorDirectMessageAndGuildText)
	case e.Allowed == (plugin.DirectMessages | plugin.GuildNewsChannels):
		desc, _ = l.Localize(channelTypeErrorDirectMessageAndGuildNews)
	default:
		desc, _ = l.Localize(channelTypeErrorFallback)
	}

	return
}

func (e *ChannelTypeError) Error() string {
	return "channel type error"
}

func (e *ChannelTypeError) Is(target error) bool {
	var typedTarget *ChannelTypeError
	if !As(target, &typedTarget) {
		return false
	}

	return e.Allowed == typedTarget.Allowed
}

// Handle handles the ChannelTypeError.
// By default it sends an error message stating the allowed channel types.
func (e *ChannelTypeError) Handle(s *state.State, ctx *plugin.Context) error {
	return HandleChannelTypeError(e, s, ctx)
}

var HandleChannelTypeError = func(ierr *ChannelTypeError, s *state.State, ctx *plugin.Context) error {
	embed := ErrorEmbed.Clone().
		WithDescription(ierr.Description(ctx.Localizer))

	_, err := ctx.ReplyEmbedBuilder(embed)
	return err
}
