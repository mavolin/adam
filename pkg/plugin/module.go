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
		GetShortDescription(l *localization.Localizer) string
		// GetLongDescription returns an option long description of the module.
		GetLongDescription(l *localization.Localizer) string
		// IsHidden specifies whether this module will be hidden from the help
		// page.
		//
		// If set to true, all submodules and subcommands will be hidden as
		// well.
		IsHidden() bool
		// GetDefaultChannelTypes returns the ChannelTypes required to use this
		// module.
		//
		// Commands can overwrite this, by setting a custom ChannelTypes.
		GetDefaultChannelTypes() ChannelTypes
		// GetDefaultBotPermissions get the permissions needed to use this
		// module.
		//
		// Commands can overwrite this, by setting their bot permissions to a
		// non-nil value.
		GetDefaultBotPermissions() *discord.Permissions
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
		// one use throttling.None.
		GetDefaultThrottler() Throttler
	}

	// RegisteredModule is the abstraction of a module as returned by a
	// Provider.
	// In contrast to the regular module abstraction, RegisteredModule will
	// return data that takes into account it's parents settings.
	RegisteredModule interface {
		// Parent returns the parent of this module.
		// It will return nil, nil, if this module is top-level.
		//
		// In any other case will always return valid data, even if error !=
		// nil.
		// It is also guaranteed that the original parent of the module, i.e.
		// the module that provides this command is included, if there is one.
		//
		// However, all runtime plugin providers that returned an error won't
		// be included, and their error will be returned wrapped in a
		// bot.RuntimePluginProviderError.
		// If multiple errors occur, a errors.MultiError filled with
		// bot.RuntimePluginProviderErrors will be returned.
		Parent() (RegisteredModule, error)
		// Identifier returns the identifier of the module.
		Identifier() Identifier

		// Name returns the name of the module.
		Name() string
		// ShortDescription returns an optional one-sentence description of the
		// module.
		ShortDescription(l *localization.Localizer) string
		// LongDescription returns an option thorough description of the
		// module.
		LongDescription(l *localization.Localizer) string
		// IsHidden specifies whether this module and all it's submodules and
		// commands should be hidden from help messages.
		IsHidden() bool

		// Commands returns the subcommands of the module.
		Commands() []RegisteredCommand
		// Modules returns the submodules of the module.
		Modules() []RegisteredModule
		// FindCommand finds the command with the given invoke inside this
		// module.
		// For example if this is the top-level 'abc' module and you want to
		// find the command with the identifier '.abc.def', you should search
		// for 'def'.
		FindCommand(invoke string) RegisteredCommand
		// FindModule finds the module with the given invoke inside this
		// module.
		// For example if this is the top-level 'abc' module and you want to
		// find the module with the identifier '.abc.def', you should search
		// for def'.
		FindModule(invoke string) RegisteredModule
	}
)
