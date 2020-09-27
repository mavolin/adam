package errors

import (
	"github.com/mavolin/adam/internal/constant"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

// ErrorEmbed is the embedutil.Builder used to create new error embeds.
// Errors may fill the description of the embed or add fields.
//
// It should be made sure, that embed builder always succeeds in building, as
// otherwise errors might not get sent.
// This means if localizing the embed, fallbacks should be defined.
var ErrorEmbed = embedutil.NewBuilder().
	WithSimpleTitlel(errorTitle).
	WithColor(constant.ErrorColor)

// InfoEmbed is the embedutil.Builder used to create new info embeds.
// Infos may fill the description of the embed or add fields.
//
// It should be made sure, that embed builder always succeeds in building, as
// otherwise errors might not get sent.
// This means if localizing the embed, fallbacks should be defined.
var InfoEmbed = embedutil.NewBuilder().
	WithSimpleTitlel(infoTitle).
	WithColor(constant.InfoColor)
