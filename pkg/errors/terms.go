package errors

import "github.com/mavolin/adam/pkg/localization"

const (
	// termError is the term used for
	termError               = "errors.error"
	termInternalDescription = "errors.internal.description"
)

var errorTitleConfig = localization.Config{
	Term: termError,
	Fallback: localization.Fallback{
		Other: "Error",
	},
}
