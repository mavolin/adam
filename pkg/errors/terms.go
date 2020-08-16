package errors

import "github.com/mavolin/adam/pkg/localization"

const (
	// termError is the title of error embeds.
	termError = "errors.error"
	// termInternalDescription is the default description of an InternalError.
	termInternalDescription = "errors.internal.description"
)

// errorTitleConfig is the localization.Config used to generate the title of an
// error embed.
var errorTitleConfig = localization.Config{
	Term: termError,
	Fallback: localization.Fallback{
		Other: "Error",
	},
}
