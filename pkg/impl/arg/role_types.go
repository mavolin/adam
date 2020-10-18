package arg

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/plugin"
)

// RoleAllowIDs is a global flag that allows you to specify whether Roles
// may be noted as Snowflakes.
//
// Defaults to true.
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
var Role = new(role)

type role struct{}

func (r role) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(roleName) // we have a fallback
	return name
}

func (r role) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(roleDescription) // we have a fallback
	return desc
}

func (r role) Parse(s *state.State, ctx *Context) (interface{}, error) {
	err := restriction.ChannelTypes(plugin.GuildChannels)(s, ctx.Context)
	if err != nil {
		return nil, err
	}

	if matches := userMentionRegexp.FindStringSubmatch(ctx.Raw); len(matches) > 1 {
		id := matches[1]

		rid, err := discord.ParseSnowflake(id)
		if err != nil { // range err
			return nil, newArgParsingErr(roleInvalidMentionArg, roleInvalidMentionFlag, ctx, nil)
		}

		role, err := s.Role(ctx.GuildID, discord.RoleID(rid))
		if err != nil {
			return nil, newArgParsingErr(roleInvalidMentionArg, roleInvalidMentionFlag, ctx, nil)
		}

		return role, nil
	}

	if !RoleAllowIDs {
		return nil, newArgParsingErr(roleInvalidMentionWithRaw, roleInvalidMentionWithRaw, ctx, nil)
	}

	rid, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgParsingErr(roleInvalidIDWithRaw, roleInvalidIDWithRaw, ctx, nil)
	}

	role, err := s.Role(ctx.GuildID, discord.RoleID(rid))
	if err != nil {
		return nil, newArgParsingErr(roleInvalidIDArg, roleInvalidIDFlag, ctx, nil)
	}

	return role, nil
}

func (r role) Default() interface{} {
	return (*discord.Role)(nil)
}
