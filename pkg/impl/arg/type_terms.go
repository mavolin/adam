package arg

import (
	"github.com/mavolin/adam/pkg/i18n"
	emojiutil "github.com/mavolin/adam/pkg/utils/emoji"
)

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
		"The `{{.used_name}}`-flag must match `{{.regexp}}`.")
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
	"args.types.switch.errors.with_content", "`{{.name}}` is a Switch flag and cannot be used with content.")

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
		"{{.raw}} is too small, try giving the `{{.used_name}}`-flag a larger number.")

	numberOverRangeErrorArg = i18n.NewFallbackConfig(
		"args.types.number.errors.over_range.arg",
		"{{.raw}} is too large, try using a smaller number as argument {{.position}}.")
	numberOverRangeErrorFlag = i18n.NewFallbackConfig(
		"args.types.integer.errors.over_range.flag",
		"{{.raw}} is a too large, try giving the `{{.used_name}}`-flag a smaller number.")

	numberBelowMinErrorArg = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.arg",
		"Argument {{.position}} must be larger or equal to {{.min}}.")
	numberBelowMinErrorFlag = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.flag",
		"The `{{.used_name}}`-flag must be larger or equal to {{.min}}.")

	numberAboveMaxErrorArg = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.arg",
		"Argument {{.position}} must be smaller or equal to {{.max}}.")
	numberAboveMaxErrorFlag = i18n.NewFallbackConfig(
		"args.types.number.errors.below_min.flag",
		"The `{{.used_name}}`-flag must be smaller or equal to {{.max}}.")
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
		"The text used in the `{{.used_name}}`-flag must be at least {{.min}} characters long.")

	textAboveMaxLengthErrorArg = i18n.NewFallbackConfig(
		"args.types.text.errors.above_max_length.arg",
		"The text in argument {{.position}} may not be longer than {{.max}} characters.")
	textAboveMaxLengthErrorFlag = i18n.NewFallbackConfig(
		"args.types.text.errors.above_max_length.flag",
		"The text used in the `{{.used_name}}`-flag may not be longer than {{.max}} characters.")
)

// =============================================================================
// Link
// =====================================================================================

// ================================ Meta Data ================================

var (
	linkName        = i18n.NewFallbackConfig("args.types.link.name", "Link")
	linkDescription = i18n.NewFallbackConfig("args.types.link.description", "A link to something on the web.")
)

// ================================ Errors ================================

var (
	linkInvalidErrorArg = i18n.NewFallbackConfig(
		"args.types.link.errors.invalid.arg", "Argument {{.position}} must be a valid link.")
	linkInvalidErrorFlag = i18n.NewFallbackConfig(
		"args.types.link.errors.invalid.flag", "The `{{.used.name}}`-flag must be a valid link.")
)

// =============================================================================
// ID
// =====================================================================================

// ================================ Meta Data ================================

var (
	idName        = i18n.NewFallbackConfig("args.types.id.name", "ID")
	idDescription = i18n.NewFallbackConfig("args.types.id.name", "The id of something.")
)

// ================================ Errors ================================

var (
	idBelowMinLengthErrorArg = i18n.NewFallbackConfig(
		"args.types.id.errors.below_min_length.arg",
		"Argument {{.position}} must be at least {{.min}} characters long.")
	idBelowMinLengthErrorFlag = i18n.NewFallbackConfig(
		"args.types.id.errors.below_min_length.flag",
		"The `{{.used_name}}`-flag must be at least {{.min}} characters long.")

	idAboveMaxLengthErrorArg = i18n.NewFallbackConfig(
		"args.types.id.errors.above_max_length.arg",
		"Argument {{.position}} may not be longer than {{.max}} characters.")
	idAboveMaxLengthErrorFlag = i18n.NewFallbackConfig(
		"args.types.id.errors.above_max_length.flag",
		"The `{{.used_name}}`-flag may not be longer than {{.max}} characters.")

	idNotANumberErrorArg = i18n.NewFallbackConfig(
		"args.types.id.errors.not_a_number.arg",
		"Argument {{.position}} must be a number.")
	idNotANumberErrorFlag = i18n.NewFallbackConfig(
		"args.types.id.errors.not_a_number.flag",
		"The `{{.used_name}}`-flag must be a number.")
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
		"args.types.choice.errors.invalid.flag", "`{{.raw}}` is not a valid choice for the `{{.used_name}}`-flag.")
)

