package arg

import (
	"regexp"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/plugin"
)

// MemberAllowIDs is a global flag that defines whether Members may also be
// noted as plain Snowflakes.
var MemberAllowIDs = true

// =============================================================================
// User
// =====================================================================================

// User is the Type used to specify users globally.
// The User doesn't have to be on the same guild as the invoking one.
// In contrast to member, this can also be used in direct messages.
// A User can either be a mention, or an id.
//
// Gp type: *discord.User
var User plugin.ArgType = new(user)

type user struct{}

func (u user) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(userName) // we have a fallback
	return name
}

func (u user) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(userDescription) // we have a fallback
	return desc
}

var userMentionRegexp = regexp.MustCompile(`^<@!?(?P<id>\d+)>$`)

func (u user) Parse(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	if matches := userMentionRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		rawID := matches[1]

		id, err := discord.ParseSnowflake(rawID)
		if err != nil { // range err
			return nil, newArgumentError2(userInvalidMentionErrorArg, userInvalidMentionErrorFlag, ctx, nil)
		}

		for _, m := range ctx.Mentions {
			if m.ID == discord.UserID(id) {
				return &m.User, nil
			}
		}

		user, err := s.User(discord.UserID(id))
		if err != nil {
			return nil, newArgumentError2(userInvalidMentionErrorArg, userInvalidMentionErrorFlag, ctx, nil)
		}

		return user, nil
	}

	id, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgumentError(userInvalidError, ctx, nil)
	}

	user, err := s.User(discord.UserID(id))
	if err != nil {
		return nil, newArgumentError(userIDInvalidError, ctx, nil)
	}

	return user, nil
}

func (u user) GetDefault() interface{} {
	return (*discord.User)(nil)
}

// =============================================================================
// Member
// =====================================================================================

// Member is the Type used for members of a guild.
// It will always return an error, if the Command is called in a direct
// message.
//
// A Member can either be a mention of a member, or, if enabled, an id of a
// guild member.
//
// Go type: *discord.Member
var Member plugin.ArgType = new(member)

type member struct{}

func (m member) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(memberName) // we have a fallback
	return name
}

func (m member) GetDescription(l *i18n.Localizer) string {
	if MemberAllowIDs {
		desc, err := l.Localize(memberDescriptionWithIDs)
		if err == nil {
			return desc
		}
	}

	desc, _ := l.Localize(memberDescriptionNoIDs) // we have a fallback
	return desc
}

func (m member) Parse(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	err := restriction.ChannelTypes(plugin.GuildChannels)(s, ctx.Context)
	if err != nil {
		return nil, err
	}

	if matches := userMentionRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		rawID := matches[1]

		id, err := discord.ParseSnowflake(rawID)
		if err != nil { // range err
			return nil, newArgumentError2(userInvalidMentionErrorArg, userInvalidMentionErrorFlag, ctx, nil)
		}

		for _, m := range ctx.Mentions {
			if m.ID == discord.UserID(id) {
				m.Member.User = m.User
				return m.Member, nil
			}
		}

		member, err := s.Member(ctx.GuildID, discord.UserID(id))
		if err != nil {
			return nil, newArgumentError2(userInvalidMentionErrorArg, userInvalidMentionErrorFlag, ctx, nil)
		}

		return member, nil
	}

	if !MemberAllowIDs {
		return nil, newArgumentError(userInvalidMentionWithRawError, ctx, nil)
	}

	id, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgumentError(userInvalidError, ctx, nil)
	}

	member, err := s.Member(ctx.GuildID, discord.UserID(id))
	if err != nil {
		return nil, newArgumentError(userIDInvalidError, ctx, nil)
	}

	return member, nil
}

func (m member) GetDefault() interface{} {
	return (*discord.Member)(nil)
}
