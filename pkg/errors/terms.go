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
		"errors.insufficient_bot_permissions.missing_permissions.name", "Missing Permissions")

	argumentParsingReasonFieldName = localization.NewFallbackConfig("errors.argument_parsing.reason.name", "Reason")

	errorIDFooter = localization.NewFallbackConfig("errors.error_id", "Error-ID: {{.error_id}}")
)

type errorIDFooterPlaceholders struct {
	ErrorID string
}

type insufficientBotPermissionsDescSinglePlaceholders struct {
	MissingPermission string
}
