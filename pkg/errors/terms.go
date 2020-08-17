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

	insufficientBotPermissionsDesc = localization.NewFallbackConfig(
		"errors.insufficient_bot_permissions.description", "I need the following permissions to execute the command:")
	insufficientBotPermissionPermissionField = localization.NewFallbackConfig(
		"errors.insufficient_bot_permissions.permission", "Permission")

	errorIDFooter = localization.NewFallbackConfig("errors.error_id", "Error-ID: {{.error_id}}")
)

type errorIDPlaceholders struct {
	ErrorID string
}
