package i18n

import (
	"errors"
	"reflect"

	"github.com/iancoleman/strcase"

	"github.com/mavolin/adam/internal/errorutil"
)

// ErrPlaceholders gets returned, if the type of Placeholders is
// invalid.
var ErrPlaceholders = errors.New("i18n: placeholders must be of type map[string]string or struct")

// EmptyConfig is the Config used as a replacement for an empty string.
var EmptyConfig = new(Config)

// Config is a data struct that contains all information needed to create
// a localized message.
// However, Configs don't necessarily produce localized text:
// If the Term is empty, or the Func returns with an error, the Config will
// resort to its Fallback.
type Config struct {
	// Term is the key of the translation.
	//
	// Term may be empty in which case Fallback will always be used.
	Term Term
	// Placeholders contains the placeholder data.
	// This can either be a map[string]string or a struct.
	//
	// Structs
	//
	// If you use a struct the Localizer will to convert the
	// name of the fields to snake_case.
	// However, if you want to use a custom name for the keys, you
	// can use the `i18n:"myname"` struct tag.
	Placeholders interface{}
	// Plural is a number or a string containing such, that is used to
	// identify if the message should be pluralized or not.
	//
	// If nil, the other message should be used.
	Plural interface{}

	// Fallback is the fallback used if the Func is nil or the
	// Func returned an error.
	Fallback Fallback
}

// NewTermConfig is a utility function that can be used to inline term-only
// Configs.
func NewTermConfig(term Term) *Config {
	return term.AsConfig()
}

// NewStaticConfig returns a static *Config that always produces s, even
// if '{{' and '}}' are used.
func NewStaticConfig(s string) *Config {
	c := NewFallbackConfig("", s)
	c.Fallback.literal = true

	return c
}

// NewFallbackConfig is a utility function that can be used to inline
// term-only Configs with a fallback.
func NewFallbackConfig(term Term, fallback string) *Config {
	c := &Config{
		Term:     term,
		Fallback: Fallback{Other: fallback},
	}

	return c
}

// WithPlaceholders returns a copy of the Config with the passed placeholders
// set.
func (c Config) WithPlaceholders(placeholders interface{}) *Config {
	c.Placeholders = placeholders
	return &c
}

// WithPlural returns a copy of the Config with the passed plural set.
func (c Config) WithPlural(plural interface{}) *Config {
	c.Plural = plural
	return &c
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
		return nil, errorutil.WithStack(ErrPlaceholders)
	}

	placeholders := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		tv := t.Field(i)

		if !fv.CanInterface() {
			continue
		}

		key := tv.Tag.Get("i18n")

		if key == "" {
			key = strcase.ToSnake(tv.Name)
		}

		placeholders[key] = fv.Interface()
	}

	return placeholders, nil
}

// Term is a unique identifier of a localized message.
type Term string

// AsConfig wraps the term in a Config.
func (t Term) AsConfig() *Config {
	return &Config{Term: t}
}

// Fallback is the English fallback used if a translation is not available.
// The message is created using go's text/template system.
// `{{` and `}}` are used as delimiters, respectively.
type Fallback struct {
	// One is the singular form of the fallback message, if there is any.
	One string
	// Other is the plural form of the fallback message.
	// This is also the default form, meaning if no pluralization is needed
	// this field should be used.
	Other string

	literal bool
}

// genTranslation attempts to generate the translation using the passed
// Placeholders and the passed plural data.
func (f *Fallback) genTranslation(placeholderData map[string]interface{}, plural interface{}) (string, error) {
	if plural != nil { // we have plural information
		if isOne, err := isOne(plural); err != nil { // attempt to check if plural is == 1
			return "", err
		} else if isOne {
			return fillTemplate(f.One, placeholderData)
		}
	}

	if f.literal {
		return f.Other, nil
	}

	// no plural information or plural was != 1
	return fillTemplate(f.Other, placeholderData)
}
