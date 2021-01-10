package plugin

import "github.com/mavolin/adam/pkg/i18n"

// ================================ BotPermissionsError ================================

var (
	insufficientPermissionsDefault = i18n.NewFallbackConfig(
		"errors.insufficient_permissions.default",
		"I don't have sufficient permission to execute this command.")

	insufficientPermissionsDescSingle = i18n.NewFallbackConfig(
		"errors.insufficient_permissions.description.single",
		"It seems as if I don't have sufficient permissions to run this command. Please give me the "+
			`"{{.missing_permission}}" permission and try again.`)

	insufficientPermissionsDescMulti = i18n.NewFallbackConfig(
		"errors.insufficient_permissions.description.multi",
		"It seems as if I don't have sufficient permissions to run this command. Please give me the following "+
			"permissions and try again:")

	insufficientPermissionsMissingPermissionsFieldName = i18n.NewFallbackConfig(
		"errors.insufficient_permissions.fields.missing_permissions.name",
		"Missing Permissions")
)

type insufficientBotPermissionsDescSinglePlaceholders struct {
	MissingPermission string
}

// ================================ ChannelTypeError ================================

var (
	channelTypeErrorGuildText = i18n.NewFallbackConfig(
		"errors.channel_type.description.guild_text",
		"You can only use this command in a regular text channel.")

	channelTypeErrorGuildNews = i18n.NewFallbackConfig(
		"errors.channel_types.description.guild_news.",
		"You can only use this command in an announcement channel.")

	channelTypeErrorDirectMessage = i18n.NewFallbackConfig(
		"errors.channel_types.description.direct_message.",
		"You can only use this command in a direct message.")

	channelTypeErrorGuild = i18n.NewFallbackConfig(
		"errors.channel_types.description.guild.",
		"You can only use this command in a server.")

	channelTypeErrorDirectMessageAndGuildText = i18n.NewFallbackConfig(
		"errors.channel_types.description.direct_message_and_guild_text.",
		"You can only use this command in a direct message or a regular text channel.")

	channelTypeErrorDirectMessageAndGuildNews = i18n.NewFallbackConfig(
		"errors.channel_types.description.direct_message_and_guild_news.",
		"You can only use this command in a direct message or an announcement channel.")

	channelTypeErrorFallback = i18n.NewFallbackConfig(
		"errors.channel_type.description.fallback",
		"Ypu can't use this command in this type of channel.")
)

// ================================ RestrictionError ================================

var defaultRestrictionDesc = i18n.NewFallbackConfig(
	"errors.restriction.description.default",
	"👮 You are not allowed to use this command.")
