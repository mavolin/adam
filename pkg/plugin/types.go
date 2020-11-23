package plugin

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"
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
	switch target { //nolint: exhaustive // other types handled in default
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

type (
	// RestrictionFunc is the function used to determine if a user is authorized
	// to use a command or module.
	//
	// Implementations can be found in impl/restriction.
	RestrictionFunc func(s *state.State, ctx *Context) error

	// RestrictionErrorWrapper is the interface used to wrap errors returned by a
	// RestrictionFunc.
	// If the RestrictionFunc of a plugin returns an error, that implements this,
	// It will call Wrap() to properly wrap the error.
	RestrictionErrorWrapper interface {
		// Wrap wraps the error returned by the RestrictionFunc.
		Wrap(*state.State, *Context) error
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
	// The returned error should be of type errors.ThrottlingError.
	//
	// The returned function will be called, if the command exits with an
	// error and that error is not ignored as defined via the bots Options.
	Check(ctx *Context) (func(), error)
}
