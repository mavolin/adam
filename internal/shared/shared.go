// Package shared provides place to share variables used internally.
package shared

import (
	"github.com/mavolin/adam/internal/embedbuilder"
)

const Whitespace = " \n"

// ErrorEmbed is the global error embed template.
// See errors.SetErrorEmbed and errors.NewErrorEmbed for more information.
var ErrorEmbed = embedbuilder.New().
	WithTitlel(errorTitle).
	WithColor(0xff5a5a)

// InfoEmbed is the global info embed template
// See errors.SetInfoEmbed and errors.NewInfoEmbed for more information.
var InfoEmbed = embedbuilder.New().
	WithTitlel(infoTitle).
	WithColor(0x6eb7b1)
