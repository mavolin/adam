package constant

import "github.com/diamondburned/arikawa/discord"

// DMPermissions are the permissions that are granted in a private channel.
var DMPermissions = discord.PermissionAddReactions | discord.PermissionViewChannel | discord.PermissionSendMessages |
	discord.PermissionSendTTSMessages | discord.PermissionEmbedLinks | discord.PermissionAttachFiles |
	discord.PermissionReadMessageHistory | discord.PermissionUseExternalEmojis
