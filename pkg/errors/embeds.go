package errors

import (
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

// ErrorEmbed is the embedutil.Builder used to create new error embeds.
// Errors may fill the description of the Embed or add fields.
//
// It should be made sure, that Embed builder always succeeds in building, as
// otherwise errors might not get sent.
// This means if localizing the Embed, fallbacks should be defined.
var ErrorEmbed = embedutil.NewBuilder().
	WithSimpleTitlel(errorTitle).
	WithColor(0xff5a5a)

// InfoEmbed is the embedutil.Builder used to create new info embeds.
// Infos may fill the description of the Embed or add fields.
//
// It should be made sure, that Embed builder always succeeds in building, as
// otherwise errors might not get sent.
// This means if localizing the Embed, fallbacks should be defined.
var InfoEmbed = embedutil.NewBuilder().
	WithSimpleTitlel(infoTitle).
	WithColor(0x6eb7b1)
