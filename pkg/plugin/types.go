package plugin

import (
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"
)

// ChannelType is an enum used to specify in which channel types the command
// may be executed.
// It is bit-shifted to allow for combinations of different channel types.
type ChannelType uint8

const (
	// GuildText is the ChannelType of a regular guild text channel (0).
	GuildText ChannelType = 1 << iota
	// GuildNews is the ChannelType of a news channel (5).
	GuildNews
	// DirectMessage is the ChannelType of a private chat (1).
	DirectMessage

	// Combinations

	// All is a combination of all ChannelTypes.
	All = DirectMessage | Guild
	// Guild is a combination of all ChannelTypes used in guilds, i.e.
	// GuildText and GuildNews.
	Guild = GuildText | GuildNews
)

// Has checks if the passed discord.ChannelType is found in the ChannelType.
func (t ChannelType) Has(target discord.ChannelType) bool {
	switch target {
	case discord.GuildText:
		return t&GuildText == GuildText
	case discord.DirectMessage:
		return t&DirectMessage == DirectMessage
	case discord.GuildNews:
		return t&GuildNews == GuildNews
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
