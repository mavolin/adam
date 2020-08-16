package errors

import "github.com/mavolin/adam/pkg/localization"

const (
	// termErrorTitle is the title of an error message.
	termErrorTitle = "errors.title"
	// termInternalDescription is the default description of an InternalError.
	termInternalDescription = "errors.internal.description"
	// termInfoTitle is the title of an info message.
	termInfoTitle = "info.title"
)

var (
	// errorTitleConfig is the localization.Config used to generate the title
	// of an error message.
	errorTitleConfig = localization.Config{
		Term: termErrorTitle,
		Fallback: localization.Fallback{
			Other: "Error",
		},
	}

	// infoTitleConfig is the localization.Config used to generate the title
	// of an info message.
	infoTitleConfig = localization.Config{
		Term: termInfoTitle,
		Fallback: localization.Fallback{
			Other: "Info",
		},
	}
)
