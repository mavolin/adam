package errors

import "github.com/mavolin/adam/pkg/localization"

var (
	errorTitle = localization.NewFallbackConfig("error.title", "Error")
	infoTitle  = localization.NewFallbackConfig("info.title", "Info")

	defaultInternalDesc = localization.NewFallbackConfig("errors.internal.description.default",
		"Oh no! Something went wrong and I couldn't finish executing your command. I've informed my team and they'll "+
			"get on fixing the bug asap.")

	defaultRestrictionDesc = localization.NewFallbackConfig("errors.restriction.description.default",
		"👮 You are not allowed to use this command.")

	insufficientBotPermissionsDescSingle = localization.NewFallbackConfig(
		"errors.insufficient_bot_permissions.description.single",
		"It seems as if I don't have sufficient permissions to run this command. Please give me the "+
			`"{{.missing_permission}}" permission and try again.`)
	insufficientBotPermissionsDescMulti = localization.NewFallbackConfig(
		"errors.insufficient_bot_permissions.description.multi",
		"It seems as if I don't have sufficient permissions to run this command. Please give me the following "+
			"permissions and try again.")
	insufficientBotPermissionMissingMissingPermissionsFieldName = localization.NewFallbackConfig(
		"errors.insufficient_bot_permissions.fields.missing_permissions.name",
		"Missing Permissions")

	argumentParsingReasonFieldName = localization.NewFallbackConfig("errors.argument_parsing.reason.name", "Reason")

	channelTypeErrorGuildText = localization.NewFallbackConfig(
		"errors.channel_type.description.guild_text",
		"You must use this command in a regular text channel.")
	channelTypeErrorGuildNews = localization.NewFallbackConfig(
		"errors.channel_types.description.guild_news.",
		"You must use this command in an announcement channel.")
	channelTypeErrorDirectMessage = localization.NewFallbackConfig(
		"errors.channel_types.description.direct_message.",
		"You must use this command in a direct message.")
	channelTypeErrorGuild = localization.NewFallbackConfig(
		"errors.channel_types.description.guild.",
		"You must use this command in a server.")
	channelTypeErrorDirectMessageAndGuildText = localization.NewFallbackConfig(
		"errors.channel_types.description.direct_message_and_guild_text.",
		"You must use this command in a direct message or a regular text channel.")
	channelTypeErrorDirectMessageAndGuildNews = localization.NewFallbackConfig(
		"errors.channel_types.description.direct_message_and_guild_news.",
		"You must use this command in a direct message or a announcement channel.")
	channelTypeErrorFallback = localization.NewFallbackConfig(
		"errors.channel_type.description.fallback",
		"Ypu can't use this command in this type of channel.")
)

type insufficientBotPermissionsDescSinglePlaceholders struct {
	MissingPermission string
}
