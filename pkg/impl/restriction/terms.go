package restriction

import "github.com/mavolin/adam/pkg/localization"

var (
	anyMessageHeader = localization.NewFallbackConfig(
		"restrictions.compare.any.header",
		"You need to fulfill at least one of these requirements to execute the command:")
	anyMessageInline = localization.NewFallbackConfig(
		"restrictions.compare.any.inline",
		"You need to fulfill at least one of these requirements:")

	allMessageHeader = localization.NewFallbackConfig(
		"restrictions.compare.all.header",
		"You need to fulfill all of these requirements to execute the command:")
	allMessageInline = localization.NewFallbackConfig(
		"restrictions.compare.all.inline",
		"You need to fulfill all of these requirements:")

	notNSFWChannelError = localization.NewFallbackConfig(
		"restrictions.nsfw.errors.not_nsfw.description",
		"This command must be invoked in an NSFW channel.")

	notOwnerError = localization.NewFallbackConfig(
		"restrictions.owner.errors.not_owner.description",
		"You must be the owner of the server to invoke this command.")

	notBotOwnerError = localization.NewFallbackConfig(
		"restrictions.bot_owner.errors.not_bot_owner.description",
		"You must be the owner of the bot to invoke this command.")

	missingRoleError = localization.NewFallbackConfig(
		"restrictions.roles.errors.missing_role",
		"You need the {{.role}} role to use this command.")
	missingRolesAllError = localization.NewFallbackConfig(
		"restrictions.roles.errors.missing_roles.all",
		"You need these roles to use this command:")
	missingRolesAnyError = localization.NewFallbackConfig(
		"restrictions.roles.errors.missing_roles.any",
		"You need at least one of these roles to use this command:")

	blockedChannelErrorSingle = localization.NewFallbackConfig(
		"restrictions.channels.errors.blocked_channel.single",
		"You can only use this command in {{.channel}}.")
	blockedChannelErrorMulti = localization.NewFallbackConfig(
		"restrictions.channels.errors.blocked_channel.multi",
		"You must use this command in one of these channels:")

	insufficientUserPermissionsDescSingle = localization.NewFallbackConfig(
		"errors.insufficient_user_permissions.description.single",
		`You need the "{{.missing_permission}}" permission to use this command.`)
	insufficientUserPermissionsDescMulti = localization.NewFallbackConfig(
		"errors.insufficient_user_permissions.description.multi",
		"You need these permissions to use this command:")
)

type (
	missingRoleErrorPlaceholders struct {
		Role string
	}

	blockedChannelErrorSinglePlaceholders struct {
		Channel string
	}

	insufficientUserPermissionsDescSinglePlaceholders struct {
		MissingPermission string
	}
)
