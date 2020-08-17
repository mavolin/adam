package localization

import (
	"errors"
	"reflect"

	"github.com/iancoleman/strcase"
)

var ErrInvalidPlaceholders = errors.New("placeholders must be of type map[string]string or struct")

// Config is a data struct that contains all information needed to create
// a localized message.
type Config struct {
	// Term is the key of the translation.
	Term string
	// Placeholders contains the placeholder data.
	// This can either be a map[string]string or a struct (see section
	// Structs for further info).
	//
	// Structs
	//
	// If you use a struct the Localizer will to convert the
	// name of the fields to snake_case.
	// However, if you want to use a custom name for the keys, you
	// can use the `localization:"myname"` struct tag.
	Placeholders interface{}
	// Plural is a number or a string containing such, that is used to
	// identify if the message should be pluralized or not.
	//
	// If nil, the other message should be used.
	Plural interface{}

	// Fallback is the fallback used if the LangFunc is nil or the
	// LangFunc returned an error.
	Fallback Fallback
}

// Term is a utility function that can be used to inline term-only Configs.
func Term(term string) Config {
	return Config{
		Term: term,
	}
}

// NewFallbackConfig is a utility function that can be used to inline
// term-only Configs with a fallback.
func NewFallbackConfig(term, fallback string) Config {
	return Config{
		Term: term,
		Fallback: Fallback{
			Other: fallback,
		},
	}
}

// WithPlaceholders returns a copy of the Config with the passed Placeholders
// set.
func (c Config) WithPlaceholders(placeholders interface{}) Config {
	c.Placeholders = placeholders
	return c
}

func (c Config) placeholdersToMap() (map[string]interface{}, error) {
	if c.Placeholders == nil {
		return nil, nil
	}

	if v, ok := c.Placeholders.(map[string]interface{}); ok {
		return v, nil
	}

	v := reflect.ValueOf(c.Placeholders)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	if v.Kind() != reflect.Struct {
		return nil, withStack(ErrInvalidPlaceholders)
	}

	placeholders := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		tv := t.Field(i)

		if !fv.CanInterface() {
			continue
		}

		key := tv.Tag.Get("localization")

		if key == "" {
			key = strcase.ToSnake(tv.Name)
		}

		placeholders[key] = fv.Interface()
	}

	return placeholders, nil
}

// Fallback is the English fallback used if a translation is not available.
// The message is created using go's text/template system and left and
// right delimiters are {{ and }} respectively.
type Fallback struct {
	// One is the singular form of the fallback message, if there is any.
	One string
	// Other is the plural form of the fallback message.
	// This is also the default form, meaning if no pluralization is needed
	// this field should be used.
	Other string
}

// genTranslation attempts to generate the translation using the passed
// Placeholders and the passed plural data.
func (f Fallback) genTranslation(placeholderData map[string]interface{}, plural interface{}) (string, error) {
	if plural != nil { // we have plural information
		if isOne, err := isOne(plural); err != nil { // attempt to check if plural is == 1
			return "", err
		} else if isOne {
			return fillTemplate(f.One, placeholderData)
		}
	}

	// no plural information or plural was != 1

	return fillTemplate(f.Other, placeholderData)
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
	placeholders, err := c.placeholdersToMap()
	if err != nil {
		return c.Term, err
	}

	if l.f != nil { // try the user-defined translator first, if there is one
		s, err = l.f(c.Term, placeholders, c.Plural)
		if err == nil {
			return
		}
	}

	// otherwise use fallback if there is;
	// checking other suffices as it will always be set if there is a fallback
	if c.Fallback.Other == "" {
		return c.Term, NewNoTranslationGeneratedError(c.Term)
	}

	s, err = c.Fallback.genTranslation(placeholders, c.Plural)
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
