package errors

import (
	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/internal/embedbuilder"
	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/i18n"
)

var (
	errorEmbed = embedbuilder.New().
			WithTitlel(shared.ErrorTitle).
			WithColor(0xff5a5a)

	infoEmbed = embedbuilder.New().
			WithTitlel(shared.InfoTitle).
			WithColor(0x6eb7b1)
)

// This may be a bit ugly, but it allows two things:
// a) Package msgbuilder can import package plugin, and
// b) plugin-related errors can actually be provided by package plugin, and not
// error.
func init() {
	shared.NewErrorEmbed = func(l *i18n.Localizer, desc string) (discord.Embed, error) {
		return NewErrorEmbed().
			WithDescription(desc).
			Build(l)
	}

	shared.NewInfoEmbed = func(l *i18n.Localizer, desc string) (discord.Embed, error) {
		return NewInfoEmbed().
			WithDescription(desc).
			Build(l)
	}
}

// SetErrorEmbed updates the *msgbuilder.EmbedBuilder used to create new error
// embeds.
//
// It should be made sure that EmbedBuilder always succeeds in building, as
// otherwise errors might not get sent.
// This means if localizing the Embed, fallbacks should be defined.
//
// SetErrorEmbed is not safe for concurrent use and should not be called after
// the bot has been started.
func SetErrorEmbed(b *embedbuilder.Builder) {
	errorEmbed = b
}

// NewErrorEmbed creates a new *msgbuilder.EmbedBuilder that can be used to
// build error embeds.
func NewErrorEmbed() *embedbuilder.Builder {
	return errorEmbed.Clone()
}

// SetInfoEmbed updates the *msgbuilder.EmbedBuilder used to create new info
// embeds.
//
// It should be made sure that Embed builder always succeeds in building, as
// otherwise errors might not get sent.
// This means if localizing the Embed, fallbacks should be defined.
//
// SetInfoEmbed is not safe for concurrent use and should not be called after
// the bot has been started.
func SetInfoEmbed(b *embedbuilder.Builder) {
	infoEmbed = b
}

// NewInfoEmbed creates a new *msgbuilder.EmbedBuilder that can be used to
// build info embeds.
func NewInfoEmbed() *embedbuilder.Builder {
	return infoEmbed.Clone()
}
