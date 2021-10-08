package permutil

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/stretchr/testify/assert"

	mocki18n "github.com/mavolin/adam/internal/mock/i18n"
)

func TestPermissionNames(t *testing.T) {
	t.Parallel()

	expect := []string{"Ban Members", "View Channel"}

	perms := discord.PermissionBanMembers | discord.PermissionViewChannel

	l := mocki18n.NewLocalizer(t).
		On("permission.ban_members", "Ban Members").
		On("permission.view_channel", "View Channel").
		Build()

	actual := Names(l, perms)
	assert.Equal(t, expect, actual)
}
