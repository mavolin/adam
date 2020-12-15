package arg

import (
	"regexp"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/plugin"
)

// RoleAllowIDs is a global flag that defines whether Roles may also be noted
// as plain Snowflakes.
var RoleAllowIDs = true

// =============================================================================
// Role
// =====================================================================================

// Role is the Type used for roles.
// A role can either be a role mention or the id of the role.
//
// It will return an error if used on a guild.
//
// Go type: *discord.Role
var Role Type = new(role)

type role struct{}

func (r role) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(roleName) // we have a fallback
	return name
}

func (r role) Description(l *i18n.Localizer) string {
	if RoleAllowIDs {
		desc, err := l.Localize(roleDescriptionWithID)
		if err == nil {
			return desc
		}
	}

	desc, _ := l.Localize(roleDescriptionNoID) // we have a fallback
	return desc
}

var roleMentionRegexp = regexp.MustCompile(`^<@&(?P<id>\d+)>$`)

func (r role) Parse(s *state.State, ctx *Context) (interface{}, error) {
	err := restriction.ChannelTypes(plugin.GuildChannels)(s, ctx.Context)
	if err != nil {
		return nil, err
	}

	if matches := roleMentionRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		rawID := matches[1]

		id, err := discord.ParseSnowflake(rawID)
		if err != nil { // range err
			return nil, newArgParsingErr2(roleInvalidMentionErrorArg, roleInvalidMentionErrorFlag, ctx, nil)
		}

		role, err := s.Role(ctx.GuildID, discord.RoleID(id))
		if err != nil {
			return nil, newArgParsingErr2(roleInvalidMentionErrorArg, roleInvalidMentionErrorFlag, ctx, nil)
		}

		return role, nil
	}

	if !RoleAllowIDs {
		return nil, newArgParsingErr(roleInvalidMentionWithRawError, ctx, nil)
	}

	id, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgParsingErr(roleInvalidError, ctx, nil)
	}

	role, err := s.Role(ctx.GuildID, discord.RoleID(id))
	if err != nil {
		return nil, newArgParsingErr(roleIDInvalidError, ctx, nil)
	}

	return role, nil
}

func (r role) Default() interface{} {
	return (*discord.Role)(nil)
}
