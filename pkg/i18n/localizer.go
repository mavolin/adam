package i18n

import (
	"errors"

	"github.com/mavolin/adam/internal/errorutil"
)

// ErrNilConfig is the error returned if a nil config is given to
// Localizer.Localize.
var ErrNilConfig = errors.New("i18n: cannot translate nil Config")

// Localizer is a translator for a specific language.
// It provides multiple utility functions and wraps a Func.
//
// The zero value of a Localizer is a fallback localizer.
type Localizer struct {
	// f is the Func used to create translations.
	f Func
	// Lang is the language the Localizer is translating to.
	// This does not account for possible fallbacks being used, because
	// the wanted language was not available.
	//
	// It is unique to every language and dialect.
	//
	// If Lang is empty, the localizer is a fallback localizer.
	Lang string

	// defaultPlaceholders is a list of placeholders that is automatically
	// added to every config.
	defaultPlaceholders map[string]interface{}
}

// NewLocalizer creates a new Localizer for the passed language that generates
// text using the passed Func.
//
// Lang must be unique to every language and dialect used.
func NewLocalizer(lang string, f Func) *Localizer {
	return &Localizer{f: f, Lang: lang}
}

// NewFallbackLocalizer creates a new *Localizer that always uses the fallback
// messages.
func NewFallbackLocalizer() *Localizer {
	return new(Localizer)
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
	if c == nil {
		return "", errorutil.WithStack(ErrNilConfig)
	}

	placeholders, err := c.placeholdersToMap()
	if err != nil {
		return "", err
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
			return s, nil
		}
	}

	// otherwise use fallback if there is;
	// checking other suffices as it will always be set if there is a fallback
	if len(c.Fallback.Other) == 0 && !c.Fallback.empty {
		return "", newNoTranslationGeneratedError(c.Term)
	}

	return c.Fallback.genTranslation(placeholders, c.Plural)
}

// LocalizeTerm is a short for
//		l.Localize(&i18n.Config{
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
