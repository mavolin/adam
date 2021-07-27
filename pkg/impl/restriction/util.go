package restriction

import (
	"fmt"
	"sort"

	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/adam/pkg/plugin"
)

// assertChannelTypes asserts that the command with the passed context
// is used in the passed channel types.
// If that is not the case a *plugin.RestrictionError generated through
// newChannelTypesError is returned.
//
// assertChannelTypes will also silently report errors in some cases.
func assertChannelTypes(ctx *plugin.Context, allowed plugin.ChannelTypes) error {
	ok, err := allowed.Check(ctx)
	if err != nil {
		return err
	} else if ok {
		return nil
	}

	remaining := ctx.InvokedCommand.ChannelTypes() & allowed
	if remaining == 0 { // no channel types remaining
		// there is no need to prevent execution, as another restriction
		// may permit it, still we should capture this
		ctx.HandleErrorSilently(fmt.Errorf("restriction: need channel types %s, but command only allows %s",
			allowed, ctx.InvokedCommand.ChannelTypes()))

		return plugin.DefaultFatalRestrictionError
	}

	fatal := (ctx.GuildID == 0 && remaining&plugin.DirectMessages == 0) ||
		(ctx.GuildID != 0 && remaining == plugin.DirectMessages)

	return newChannelTypesError(remaining, ctx.Localizer, fatal)
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
