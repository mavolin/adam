package errors

import (
	"github.com/mavolin/adam/internal/embedbuilder"
	"github.com/mavolin/adam/internal/shared"
)

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
	shared.ErrorEmbed = b
}

// NewErrorEmbed creates a new *msgbuilder.EmbedBuilder that can be used to
// build error embeds.
func NewErrorEmbed() *embedbuilder.Builder {
	return shared.ErrorEmbed.Clone()
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
	shared.InfoEmbed = b
}

// NewInfoEmbed creates a new *msgbuilder.EmbedBuilder that can be used to
// build info embeds.
func NewInfoEmbed() *embedbuilder.Builder {
	return shared.InfoEmbed.Clone()
}
