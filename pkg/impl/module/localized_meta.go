package module

import (
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// LocalizedMeta is the localized, implementation of the plugin.ModuleMeta
// interface.
type LocalizedMeta struct {
	// Name is the name of the module.
	// It may not contain whitespace or dots.
	Name string
	// ShortDescription is an optional short description of the module.
	ShortDescription *i18n.Config
	// LongDescription is an optional long description of the module.
	LongDescription *i18n.Config
}

var _ plugin.ModuleMeta = LocalizedMeta{}

func (m LocalizedMeta) GetName() string {
	return m.Name
}

func (m LocalizedMeta) GetShortDescription(l *i18n.Localizer) string {
	desc, err := l.Localize(m.ShortDescription)
	if err != nil {
		return ""
	}

	return desc
}

func (m LocalizedMeta) GetLongDescription(l *i18n.Localizer) string {
	desc, err := l.Localize(m.LongDescription)
	if err != nil {
		return ""
	}

	return desc
}
