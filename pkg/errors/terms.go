package errors

import "github.com/mavolin/adam/pkg/localization"

const (
	// termErrorTitle is the title of an error message.
	termErrorTitle = "errors.title"
	// termInternalDefaultDescription is the default description of an
	// InternalError.
	termInternalDefaultDescription = "errors.internal.description.default"
	// termRestrictionDefaultDescription is the default description of an
	// RestrictionError.
	termRestrictionDefaultDescription = "errors.restriction.description.default"
	// termErrorID is the error id footer of an InternalError.
	termErrorID = "errors.error_id"
	// termInfoTitle is the title of an info message.
	termInfoTitle = "info.title"
)

var (
	// errorTitleConfig is the localization.Config used to generate the title
	// of an error message.
	errorTitleConfig = localization.QuickFallbackConfig(termErrorTitle, "Error")

	// defaultInternalDescConfig is the localization.Config used by default as
	// description for an InternalError.
	defaultInternalDescConfig = localization.QuickFallbackConfig(termInternalDefaultDescription,
		"Oh no! Something went wrong and I couldn't finish executing your command. I've informed my team and they'll "+
			"get on fixing the bug asap.")

	defaultRestrictionDescConfig = localization.QuickFallbackConfig(termRestrictionDefaultDescription,
		"👮 You are not allowed to use this command.")

	// infoTitleConfig is the localization.Config used to generate the title
	// of an info message.
	infoTitleConfig = localization.QuickFallbackConfig(termInfoTitle, "Info")
)

// errorIDPlaceholders is the placeholders struct for the errors.error_id
// config.
type errorIDPlaceholders struct {
	// ErrorID is the id of the error.
	ErrorID string
}
