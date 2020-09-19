package localization

import (
	"fmt"

	"github.com/mavolin/adam/internal/errorutil"
)

// NoTranslationGeneratedError gets returned if the Localizer is unable
// to produce a translation with the data available, i.e. if neither the
// underlying LangFunc, nor the Fallback return a non-error value.
type NoTranslationGeneratedError struct {
	Term Term
	s    errorutil.Stack
}

// NewNoTranslationGeneratedError creates a new NoTranslationGeneratedError
// for the passed term.
func NewNoTranslationGeneratedError(term Term) *NoTranslationGeneratedError {
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

func (e *NoTranslationGeneratedError) StackTrace() []uintptr { return e.s }

// Is checks if the error matches the passed error.
func (e *NoTranslationGeneratedError) Is(target error) bool {
	casted, ok := target.(*NoTranslationGeneratedError)
	if !ok {
		return false
	}

	return e.Term == casted.Term
}

// stackError is a copy of errors.InternalError to prevent an import cycle.
type stackError struct {
	cause error
	s     errorutil.Stack
}

func withStack(err error) error {
	if err == nil {
		return nil
	}

	return &stackError{
		cause: err,
		s:     errorutil.GenerateStackTrace(1),
	}
}

func (s *stackError) Error() string { return s.cause.Error() }
func (s *stackError) Unwrap() error { return s.cause }
