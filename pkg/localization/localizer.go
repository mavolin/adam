package localization

type (
	// Config is a data struct that contains all information needed to create
	// a localized message.
	Config struct {
		// Term is the key of the translation.
		Term string
		// Placeholders are the filled placeholders of the translation.
		Placeholders Placeholders
		// Plural is a number or a string containing such, that is used to
		// identify if the message should be pluralized or not.
		//
		// If nil, the other message should be used.
		Plural interface{}

		// Fallback is the fallback used if the LangFunc is nil or the
		// LangFunc returned an error.
		Fallback Fallback
	}

	// Placeholders is the type used for placeholders.
	Placeholders map[string]interface{}

	// Fallback is the English fallback used if a translation is not available.
	// The message is created using go's text/template system and left and
	// right delimiters are {{ and }} respectively.
	Fallback struct {
		// One is the singular form of the fallback message, if there is any.
		One string
		// Other is the plural form of the fallback message.
		// This is also the default form, meaning if no pluralization is needed
		// this field should be used.
		Other string
	}
)

// genTranslation attempts to generate the translation using the passed
// Placeholders and the passed plural data.
func (f Fallback) genTranslation(placeholders Placeholders, plural interface{}) (string, error) {
	if plural != nil { // we have plural information
		if isOne, err := isOne(plural); err != nil { // attempt to check if plural is == 1
			return "", err
		} else if isOne {
			s, err := fillTemplate(f.One, placeholders)
			return s, err
		}
	}

	// no plural information or plural was != 1

	s, err := fillTemplate(f.Other, placeholders)
	return s, err
}

// Localizer is a translator for a specific language.
// It provides multiple utility functions and wraps a LangFunc.
type Localizer struct {
	// f is the LangFunc used to create translations.
	f LangFunc
	// Lang is the language the Localizer is translating to.
	// This does not account for possible fallbacks being used, because
	// the wanted language was not available.
	Lang string
}

// Localize generates a localized message using the passed config.
// c.Term must be set.
func (l *Localizer) Localize(c Config) (s string, err error) {
	if l.f != nil { // try the user-defined translator first, if there is one
		s, err = l.f(c.Term, c.Placeholders, c.Plural)
		if err == nil {
			return
		}
	}

	// otherwise use fallback if there is;
	// checking other suffices as it will always be set if there is a fallback
	if c.Fallback.Other == "" {
		return c.Term, &NoTranslationGeneratedError{
			Term: c.Term,
		}
	}

	s, err = c.Fallback.genTranslation(c.Placeholders, c.Plural)
	if err != nil {
		return c.Term, err
	}

	return
}

// LocalizeTerm is a short for
//		l.Localize(localization.Config{
//			Term: term,
//		})
func (l *Localizer) LocalizeTerm(term string) (string, error) { return l.Localize(Config{Term: term}) }

// MustLocalize is the same as Localize, but it panics if there is an error.
func (l *Localizer) MustLocalize(c Config) string {
	s, err := l.Localize(c)
	if err != nil {
		panic(err)
	}

	return s
}

// MustLocalizeTerm is the same as LocalizeTerm, but it panics if there is an
// error.
func (l *Localizer) MustLocalizeTerm(term string) string {
	s, err := l.LocalizeTerm(term)
	if err != nil {
		panic(err)
	}

	return s
}
