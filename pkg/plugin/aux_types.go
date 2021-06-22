package plugin

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"
)

// ChannelTypes is an enum used to specify in which channel types the command
// may be executed.
// It is bit-shifted to allow for combinations of different channel types.
type ChannelTypes uint8

const (
	// GuildTextChannels is the ChannelTypes of a regular guild text channel
	// (0).
	GuildTextChannels ChannelTypes = 1 << iota
	// GuildNewsChannels is the ChannelTypes of a news channel (5).
	GuildNewsChannels
	// DirectMessages is the ChannelTypes of a private chat (1).
	DirectMessages

	// Combinations

	// AllChannels is a combination of all ChannelTypes.
	AllChannels = DirectMessages | GuildChannels
	// GuildChannels is a combination of all ChannelTypes used in guilds, i.e.
	// GuildTextChannels and GuildNewsChannels.
	GuildChannels = GuildTextChannels | GuildNewsChannels
)

// Has checks if the passed discord.ChannelType is found in the ChannelTypes.
func (t ChannelTypes) Has(target discord.ChannelType) bool {
	switch target { //nolint:exhaustive // other types handled in default
	case discord.GuildText:
		return t&GuildTextChannels == GuildTextChannels
	case discord.DirectMessage:
		return t&DirectMessages == DirectMessages
	case discord.GuildNews:
		return t&GuildNewsChannels == GuildNewsChannels
	default:
		return false
	}
}

// Check checks if the ChannelTypes match the channel type of the invoking
// channel.
// It tries to avoid a call to Channel.
func (t ChannelTypes) Check(ctx *Context) (bool, error) {
	if t&AllChannels == AllChannels { // we match all channel types
		return true, nil
	} else if t&AllChannels == 0 { // we match no valid channel types
		return false, nil
	}

	if ctx.GuildID == 0 { // we are in a dm and...
		if t&DirectMessages == DirectMessages { // ... allow them
			return true, nil
		}

		// ... don't allow them

		return false, nil
	}

	// we are in a guild channel and...

	// ... allow all types of guild channels
	if t&GuildChannels == GuildChannels {
		return true, nil

		// ... allow one of the guild channels, we just don't know if we match
		// the right one
	} else if t&GuildChannels != 0 {
		// so we have to check
		c, err := ctx.Channel()
		if err != nil {
			return false, err
		}

		return t.Has(c.Type), nil
	}

	// ... don't allow guild channels
	return false, nil
}

func (t ChannelTypes) String() string {
	switch {
	case t&AllChannels == AllChannels:
		return "all channels"
	case t&GuildChannels == GuildChannels:
		return "guild channels"
	case t&GuildTextChannels == GuildTextChannels:
		return "guild text channels"
	case t&GuildNewsChannels == GuildNewsChannels:
		return "guild news channels"
	case t&DirectMessages == DirectMessages:
		return "direct messages"
	default:
		return ""
	}
}

type (
	// RestrictionFunc is the function used to determine if a user is
	// authorized to use a command or module.
	//
	// Implementations can be found in impl/restriction.
	RestrictionFunc func(*state.State, *Context) error

	// RestrictionErrorWrapper is the interface used to wrap errors returned by
	// a RestrictionFunc.
	// If the RestrictionFunc of a plugin returns an error, that implements
	// this, It will call Wrap() to properly wrap the error.
	RestrictionErrorWrapper interface {
		// Wrap wraps the error returned by the RestrictionFunc.
		Wrap(*state.State, *Context) error
		error
	}
)

// Throttler is used to create cooldowns for commands.
//
// Implementations can be found in impl/throttler.
type Throttler interface {
	// Check checks if the command may be executed and increments the counter
	// if so.
	// It returns non-nil, nil if the command may be executed and nil, non-nil
	// if the command is throttled.
	// The returned error should be of type *plugin.ThrottlingError.
	//
	// If the returned function gets called, the command invoke should not be
	// counted, e.g. if a Command returns with an error.
	// This will be the case, if the ThrottlerCancelChecker function in the bot's
	// Options returns true.
	//
	// Note that the Throttler will be called before non-default bot
	// middlewares are run.
	// Therefore, only context data set through event handlers will be
	// available.
	Check(*state.State, *Context) (func(), error)
}
