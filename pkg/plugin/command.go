package plugin

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

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
		//	• *embedutil.Builder
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
		// If Invoke returns an error it will be handed to the error handler
		// of the bot.
		//
		// As a special case if both return values are non-nil, both the
		// reply and the error will be handled.
		// Any errors that occur when sending the reply will be silently
		// handled.
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
		// GetShortDescription returns an optional short description
		// of the command.
		GetShortDescription(l *i18n.Localizer) string
		// GetLongDescription returns an optional long description of the
		// command.
		GetLongDescription(l *i18n.Localizer) string
		// GetExampleArgs returns optional example arguments of the command.
		GetExampleArgs(l *i18n.Localizer) []string

		// GetArgs returns the ArgConfig of the command.
		//
		// If this is nil, the command will accept no arguments and flags.
		GetArgs() ArgConfig

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
		//
		// If the RestrictionFunc returns an error that implements
		// RestrictionErrorWrapper, it will be wrapped accordingly.
		IsRestricted(s *state.State, ctx *Context) error
		// GetThrottler returns the Throttler for the command.
		GetThrottler() Throttler
	}
)

// Reply is used to send a reply, if returned as first return value of a
// Command.Invoke call.
type Reply interface {
	// SendReply sends the reply using the passed state.State and Context.
	SendReply(s *state.State, ctx *Context) error
}
