// Package shared provides variables and constants used across multiple
// packages.
package shared

import (
	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/pkg/i18n"
)

const Whitespace = " \n"

// ErrorEmbedTemplate is the global error embed template.
// See errors.SetErrorEmbed and errors.NewErrorEmbed for more information.
var ErrorEmbedTemplate = func(l *i18n.Localizer) discord.Embed {
	title, _ := l.Localize(errorTitle) // we have a fallback
	return discord.Embed{Title: title, Color: 0xff5a5a}
}

// InfoEmbedTemplate is the global info embed template.
// See errors.SetInfoEmbed and errors.NewInfoEmbed for more information.
var InfoEmbedTemplate = func(l *i18n.Localizer) discord.Embed {
	title, _ := l.Localize(infoTitle) // we have a fallback
	return discord.Embed{Title: title, Color: 0x6eb7b1}
}
