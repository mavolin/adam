package plugin

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
)

type (
	// Command is the abstraction of a command.
	//
	// Defaults for simple commands can be found in impl/command.
	Command interface {
		CommandMeta

		// Invoke calls the command.
		//
		// Possible first return values are
		//	• uint, uint8, uint16, uint32, uint64
		//	• int, int8, int16, int32, int64
		// 	• float32, float64
		//	• string
		//	• discord.Embed
		//	• *embedutil.Builder
		//	• api.SendMessageData
		//	• localization.Term
		//	• localization.Config
		//	• any type implementing Response
		//	• nil for no response
		//
		// Error Handling
		//
		// If Invoke returns an error it will be handed to the error handler
		// of the bot.
		// As a special case if both return values are non-nil, both the
		// response and the error will be handled.
		//
		// Panic Handling
		//
		// Similarly, if Invoke panics the panic will be handled by the
		// PanicHandler of the executing bot.
		Invoke(s *state.State, ctx *Context) (interface{}, error)
	}

	// CommandMeta is the abstraction of the Command's meta data.
	//
	// Default implementations can be found in impl/command.
	CommandMeta interface {
		// GetName gets the name of the command.
		GetName() string
		// GetAliases returns the optional aliases of the command.
		GetAliases() []string
		// GetArgs returns the ArgConfig of the command.
		//
		// If this is nil, the command will accept no arguments and flags.
		GetArgs() ArgConfig
		// GetShortDescription returns an optional one-sentence description
		// of the command.
		GetShortDescription(l *localization.Localizer) string
		// GetLongDescription returns an optional long description of the
		// command.
		GetLongDescription(l *localization.Localizer) string
		// GetExamples returns optional example usages of the command.
		GetExamples(l *localization.Localizer) []string
		// IsHidden specifies whether this command will be hidden in the help
		// page.
		IsHidden() bool
		// GetChannelType returns the ChannelTypes this command may be invoked
		// in.
		//
		// Setting this overrides ChannelTypes defined by the parent.
		//
		// If this 0, the parents ChannelTypes will be used.
		GetChannelTypes() ChannelTypes
		// GetBotPermissions gets the permissions the bot needs to execute this
		// command.
		// If the bot lacks one ore more permissions command execution will
		// stop with an errors.InsufficientPermissionsError.
		//
		// Setting this to a non-nil value overrides bot permissions defined by
		// parents.
		//
		// Note that that direct messages may also pass this, if the passed
		// permissions only require constant.DMPermissions.
		GetBotPermissions() *discord.Permissions
		// GetRestrictionFunc checks if the user is restricted from using the
		// command.
		//
		// Setting this will override restrictions defined by the parent.
		//
		// If they are restricted, a errors.RestrictionError should be
		// returned.
		//
		// If the RestrictionFunc returns an error that implements
		// RestrictionErrorWrapper, it will be properly wrapped.
		GetRestrictionFunc() RestrictionFunc
		// GetThrottler returns the Throttler for the command.
		//
		// Setting this will override the Throttler defined by the parent.
		// To remove a Throttler defined by a parent without defining a new
		// one use throttling.None.
		GetThrottler() Throttler
	}

	// RegisteredCommand is the abstraction of a command as returned by a
	// Provider.
	// In contrast to the regular command abstraction, RegisteredCommand will
	// return data that takes into account it's parents settings.
	RegisteredCommand interface {
		// Parent returns the parent of this command.
		// It will return nil, nil, if this command is top-level.
		//
		// In any other case will always return valid data, even if error !=
		// nil.
		// It is also  guaranteed that the original parent of the command, i.e.
		// the module that provides this command is included, if there is one.
		//
		// However, all runtime plugin providers that returned an error won't
		// be included, and their error will be returned wrapped in a
		// bot.RuntimePluginProviderError.
		// If multiple errors occur, a errors.MultiError filled with
		// bot.RuntimePluginProviderErrors will be returned.
		Parent() (RegisteredModule, error)
		// Identifier returns the identifier of the command.
		Identifier() Identifier

		// Name returns the name of the command.
		Name() string
		// Aliases returns the optional aliases of the command.
		Aliases() []string
		// Args returns the argument configuration of the command.
		//
		// If this is nil, the command accepts no arguments.
		Args() ArgConfig
		// ShortDescription returns an optional one-sentence description of the
		// command.
		ShortDescription(l *localization.Localizer) string
		// LongDescription returns an optional thorough description of the
		// command.
		LongDescription(l *localization.Localizer) string
		// Examples returns optional examples for the command.
		Examples(l *localization.Localizer) []string
		// IsHidden specifies whether to show this command in the help.
		IsHidden() bool
		// ChannelTypes returns the ChannelTypes this command can be run in.
		//
		// If the command itself did not define some, ChannelTypes will return
		// the ChannelTypes of the closest parent that has defaults defined.
		ChannelTypes() ChannelTypes
		// BotPermissions returns the permissions this command needs to
		// execute.
		// If the command itself did not define some, BotPermissions will
		// return the permissions of the closest parent that has a default
		// defined.
		BotPermissions() discord.Permissions
		// IsRestricted returns whether or not this command is restricted.
		IsRestricted(*state.State, *Context) error
		// Throttler returns the Throttler of this command.
		//
		// If the command itself did not define one, Throttler will return the
		// Throttler of the closest parent.
		Throttler() Throttler

		// Invoke invokes the command.
		// See Command.Invoke for more details.
		Invoke(*state.State, *Context) (interface{}, error)
	}
)

// Response is the interface that a type can implement, to be used as a return
// value of a Command.Invoke call.
type Response interface {
	// Send sends the response using the passed state.State and Context.
	Send(s *state.State, ctx *Context) error
}
