package errors

import "github.com/mavolin/adam/pkg/i18n"

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
