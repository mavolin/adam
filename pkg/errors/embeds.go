package errors

import (
	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/i18n"
)

// SetErrorEmbedTemplate updates the template function used to create new error
// embeds.
//
// The returned embed must have a title.
// However, that title may be overwritten by the caller.
//
// SetErrorEmbed is not safe for concurrent use and should not be called after
// the bot has been started.
func SetErrorEmbedTemplate(tmpl func(*i18n.Localizer) discord.Embed) {
	shared.ErrorEmbedTemplate = tmpl
}

// NewErrorEmbed creates a new discord.Embed that can be used to build error
// embeds.
func NewErrorEmbed(l *i18n.Localizer) discord.Embed {
	return shared.ErrorEmbedTemplate(l)
}

// SetInfoEmbedTemplate updates the template function used to create new info
// embeds.
//
// The returned embed must have a title.
// However, that title may be overwritten by the caller.
//
// SetInfoEmbedTemplate is not safe for concurrent use and should not be called
// after the bot has been started.
func SetInfoEmbedTemplate(tmpl func(localizer *i18n.Localizer) discord.Embed) {
	shared.InfoEmbedTemplate = tmpl
}

// NewInfoEmbed creates a new discord.Embed that can be used to build info
// embeds.
func NewInfoEmbed(l *i18n.Localizer) discord.Embed {
	return shared.InfoEmbedTemplate(l)
}