// =============================================================================
// Emoji
// =====================================================================================

// ================================ Meta Data ================================

var (
	emojiName        = i18n.NewFallbackConfig("args.types.emoji.name", "Emoji")
	emojiDescription = i18n.NewFallbackConfig("args.types.emoji.description", "An emoji. "+emojiutil.Ghost)
)

// ================================ Errors ================================

var (
	emojiCustomEmojiInDMError = i18n.NewFallbackConfig(
		"args.types.emoji.errors.custom_emoji_in_dm", "You can't use custom emojis in DMs.")

	emojiOnlyUnicodeErrorArg = i18n.NewFallbackConfig(
		"args.types.emoji.errors.only_unicode.arg",
		emojiutil.Prohibited+" You can only use default emojis as argument {{.position}}.")
	emojiOnlyUnicodeErrorFlag = i18n.NewFallbackConfig(
		"args.types.emoji.errors.only_unicode.flag",
		emojiutil.Prohibited+" You can only use default emojis as `{{.used_name}}`-flag.")

	emojiInvalidError = i18n.NewFallbackConfig(
		"args.types.emoji.errors.invalid",
		emojiutil.DesktopComputer+emojiutil.Collision+" {{.raw}} is not an emoji.")

	emojiNoAccessError = i18n.NewFallbackConfig(
		"args.types.user.errors.no_access",
		"{{.raw}} is not a valid emoji or I'm unable to access it. "+
			"Make sure to only use emojis from this server.")
)

// =============================================================================
// EmojiID
// =====================================================================================

// ================================ Errors ================================

var (
	emojiIDNoAccessErrorArg = i18n.NewFallbackConfig(
		"args.types.emoji_id.errors.no_access.arg",
		"Argument {{.position}} is not a valid emoji id or I'm unable to access the emoji it belongs to. "+
			"Make sure to only use emojis from this server.")
	emojiIDNoAccessErrorFlag = i18n.NewFallbackConfig(
		"args.types.emoji_id.errors.no_access.flag",
		"The `{{.used_name}}`-flag contains no valid emoji id or I'm unable to access the emoji it belongs to. "+
			"Make sure to only use emojis from this server.")
)

// =============================================================================
// Member
// =====================================================================================

// ================================ Meta Data ================================

var (
	memberName             = i18n.NewFallbackConfig("args.types.member.name", "Member")
	memberDescriptionNoIDs = i18n.NewFallbackConfig(
		"args.types.member.description.no_ids", "A member is a mention of a user in a server. For example @Wumpus.")
	memberDescriptionWithIDs = i18n.NewFallbackConfig(
		"args.types.member.description.with_ids",
		"A member is either a mention of a user or their id. For example @Wumpus or 123456789098765432.")
)

// =============================================================================
// MemberID
// =====================================================================================

// ================================ Meta Data ================================

var (
	memberIDName        = i18n.NewFallbackConfig("args.types.member_id.name", "Member ID")
	memberIDDescription = i18n.NewFallbackConfig(
		"args.types.member_id.description", "The id of a server member. For example 123456789098765432.")
)

// =============================================================================
// User
// =====================================================================================

// ================================ Meta Data ================================

var (
	userName        = i18n.NewFallbackConfig("args.types.user.name", "User")
	userDescription = i18n.NewFallbackConfig(
		"args.types.user.description",
		"A user is either a mention of a user or their id. "+
			"The command doesn't need to be invoked on the server the user is on.")
)

// ================================ Errors ================================

var (
	userInvalidMentionWithRaw = i18n.NewFallbackConfig(
		"args.types.user.errors.invalid_mention_with_raw", "{{.raw}} is not a valid mention.")

	userInvalidMentionArg = i18n.NewFallbackConfig(
		"args.types.user.errors.invalid_mention.arg",
		"The mention in argument {{.position}} is invalid. Make sure the user is still on the server.")
	userInvalidMentionFlag = i18n.NewFallbackConfig(
		"args.types.user.errors.invalid_mention.flag",
		"The mention in the `{{.used_name}}`-flag is invalid. Make sure the user is still on the server.")
)

