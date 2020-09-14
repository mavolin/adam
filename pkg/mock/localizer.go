package mock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/localization"
)

type Localizer struct {
	t *testing.T

	def   string
	on    map[localization.Term]string
	errOn map[localization.Term]struct{}
}

// NewLocalizer creates a new Localizer.
// If a term is not found, Localizer will panic.
func NewLocalizer(t *testing.T) *Localizer {
	return &Localizer{
		t:     t,
		on:    make(map[localization.Term]string),
		errOn: make(map[localization.Term]struct{}),
	}
}

// NewLocalizer creates a new Localizer using the passed default.
// If a term is not found, Localizer will return the default value.
func NewLocalizerWithDefault(def string) *Localizer {
	return &Localizer{
		def:   def,
		on:    make(map[localization.Term]string),
		errOn: make(map[localization.Term]struct{}),
	}
}

// NewNoOpLocalizer returns a localization.Localizer that returns only errors.
func NewNoOpLocalizer() *localization.Localizer {
	m := localization.NewManager(func(lang string) localization.LangFunc {
		return func(_ localization.Term, _ map[string]interface{}, _ interface{}) (string, error) {
			return "", errors.New("error")
		}
	})

	return m.Localizer("")
}

// On adds the passed response for the passed term to the localizer.
func (l *Localizer) On(term localization.Term, response string) *Localizer {
	l.on[term] = response
	return l
}

// ErrorOn returns an error whenever the passed term is requested.
func (l *Localizer) ErrorOn(term localization.Term) *Localizer {
	l.errOn[term] = struct{}{}
	return l
}

// Clone creates a copy of the localizer.
func (l *Localizer) Clone(t *testing.T) *Localizer {
	on := make(map[localization.Term]string, len(l.on))
	errOn := make(map[localization.Term]struct{}, len(l.on))

	for k, v := range l.on {
		on[k] = v
	}

	for k := range l.errOn {
		errOn[k] = struct{}{}
	}

	return &Localizer{
		t:     t,
		def:   l.def,
		on:    on,
		errOn: errOn,
	}
}

// Build build the localizer.
func (l *Localizer) Build() *localization.Localizer {
	m := localization.NewManager(func(lang string) localization.LangFunc {
		return func(term localization.Term, _ map[string]interface{}, _ interface{}) (string, error) {
			r, ok := l.on[term]
			if ok {
				return r, nil
			}

			_, ok = l.errOn[term]
			if ok {
				return r, errors.New("error")
			}

			if l.def == "" {
				assert.Failf(l.t, "unexpected call to Localize", "unknown term %s", term)

				return string(term), errors.New("unknown term")
			}

			return l.def, nil
		}
	})

	return m.Localizer("")
}
