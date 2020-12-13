package i18n

// Localizer is a translator for a specific language.
// It provides multiple utility functions and wraps a LangFunc.
type Localizer struct {
	// f is the LangFunc used to create translations.
	f LangFunc
	// Lang is the language the Localizer is translating to.
	// This does not account for possible fallbacks being used, because
	// the wanted language was not available.
	Lang string

	// defaultPlaceholders is a list of placeholders that is automatically
	// added to every config.
	defaultPlaceholders map[string]interface{}
}

// WithPlaceholder adds the passed default placeholder to the Localizer.
func (l *Localizer) WithPlaceholder(key string, val interface{}) {
	if l.defaultPlaceholders == nil {
		l.defaultPlaceholders = map[string]interface{}{key: val}

		return
	}

	l.defaultPlaceholders[key] = val
}

// WithPlaceholders adds the passed default placeholders to the
// Localizer.
func (l *Localizer) WithPlaceholders(p map[string]interface{}) {
	if l.defaultPlaceholders == nil {
		l.defaultPlaceholders = p
		return
	}

	for k, v := range p {
		l.defaultPlaceholders[k] = v
	}
}

// Localize generates a localized message using the passed config.
// c.NewTermConfig must be set.
func (l *Localizer) Localize(c *Config) (s string, err error) {
	placeholders, err := c.placeholdersToMap()
	if err != nil {
		return string(c.Term), err
	}

	if placeholders == nil && len(l.defaultPlaceholders) > 0 {
		placeholders = make(map[string]interface{}, len(l.defaultPlaceholders))
	}

	for k, v := range l.defaultPlaceholders {
		if _, ok := placeholders[k]; ok {
			continue
		}

		//goland:noinspection GoNilness // see if above for-loop
		placeholders[k] = v
	}

	if l.f != nil { // try the user-defined translator first, if there is one
		s, err = l.f(c.Term, placeholders, c.Plural)
		if err == nil {
			return s, err
		}
	}

	// otherwise use fallback if there is;
	// checking other suffices as it will always be set if there is a fallback
	if c.Fallback.Other == "" {
		return string(c.Term), NewNoTranslationGeneratedError(c.Term)
	}

	s, err = c.Fallback.genTranslation(placeholders, c.Plural)
	if err != nil {
		return string(c.Term), err
	}

	return s, err
}

// LocalizeTerm is a short for
//		l.Localize(i18n.Config{
//			Term: term,
//		})
func (l *Localizer) LocalizeTerm(term Term) (string, error) {
	return l.Localize(term.AsConfig())
}

// MustLocalize is the same as Localize, but it panics if there is an error.
func (l *Localizer) MustLocalize(c *Config) string {
	s, err := l.Localize(c)
	if err != nil {
		panic(err)
	}

	return s
}

// MustLocalizeTerm is the same as LocalizeTerm, but it panics if there is an
// error.
func (l *Localizer) MustLocalizeTerm(term Term) string {
	return l.MustLocalize(term.AsConfig())
}
