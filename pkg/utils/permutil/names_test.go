package permutil

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
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

	l := newMockedLocalizer(t).
		on("permissions.ban_members", "Ban Members").
		on("permissions.view_channel", "View Channel").
		build()

	actual := Namesl(perms, l)
	assert.Equal(t, expect, actual)
}
