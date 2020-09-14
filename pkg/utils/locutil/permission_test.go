package locutil

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/utils/mock"
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
		NewLocalizer(t).
		On("permissions.ban_members", "Ban Members").
		On("permissions.view_channel", "View Channel").
		Build()

	actual := PermissionNamesl(perms, l)
	assert.Equal(t, expect, actual)
}
