package plugin

import (
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"
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
	switch target {
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

// RestrictionFunc is the function used to determine if a user is authorized
// to use a command or module.
//
// Implementations can be found in impl/restriction.
type RestrictionFunc func(s *state.State, ctx *Context) error

// ThrottlingOptions is used to create cooldowns for commands.
// Throttling is applied on a per-user basis.
type ThrottlingOptions struct {
	// MaxInvokes specifies the inclusive maximum amount of invokes within
	// the given Timeframe
	MaxInvokes uint
	// Duration is the time.Duration where the MaxInvokes level is measured.
	Duration time.Duration
}

// RestrictionErrorWrapper is the interface used to wrap errors returned by a
// RestrictionFunc.
// If the RestrictionFunc of a plugin returns an error, that implements this,
// It will call Wrap() to properly wrap the error.
type RestrictionErrorWrapper interface {
	// Wrap wraps the error returned by the RestrictionFunc.
	Wrap(*state.State, *Context) error
}
