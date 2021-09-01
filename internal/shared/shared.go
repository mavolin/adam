// Package shared provides place to share variables used internally.
package shared

import (
	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/internal/embedbuilder"
	"github.com/mavolin/adam/pkg/i18n"
)

const Whitespace = " \n"

// NewErrorEmbed creates a new error embed.
// It is initialized in package error.
// See errors.SetErrorEmbed and errors.NewErrorEmbed for more information.
var NewErrorEmbed = func(l *i18n.Localizer, desc string) (discord.Embed, error) {
	return embedbuilder.New().
		WithTitlel(ErrorTitle).
		WithColor(0xff5a5a).
		WithDescription(desc).
		Build(l)
}

// NewInfoEmbed creates a new info embed.
// It is initialized in package error.
// See errors.SetInfoEmbed and errors.NewInfoEmbed for more information.
var NewInfoEmbed = func(l *i18n.Localizer, desc string) (discord.Embed, error) {
	return embedbuilder.New().
		WithTitlel(InfoTitle).
		WithColor(0x6eb7b1).
		WithDescription(desc).
		Build(l)
}
