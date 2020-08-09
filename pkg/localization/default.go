package localization

// DefaultFunc is the Func that can be used to retrieve the DefaultTranslator.
func DefaultFunc(string) Translator { return DefaultTranslator }

// DefaultTranslator is the default translator.
// It will be used if the user doesn't localize, but is also a suitable
// fallback for custom Translator implementations.
var DefaultTranslator = new(defaultTranslator)

type defaultTranslator struct{}
