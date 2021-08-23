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

		// GetCommands returns the subcommands of the module.
		GetCommands() []Command
		// GetModules returns the submodules of the module.
		GetModules() []Module
	}

	// ModuleMeta is the abstraction of the Module's meta data.
	//
	// Default implementations can be found in impl/module.
	ModuleMeta interface {
		// GetName returns the name of the module.
		// It may not contain whitespace or dots.
		GetName() string
		// GetShortDescription returns an optional one-sentence description of
		// the module.
		GetShortDescription(l *i18n.Localizer) string
		// GetLongDescription returns an option long description of the module.
		GetLongDescription(l *i18n.Localizer) string
	}
)

type (
	// ResolvedModule is a resolved module as returned by a Provider.
	// In contrast to the regular Module abstraction, ResolvedModule's plugins
	// reflect the plugins provided by all modules with the same ID, i.e. a
	// plugin with the same name provided through different bot.PluginProvider.
	ResolvedModule interface {
		// Parent is the parent of this module.
		// If the module is top-level Parent will be nil.
		Parent() ResolvedModule
		// Sources contains the Modules this module is based upon.
		// If the module is top-level, this will be empty.
		Sources() []SourceModule
		// ID is the identifier of the module.
		ID() ID
		// Name is the name of the module.
		Name() string
		// ShortDescription returns an optional one-sentence description of the
		// module.
		ShortDescription(l *i18n.Localizer) string
		// LongDescription returns an option thorough description of the
		// module.
		//
		// If the module only provides a short description, it will be returned
		// instead.
		LongDescription(l *i18n.Localizer) string

		// IsHidden specifies if all Sources are hidden.
		// A source module is considered hidden if all of it's commands and
		// modules are hidden as well.
		IsHidden() bool

		// Commands are the subcommands of the module.
		// They are sorted in ascending order by name.
		Commands() []ResolvedCommand
		// Modules are the submodules of the module.
		// They are sorted in ascending order by name.
		Modules() []ResolvedModule

		// FindCommand finds the command with the given name inside this module.
		// A name can either be the actual name of a command, or an alias.
		//
		// If there is no command with the given name, nil will be returned.
		FindCommand(name string) ResolvedCommand
		// FindModule finds the module with the given name inside the module.
		//
		// If there is no module with the given name, nil will be returned.
		FindModule(name string) ResolvedModule
	}

	// SourceModule contains the parent Modules of a ResolvedModule.
	SourceModule struct {
		// SourceName is the name of the plugin source that
		// provided the module.
		SourceName string
		// Modules contains the parents of the ResolvedModule.
		// They are sorted in ascending order from the most distant to the
		// closest parent.
		Modules []Module
	}
)