// =============================================================================
// UserID
// =====================================================================================

// ================================ Meta Data ================================

var (
	userIDName        = i18n.NewFallbackConfig("args.types.user_id.name", "User ID")
	userIDDescription = i18n.NewFallbackConfig(
		"args.types.user_id.description", "The id of a user. For example 123456789098765432.")
)

// ================================ Errors ================================

var (
	userInvalidIDWithRaw = i18n.NewFallbackConfig(
		"args.types.user.errors.invalid_id_with_raw", "{{.raw}} is not a valid user id.")

	userInvalidIDArg = i18n.NewFallbackConfig(
		"args.types.user.errors.invalid_id.arg", "The user id in argument {{.position}} is invalid.")
	userInvalidIDFlag = i18n.NewFallbackConfig(
		"args.types.user.errors.invalid_id.flag", "The user id in the `{{.used_name}}`-flag is invalid.")
)

// =============================================================================
// Role
// =====================================================================================

// ================================ Meta Data ================================

var (
	roleName        = i18n.NewFallbackConfig("args.types.role.name", "Role")
	roleDescription = i18n.NewFallbackConfig(
		"args.types.role.description",
		"A role mention or an id of a role. For example @WumpusGang or 123456789098765432.")
)

// ================================ Errors ================================

var (
	roleInvalidMentionWithRaw = i18n.NewFallbackConfig(
		"args.types.role.errors.invalid_mention_with_raw", "{{.raw}} is not a valid role mention.")

	roleInvalidMentionArg = i18n.NewFallbackConfig(
		"args.types.role.errors.invalid_mention.arg",
		"The mention in argument {{.position}} is invalid. Make sure the role exists.")
	roleInvalidMentionFlag = i18n.NewFallbackConfig(
		"args.types.role.errors.invalid_mention.flag",
		"The mention in the `{{.used_name}}`-flag is invalid. Make sure the role exists.")
)

// =============================================================================
// RoleID
// =====================================================================================

// ================================ Meta Data ================================

var (
	roleIDName        = i18n.NewFallbackConfig("args.types.role_id.name", "Role ID")
	roleIDDescription = i18n.NewFallbackConfig(
		"args.types.role_id.description", "The id of a role. For example 123456789098765432")
)

// ================================ Errors ================================

var (
	roleInvalidIDWithRaw = i18n.NewFallbackConfig(
		"args.types.role.errors.invalid_id_with_raw", "{{.raw}} is not a valid role id.")

	roleInvalidIDArg = i18n.NewFallbackConfig(
		"args.types.role.errors.invalid_id.arg", "The role id in argument {{.position}} is invalid.")
	roleInvalidIDFlag = i18n.NewFallbackConfig(
		"args.types.role.errors.invalid_id.flag", "The role id in the `{{.used_name}}`-flag is invalid.")
)

// =============================================================================
// RegularExpression
// =====================================================================================

// ================================ Meta Data ================================

var (
	regexpName        = i18n.NewFallbackConfig("args.types.regular_expression.name", "Regular Expression")
	regexpDescription = i18n.NewFallbackConfig(
		"args.types.regular_expression.description",
		"A regular expression is a regular expression following the RE2/Golang-flavor. "+
			"It can be used to macht text that follows user-defined rules.")
)

// ================================ Errors ================================

