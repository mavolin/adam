package plugin

import (
	"github.com/mavolin/adam/pkg/localization"
)

type mockLocalizer struct {
	def      string
	onReturn map[localization.Term]string
}

// newMockedLocalizer creates a new mockLocalizer.
// If a term is not found, mockLocalizer will panic.
func newMockedLocalizer() *mockLocalizer {
	return &mockLocalizer{
		onReturn: make(map[localization.Term]string),
	}
}

// newMockedLocalizer creates a new mockLocalizer using the passed default.
// If a term is not found, mockLocalizer will return the default value.
func newMockedLocalizerWithDefault(def string) *mockLocalizer {
	return &mockLocalizer{
		def:      def,
		onReturn: make(map[localization.Term]string),
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

			if l.def == "" {
				panic("unexpected localization requested for term " + term)
			}

			return l.def, nil
		}
	})

	return m.Localizer("")
}
