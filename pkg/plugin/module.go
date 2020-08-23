package plugin

import (
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/localization"
)

type (
	// Module is the abstraction of a module.
	//
	// A default for a simple module can be found in impl/module.
	Module interface {
		// GetMeta returns the meta information of the module.
		Meta() ModuleMeta

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
		GetShortDescription(l *localization.Localizer) (string, error)
		// GetLongDescription returns an option long description of the module.
		GetLongDescription(l *localization.Localizer) (string, error)
		// IsHidden specifies whether this module will be hidden from the help
		// page.
		//
		// If set to true, all submodules and subcommands will be hidden as
		// well.
		IsHidden() bool
		// GetChannelTypes returns the ChannelTypes required to use this module.
		//
		// Commands can overwrite this, by setting a custom ChannelTypes.
		GetChannelTypes() ChannelTypes
		// GetBotPermissions get the permissions needed to use this module.
		//
		// Commands can overwrite this, by setting custom BotPermissions.
		GetBotPermissions() discord.Permissions
		// IsRestricted checks if the user calling the command is restricted
		// from using this module.
		//
		// Commands can overwrite this, by returning a non-nil RestrictionFunc.
		GetRestrictionFunc() RestrictionFunc
		// GetThrottling returns the ThrottlingOptions for the module.
		// This defines how often all commands and submodules in this module
		// together may be used.
		//
		// If either of the fields in ThrottlingOptions is zero value, the
		// module won't be throttled.
		GetThrottling() ThrottlingOptions
	}
)
