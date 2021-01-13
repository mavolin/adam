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

	// Hidden specifies whether this module should be hidden from the help
	// message.
	//
	// All subcommands and submodules will be hidden as well.
	Hidden bool

	// DefaultChannelTypes are the plugin.ChannelTypes the used as default, for
	// all submodule and -commands that don't specify some.
	//
	// If this is not set, the channel types of the parent will be used.
	DefaultChannelTypes plugin.ChannelTypes
	// DefaultRestrictions are the restrictions used as default, for all
	// submodule and -commands that don't specify some.
	//
	// If this is not set, the restrictions of the parent will be used.
	DefaultRestrictions plugin.RestrictionFunc
	// DefaultThrottler is the plugin.Throttler used as default, for all
	// submodule and -commands that don't specify some.
	//
	// If this is not set, the throttler of the parent will be used.
	DefaultThrottler plugin.Throttler
}

var _ plugin.ModuleMeta = LocalizedMeta{}

func (m LocalizedMeta) GetName() string { return m.Name }

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

	if len(desc) > 0 {
		return desc
	}

	return m.GetShortDescription(l)
}

func (m LocalizedMeta) IsHidden() bool                              { return m.Hidden }
func (m LocalizedMeta) GetDefaultChannelTypes() plugin.ChannelTypes { return m.DefaultChannelTypes }

func (m LocalizedMeta) GetDefaultRestrictionFunc() plugin.RestrictionFunc {
	return m.DefaultRestrictions
}

func (m LocalizedMeta) GetDefaultThrottler() plugin.Throttler { return m.DefaultThrottler }
