package i18nutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/i18n"
)

// mockLocalizer is a copy of mock.Localizer, used to prevent import cycles.
type mockLocalizer struct {
	t *testing.T

	def      string
	onReturn map[i18n.Term]string
	errOn    map[i18n.Term]struct{}
}

func newMockedLocalizer(t *testing.T) *mockLocalizer {
	return &mockLocalizer{
		t:        t,
		onReturn: make(map[i18n.Term]string),
		errOn:    make(map[i18n.Term]struct{}),
	}
}

func (l *mockLocalizer) on(term i18n.Term, response string) *mockLocalizer {
	l.onReturn[term] = response
	return l
}

func (l *mockLocalizer) build() *i18n.Localizer {
	m := i18n.NewManager(func(lang string) i18n.LangFunc {
		return func(term i18n.Term, _ map[string]interface{}, _ interface{}) (string, error) {
			r, ok := l.onReturn[term]
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
