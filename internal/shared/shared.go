// Package shared provides place to share variables used internally;
package shared

import "github.com/mavolin/adam/pkg/utils/embedutil"

const Whitespace = " \n"

// ErrorEmbed is the shared errors *embedutil.Builder.
// See errors.SetErrorEmbed and errors.NewErrorEmbed for more information.
var ErrorEmbed = embedutil.NewBuilder().
	WithTitlel(errorTitle).
	WithColor(0xff5a5a)

// InfoEmbed is the shared info *embedutil.Builder.
// See errors.SetInfoEmbed and errors.NewInfoEmbed for more information.
var InfoEmbed = embedutil.NewBuilder().
	WithTitlel(infoTitle).
	WithColor(0x6eb7b1)
