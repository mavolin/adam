package localization

type (
	// Func is the function used to retrieve a Translator.
	// If there is no data for the passed language or the passed language is
	// invalid, it should either use to a default language or fall back to the
	// DefaultTranslator.
	// However, it should never return nil.
	Func func(lang string) Translator

	// Translator is the interface used to translate various messages, as found
	// in the library.
	// An example of an implementation may be found in pkg/impl/translator.
	Translator interface{}
)