var (
	regexpInvalidErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid.arg",
		"The regular expression in argument {{.position}} is invalid.")
	regexpInvalidErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid.flag",
		"The regular expression you used in the `{{.used_name}}`-flag is invalid.")

	regexpInvalidCharClassErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_character_class.arg",
		"The regular expression in argument {{.position}} uses an invalid character class:\n```\n{{.expression}}\n```")
	regexpInvalidCharClassErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_character_class.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses an invalid character class:"+
			"\n```\n{{.expression}}\n```")

	regexpInvalidCharRangeErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_character_range.arg",
		"The regular expression in argument {{.position}} uses an invalid character class range:\n"+
			"```\n{{.expression}}\n```")
	regexpInvalidCharRangeErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_character_range.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses an invalid character class range:\n"+
			"```\n{{.expression}}\n```")

	regexpInvalidEscapeErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_escape.arg",
		"The regular expression in argument {{.position}} uses an invalid escape sequence:\n```\n{{.expression}}\n```")
	regexpInvalidEscapeErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_escape.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses an invalid escape sequence\n"+
			"```\n{{.expression}}\n```")

	regexpInvalidNamedCaptureErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_named_capture.arg",
		"The regular expression in argument {{.position}} uses an invalid named capture:\n```\n{{.expression}}\n```")
	regexpInvalidNamedCaptureErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_named_capture.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses an invalid named capture:\n"+
			"```\n{{.expression}}\n```")

	regexpInvalidPerlOpErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_perl_operation.arg",
		"The regular expression in argument {{.position}} uses invalid or unsupported Perl syntax:\n"+
			"```\n{{.expression}}\n```")
	regexpInvalidPerlOpErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_perl_operation.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses invalid or unsupported Perl syntax:\n"+
			"```\n{{.expression}}\n```")

	regexpInvalidRepeatOpErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_repeat_operation.arg",
		"The regular expression in argument {{.position}} has two consecutive `+`.")
	regexpInvalidRepeatOpErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_repeat_operation.flag",
		"The regular expression you used in the `{{.used_name}}`-flag has two consecutive `+`.")

	regexpInvalidRepeatSizeErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_repeat_size.arg",
		"The regular expression in argument {{.position}} uses an invalid invalid repeat count:\n"+
			"```\n{{.expression}}\n```")
	regexpInvalidRepeatSizeErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_repeat_size.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses an invalid invalid repeat count:\n"+
			"```\n{{.expression}}\n```")

	regexpInvalidUTF8ErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_utf8.arg",
		"The regular expression in argument {{.position}} uses invalid UTF-8:\n```\n{{.expression}}\n```")
	regexpInvalidUTF8ErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.invalid_utf8.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses invalid UTF-8:\n```\n{{.expression}}\n```")

	regexpMissingBracketErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.missing_bracket.arg",
		"The regular expression in argument {{.position}} is missing a closing `]`:\n```\n{{.expression}}\n```")
	regexpMissingBracketErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.missing_bracket.flag",
		"The regular expression you used in the `{{.used_name}}`-flag is missing a closing `]`:\n"+
			"```\n{{.expression}}\n```")

	regexpMissingParenErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.missing_parentheses.arg",
		"The regular expression in argument {{.position}} is missing a closing `)`:\n```\n{{.expression}}\n```")
	regexpMissingParenErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.missing_parentheses.flag",
		"The regular expression you used in the `{{.used_name}}`-flag is missing a closing `)`:\n"+
			"```\n{{.expression}}\n```")

	regexpMissingRepeatArgErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.missing_repeat_argument.arg",
		"The regular expression in argument {{.position}} is missing an argument to the repetition operator:\n"+
			"```\n{{.expression}}\n```")
	regexpMissingRepeatArgErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.missing_repeat_argument.flag",
		"The regular expression you used in the `{{.used_name}}`-flag is missing an argument to the repetition "+
			"operator:\n```\n{{.expression}}\n```")

	regexpTrailingBackslashErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.trailing_backslash.arg",
		"The regular expression in argument {{.position}} has a trailing backlash at the end of the expression.")
	regexpTrailingBackslashErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.trailing_backslash.flag",
		"The regular expression you used in the `{{.used_name}}`-flag has a trailing backlash at the end of the "+
			"expression.")

	regexpUnexpectedParenErrorArg = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.unexpected_parentheses.arg",
		"The regular expression in argument {{.position}} has an unexpected `)`:\n```\n{{.expression}}\n```")
	regexpUnexpectedParenErrorFlag = i18n.NewFallbackConfig(
		"args.types.regular_expression.errors.unexpected_parentheses.flag",
		"The regular expression you used in the `{{.used_name}}`-flag has an unexpected `)`:"+
			"\n```\n{{.expression}}\n```")
)
