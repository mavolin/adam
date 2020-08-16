package discordutil

import (
	. "github.com/diamondburned/arikawa/discord"

	. "github.com/mavolin/adam/pkg/localization"
)

var permissionConfigs = map[Permissions]Config{
	PermissionCreateInstantInvite: QuickFallbackConfig("permissions.create_instant_invite", "Create Invite"),
	PermissionKickMembers:         QuickFallbackConfig("permissions.kick_members", "Kick Members"),
	PermissionBanMembers:          QuickFallbackConfig("permissions.ban_members", "Ban Members"),
	PermissionAdministrator:       QuickFallbackConfig("permissions.administrator", "Administrator"),
	PermissionManageChannels:      QuickFallbackConfig("permissions.manage_channels", "Manage Channels"),
	PermissionManageGuild:         QuickFallbackConfig("permissions.manage_guild", "Manage Server"),
	PermissionAddReactions:        QuickFallbackConfig("permissions.add_reactions", "Add Reactions"),
	PermissionViewAuditLog:        QuickFallbackConfig("permissions.view_audit_log", "View Audit Log"),
	PermissionPrioritySpeaker:     QuickFallbackConfig("permissions.priority_speaker", "Priority Speaker"),
	PermissionStream:              QuickFallbackConfig("permissions.stream", "Video"),
	PermissionViewChannel:         QuickFallbackConfig("permissions.view_channel", "View Channel"),
	PermissionSendMessages:        QuickFallbackConfig("permissions.send_messages.text", "Send Messages"),
	PermissionSendTTSMessages:     QuickFallbackConfig("permissions.send_tts_messages", "Send TTS Messages"),
	PermissionManageMessages:      QuickFallbackConfig("permissions.manage_messages", "Manage Messages"),
	PermissionEmbedLinks:          QuickFallbackConfig("permissions.embed_links", "Embed Links"),
	PermissionAttachFiles:         QuickFallbackConfig("permissions.attach_files", "Attach Files"),
	PermissionReadMessageHistory:  QuickFallbackConfig("permissions.read_message_history", "Read Message History"),
	PermissionMentionEveryone: QuickFallbackConfig("permissions.mention_everyone",
		"Mention everyone, here and All Roles"),
	PermissionUseExternalEmojis: QuickFallbackConfig("permissions.use_external_emojis", "Use External Emojis"),
	//
	PermissionConnect:         QuickFallbackConfig("permissions.connect", "Connect"),
	PermissionSpeak:           QuickFallbackConfig("permissions.speak", "Speak"),
	PermissionMuteMembers:     QuickFallbackConfig("permissions.mute_members", "Mute Members"),
	PermissionDeafenMembers:   QuickFallbackConfig("permissions.deafen_members", "Deafen Members"),
	PermissionMoveMembers:     QuickFallbackConfig("permissions.move_members", "Move Members"),
	PermissionUseVAD:          QuickFallbackConfig("permissions.use_vad", "Use Voice Activity"),
	PermissionChangeNickname:  QuickFallbackConfig("permissions.change_nickname", "Change Nickname"),
	PermissionManageNicknames: QuickFallbackConfig("permissions.manage_nicknames", "Manage Nicknames"),
	PermissionManageRoles:     QuickFallbackConfig("permissions.manage_roles", "Manage Roles"),
	PermissionManageWebhooks:  QuickFallbackConfig("permissions.manage_webhooks", "Manage Webhooks"),
	PermissionManageEmojis:    QuickFallbackConfig("permissions.manage_emojis", "Manage Emojis"),
}
