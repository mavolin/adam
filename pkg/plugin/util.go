package plugin

import "github.com/diamondburned/arikawa/discord"

// Permissions is an utility used to create a pointer to a discord.Permissions.
func Permissions(perms discord.Permissions) *discord.Permissions { return &perms }
