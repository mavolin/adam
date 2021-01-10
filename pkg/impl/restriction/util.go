package restriction

import (
	"sort"

	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/adam/pkg/plugin"
)

// assertChannelTypes asserts that the command with the passed context
// is used in the passed channel types.
//
// assertChannelTypes will also silently report errors in some cases.
func assertChannelTypes(ctx *plugin.Context, allowed plugin.ChannelTypes, noRemainingError error) error {
	ok, err := allowed.Check(ctx)
	if err != nil {
		return err
	} else if ok {
		return nil
	}

	remaining := ctx.InvokedCommand.ChannelTypes & allowed
	if remaining == 0 { // no channel types remaining
		// there is no need to prevent execution, as another restriction
		// may permit it, still we should capture this
		ctx.HandleErrorSilent(noRemainingError)

		return plugin.DefaultFatalRestrictionError
	}

	fatal := false

	if ctx.GuildID == 0 && remaining&plugin.DirectMessages == 0 {
		fatal = true
	} else if ctx.GuildID != 0 && remaining == plugin.DirectMessages {
		fatal = true
	}

	return newInvalidChannelTypeError(remaining, ctx.Localizer, fatal)
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

	roles = append(roles, r) // make space for another element

	// only insert if r wasn't supposed to go to the end
	if i < len(roles) {
		copy(roles[i+1:], roles[i:])
		roles[i] = r
	}

	return roles
}
