package i18n

import (
	"errors"
	"fmt"

	"github.com/mavolin/adam/internal/errorutil"
)

// LocalizationError gets returned if the Localizer is unable
// to produce a translation with the data available, i.e. if neither the
// underlying Func, nor the Fallback return a non-error value.
type LocalizationError struct {
	Term Term
	s    errorutil.StackTrace
}

// newLocalizationError creates a new LocalizationError for the passed term.
func newLocalizationError(term Term) *LocalizationError {
	return &LocalizationError{
		Term: term,
		s:    errorutil.GenerateStackTrace(1),
	}
}

// Error generates a error message.
func (e *LocalizationError) Error() string {
	if e.Term != "" {
		return fmt.Sprintf("i18n: unable to generate a translation for term '%s'", e.Term)
	}

	return "i18n: unable to generate a translation"
}

func (e *LocalizationError) StackTrace() errorutil.StackTrace { return e.s }

// Is checks if the error matches the passed error.
func (e *LocalizationError) Is(target error) bool {
	var typedTarget *LocalizationError
	if !errors.As(target, &typedTarget) {
		return false
	}

	return e.Term == typedTarget.Term
}
