package arg

import "github.com/mavolin/adam/pkg/i18n"

// =============================================================================
// Switch
// =====================================================================================

var (
	// ================================ Meta Data ================================

	switchName        = i18n.NewFallbackConfig("args.types.switch.name", "Switch")
	switchDescription = i18n.NewFallbackConfig(
		"args.types.switch.description",
		"Used to turn on a feature of a command. Only used with flags.")

	// ================================ Errors ================================

	switchWithContentError = i18n.NewFallbackConfig(
		"args.types.switch.errors.with_content", "`{{.name}}` is a Switch flag and cannot be used with content.")
)

type (
	// ================================ Errors ================================
	switchWithContentErrorPlaceholders struct {
		Name string
	}
)
