package errors

import (
	"github.com/mavolin/adam/internal/constant"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

// ErrorEmbed is the embedutil.Builder used to create new error embeds.
// Errors may fill the description of the embed or add fields.
var ErrorEmbed = embedutil.NewBuilder().
	WithSimpleTitlel(errorTitle).
	WithColor(constant.ErrorColor)

// InfoEmbed is the embedutil.Builder used to create new info embeds.
// Errors may fill the description of the embed or add fields.
var InfoEmbed = embedutil.NewBuilder().
	WithSimpleTitlel(infoTitle).
	WithColor(constant.InfoColor)
