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

// =============================================================================
// ID
// =====================================================================================

// ================================ Meta Data ================================

var (
	idName        = i18n.NewFallbackConfig("args.types.id.name", "ID")
	idDescription = i18n.NewFallbackConfig("args.types.id.name", "The ID of something.")
)

// ================================ Errors ================================

var (
	idBelowMinLengthErrorArg = i18n.NewFallbackConfig(
		"args.types.id.errors.below_min_length.arg",
		"Argument {{.position}} must be at least {{.min}} characters long.")
	idBelowMinLengthErrorFlag = i18n.NewFallbackConfig(
		"args.types.id.errors.below_min_length.flag",
		"The `-{{.used_name}}`-flag must be at least {{.min}} characters long.")

	idAboveMaxLengthErrorArg = i18n.NewFallbackConfig(
		"args.types.id.errors.above_max_length.arg",
		"Argument {{.position}} may not be longer than {{.max}} characters.")
	idAboveMaxLengthErrorFlag = i18n.NewFallbackConfig(
		"args.types.id.errors.above_max_length.flag",
		"The `-{{.used_name}}`-flag may not be longer than {{.max}} characters.")

	idNotANumberErrorArg = i18n.NewFallbackConfig(
		"args.types.id.errors.not_a_number.arg",
		"Argument {{.position}} must be a number.")
	idNotANumberErrorFlag = i18n.NewFallbackConfig(
		"args.types.id.errors.not_a_number.flag",
		"The `-{{.used_name}}`-flag must be a number.")
)

// =============================================================================
// Choice
// =====================================================================================

// ================================ Meta Data ================================

var (
	choiceName        = i18n.NewFallbackConfig("args.types.choice.name", "Choice")
	choiceDescription = i18n.NewFallbackConfig(
		"args.types.choice.Name",
		"A choice is a list of elements from which you can to pick one. "+
			"Refer to the help of the command to see all possible choices.")
)

// ================================ Error ================================

var (
	choiceInvalidErrorArg = i18n.NewFallbackConfig(
		"args.types.choice.errors.invalid.arg", "`{{.raw}}` is not a valid choice for argument {{.position}}.")
	choiceInvalidErrorFlag = i18n.NewFallbackConfig(
		"args.types.choice.errors.invalid.flag", "`{{.raw}}` is not a valid choice for the `-{{.{{.used_name}}`-flag.")
)

// =============================================================================
// Member
// =====================================================================================

// ================================ Meta Data ================================

var (
	memberName             = i18n.NewFallbackConfig("args.types.member.name", "Member")
	memberDescriptionNoIDs = i18n.NewFallbackConfig("args.types.member.description.no_ids",
		"A member is a mention of a user in a server. For example @Wumpus.")
	memberDescriptionWithIDs = i18n.NewFallbackConfig("args.types.member.description.with_ids",
		"A member is either a mention of a user in a server or their id. For example @Wumpus or 123456789098765432.")
)

// =============================================================================
// Users
// =====================================================================================

// ================================ Errors ================================

var (
	userInvalidMentionNoDigits = i18n.NewFallbackConfig(
		"args.types.user.errors.invalid_mention_no_digits.arg",
		"{{.raw}} is not a valid mention.")

	userInvalidMentionArg = i18n.NewFallbackConfig(
		"args.types.user.errors.invalid_mention.arg",
		"The mention in argument {{.position}} is invalid. Make sure the user is still on the server.")
	userInvalidMentionFlag = i18n.NewFallbackConfig(
		"args.types.user.errors.invalid_mention.flag",
		"The mention for the `-{{.used_name}}`-flag is invalid. Make sure the user is still on the server.")
)
