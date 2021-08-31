// Package shared provides place to share variables used internally;
package shared

import "github.com/mavolin/adam/pkg/utils/msgbuilder"

const Whitespace = " \n"

// ErrorEmbed is the shared errors *msgbuilder.EmbedBuilder.
// See errors.SetErrorEmbed and errors.NewErrorEmbed for more information.
var ErrorEmbed = msgbuilder.NewEmbed().
	WithTitlel(errorTitle).
	WithColor(0xff5a5a)

// InfoEmbed is the shared info *msgbuilder.EmbedBuilder.
// See errors.SetInfoEmbed and errors.NewInfoEmbed for more information.
var InfoEmbed = msgbuilder.NewEmbed().
	WithTitlel(infoTitle).
	WithColor(0x6eb7b1)
