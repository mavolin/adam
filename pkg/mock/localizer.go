package mock

import "github.com/mavolin/adam/pkg/localization"

type Localizer struct {
	def string
	on  map[string]string
}

// NewLocalizer creates a new Localizer.
// If a term is not found, Localizer will panic.
func NewLocalizer() *Localizer {
	return &Localizer{
		on: make(map[string]string),
	}
}

// NewLocalizer creates a new Localizer using the passed default.
// If a term is not found, Localizer will return the default value.
func NewLocalizerWithDefault(def string) *Localizer {
	return &Localizer{
		def: def,
		on:  make(map[string]string),
	}
}

func (l *Localizer) On(term, response string) *Localizer {
	l.on[term] = response
	return l
}

func (l *Localizer) Clone() *Localizer {
	cp := make(map[string]string, len(l.on))

	for k, v := range l.on {
		cp[k] = v
	}

	return &Localizer{
		def: l.def,
		on:  cp,
	}
}

func (l *Localizer) Build() *localization.Localizer {
	m := localization.NewManager(func(lang string) localization.LangFunc {
		return func(term string, _ map[string]interface{}, _ interface{}) (string, error) {
			r, ok := l.on[term]
			if !ok {
				if l.def == "" {
					panic("unexpected localization for term " + term)
				}

				return l.def, nil
			}

			return r, nil
		}
	})

	return m.Localizer("")
}
