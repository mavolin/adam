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

		// IsHidden specifies whether this module will be hidden from the help
		// page.
		IsHidden() bool
		// GetDefaultChannelTypes returns the ChannelTypes required to use this
		// module.
		//
		// Commands can overwrite this, by setting a custom ChannelTypes.
		GetDefaultChannelTypes() ChannelTypes
		// IsRestricted checks if the user calling the command is restricted
		// from using this module.
		// If the bot lacks one ore more permissions command execution will
		// stop with an errors.InsufficientPermissionsError.
		//
		// Commands can overwrite this, by returning a non-nil RestrictionFunc.
		// To remove a RestrictionFunc defined by a parent without defining a
		// new one use restriction.None.
		//
		// Note that that direct messages may also pass this, if the passed
		// permissions only require constant.DMPermissions.
		GetDefaultRestrictionFunc() RestrictionFunc
		// GetDefaultThrottler returns the Throttler for the module.
		// The throttler is used for all subcommands and submodules of the
		// module.
		// However, a command or module can overwrite this, by setting its own
		// Throttler.
		//
		// To remove a Throttler defined by a parent without defining a new
		// one use throttler.None.
		GetDefaultThrottler() Throttler
	}
)
