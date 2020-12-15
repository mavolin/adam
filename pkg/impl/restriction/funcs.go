package restriction

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/permutil"
)

var (
	// ErrNotNSFWChannel is the error returned by NSFW if the command is not
	// executed in an NSFW channel.
	ErrNotNSFWChannel = errors.NewRestrictionErrorl(notNSFWChannelError)
	// ErrNotGuildOwner is the error returned by GuildOwner if the command is
	// not executed by the guild owner.
	ErrNotGuildOwner = errors.NewFatalRestrictionErrorl(notOwnerError)
	// ErrNotBotOwner is the error returned by BotOwner if the command is not
	// executed by the bot owner.
	ErrNotBotOwner = errors.NewFatalRestrictionErrorl(notBotOwnerError)
)

// NSFW asserts that a command is executed in an NSFW channel.
// It fails if the command is used in a direct message.
func NSFW(_ *state.State, ctx *plugin.Context) error {
	err := assertChannelTypes(ctx, plugin.GuildChannels,
		errors.NewWithStack("restriction: invalid assertion NSFW for DM-only command"))
	if err != nil {
		return err
	}

	c, err := ctx.Channel()
	if err != nil {
		return err
	}

	if c.NSFW {
		return nil
	}

	return ErrNotNSFWChannel
}

var _ plugin.RestrictionFunc = NSFW

// GuildOwner asserts that a command is executed by the guild owner.
// It fails if the command is used in a direct message.
func GuildOwner(_ *state.State, ctx *plugin.Context) error {
	err := assertChannelTypes(ctx, plugin.GuildChannels,
		errors.NewWithStack("restriction: invalid assertion GuildOwner for DM-only command"))
	if err != nil {
		return err
	}

	g, err := ctx.Guild()
	if err != nil {
		return err
	}

	if g.OwnerID == ctx.Author.ID {
		return nil
	}

	return ErrNotGuildOwner
}

var _ plugin.RestrictionFunc = GuildOwner

// BotOwner asserts that a command is executed by a bot owner.
func BotOwner(_ *state.State, ctx *plugin.Context) error {
	if ctx.IsBotOwner() {
		return nil
	}

	return ErrNotBotOwner
}

var _ plugin.RestrictionFunc = BotOwner

// Users creates a plugin.RestrictionFunc that defines a set of users that may
// use a command.
// It returns a errors.DefaultRestrictionError if the author isn't one of them.
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

		return errors.DefaultFatalRestrictionError
	}
}

// MustAllRoles asserts that the user has all of the passed roles or is able
// to assign themself any of the passed roles.
// You can mix roles from different guilds, roles that aren't available in a
// guild are ignored.
// However, the guild the command was invoked in must have at least one of the
// passed roles.
// This effectively means, that only guilds whose roles are included, are able
// to use the command at all.
//
// It fails if the command is used in a direct message.
func AllRoles(allowed ...discord.RoleID) plugin.RestrictionFunc { //nolint:gocognit
	return func(_ *state.State, ctx *plugin.Context) error {
		if len(allowed) == 0 {
			return nil
		}

		if ctx.GuildID == 0 {
			err := assertChannelTypes(ctx, plugin.GuildChannels,
				errors.NewWithStack("restriction: invalid assertion AllRoles for DM-only command"))
			if err != nil {
				return err
			}
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

			return errors.DefaultFatalRestrictionError
		}

		if canManageRole(missingRoles[len(missingRoles)-1], g, ctx.Member) {
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
func MustAllRoles(allowed ...discord.RoleID) plugin.RestrictionFunc { //nolint:gocognit
	return func(_ *state.State, ctx *plugin.Context) error {
		if len(allowed) == 0 {
			return nil
		}

		if ctx.GuildID == 0 {
			err := assertChannelTypes(ctx, plugin.GuildChannels,
				errors.NewWithStack("restriction: invalid assertion MustAllRoles for DM-only command"))
			if err != nil {
				return err
			}
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

			return errors.DefaultFatalRestrictionError
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

		if ctx.GuildID == 0 {
			err := assertChannelTypes(ctx, plugin.GuildChannels,
				errors.NewWithStack("restriction: invalid assertion MustAllRoles for DM-only command"))
			if err != nil {
				return err
			}
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
			return errors.DefaultFatalRestrictionError
		}

		if canManageRole(missingRoles[0], g, ctx.Member) {
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

		if ctx.GuildID == 0 {
			err := assertChannelTypes(ctx, plugin.GuildChannels,
				errors.NewWithStack("restriction: invalid assertion MustAllRoles for DM-only command"))
			if err != nil {
				return err
			}
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
			return errors.DefaultFatalRestrictionError
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
			return errors.DefaultFatalRestrictionError
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
			return errors.DefaultFatalRestrictionError
		}

		return newChannelsError(missingIDs, ctx.Localizer)
	}
}

// ChannelType asserts that a command is executed in a channel of an allowed
// type.
//
// Note that the resulting plugin.RestrictionFunc won't return a
// errors.InvalidChannelTypeError but a errors.RestrictionError.
func ChannelTypes(allowed plugin.ChannelTypes) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		return assertChannelTypes(ctx, allowed,
			errors.NewWithStack("restriction: invalid assertion ChannelTypes for command with opposite channel types"))
	}
}

// BotPermissions asserts that the bot has the passed permissions.
// When using this, the commands bot permissions should be set to
// plugin.NoPermissions.
//
// Note that direct messages may also pass this, if the passed permissions
// only require constant.DMPermissions.
//
// Also note that the resulting plugin.RestrictionFunc won't return a
// errors.InsufficientPermissionsError but a errors.RestrictionError.
func BotPermissions(required discord.Permissions) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		if required == 0 {
			return nil
		}

		if ctx.GuildID == 0 {
			if permutil.DMPermissions.Has(required) {
				return nil
			}

			return assertChannelTypes(ctx, plugin.GuildChannels,
				errors.NewWithStack("restriction: invalid assertion BotPermissions with guild only permissions for "+
					"DM command"))
		}

		actual, err := ctx.SelfPermissions()
		if err != nil {
			return err
		}

		missing := (actual & required) ^ required
		if missing == 0 {
			return nil
		}

		return newInsufficientBotPermissionsError(missing, ctx.Localizer)
	}
}

// UserPermissions asserts that the invoking user has the passed permissions.
//
// Note that direct messages may also pass this, if the passed permissions
// only require constant.DMPermissions.
func UserPermissions(perms discord.Permissions) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		if perms == 0 {
			return nil
		}

		if ctx.GuildID == 0 {
			if permutil.DMPermissions.Has(perms) {
				return nil
			}

			return assertChannelTypes(ctx, plugin.GuildChannels,
				errors.NewWithStack("restriction: invalid assertion UserPermissions with guild only permissions for "+
					"DM-only command"))
		}

		actual, err := ctx.UserPermissions()
		if err != nil {
			return err
		}

		missing := (actual & perms) ^ perms
		if missing == 0 {
			return nil
		}

		return newInsufficientUserPermissionsError(missing, ctx.Localizer)
	}
}
