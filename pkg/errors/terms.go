package errors

import "github.com/mavolin/adam/pkg/localization"

// termErrorID is the error id footer of an InternalError.
const termErrorID = "errors.error_id"

var (
	// errorTitleConfig is the localization.Config used to generate the title
	// of an error message.
	errorTitleConfig = localization.QuickFallbackConfig("errors.title", "Error")

	// defaultInternalDescConfig is the localization.Config used by default as
	// description for an InternalError.
	defaultInternalDescConfig = localization.QuickFallbackConfig("errors.internal.description.default",
		"Oh no! Something went wrong and I couldn't finish executing your command. I've informed my team and they'll "+
			"get on fixing the bug asap.")

	defaultRestrictionDescConfig = localization.QuickFallbackConfig("errors.restriction.description.default",
		"👮 You are not allowed to use this command.")

	// infoTitleConfig is the localization.Config used to generate the title
	// of an info message.
	infoTitleConfig = localization.QuickFallbackConfig("info.title", "Info")
)

// errorIDPlaceholders is the placeholders struct for the errors.error_id
// config.
type errorIDPlaceholders struct {
	// ErrorID is the id of the error.
	ErrorID string
}
