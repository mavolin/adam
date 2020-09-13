package restriction

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/internal/constant"
	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

var (
	// ErrNotNSFWChannel is the error returned by NSFW if the command is not
	// executed in an NSFW channel.
	ErrNotNSFWChannel = errors.NewRestrictionErrorl(notNSFWChannelError)
	// ErrNotGuildOwner is the error returned by GuildOwner, if the command is
	// not executed by the guild owner.
	ErrNotGuildOwner = errors.NewFatalRestrictionErrorl(notOwnerError)
	// ErrNotBotOwner is the error returned by BotOwner, if the command is not
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

// BotOwner asserts that a command is executed by a bot owner.
func BotOwner(_ *state.State, ctx *plugin.Context) error {
	if ctx.IsBotOwner() {
		return nil
	}

	return ErrNotBotOwner
}

// Users creates a plugin.RestrictionFunc that defines a set of users that may
// use a command.
// It returns a errors.DefaultRestrictionError if the author isn't one of them.
func Users(userIDs ...discord.UserID) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		if len(userIDs) == 0 {
			return nil
		}

		for _, id := range userIDs {
			if id == ctx.Author.ID {
				return nil
			}
		}

		return errors.DefaultFatalRestrictionError
	}
}

// AllRoles asserts that the user has all of the passed roles.
// You can mix roles from different guilds, roles that aren't available in a
// guild are ignored.
// However, the guild the command was invoked in must have at least one of the
// passed roles.
// This effectively means, that only guilds whose roles are included, are able
// to use the command at all.
//
// It fails if the command is used in a direct message.
func AllRoles(roleIDs ...discord.RoleID) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		if len(roleIDs) == 0 {
			return nil
		}

		if ctx.GuildID == 0 {
			err := assertChannelTypes(ctx, plugin.GuildChannels,
				errors.NewWithStack("restriction: invalid assertion AllRoles for DM-only command"))
			if err != nil {
				return err
			}
		}

		missingIDs := make([]discord.RoleID, 0, len(roleIDs))

	RoleIDs:
		for _, targetID := range roleIDs {
			for _, id := range ctx.Member.RoleIDs {
				if targetID == id {
					continue RoleIDs
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

	MissingIDs:
		for i := 0; i < len(missingIDs); i++ {
			id := missingIDs[i]

			for _, role := range g.Roles {
				if role.ID == id {
					continue MissingIDs
				}
			}

			// role is not in the guild, remove it, but preserve hierarchy
			missingIDs = append(missingIDs[:i], missingIDs[i+1:]...)
			i--
		}

		if len(missingIDs) == 0 { // no roles missing from this guild
			// check if this guild even has a role in our checklist
			for _, id := range roleIDs {
				for _, r := range g.Roles {
					if id == r.ID {
						return nil
					}
				}
			}

			return errors.DefaultFatalRestrictionError
		}

		return newAllMissingRolesError(missingIDs, ctx.Localizer)
	}
}

// AnyRole asserts that the invoking user has at least one of the passed
// roles.
//
// It fails if the command is used in a direct message.
func AnyRole(roleIDs ...discord.RoleID) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if len(roleIDs) == 0 {
			return nil
		}

		if ctx.GuildID == 0 {
			err := assertChannelTypes(ctx, plugin.GuildChannels,
				errors.NewWithStack("restriction: invalid assertion AllRoles for DM-only command"))
			if err != nil {
				return err
			}
		}

		for _, targetID := range roleIDs {
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

		missingIDs := make([]discord.RoleID, 0, len(roleIDs))

	RoleIDs:
		for _, id := range roleIDs {
			for _, role := range g.Roles {
				if role.ID == id {
					missingIDs = append(missingIDs, role.ID)
					continue RoleIDs
				}
			}
		}

		if len(missingIDs) == 0 { // none of the roles are from this guild
			return errors.DefaultFatalRestrictionError
		}

		return newAnyMissingRolesError(missingIDs, ctx.Localizer)
	}
}

// Channels asserts that a command is executed in one of the passed channels.
func Channels(channelIDs ...discord.ChannelID) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if len(channelIDs) == 0 {
			return nil
		}

		for _, id := range channelIDs {
			if id == ctx.ChannelID {
				return nil
			}
		}

		if ctx.GuildID == 0 {
			return errors.DefaultFatalRestrictionError
		}

		channels, err := s.Channels(ctx.GuildID)
		if err != nil {
			return err
		}

		g, err := ctx.Guild()
		if err != nil {
			return err
		}

		missingIDs := make([]discord.ChannelID, 0, len(channelIDs))

	ChannelIDs:
		for _, targetID := range channelIDs {
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
// errors.InsufficientBotPermissionsError but a errors.RestrictionError.
func BotPermissions(required discord.Permissions) plugin.RestrictionFunc {
	return func(_ *state.State, ctx *plugin.Context) error {
		if required == 0 {
			return nil
		}

		if ctx.GuildID == 0 {
			if constant.DMPermissions.Has(required) {
				return nil
			}

			return assertChannelTypes(ctx, plugin.GuildChannels,
				errors.NewWithStack("restriction: invalid assertion BotPermissions with guild only permissions for "+
					"DM-only command"))
		}

		g, err := ctx.Guild()
		if err != nil {
			return err
		}

		c, err := ctx.Channel()
		if err != nil {
			return err
		}

		s, err := ctx.Self()
		if err != nil {
			return err
		}

		actual := discord.CalcOverwrites(*g, *c, *s)

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
			if constant.DMPermissions.Has(perms) {
				return nil
			}

			return assertChannelTypes(ctx, plugin.GuildChannels,
				errors.NewWithStack("restriction: invalid assertion UserPermissions with guild only permissions for "+
					"DM-only command"))
		}

		g, err := ctx.Guild()
		if err != nil {
			return err
		}

		c, err := ctx.Channel()
		if err != nil {
			return err
		}

		actual := discord.CalcOverwrites(*g, *c, *ctx.Member)

		missing := (actual & perms) ^ perms
		if missing == 0 {
			return nil
		}

		return newInsufficientUserPermissionsError(missing, ctx.Localizer)
	}
}
