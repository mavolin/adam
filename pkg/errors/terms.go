package errors

import "github.com/mavolin/adam/pkg/localization"

const (
	// termErrorTitle is the title of an error message.
	termErrorTitle = "errors.title"
	// termInternalDescription is the default description of an InternalError.
	termInternalDescription = "errors.internal.description"
	// termErrorID is the error id footer of an InternalError.
	termErrorID = "errors.error_id"
	// termInfoTitle is the title of an info message.
	termInfoTitle = "info.title"
)

var (
	// errorTitleConfig is the localization.Config used to generate the title
	// of an error message.
	errorTitleConfig = localization.QuickFallbackConfig(termErrorTitle, "Error")

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
