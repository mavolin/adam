package plugin

import (
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
)

// ChannelTypes is an enum used to specify in which channel types the command
// may be executed.
// It is bit-shifted to allow for combinations of different channel types.
type ChannelTypes uint8

const (
	// GuildText is the ChannelTypes of a regular guild text channel (0).
	GuildText ChannelTypes = 1 << iota
	// GuildNews is the ChannelTypes of a news channel (5).
	GuildNews
	// DirectMessage is the ChannelTypes of a private chat (1).
	DirectMessage

	// Combinations

	// All is a combination of all ChannelTypes.
	All = DirectMessage | Guild
	// Guild is a combination of all ChannelTypes used in guilds, i.e.
	// GuildText and GuildNews.
	Guild = GuildText | GuildNews
)

// Has checks if the passed discord.ChannelType is found in the ChannelTypes.
func (t ChannelTypes) Has(target discord.ChannelType) bool {
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

// Names returns the names of the ChannelTypes or nil if ChannelTypes is 0
// or invalid.
func (t ChannelTypes) Names(l *localization.Localizer) (s []string) {
	if t&GuildText == GuildText {
		// we can ignore the error, as there is a fallback
		t, _ := l.Localize(guildTextType)
		s = append(s, t)
	}

	if t&GuildNews == GuildNews {
		// we can ignore the error, as there is a fallback
		t, _ := l.Localize(guildNewsType)
		s = append(s, t)
	}

	if t&DirectMessage == DirectMessage {
		// we can ignore the error, as there is a fallback
		t, _ := l.Localize(directMessageType)
		s = append(s, t)
	}

	return
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
