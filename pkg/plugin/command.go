package plugin

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

type (
	// Command is the abstraction of a command.
	//
	// Defaults for simple commands can be found in impl/command.
	Command interface {
		CommandMeta

		// Invoke calls the command.
		// The first return value is the reply sent to the user in the channel
		// they invoked the command in.
		//
		// Possible first return values are:
		//	• uint, uint8, uint16, uint32, uint64
		//	• int, int8, int16, int32, int64
		// 	• float32, float64
		//	• string
		//	• discord.Embed, *discord.Embed
		//	• *msgbuilder.Builder
		//	• api.SendMessageData
		//	• i18n.Term
		//	• *i18n.Config
		//	• any type implementing Reply
		//	• nil for no reply
		//
		// All other values will be captured through a *bot.ReplyTypeError.
		//
		// Error Handling
		//
		// If Invoke returns an error it will be handed down the middleware
		// chain until it reaches the bot's ErrorHandler.
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
		// It may not contain whitespace or dots.
		GetName() string
		// GetAliases returns the optional aliases of the command.\
		// They may not contain whitespace or dots.
		GetAliases() []string
		// GetShortDescription returns an optional short description
		// of the command.
		GetShortDescription(l *i18n.Localizer) string
		// GetLongDescription returns an optional long description of the
		// command.
		GetLongDescription(l *i18n.Localizer) string

		// GetArgs returns the ArgConfig of the command.
		//
		// If this is nil, the command will accept no arguments and flags.
		GetArgs() ArgConfig
		// GetArgParser returns the optional custom ArgParser of the command.
		GetArgParser() ArgParser
		// GetExampleArgs returns optional example arguments of the command.
		GetExampleArgs(l *i18n.Localizer) ExampleArgs

		// IsHidden specifies whether this command will be hidden in the help
		// page.
		IsHidden() bool
		// GetChannelTypes returns the ChannelTypes this command may be invoked
		// in.
		// If this is 0, AllChannels will be used.
		GetChannelTypes() ChannelTypes
		// GetBotPermissions gets the permissions the bot needs to execute this
		// command.
		// If the bot lacks one ore more permissions command execution will
		// stop with an errors.InsufficientPermissionsError.
		//
		// Note that that direct messages may also pass this, if the passed
		// permissions only require permutil.DMPermissions.
		GetBotPermissions() discord.Permissions
		// IsRestricted checks if the user is restricted from using the
		// command.
		//
		// If they are restricted, a *plugin.RestrictionError should be
		// returned.
		IsRestricted(s *state.State, ctx *Context) error
		// GetThrottler returns the Throttler for the command.
		GetThrottler() Throttler
	}

	// ExampleArgs is a struct containing a set of exemplary arguments and
	// flags.
	// They are formatted using their ArgParser's FormatArgs method.
	ExampleArgs []struct {
		// Args contains the example arguments.
		Args []string
		// Flags is a map of exemplary flags.
		Flags map[string]string
	}
)

func (a ExampleArgs) BaseType(*i18n.Localizer) ExampleArgs { return a }

// ResolvedCommand is a resolved command as returned by a Provider.
type ResolvedCommand interface {
	// Parent returns the parent of this command.
	// The returned ResolvedModule may not consists of all modules that share
	// the same namespace, if some plugin sources are unavailable.
	// Check PluginProvider.UnavailableProviders() to check if that is the
	// case.
	//
	// In any case the module will contain the built-in module and the module
	// that provides the command.
	Parent() ResolvedModule
	// SourceName returns the name of the plugin source that provides the
	// command.
	//
	// If the command is built-in, ProviderName will be set to BuiltInSource.
	SourceName() string
	// Source returns the original Command this command is based on.
	Source() Command
	// SourceParents returns the original parent Modules in ascending order
	// from the most distant to the closest parent.
	//
	// If the command is top-level, SourceParents will return nil.
	SourceParents() []Module

	// ID returns the identifier of the command.
	ID() ID
	// Name returns the name of the command.
	Name() string
	// Aliases returns the optional aliases of the command.
	Aliases() []string
	// ShortDescription returns an optional brief description of the command.
	ShortDescription(*i18n.Localizer) string
	// LongDescription returns an optional long description of the command.
	//
	// If the command only provides a short description, it will be returned
	// instead.
	LongDescription(*i18n.Localizer) string

	// Args returns the argument configuration of the command.
	//
	// If this is nil, the command accepts no arguments.
	Args() ArgConfig
	// ArgParser returns ArgParser of the command.
	// In contrast to Command's GetArgParser equivalent, ArgParser will return
	// the global ArgParser if the command did not define one.
	ArgParser() ArgParser

	// ExampleArgs returns optional example arguments of the command.
	ExampleArgs(*i18n.Localizer) ExampleArgs
	// Examples returns the command's example arguments prefixed with their
	// invoke.
	// Invoke and example arguments are separated by a space.
	Examples(*i18n.Localizer) []string

	// IsHidden returns whether to show this command in the help.
	IsHidden() bool
	// ChannelTypes are the ChannelTypes this command can be run in.
	//
	// If the command itself did not define some, ChannelTypes will be
	// AllChannels.
	ChannelTypes() ChannelTypes
	// BotPermissions returns the permissions the command needs to execute.
	BotPermissions() discord.Permissions
	// IsRestricted checks whether or not this command is restricted.
	//
	// If the RestrictionFunc returns an error that implements
	// RestrictionErrorWrapper, it will be wrapped accordingly.
	IsRestricted(*state.State, *Context) error
	// Throttler returns the Throttler of this command.
	Throttler() Throttler

	// Invoke invokes the command.
	// See Command.Invoke for more details.
	Invoke(*state.State, *Context) (interface{}, error)
}
