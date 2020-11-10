package arg

import (
	"regexp"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/plugin"
)

// TextChannelAllowIDs is a global flag that defines whether TextChannels may
// also be noted as plain Snowflakes.
var TextChannelAllowIDs = false

// =============================================================================
// TextChannel
// =====================================================================================

// TextChannel is the Type used for guild text channels and news channels.
// The channel must be on the same guild as the invoking one.
//
// TextChannel will always fail if used in a direct message.
//
// Go type: *discord.Channel
var TextChannel Type = new(textChannel)

type textChannel struct{}

func (t textChannel) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(textChannelName) // we have a fallback
	return name
}

func (t textChannel) Description(l *i18n.Localizer) string {
	if TextChannelAllowIDs {
		desc, err := l.Localize(textChannelDescriptionWithID)
		if err == nil {
			return desc
		}
	}

	desc, _ := l.Localize(textChannelDescriptionNoID) // we have a fallback
	return desc
}

var textChannelMentionRegexp = regexp.MustCompile(`^<#(?P<id>\d+)>$`)

func (t textChannel) Parse(s *state.State, ctx *Context) (interface{}, error) {
	err := restriction.ChannelTypes(plugin.GuildChannels)(s, ctx.Context)
	if err != nil {
		return nil, err
	}

	if matches := textChannelMentionRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		rawID := matches[1]

		id, err := discord.ParseSnowflake(rawID)
		if err != nil { // range err
			return nil, newArgParsingErr2(textChannelInvalidMentionErrorArg, textChannelInvalidMentionErrorFlag, ctx, nil)
		}

		c, err := s.Channel(discord.ChannelID(id))
		if err != nil {
			return nil, newArgParsingErr2(textChannelInvalidMentionErrorArg, textChannelInvalidMentionErrorFlag, ctx, nil)
		}

		if c.GuildID != ctx.GuildID {
			return nil, newArgParsingErr(textChannelGuildNotMatchingError, ctx, nil)
		} else if c.Type != discord.GuildText && c.Type != discord.GuildNews {
			return nil, newArgParsingErr(textChannelInvalidTypeError, ctx, nil)
		}

		return c, nil
	}

	if !TextChannelAllowIDs {
		return nil, newArgParsingErr(textChannelInvalidMentionWithRawError, ctx, nil)
	}

	id, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgParsingErr(textChannelInvalidError, ctx, nil)
	}

	c, err := s.Channel(discord.ChannelID(id))
	if err != nil {
		return nil, newArgParsingErr(channelIDInvalidError, ctx, nil)
	}

	if c.GuildID != ctx.GuildID {
		return nil, newArgParsingErr(textChannelIDGuildNotMatchingError, ctx, nil)
	} else if c.Type != discord.GuildText && c.Type != discord.GuildNews {
		return nil, newArgParsingErr(textChannelIDInvalidTypeError, ctx, nil)
	}

	return c, nil
}

func (t textChannel) Default() interface{} {
	return (*discord.Channel)(nil)
}
