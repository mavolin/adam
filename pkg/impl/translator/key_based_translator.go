package translator

import "github.com/mavolin/adam/pkg/localization"

type (
	// FieldTranslatorFunc is a language specific function for localizing
	// messages.
	FieldTranslatorFunc func(c Config) (string, error)

	// Config is a data struct that contains all information needed to create
	// a localized message.
	Config struct {
		// Term is the unique key of the translation.
		Term string
		// Placeholders are the filled placeholders of the translation.
		Placeholders map[string]interface{}
		// PluralKey points to a number within the Placeholders map, that
		// specifies whether or not the message shall be pluralized.
		//
		// If the string is empty, the message doesn't need to be pluralized.
		PluralKey string
	}
)

// NewFieldBasedFunc creates a localization.Func for a field-based Translator.
// The passed function is used to retrieve the appropriate
// FieldTranslatorFunc for the language.
//
// If there is no localization data for the passed language or the language
// is invalid, it should return a default.
// If nil is returned, the translator will fall back to the
// localization.DefaultTranslator.
func NewFieldBasedFunc(f func(lang string) FieldTranslatorFunc) localization.Func {
	return func(lang string) localization.Translator {
		return &fieldBasedTranslator{t: f(lang)}
	}
}

// fieldBasedTranslator is the actual implementation of a field-based translator.
type fieldBasedTranslator struct {
	// t is the FieldTranslatorFunc for the language, that is called to
	// translate.
	t FieldTranslatorFunc
}
