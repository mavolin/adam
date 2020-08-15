package localization

// Manager manages the distribution of Localizers.
type Manager struct {
	f Func
}

// NewManager creates a new Manager using the passed Func.
//
// If the passed Func is nil, Localizers generated from this Manager will
// always use their fallback messages.
func NewManager(f Func) *Manager {
	if f == nil {
		f = func(lang string) LangFunc { return nil }
	}

	return &Manager{
		f: f,
	}
}

// Localizer returns a localizer for the passed language.
func (m *Manager) Localizer(lang string) *Localizer {
	return &Localizer{
		f:    m.f(lang),
		Lang: lang,
	}
}
