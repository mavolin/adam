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
		"args.types.switch.errors.with_content", "`-{{.name}}` is a Switch flag and cannot be used with content.")
)

type (
	// ================================ Errors ================================
	switchWithContentErrorPlaceholders struct {
		Name string
	}
)

// =============================================================================
// Integer
// =====================================================================================

var (
	// ================================ Meta Data ================================

	integerName        = i18n.NewFallbackConfig("args.types.integer.name", "Integer")
	integerDescription = i18n.NewFallbackConfig(
		"args.types.integer.description",
		"A whole number.")

	// ================================ Errors ================================

	integerSyntaxError = i18n.NewFallbackConfig(
		"args.types.integer.errors.syntax.argument", "{{.raw}} is not a number.")

	integerUnderRangeErrorArg = i18n.NewFallbackConfig(
		"args.types.integer.errors.under_range.argument",
		"{{.raw}} is too small, try using a larger number as argument {{.postion}}.")
	integerUnderRangeErrorFlag = i18n.NewFallbackConfig(
		"args.types.integer.errors.under_range.flag",
		"{{.raw}} is too small, try giving the `-{{.used_name}}`-flag a larger number.")

	integerOverRangeErrorArg = i18n.NewFallbackConfig(
		"args.types.integer.errors.over_range.argument",
		"{{.raw}} is too large, try using a smaller number as argument {{.postion}}.")
	integerOverRangeErrorFlag = i18n.NewFallbackConfig(
		"args.types.integer.errors.over_range.flag",
		"{{.raw}} is a bit too large, try giving the `-{{.used_name}}`-flag a smaller number.")

	integerBelowMinErrorArg = i18n.NewFallbackConfig(
		"args.types.integer.errors.below_min.argument",
		"Argument {{.position}} must be larger or equal to {{.min}}.")
	integerBelowMinErrorFlag = i18n.NewFallbackConfig(
		"args.types.integer.errors.below_min.flag",
		"The `-{{.used_name}}`-flag must be larger or equal to {{.min}}.")

	integerAboveMaxErrorArg = i18n.NewFallbackConfig(
		"args.types.integer.errors.below_min.argument",
		"Argument {{.position}} must be smaller or equal to {{.max}}.")
	integerAboveMaxErrorFlag = i18n.NewFallbackConfig(
		"args.types.integer.errors.below_min.flag",
		"The `-{{.used_name}}`-flag must be smaller or equal to {{.max}}.")
)
