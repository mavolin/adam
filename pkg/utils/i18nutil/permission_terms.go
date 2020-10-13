package i18nutil

import (
	. "github.com/diamondburned/arikawa/discord"

	. "github.com/mavolin/adam/pkg/i18n"
)

var permissionConfigs = map[Permissions]Config{
	PermissionCreateInstantInvite: NewFallbackConfig("permissions.create_instant_invite", "Create Invite"),
	PermissionKickMembers:         NewFallbackConfig("permissions.kick_members", "Kick Members"),
	PermissionBanMembers:          NewFallbackConfig("permissions.ban_members", "Ban Members"),
	PermissionAdministrator:       NewFallbackConfig("permissions.administrator", "Administrator"),
	PermissionManageChannels:      NewFallbackConfig("permissions.manage_channels", "Manage Channels"),
	PermissionManageGuild:         NewFallbackConfig("permissions.manage_guild", "Manage Server"),
	PermissionAddReactions:        NewFallbackConfig("permissions.add_reactions", "Add Reactions"),
	PermissionViewAuditLog:        NewFallbackConfig("permissions.view_audit_log", "View Audit Log"),
	PermissionPrioritySpeaker:     NewFallbackConfig("permissions.priority_speaker", "Priority Speaker"),
	PermissionStream:              NewFallbackConfig("permissions.stream", "Video"),
	PermissionViewChannel:         NewFallbackConfig("permissions.view_channel", "View Channel"),
	PermissionSendMessages:        NewFallbackConfig("permissions.send_messages.text", "Send Messages"),
	PermissionSendTTSMessages:     NewFallbackConfig("permissions.send_tts_messages", "Send TTS Messages"),
	PermissionManageMessages:      NewFallbackConfig("permissions.manage_messages", "Manage Messages"),
	PermissionEmbedLinks:          NewFallbackConfig("permissions.embed_links", "Embed Links"),
	PermissionAttachFiles:         NewFallbackConfig("permissions.attach_files", "Attach Files"),
	PermissionReadMessageHistory:  NewFallbackConfig("permissions.read_message_history", "Read Message History"),
	PermissionMentionEveryone: NewFallbackConfig("permissions.mention_everyone",
		"Mention everyone, here and All Roles"),
	PermissionUseExternalEmojis: NewFallbackConfig("permissions.use_external_emojis", "Use External Emojis"),
	//
	PermissionConnect:         NewFallbackConfig("permissions.connect", "Connect"),
	PermissionSpeak:           NewFallbackConfig("permissions.speak", "Speak"),
	PermissionMuteMembers:     NewFallbackConfig("permissions.mute_members", "Mute Members"),
	PermissionDeafenMembers:   NewFallbackConfig("permissions.deafen_members", "Deafen Members"),
	PermissionMoveMembers:     NewFallbackConfig("permissions.move_members", "Move Members"),
	PermissionUseVAD:          NewFallbackConfig("permissions.use_vad", "Use Voice Activity"),
	PermissionChangeNickname:  NewFallbackConfig("permissions.change_nickname", "Change Nickname"),
	PermissionManageNicknames: NewFallbackConfig("permissions.manage_nicknames", "Manage Nicknames"),
	PermissionManageRoles:     NewFallbackConfig("permissions.manage_roles", "Manage Roles"),
	PermissionManageWebhooks:  NewFallbackConfig("permissions.manage_webhooks", "Manage Webhooks"),
	PermissionManageEmojis:    NewFallbackConfig("permissions.manage_emojis", "Manage Emojis"),
}
