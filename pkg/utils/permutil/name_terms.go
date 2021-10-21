package permutil

import (
	. "github.com/diamondburned/arikawa/v3/discord" //nolint:revive // make this file readable

	. "github.com/mavolin/adam/pkg/i18n" //nolint:revive // make this file readable
)

var permissionConfigs = map[Permissions]*Config{
	PermissionCreateInstantInvite: NewFallbackConfig("permission.create_instant_invite", "Create Invite"),
	PermissionKickMembers:         NewFallbackConfig("permission.kick_members", "Kick Members"),
	PermissionBanMembers:          NewFallbackConfig("permission.ban_members", "Ban Members"),
	PermissionAdministrator:       NewFallbackConfig("permission.administrator", "Administrator"),
	PermissionManageChannels:      NewFallbackConfig("permission.manage_channels", "Manage Channels"),
	PermissionManageGuild:         NewFallbackConfig("permission.manage_guild", "Manage Server"),
	PermissionAddReactions:        NewFallbackConfig("permission.add_reactions", "AddSource Reactions"),
	PermissionViewAuditLog:        NewFallbackConfig("permission.view_audit_log", "View Audit Log"),
	PermissionPrioritySpeaker:     NewFallbackConfig("permission.priority_speaker", "Priority Speaker"),
	PermissionStream:              NewFallbackConfig("permission.stream", "Video"),
	PermissionViewChannel:         NewFallbackConfig("permission.view_channel", "View Channel"),
	PermissionSendMessages:        NewFallbackConfig("permission.send_messages.text", "Send Messages"),
	PermissionSendTTSMessages:     NewFallbackConfig("permission.send_tts_messages", "Send TTS Messages"),
	PermissionManageMessages:      NewFallbackConfig("permission.manage_messages", "Manage Messages"),
	PermissionEmbedLinks:          NewFallbackConfig("permission.embed_links", "Embed Links"),
	PermissionAttachFiles:         NewFallbackConfig("permission.attach_files", "Attach Files"),
	PermissionReadMessageHistory:  NewFallbackConfig("permission.read_message_history", "Read Message History"),
	PermissionMentionEveryone: NewFallbackConfig("permission.mention_everyone",
		"Mention everyone, here and All Roles"),
	PermissionUseExternalEmojis: NewFallbackConfig("permission.use_external_emojis", "Use External Emojis"),
	//
	PermissionConnect:         NewFallbackConfig("permission.connect", "Connect"),
	PermissionSpeak:           NewFallbackConfig("permission.speak", "Speak"),
	PermissionMuteMembers:     NewFallbackConfig("permission.mute_members", "Mute Members"),
	PermissionDeafenMembers:   NewFallbackConfig("permission.deafen_members", "Deafen Members"),
	PermissionMoveMembers:     NewFallbackConfig("permission.move_members", "Move Members"),
	PermissionUseVAD:          NewFallbackConfig("permission.use_vad", "Use Voice Activity"),
	PermissionChangeNickname:  NewFallbackConfig("permission.change_nickname", "Change Nickname"),
	PermissionManageNicknames: NewFallbackConfig("permission.manage_nicknames", "Manage Nicknames"),
	PermissionManageRoles:     NewFallbackConfig("permission.manage_roles", "Manage Roles"),
	PermissionManageWebhooks:  NewFallbackConfig("permission.manage_webhooks", "Manage Webhooks"),
	PermissionManageEmojisAndStickers: NewFallbackConfig("permission.manage_emojis_and_stickers",
		"Manage Emojis and Stickers"),
	PermissionUseSlashCommands:     NewFallbackConfig("permission.use_slash_commands", "Use Slash Commands"),
	PermissionRequestToSpeak:       NewFallbackConfig("permission.request_to_speak", "Request to Speak"),
	PermissionManageThreads:        NewFallbackConfig("permission.manage_threads", "Manage Threads"),
	PermissionCreatePublicThreads:  NewFallbackConfig("permission.create_public_threads", "Create Public Threads"),
	PermissionCreatePrivateThreads: NewFallbackConfig("permission.create_private_threads", "Create Private Threads"),
	PermissionUseExternalStickers:  NewFallbackConfig("permission.use_external_stickers", "Use External Stickers"),
	PermissionSendMessagesInThreads: NewFallbackConfig("permission.send_messages_in_threads",
		"Send Messages in Threads"),
	PermissionStartEmbeddedActivities: NewFallbackConfig("permission.start_embedded_activities",
		"Start Embedded Activities"),
}
