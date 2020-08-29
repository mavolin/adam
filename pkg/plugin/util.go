package plugin

import "github.com/diamondburned/arikawa/discord"

// NoPermissions defines explicitly that a command or module requires no
// permissions.
var NoPermissions = Permissions(0)

// Permissions is an utility used to create a pointer to a discord.Permissions.
func Permissions(perms discord.Permissions) *discord.Permissions { return &perms }
