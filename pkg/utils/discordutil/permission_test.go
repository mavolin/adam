package discordutil

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/mock"
)

func TestPermissionNames(t *testing.T) {
	expect := []string{"Administrator", "Video"}

	perms := discord.PermissionAdministrator | discord.PermissionStream

	actual := PermissionNames(perms)
	assert.Equal(t, expect, actual)
}

func TestPermissionNamesl(t *testing.T) {
	expect := []string{"Ban Members", "View Channel"}

	perms := discord.PermissionBanMembers | discord.PermissionViewChannel

	l := mock.
		NewLocalizer().
		On("permissions.ban_members", "Ban Members").
		On("permissions.view_channel", "View Channel").
		Build()

	actual := PermissionNamesl(perms, l)
	assert.Equal(t, expect, actual)
}

func TestPermissionList(t *testing.T) {
	// maps don't have a deterministic order, so actual could be any one of these
	expect := "Ban Members, Manage Nicknames and View Channel"

	perms := discord.PermissionBanMembers | discord.PermissionManageNicknames | discord.PermissionViewChannel

	actual := PermissionList(perms)
	assert.Equal(t, expect, actual)
}

func TestPermissionListl(t *testing.T) {
	// maps don't have a deterministic order, so actual could be any one of these
	expect := "Ban Members, Manage Nicknames and View Channel"

	perms := discord.PermissionBanMembers | discord.PermissionManageNicknames | discord.PermissionViewChannel

	l := mock.
		NewLocalizer().
		On("lang.lists.last_separator", " and ").
		On("permissions.ban_members", "Ban Members").
		On("permissions.manage_nicknames", "Manage Nicknames").
		On("permissions.view_channel", "View Channel").
		On("lang.lists.default_separator", ", ").
		Build()

	actual := PermissionListl(perms, l)
	assert.Equal(t, expect, actual)
}
