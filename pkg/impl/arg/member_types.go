package arg

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/plugin"
)

// MemberAllowIDs is a global flag that allows you to specify whether Members
// may be noted as Snowflakes.
//
// Defaults to true.
var MemberAllowIDs = true

// =============================================================================
// Member
// =====================================================================================

// Member is the Type used for members of a guild.
// It will always return an error, if the command is called in a direct
// message.
//
// A Member can either be a mention of a member, or, if enabled, an ID
// of a guild member.
//
// Go type: *discord.Member
var Member = member{}

type member struct{}

func (m member) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(memberName) // we have a fallback
	return name
}

func (m member) Description(l *i18n.Localizer) string {
	if MemberAllowIDs {
		desc, err := l.Localize(memberDescriptionWithIDs)
		if err == nil {
			return desc
		}
	}

	desc, _ := l.Localize(memberDescriptionNoIDs) // we have a fallback
	return desc
}

func (m member) Parse(s *state.State, ctx *Context) (interface{}, error) {
	err := restriction.ChannelTypes(plugin.GuildChannels)(s, ctx.Context)
	if err != nil {
		return nil, err
	}

	if userMentionRegexp.MatchString(ctx.Raw) {
		id := userMentionRegexp.FindStringSubmatch(ctx.Raw)[1]

		mid, err := discord.ParseSnowflake(id)
		if err != nil { // range err
			return nil, newArgParsingErr(userInvalidMentionArg, userInvalidMentionFlag, ctx, nil)
		}

		member, err := s.Member(ctx.GuildID, discord.UserID(mid))
		if err != nil {
			return nil, newArgParsingErr(userInvalidMentionArg, userInvalidMentionFlag, ctx, nil)
		}

		return member, nil
	}

	if !MemberAllowIDs {
		return nil, newArgParsingErr(userInvalidMentionWithRaw, userInvalidMentionWithRaw, ctx, nil)
	}

	mid, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgParsingErr(userInvalidIDWithRaw, userInvalidIDWithRaw, ctx, nil)
	}

	member, err := s.Member(ctx.GuildID, discord.UserID(mid))
	if err != nil {
		return nil, newArgParsingErr(userInvalidIDArg, userInvalidIDFlag, ctx, nil)
	}

	return member, nil
}

func (m member) Default() interface{} {
	return (*discord.Member)(nil)
}

// =============================================================================
// MemberID
// =====================================================================================

// MemberID is the same as a Member, but it only accepts IDs.
//
// Go type: *discord.Member
var MemberID = memberID{}

type memberID struct{}

func (m memberID) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(memberIDName) // we have a fallback
	return name
}

func (m memberID) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(memberIDDescription) // we have a fallback
	return desc
}

func (m memberID) Parse(s *state.State, ctx *Context) (interface{}, error) {
	err := restriction.ChannelTypes(plugin.GuildChannels)(s, ctx.Context)
	if err != nil {
		return nil, err
	}

	mid, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgParsingErr(userInvalidIDWithRaw, userInvalidIDWithRaw, ctx, nil)
	}

	member, err := s.Member(ctx.GuildID, discord.UserID(mid))
	if err != nil {
		return nil, newArgParsingErr(userInvalidIDArg, userInvalidIDFlag, ctx, nil)
	}

	return member, nil
}

func (m memberID) Default() interface{} {
	return (*discord.Member)(nil)
}
