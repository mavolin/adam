package restriction

import (
	"sort"

	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/pluginutil"
)

// assertChannelTypes asserts that the command with the passed context
// is used in the passed channel types.
//
// assertChannelTypes will also silently report errors in some cases.
func assertChannelTypes(ctx *plugin.Context, assertTypes plugin.ChannelTypes, noRemainingError error) error {
	if assertTypes&plugin.AllChannels == plugin.AllChannels {
		return nil
	}

	if ctx.GuildID == 0 { // we are in a DM
		// we assert a DM
		if assertTypes&plugin.DirectMessages == plugin.DirectMessages {
			return nil
		}
		// no DM falls through
	} else { // we are in a guild
		// we assert all guild channels
		if assertTypes&plugin.GuildChannels == plugin.GuildChannels {
			return nil

			// we assert something other than all guild channels
		} else if !(assertTypes&plugin.GuildChannels == 0) {
			c, err := ctx.Channel()
			if err != nil {
				return err
			}

			if assertTypes.Has(c.Type) {
				return nil
			}
		}
		// not all guild types falls through
	}

	channelTypes, err := pluginutil.ChannelTypes(ctx.CommandIdentifier, ctx.Provider)
	if err != nil {
		return err
	}

	allowed := channelTypes & assertTypes
	if allowed == 0 { // no channel types remaining
		// there is no need to prevent execution, as another restriction
		// may permit it, still we should capture this
		ctx.HandleErrorSilent(noRemainingError)

		return errors.DefaultFatalRestrictionError
	}

	return newInvalidChannelTypeError(allowed, ctx.Localizer, true)
}

// canMangeRole checks if the passed member of the passed guild is able to
// modify the passed role.
func canManageRole(target discord.Role, g *discord.Guild, m *discord.Member) bool {
RoleIDs:
	for _, id := range m.RoleIDs {
		for _, r := range g.Roles {
			if r.ID == id {
				if r.Position > target.Position {
					goto Found
				}

				continue RoleIDs
			}
		}
	}

	return false

Found:
	// manage roles can't be set on a channel level, we can just pass an empty
	// channel
	perms := discord.CalcOverwrites(*g, discord.Channel{}, *m)

	return perms.Has(discord.PermissionManageRoles)
}

// insertRoleSorted inserts the passed discord.Role into the passed slice of
// discord.Roles, while keeping the order.
//
// This assumes that roles is sorted in ascending order by position.
func insertRoleSorted(r discord.Role, roles []discord.Role) []discord.Role {
	i := sort.Search(len(roles), func(i int) bool {
		return roles[i].Position >= r.Position
	})

	// keep missingRoles sorted
	if i == len(roles) {
		roles = append(roles, r)
	} else {
		copy(roles[i+1:], roles[i:])
		roles[i] = r
	}

	return roles
}
