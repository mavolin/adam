package errors

import (
	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

// SetErrorEmbed updates the *embedutil.Builder used to create new error
// embeds.
//
// It should be made sure, that Builder always succeeds in building, as
// otherwise errors might not get sent.
// This means if localizing the Embed, fallbacks should be defined.
//
// SetErrorEmbed is not safe for concurrent use and should not be called after
// the bot has been started.
func SetErrorEmbed(b *embedutil.Builder) {
	shared.ErrorEmbed = b
}

// NewErrorEmbed creates a new *embedutil.Builder that can be used to build
// error embeds.
func NewErrorEmbed() *embedutil.Builder {
	return shared.ErrorEmbed.Clone()
}

// SetInfoEmbed updates the *embedutil.Builder used to create new info embeds.
//
// It should be made sure, that Embed builder always succeeds in building, as
// otherwise errors might not get sent.
// This means if localizing the Embed, fallbacks should be defined.
//
// SetInfoEmbed is not safe for concurrent use and should not be called after
// the bot has been started.
func SetInfoEmbed(b *embedutil.Builder) {
	shared.ErrorEmbed = b
}

// NewInfoEmbed creates a new *embedutil.Builder that can be used to build
// info embeds.
func NewInfoEmbed() *embedutil.Builder {
	return shared.InfoEmbed.Clone()
}
