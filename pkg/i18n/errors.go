package i18n

import (
	"errors"
	"fmt"

	"github.com/mavolin/adam/internal/errorutil"
)

// NoTranslationGeneratedError gets returned if the Localizer is unable
// to produce a translation with the data available, i.e. if neither the
// underlying Func, nor the Fallback return a non-error value.
type NoTranslationGeneratedError struct {
	Term Term
	s    errorutil.StackTrace
}

// NewNoTranslationGeneratedError creates a new NoTranslationGeneratedError
// for the passed term.
func newNoTranslationGeneratedError(term Term) *NoTranslationGeneratedError {
	return &NoTranslationGeneratedError{
		Term: term,
		s:    errorutil.GenerateStackTrace(1),
	}
}

// Error generates a error message.
func (e *NoTranslationGeneratedError) Error() string {
	if e.Term != "" {
		return fmt.Sprintf("unable to generate a translation for term '%s'", e.Term)
	}

	return "unable to generate a translation"
}

func (e *NoTranslationGeneratedError) StackTrace() errorutil.StackTrace { return e.s }

// Is checks if the error matches the passed error.
func (e *NoTranslationGeneratedError) Is(target error) bool {
	var typedTarget *NoTranslationGeneratedError
	if !errors.As(target, &typedTarget) {
		return false
	}

	return e.Term == typedTarget.Term
}
