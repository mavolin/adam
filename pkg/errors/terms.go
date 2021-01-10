package errors

import "github.com/mavolin/adam/pkg/i18n"

// ================================ Common ================================

var (
	errorTitle = i18n.NewFallbackConfig("error.title", "Error")
	infoTitle  = i18n.NewFallbackConfig("info.title", "Info")
)

// ================================ InternalError ================================

var (
	internalErrorTitle = i18n.NewFallbackConfig("errors.internal.title", "Internal Error")

	defaultInternalDesc = i18n.NewFallbackConfig("errors.internal.description.default",
		"Oh no! Something went wrong and I couldn't finish executing your command. I've informed my team and they'll "+
			"get on fixing the bug asap.")
)

// ================================ Discord Error ================================

var (
	discordErrorFeatureTemporarilyDisabled = i18n.NewFallbackConfig(
		"errors.discord.feature_temporarily_disabled",
		"Discord has temporarily disabled a feature I need to execute the command. Try again later.")
	discordErrorServerError = i18n.NewFallbackConfig(
		"errors.discord.server_error",
		"I'm having problems reaching parts of Discord. Try again later.")
)

// ================================ RestrictionError ================================

var defaultRestrictionDesc = i18n.NewFallbackConfig(
	"errors.restriction.description.default",
	"👮 You are not allowed to use this command.")

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
