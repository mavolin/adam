package i18n

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/i18n"
)

type Localizer struct {
	t *testing.T

	def   string
	on    map[i18n.Term]string
	errOn map[i18n.Term]struct{}
}

// NewLocalizer creates a new Localizer.
// If a term is not found, Localizer will panic.
//nolint:thelper
func NewLocalizer(t *testing.T) *Localizer {
	return &Localizer{
		t:     t,
		on:    make(map[i18n.Term]string),
		errOn: make(map[i18n.Term]struct{}),
	}
}

// NewLocalizerWithDefault creates a new Localizer using the passed default.
// If a term is not found, Localizer will return the default value.
//nolint:thelper
func NewLocalizerWithDefault(t *testing.T, def string) *Localizer {
	return &Localizer{
		t:     t,
		def:   def,
		on:    make(map[i18n.Term]string),
		errOn: make(map[i18n.Term]struct{}),
	}
}

// On adds the passed response for the passed term to the localizer.
func (l *Localizer) On(term i18n.Term, response string) *Localizer {
	l.on[term] = response
	return l
}

// ErrorOn returns an error whenever the passed term is requested.
func (l *Localizer) ErrorOn(term i18n.Term) *Localizer {
	l.errOn[term] = struct{}{}
	return l
}

// Clone creates a clone of the localizer.
//nolint:thelper
func (l *Localizer) Clone(t *testing.T) *Localizer {
	on := make(map[i18n.Term]string, len(l.on))
	errOn := make(map[i18n.Term]struct{}, len(l.on))

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

// Build builds the localizer.
func (l *Localizer) Build() *i18n.Localizer {
	return i18n.NewLocalizer("dev", func(term i18n.Term, _ map[string]interface{}, _ interface{}) (string, error) {
		l.t.Helper()

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
	})
}
