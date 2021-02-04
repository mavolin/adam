package plugin

import (
	"github.com/mavolin/adam/pkg/i18n"
)

type (
	// Module is the abstraction of a module.
	//
	// A default for a simple module can be found in impl/module.
	Module interface {
		ModuleMeta

		// Commands returns the subcommands of the module.
		Commands() []Command
		// Modules returns the submodules of the module.
		Modules() []Module
	}

	// ModuleMeta is the abstraction of the Module's meta data.
	//
	// Default implementations can be found in impl/module.
	ModuleMeta interface {
		// GetName returns the name of the module
		GetName() string
		// GetShortDescription returns an optional one-sentence description of
		// the module.
		GetShortDescription(l *i18n.Localizer) string
		// GetLongDescription returns an option long description of the module.
		GetLongDescription(l *i18n.Localizer) string
	}
)
