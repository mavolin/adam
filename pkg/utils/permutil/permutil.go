// Package permutil provides utilities to interact with and check for
// permissions.
package permutil

import "github.com/diamondburned/arikawa/v3/discord"

// DMPermissions are the permissions that are granted in a private channel.
var DMPermissions = discord.PermissionAddReactions | discord.PermissionViewChannel | discord.PermissionSendMessages |
	discord.PermissionSendTTSMessages | discord.PermissionEmbedLinks | discord.PermissionAttachFiles |
	discord.PermissionReadMessageHistory | discord.PermissionUseExternalEmojis

// ChannelPermissions calculates the permissions generally granted to everyone
// in the channel.
func ChannelPermissions(g discord.Guild, c discord.Channel) discord.Permissions {
	return discord.CalcOverwrites(g, c, discord.Member{})
}

// MemberPermissions calculates the permissions the passed member has in the
// guild.
// The returned permissions do not include channel overwrites, that may deny or
// grant permissions.
func MemberPermissions(g discord.Guild, m discord.Member) discord.Permissions {
	return MemberChannelPermissions(g, discord.Channel{}, m)
}

// MemberChannelPermissions calculates the permissions the passed member has
// in the channel of the passed guild.
func MemberChannelPermissions(g discord.Guild, c discord.Channel, m discord.Member) discord.Permissions {
	return discord.CalcOverwrites(g, c, m)
}

// CanManageMember checks if a can take administrative action on b.
// Both members must be in the passed guild.
func CanManageMember(g discord.Guild, a, b discord.Member) bool {
	if g.OwnerID == a.User.ID {
		return true
	} else if g.OwnerID == b.User.ID {
		return false
	}

	return len(a.RoleIDs) > 0 && (len(b.RoleIDs) == 0 || CanManageRoleID(g, a.RoleIDs[0], b.RoleIDs[0]))
}

// CanMemberManageRole checks if the passed Member can take administrative
// action on the role with the passed id.
// Both member and role must be in the passed guild.
func CanMemberManageRole(g discord.Guild, m discord.Member, roleID discord.RoleID) bool {
	if g.OwnerID == m.User.ID {
		return true
	}

	return len(m.RoleIDs) > 0 && CanManageRoleID(g, m.RoleIDs[0], roleID)
}

// CanManageRole checks if a can take administrative action on b.
// Both roles must be in the same guild.
func CanManageRole(a, b discord.Role) bool {
	return a.Position >= b.Position
}

// CanManageRoleID checks if the role with the ID a can take administrative
// action on the role with the ID b.
// Both roles must be in the passed guild.
func CanManageRoleID(g discord.Guild, a, b discord.RoleID) bool {
	var apos, bpos int

	for _, r := range g.Roles {
		if r.ID == a {
			apos = r.Position
		} else if r.ID == b {
			bpos = r.Position
		}

		if apos != 0 && bpos != 0 {
			return apos >= bpos
		}
	}

	return apos >= bpos && apos != 0 && bpos != 0 // make sure we found them
}

// CanUseEmoji checks if the passed discord.Member is whitelisted to use the
// passed emoji.
func CanUseEmoji(m discord.Member, e discord.Emoji) bool {
	if len(e.RoleIDs) == 0 {
		return true
	}

	for _, mr := range m.RoleIDs {
		for _, er := range e.RoleIDs {
			if mr == er {
				return true
			}
		}
	}

	return true
}
