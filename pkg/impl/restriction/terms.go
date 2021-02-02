package restriction

import "github.com/mavolin/adam/pkg/i18n"

// ================================ Any ================================

var (
	anyMessageHeader = i18n.NewFallbackConfig(
		"restriction.comparator.any.header",
		"You need to fulfill at least one of these requirements to execute the command:")
	anyMessageInline = i18n.NewFallbackConfig(
		"restriction.comparator.any.inline",
		"You need to fulfill at least one of these requirements:")
)

// ================================ All ================================

var (
	allMessageHeader = i18n.NewFallbackConfig(
		"restriction.comparator.all.header",
		"You need to fulfill all of these requirements to execute the command:")
	allMessageInline = i18n.NewFallbackConfig(
		"restriction.comparator.all.inline",
		"You need to fulfill all of these requirements:")
)

// ================================ NSFW ================================

var notNSFWChannelError = i18n.NewFallbackConfig(
	"restriction.nsfw.error.not_nsfw.description",
	"This command must be invoked in an NSFW channel.")

// ================================ Owner ================================

var notOwnerError = i18n.NewFallbackConfig(
	"restriction.owner.error.not_owner.description",
	"You must be the owner of the server to invoke this command.")

// ================================ BotOwner ================================

var notBotOwnerError = i18n.NewFallbackConfig(
	"restriction.bot_owner.error.not_bot_owner.description",
	"You must be the owner of the bot to invoke this command.")

// ================================ Roles ================================

var (
	missingRoleError = i18n.NewFallbackConfig(
		"restriction.roles.error.missing_role",
		"You need the {{.role}} role to use this command.")
	missingRolesAllError = i18n.NewFallbackConfig(
		"restriction.roles.error.missing_roles.all",
		"You need these roles to use this command:")
	missingRolesAnyError = i18n.NewFallbackConfig(
		"restriction.roles.error.missing_roles.any",
		"You need at least one of these roles to use this command:")
)

type missingRoleErrorPlaceholders struct {
	Role string
}

// ================================ Channels ================================

var (
	blockedChannelErrorSingle = i18n.NewFallbackConfig(
		"restriction.channels.error.blocked_channel.single",
		"You can only use this command in {{.channel}}.")
	blockedChannelErrorMulti = i18n.NewFallbackConfig(
		"restriction.channels.error.blocked_channel.multi",
		"You can only use this command in one of these channels:")
)

type blockedChannelErrorSinglePlaceholders struct {
	Channel string
}

// ================================ UserPermissions ================================

var (
	userPermissionsDescSingle = i18n.NewFallbackConfig(
		"restriction.user_permissions.error.insufficient_permissions.description.single",
		`You need the "{{.missing_permission}}" permission to use this command.`)
	userPermissionsDescMulti = i18n.NewFallbackConfig(
		"restriction.user_permissions.error.insufficient_permissions.description.multi",
		"You need these permissions to use this command:")
)

type userPermissionsDescSinglePlaceholders struct {
	MissingPermission string
}
