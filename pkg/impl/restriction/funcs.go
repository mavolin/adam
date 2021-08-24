package restriction

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/permutil"
)

// NSFW asserts that a command is executed in an NSFW channel.
// It fails if the command is used in a direct message.
func NSFW(_ *state.State, ctx *plugin.Context) error {
	if ctx.GuildID == 0 {
		return plugin.NewFatalRestrictionErrorl(nsfwChannelError)
	}

	c, err := ctx.Channel()
	if err != nil {
		return err
	}

	if c.NSFW {
		return nil
	}

	return plugin.NewRestrictionErrorl(nsfwChannelError)
}

var _ plugin.RestrictionFunc = NSFW

// GuildOwner asserts that a command is executed by the guild owner.
// It fails if the command is used in a direct message.
func GuildOwner(_ *state.State, ctx *plugin.Context) error {
	if err := assertChannelTypes(ctx, plugin.GuildChannels); err != nil {
		return err
	}

	g, err := ctx.Guild()
	if err != nil {
		return err
	}

	if g.OwnerID == ctx.Author.ID {
		return nil
	}

	return plugin.NewFatalRestrictionErrorl(guildOwnerError)
}

var _ plugin.RestrictionFunc = GuildOwner

// BotOwner asserts that a command is executed by a bot owner.
func BotOwner(_ *state.State, ctx *plugin.Context) error {
	if ctx.IsBotOwner() {
		return nil
	}

	return plugin.NewFatalRestrictionErrorl(botOwnerError)
}

var _ plugin.RestrictionFunc = BotOwner

// Users creates a plugin.RestrictionFunc that defines a set of users that may
// use a command.
// It returns plugin.DefaultRestrictionError if the author isn't one of them.
func Users(allowed ...discord.UserID) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		if len(allowed) == 0 {
			return nil
		}

		for _, id := range allowed {
			if id == ctx.Author.ID {
				return nil
			}
		}

		return plugin.DefaultFatalRestrictionError
	}
}

// AllRoles asserts that the user has all of the passed roles or is able to
// assign themself all of the passed roles.
// You can mix roles from different guilds, roles that aren't available in a
// guild are ignored.
// However, the guild the command was invoked in must have at least one of the
// passed roles.
// This effectively means, that only guilds whose roles are included, are able
// to use the command at all.
//
// It fails if the command is used in a direct message.
//nolint:gocognit
func AllRoles(allowed ...discord.RoleID) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		if len(allowed) == 0 {
			return nil
		}

		if err := assertChannelTypes(ctx, plugin.GuildChannels); err != nil {
			return err
		}

		missingIDs := make([]discord.RoleID, 0, len(allowed))

		// find all missing roles
	Allowed:
		for _, targetID := range allowed {
			for _, id := range ctx.Member.RoleIDs {
				if targetID == id {
					continue Allowed
				}
			}

			missingIDs = append(missingIDs, targetID)
		}

		if len(missingIDs) == 0 {
			return nil
		}

		g, err := ctx.Guild()
		if err != nil {
			return err
		}

		missingRoles := make([]discord.Role, 0, len(missingIDs))

		// out of the missing roles, find those missing in this guild
	Missing:
		for _, id := range missingIDs {
			for _, r := range g.Roles {
				if r.ID == id { // role is in the guild
					missingRoles = insertRoleSorted(r, missingRoles)
					continue Missing
				}
			}
		}

		if len(missingRoles) == 0 { // no roles missing from this guild
			// check if this guild even has a role in our checklist
			for _, id := range allowed {
				for _, r := range g.Roles {
					if id == r.ID {
						return nil
					}
				}
			}

			return plugin.DefaultFatalRestrictionError
		}

		if permutil.CanMemberManageRole(*g, *ctx.Member, missingRoles[len(missingRoles)-1].ID) {
			return nil
		}

		return newAllMissingRolesError(missingRoles, ctx.Localizer)
	}
}

// MustAllRoles asserts that the user has all of the passed roles.
// You can mix roles from different guilds, roles that aren't available in the
// invoking guild are ignored.
// However, the guild the command was invoked in must have at least one of the
// passed roles.
// This effectively means, that only guilds whose roles are included, are able
// to use the command at all.
//
// It fails if the command is used in a direct message.
//nolint:gocognit
func MustAllRoles(allowed ...discord.RoleID) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		if len(allowed) == 0 {
			return nil
		}

		if err := assertChannelTypes(ctx, plugin.GuildChannels); err != nil {
			return err
		}

		missingIDs := make([]discord.RoleID, 0, len(allowed))

		// find all missing roles
	Allowed:
		for _, targetID := range allowed {
			for _, id := range ctx.Member.RoleIDs {
				if targetID == id {
					continue Allowed
				}
			}

			missingIDs = append(missingIDs, targetID)
		}

		if len(missingIDs) == 0 {
			return nil
		}

		g, err := ctx.Guild()
		if err != nil {
			return err
		}

		missingRoles := make([]discord.Role, 0, len(missingIDs))

		// out of the missing roles, find those missing in this guild
	Missing:
		for _, id := range missingIDs {
			for _, r := range g.Roles {
				if r.ID == id { // role is in the guild
					missingRoles = insertRoleSorted(r, missingRoles)
					continue Missing
				}
			}
		}

		if len(missingRoles) == 0 { // no roles missing from this guild
			// check if this guild even has a role in our checklist
			for _, id := range allowed {
				for _, r := range g.Roles {
					if id == r.ID {
						return nil
					}
				}
			}

			return plugin.DefaultFatalRestrictionError
		}

		return newAllMissingRolesError(missingRoles, ctx.Localizer)
	}
}

