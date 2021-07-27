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
var Role plugin.ArgType = new(role)

type role struct{}

func (r role) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(roleName) // we have a fallback
	return name
}

func (r role) GetDescription(l *i18n.Localizer) string {
	//goland:noinspection GoBoolExpressions
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

func (r role) Parse(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	err := restriction.ChannelTypes(plugin.GuildChannels)(s, ctx.Context)
	if err != nil {
		return nil, err
	}

	if matches := roleMentionRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		rawID := matches[1]

		id, err := discord.ParseSnowflake(rawID)
		if err != nil { // range err
			return nil, newArgumentError2(roleInvalidMentionErrorArg, roleInvalidMentionErrorFlag, ctx, nil)
		}

		role, err := s.Role(ctx.GuildID, discord.RoleID(id))
		if err != nil {
			return nil, newArgumentError2(roleInvalidMentionErrorArg, roleInvalidMentionErrorFlag, ctx, nil)
		}

		return role, nil
	}

	//goland:noinspection GoBoolExpressions
	if !RoleAllowIDs {
		return nil, newArgumentError(roleInvalidMentionWithRawError, ctx, nil)
	}

	id, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgumentError(roleInvalidError, ctx, nil)
	}

	role, err := s.Role(ctx.GuildID, discord.RoleID(id))
	if err != nil {
		return nil, newArgumentError(roleIDInvalidError, ctx, nil)
	}

	return role, nil
}

func (r role) GetDefault() interface{} {
	return (*discord.Role)(nil)
}
