package permutil

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/stretchr/testify/assert"

	mocki18n "github.com/mavolin/adam/internal/mock/i18n"
)

func TestPermissionNames(t *testing.T) {
	expect := []string{"Administrator", "Video"}

	perms := discord.PermissionAdministrator | discord.PermissionStream

	actual := Names(perms)
	assert.Equal(t, expect, actual)
}

func TestPermissionNamesl(t *testing.T) {
	expect := []string{"Ban Members", "View Channel"}

	perms := discord.PermissionBanMembers | discord.PermissionViewChannel

	l := mocki18n.NewLocalizer(t).
		On("permission.ban_members", "Ban Members").
		On("permission.view_channel", "View Channel").
		Build()

	actual := Namesl(perms, l)
	assert.Equal(t, expect, actual)
}
