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
		// Meta returns the meta data of the command.
		Meta() CommandMeta
		// Invoke calls the command.
		//
		// Possible first return values are
		//	• uint, uint8, uint16, uint32, uint64
		//	• int, int8, int16, int32, int64
		// 	• float32, float64
		//	• string
		//	• discord.Embed
		//	• discordutil.EmbedBuilder
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
		// IsRestricted checks if the user is restricted from using the
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
		// GetThrottlingOptions returns the ThrottlingOptions for the command.
		// If either of the fields in ThrottlingOptions is zero value, the
		// command won't be throttled.
		GetThrottlingOptions() ThrottlingOptions
	}
)

// Response is the interface that a type can implement, to be used as a return
// value of a Command.Invoke call.
type Response interface {
	// Send sends the response using the passed state.State and Context.
	Send(s *state.State, ctx *Context) error
}