// AnyRole asserts that the invoking user has at least one of the passed
// roles or has the ability to assign one of the passed roles to themself.
//
// It fails if the command is used in a direct message.
func AnyRole(allowed ...discord.RoleID) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		if len(allowed) == 0 {
			return nil
		}

		if err := assertChannelTypes(ctx, plugin.GuildChannels); err != nil {
			return err
		}

		for _, targetID := range allowed {
			for _, id := range ctx.Member.RoleIDs {
				if targetID == id {
					return nil
				}
			}
		}

		g, err := ctx.Guild()
		if err != nil {
			return err
		}

		missingRoles := make([]discord.Role, 0, len(allowed))

	Allowed:
		for _, id := range allowed {
			for _, r := range g.Roles {
				if r.ID == id {
					missingRoles = insertRoleSorted(r, missingRoles)
					continue Allowed
				}
			}
		}

		if len(missingRoles) == 0 { // none of the roles are from this guild
			return plugin.DefaultFatalRestrictionError
		}

		if permutil.CanMemberManageRole(*g, *ctx.Member, missingRoles[0].ID) {
			return nil
		}

		return newAnyMissingRolesError(missingRoles, ctx.Localizer)
	}
}

// MustAnyRole asserts that the invoking user has at least one of the passed
// roles.
//
// It fails if the command is used in a direct message.
func MustAnyRole(allowed ...discord.RoleID) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if len(allowed) == 0 {
			return nil
		}

		if err := assertChannelTypes(ctx, plugin.GuildChannels); err != nil {
			return err
		}

		for _, targetID := range allowed {
			for _, id := range ctx.Member.RoleIDs {
				if targetID == id {
					return nil
				}
			}
		}

		g, err := ctx.Guild()
		if err != nil {
			return err
		}

		missingRoles := make([]discord.Role, 0, len(allowed))

	Allowed:
		for _, id := range allowed {
			for _, r := range g.Roles {
				if r.ID == id {
					missingRoles = insertRoleSorted(r, missingRoles)
					continue Allowed
				}
			}
		}

		if len(missingRoles) == 0 { // none of the roles are from this guild
			return plugin.DefaultFatalRestrictionError
		}

		return newAnyMissingRolesError(missingRoles, ctx.Localizer)
	}
}

// Channels asserts that a command is executed in one of the passed channels.
func Channels(allowed ...discord.ChannelID) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if len(allowed) == 0 {
			return nil
		}

		for _, id := range allowed {
			if id == ctx.ChannelID {
				return nil
			}
		}

		if ctx.GuildID == 0 {
			return plugin.DefaultFatalRestrictionError
		}

		channels, err := s.Channels(ctx.GuildID)
		if err != nil {
			return errors.WithStack(err)
		}

		g, err := ctx.Guild()
		if err != nil {
			return err
		}

		missingIDs := make([]discord.ChannelID, 0, len(allowed))

	ChannelIDs:
		for _, targetID := range allowed {
			for _, c := range channels {
				if c.ID == targetID {
					overwrites := discord.CalcOverwrites(*g, c, *ctx.Member)

					// make sure we only list channels the user can see and use
					if overwrites.Has(discord.PermissionViewChannel | discord.PermissionSendMessages) {
						missingIDs = append(missingIDs, c.ID)
					}

					continue ChannelIDs
				}
			}
		}

		// the guild does not have any of the allowed channels, or the user
		// can't see them
		if len(missingIDs) == 0 {
			return plugin.DefaultFatalRestrictionError
		}

		return newChannelsError(missingIDs, ctx.Localizer)
	}
}

// ChannelTypes asserts that a command is executed in a channel of an allowed
// type.
//
// Note that the resulting plugin.RestrictionFunc won't return a
// errors.ChannelTypeError but a *plugin.RestrictionError.
func ChannelTypes(allowed plugin.ChannelTypes) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		return assertChannelTypes(ctx, allowed)
	}
}

// UserPermissions asserts that the invoking user has the passed permissions.
//
// Note that direct messages may also pass this, if the passed permissions
// only require permutil.DMPermissions.
func UserPermissions(required discord.Permissions) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		if required == 0 {
			return nil
		}

		if ctx.GuildID == 0 && permutil.DMPermissions.Has(required) {
			return nil
		}

		if err := assertChannelTypes(ctx, plugin.GuildChannels); err != nil {
			return err
		}

		actual, err := ctx.UserPermissions()
		if err != nil {
			return err
		}

		missing := (actual & required) ^ required
		if missing == 0 {
			return nil
		}

		return newUserPermissionsError(missing, ctx.Localizer)
	}
}
