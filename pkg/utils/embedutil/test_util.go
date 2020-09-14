package embedutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/localization"
)

// mockLocalizer is a copy of mock.Localizer, used to prevent import cycles.
type mockLocalizer struct {
	t *testing.T

	def      string
	onReturn map[localization.Term]string
	errOn    map[localization.Term]struct{}
}

func newMockedLocalizer(t *testing.T) *mockLocalizer {
	return &mockLocalizer{
		t:        t,
		onReturn: make(map[localization.Term]string),
		errOn:    make(map[localization.Term]struct{}),
	}
}

func (l *mockLocalizer) on(term localization.Term, response string) *mockLocalizer {
	l.onReturn[term] = response
	return l
}

func (l *mockLocalizer) build() *localization.Localizer {
	m := localization.NewManager(func(lang string) localization.LangFunc {
		return func(term localization.Term, _ map[string]interface{}, _ interface{}) (string, error) {
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
