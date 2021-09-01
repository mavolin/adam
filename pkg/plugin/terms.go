package plugin

import "github.com/mavolin/adam/pkg/i18n"

// ================================ BotPermissionsError ================================

var (
	botPermissionsDefault = i18n.NewFallbackConfig(
		"plugin.error.bot_permissions.default",
		"I don't have sufficient permission to execute this command.")

	botPermissionsDescSingle = i18n.NewFallbackConfig(
		"plugin.error.bot_permissions.description.single",
		"I don't have sufficient permissions to run this command. Please give me the "+
			`"{{.missing_permission}}" permission and try again.`)

	botPermissionsDescMulti = i18n.NewFallbackConfig(
		"plugin.error.bot_permissions.description.multi",
		"I don't have sufficient permissions to run this command. Please give me the "+
			"following permissions and try again:")

	botPermissionsMissingPermissionsFieldName = i18n.NewFallbackConfig(
		"plugin.error.bot_permissions.fields.missing_permissions.name",
		"Missing Permissions")
)

type botPermissionsDescSinglePlaceholders struct {
	MissingPermission string
}

// ================================ ChannelTypeError ================================

var (
	channelTypeErrorGuildText = i18n.NewFallbackConfig(
		"plugin.error.channel_type.guild_text",
		"You can only use this command in a regular text channel.")

	channelTypeErrorGuildNews = i18n.NewFallbackConfig(
		"plugin.error.channel_types.guild_news",
		"You can only use this command in an announcement channel.")

	channelTypeErrorDM = i18n.NewFallbackConfig(
		"plugin.error.channel_types.dm",
		"You can only use this command in a direct message.")

	channelTypeErrorThread = i18n.NewFallbackConfig(
		"plugin.error.channel_types.thread",
		"You can only use this command in a thread.")

	channelTypeErrorGuildTextAndGuildNews = i18n.NewFallbackConfig(
		"plugin.error.channel_types.guild_text_and_guild_news",
		"You can only use this command in a regular text channel or an announcement channel.")

	channelTypeErrorGuildTextAndThread = i18n.NewFallbackConfig(
		"plugin.error.channel_types.guild_text_and_thread",
		"You can only use this command in a regular text channel or a thread.")

	channelTypeErrorGuildNewsAndThread = i18n.NewFallbackConfig(
		"plugin.error.channel_types.guild_news_and_thread",
		"You can only use this command in an announcement channel or a thread.")

	channelTypeErrorGuild = i18n.NewFallbackConfig(
		"plugin.error.channel_types.guild",
		"You can only use this command in a server.")

	channelTypeErrorDMAndGuildText = i18n.NewFallbackConfig(
		"plugin.error.channel_types.dm_and_guild_text",
		"You can only use this command in a direct message or a regular text channel.")

	channelTypeErrorDMAndGuildNews = i18n.NewFallbackConfig(
		"plugin.error.channel_types.dm_and_guild_news.",
		"You can only use this command in a direct message or an announcement channel.")

	channelTypeErrorDMAndThread = i18n.NewFallbackConfig(
		"plugin.error.channel_types.dm_and_thread",
		"You can only use this command in a direct message or a thread.")

	channelTypeErrorFallback = i18n.NewFallbackConfig(
		"plugin.error.channel_type.fallback",
		"You can't use this command in this type of channel.")
)

// ================================ RestrictionError ================================

var defaultRestrictionDesc = i18n.NewFallbackConfig(
	"plugin.error.restriction.description.default",
	"You are not allowed to use this command.")
