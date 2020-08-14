package localization

import "fmt"

// NoTranslationGeneratedError gets returned if the Localizer is unable
// to produce a translation with the data available, i.e. if neither the
// underlying LangFunc, nor the Fallback return a non-error value.
type NoTranslationGeneratedError struct {
	Term string
}

// Error generates a error message.
func (e *NoTranslationGeneratedError) Error() string {
	if e.Term != "" {
		return fmt.Sprintf("unable to generate a translation for term '%s'", e.Term)
	}

	return "unable to generate a translation"
}

// Is checks if the error matches the passed error.
func (e *NoTranslationGeneratedError) Is(target error) bool {
	casted, ok := target.(*NoTranslationGeneratedError)
	if !ok {
		return false
	}

	return e.Term == casted.Term
}
