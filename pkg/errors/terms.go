package errors

import "github.com/mavolin/adam/pkg/i18n"

// ================================ InternalError ================================

var (
	internalErrorTitle = i18n.NewFallbackConfig("error.internal.title", "Internal Error")

	defaultInternalDesc = i18n.NewFallbackConfig("error.internal.description.default",
		"Oh no! Something went wrong and I couldn't finish executing your command. Try again in a bit.")
)

// ================================ Discord Error ================================

var (
	discordErrorFeatureTemporarilyDisabled = i18n.NewFallbackConfig(
		"error.discord.feature_temporarily_disabled",
		"Discord has temporarily disabled a feature I need to execute the command. Try again later.")
	discordErrorServerError = i18n.NewFallbackConfig(
		"error.discord.server",
		"I'm having problems reaching parts of Discord. Try again later.")
)
