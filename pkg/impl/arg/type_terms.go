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
// Numbers
// =====================================================================================

var (
	// ================================ Integer Meta Data ================================

	integerName        = i18n.NewFallbackConfig("args.types.integer.name", "Integer")
	integerDescription = i18n.NewFallbackConfig(
		"args.types.integer.description",
		"A whole number.")

	// ================================ Decimal Meta Data ================================

	decimalName        = i18n.NewFallbackConfig("args.types.decimal.name", "Decimal")
	decimalDescription = i18n.NewFallbackConfig(
		"args.types.decimal.description",
		"A decimal number.")
)

var (
	// ================================ Integer Errors ================================

	integerSyntaxError = i18n.NewFallbackConfig(
		"args.types.integer.errors.syntax.argument", "{{.raw}} is not an integer.")

	// ================================ Decimal Errors ================================

	decimalSyntaxError = i18n.NewFallbackConfig(
		"args.types.integer.errors.syntax.argument", "{{.raw}} is not an integer.")

	// ================================ Shared Errors ================================

	numberUnderRangeErrorArg = i18n.NewFallbackConfig(
		"args.types.number.errors.under_range.argument",
		"{{.raw}} is too small, try using a larger number as argument {{.postion}}.")
	numberUnderRangeErrorFlag = i18n.NewFallbackConfig(
		"args.types.number.errors.under_range.flag",
		"{{.raw}} is too small, try giving the `-{{.used_name}}`-flag a larger number.")

	numberOverRangeErrorArg = i18n.NewFallbackConfig(
		"args.types.number.errors.over_range.argument",
		"{{.raw}} is too large, try using a smaller number as argument {{.postion}}.")
	numberOverRangeErrorFlag = i18n.NewFallbackConfig(
		"args.types.integer.errors.over_range.flag",
		"{{.raw}} is a bit too large, try giving the `-{{.used_name}}`-flag a smaller number.")

	numberBelowMinErrorArg = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.argument",
		"Argument {{.position}} must be larger or equal to {{.min}}.")
	numberBelowMinErrorFlag = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.flag",
		"The `-{{.used_name}}`-flag must be larger or equal to {{.min}}.")

	numberAboveMaxErrorArg = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.argument",
		"Argument {{.position}} must be smaller or equal to {{.max}}.")
	numberAboveMaxErrorFlag = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.flag",
		"The `-{{.used_name}}`-flag must be smaller or equal to {{.max}}.")
)
