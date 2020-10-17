package arg

import "github.com/mavolin/adam/pkg/i18n"

// =============================================================================
// Common
// =====================================================================================

// ================================ Errors ================================

var (
	regexpNotMatchingErrorArg = i18n.NewFallbackConfig(
		"args.types.common.errors.regexp_not_matching.arg",
		"Argument {{.position}} must match `{{.regexp}}`.")
	regexpNotMatchingErrorFlag = i18n.NewFallbackConfig(
		"args.types.common.errors.regexp_not_matching.flag",
		"The `-{{.used_name}}`-flag must match `{{.regexp}}`.")
)

// =============================================================================
// Switch
// =====================================================================================

// ================================ Meta Data ================================

var (
	switchName        = i18n.NewFallbackConfig("args.types.switch.name", "Switch")
	switchDescription = i18n.NewFallbackConfig(
		"args.types.switch.description",
		"Used to turn on a feature of a command. Only used with flags.")
)

// ================================ Errors ================================

var switchWithContentError = i18n.NewFallbackConfig(
	"args.types.switch.errors.with_content", "`-{{.name}}` is a Switch flag and cannot be used with content.")

type switchWithContentErrorPlaceholders struct {
	Name string
}

// =============================================================================
// Numbers
// =====================================================================================

// ================================ Integer Meta Data ================================

var (
	integerName        = i18n.NewFallbackConfig("args.types.integer.name", "Integer")
	integerDescription = i18n.NewFallbackConfig("args.types.integer.description", "A whole number.")
)

// ================================ Decimal Meta Data ================================

var (
	decimalName        = i18n.NewFallbackConfig("args.types.decimal.name", "Decimal")
	decimalDescription = i18n.NewFallbackConfig("args.types.decimal.description", "A decimal number.")
)

// ================================ Integer Errors ================================

var integerSyntaxError = i18n.NewFallbackConfig("args.types.integer.errors.syntax", "{{.raw}} is not an integer.")

// ================================ Decimal Errors ================================

var decimalSyntaxError = i18n.NewFallbackConfig("args.types.decimal.errors.syntax", "{{.raw}} is not a decimal.")

// ================================ Shared Errors ================================

var (
	numberUnderRangeErrorArg = i18n.NewFallbackConfig(
		"args.types.number.errors.under_range.arg",
		"{{.raw}} is too small, try using a larger number as argument {{.position}}.")
	numberUnderRangeErrorFlag = i18n.NewFallbackConfig(
		"args.types.number.errors.under_range.flag",
		"{{.raw}} is too small, try giving the `-{{.used_name}}`-flag a larger number.")

	numberOverRangeErrorArg = i18n.NewFallbackConfig(
		"args.types.number.errors.over_range.arg",
		"{{.raw}} is too large, try using a smaller number as argument {{.position}}.")
	numberOverRangeErrorFlag = i18n.NewFallbackConfig(
		"args.types.integer.errors.over_range.flag",
		"{{.raw}} is a too large, try giving the `-{{.used_name}}`-flag a smaller number.")

	numberBelowMinErrorArg = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.arg",
		"Argument {{.position}} must be larger or equal to {{.min}}.")
	numberBelowMinErrorFlag = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.flag",
		"The `-{{.used_name}}`-flag must be larger or equal to {{.min}}.")

	numberAboveMaxErrorArg = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.arg",
		"Argument {{.position}} must be smaller or equal to {{.max}}.")
	numberAboveMaxErrorFlag = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.flag",
		"The `-{{.used_name}}`-flag must be smaller or equal to {{.max}}.")
)

// =============================================================================
// Text
// =====================================================================================

// ================================ Meta Data ================================

var (
	textName        = i18n.NewFallbackConfig("args.types.text.name", "Text")
	textDescription = i18n.NewFallbackConfig("args.types.text.desc", "A text. What else is there to say.")
)

// ================================ Errors ================================

var (
	textBelowMinLengthErrorArg = i18n.NewFallbackConfig(
		"args.types.text.errors.below_min_length.arg",
		"The text in argument {{.position}} must be at least {{.min}} characters long.")
	textBelowMinLengthErrorFlag = i18n.NewFallbackConfig(
		"args.types.text.errors.below_min_length.flag",
		"The text used in the `-{{.used_name}}`-flag must be at least {{.min}} characters long.")

	textAboveMaxLengthErrorArg = i18n.NewFallbackConfig(
		"args.types.text.errors.above_max_length.arg",
		"The text in argument {{.position}} may not be longer than {{.max}} characters.")
	textAboveMaxLengthErrorFlag = i18n.NewFallbackConfig(
		"args.types.text.errors.above_max_length.flag",
		"The text used in the `-{{.used_name}}`-flag may not be longer than {{.max}} characters.")
)
