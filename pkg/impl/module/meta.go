package module

import (
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// Meta is the static, unlocalized, implementation of the plugin.ModuleMeta
// interface.
type Meta struct {
	// Name is the name of the module.
	// It may not contain whitespace or dots.
	Name string
	// ShortDescription is an optional short description of the module.
	ShortDescription string
	// LongDescription is an optional long description of the module.
	LongDescription string
}

var _ plugin.ModuleMeta = Meta{}

func (m Meta) GetName() string                            { return m.Name }
func (m Meta) GetShortDescription(*i18n.Localizer) string { return m.ShortDescription }
func (m Meta) GetLongDescription(*i18n.Localizer) string  { return m.LongDescription }
