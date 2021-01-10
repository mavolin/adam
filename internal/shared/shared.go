// Package shared provides place to share variables used internally;
package shared

import "github.com/mavolin/adam/pkg/utils/embedutil"

// ErrorEmbed is the shared *embedutil.Builder.
// See errors.SetErrorEmbed and errors.NewErrorEmbed for more information.
var ErrorEmbed = embedutil.NewBuilder().
	WithSimpleTitlel(errorTitle).
	WithColor(0xff5a5a)

// InfoEmbed is the shared info *embedutil.Builder.
// See errors.SetInfoEmbed and errors.NewInfoEmbed for more information.
var InfoEmbed = embedutil.NewBuilder().
	WithSimpleTitlel(infoTitle).
	WithColor(0x6eb7b1)
