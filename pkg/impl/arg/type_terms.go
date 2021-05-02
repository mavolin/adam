package arg

import (
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/adam/pkg/i18n"
	emojiutil "github.com/mavolin/adam/pkg/utils/emoji"
)

// =============================================================================
// Common
// =====================================================================================

// ================================ Errors ================================

var (
	regexpNotMatchingErrorArg = i18n.NewFallbackConfig(
		"arg.common.error.regexp_not_matching.arg",
		"Argument {{.position}} must match `{{.regexp}}`.")
	regexpNotMatchingErrorFlag = i18n.NewFallbackConfig(
		"arg.common.error.regexp_not_matching.flag",
		"The `{{.used_name}}`-flag must match `{{.regexp}}`.")
)

// =============================================================================
// Switch
// =====================================================================================

// ================================ Meta Data ================================

var (
	switchName        = i18n.NewFallbackConfig("arg.type.switch.name", "Switch")
	switchDescription = i18n.NewFallbackConfig(
		"arg.type.switch.description",
		"Used to turn on a feature of a command. Only used for flags.")
)

// ================================ Errors ================================

var switchWithContentError = i18n.NewFallbackConfig(
	"arg.type.switch.error.with_content",
	"The `{{.used_name}}`-flag is a switch flag and cannot be used with content.")

type switchWithContentErrorPlaceholders struct {
	Name string
}

// =============================================================================
// Choice
// =====================================================================================

// ================================ Meta Data ================================

var (
	choiceName        = i18n.NewFallbackConfig("arg.type.choice.name", "Choice")
	choiceDescription = i18n.NewFallbackConfig(
		"arg.type.choice.description",
		"A list of elements from which you can to pick one.")
)

// ================================ Error ================================

var choiceInvalidError = i18n.NewFallbackConfig(
	"arg.type.choice.error.invalid", "`{{.raw}}` is not a valid choice.")

// =============================================================================
// Numbers
// =====================================================================================

// ================================ Integer Meta Data ================================

var (
	integerName        = i18n.NewFallbackConfig("arg.type.integer.name", "Integer")
	integerDescription = i18n.NewFallbackConfig("arg.type.integer.description", "A whole number.")
)

// ================================ Decimal Meta Data ================================

var (
	decimalName        = i18n.NewFallbackConfig("arg.type.decimal.name", "Decimal")
	decimalDescription = i18n.NewFallbackConfig("arg.type.decimal.description", "A decimal number.")
)

// ================================ Integer Errors ================================

var integerSyntaxError = i18n.NewFallbackConfig("arg.type.integer.error.syntax", "`{{.raw}}` is not an integer.")

// ================================ Decimal Errors ================================

var decimalSyntaxError = i18n.NewFallbackConfig("arg.type.decimal.error.syntax", "`{{.raw}}` is not a decimal.")

// ================================ Shared Errors ================================

var (
	numberBelowRangeError = i18n.NewFallbackConfig(
		"arg.type.number.error.under_range",
		"`{{.raw}}` is too small, try using a larger number.")
	numberOverRangeError = i18n.NewFallbackConfig(
		"arg.type.number.error.over_range",
		"`{{.raw}}` is too large, try using a smaller number.")

	numberBelowMinErrorArg = i18n.NewFallbackConfig(
		"arg.type.number.error.below_min.arg",
		"Argument {{.position}} may be no smaller than {{.min}}.")
	numberBelowMinErrorFlag = i18n.NewFallbackConfig(
		"arg.type.number.error.below_min.flag",
		"The `{{.used_name}}`-flag may be smaller than {{.min}}.")

	numberAboveMaxErrorArg = i18n.NewFallbackConfig(
		"arg.type.number.error.below_min.arg",
		"Argument {{.position}} may be no larger than {{.max}}.")
	numberAboveMaxErrorFlag = i18n.NewFallbackConfig(
		"arg.type.number.error.below_min.flag",
		"The `{{.used_name}}`-flag may be no larger than {{.max}}.")
)

// =============================================================================
// Duration
// =====================================================================================

// ================================ Meta Data ================================

var (
	durationName        = i18n.NewFallbackConfig("arg.type.duration.name", "Duration")
	durationDescription = i18n.NewFallbackConfig(
		"arg.type.duration.description", "A timespan. For example: `1h 3 min 4s`.\n"+
			"Available units are `ms` for milliseconds, `s` for seconds, `min` for minutes, `h` for hours, "+
			"`d` for days, `w` for weeks, `m` for months (30 days), and `y` for years.")
)

// ================================ Error ================================

var (
	durationInvalidError = i18n.NewFallbackConfig(
		"arg.type.duration.error.invalid", "`{{.raw}}` is not a valid duration.")

	durationSizeErrorArg = i18n.NewFallbackConfig(
		"arg.type.duration.error.size.arg", "The duration in argument {{.position}} is too large.")
	durationSizeErrorFlag = i18n.NewFallbackConfig(
		"arg.type.duration.error.size.flag", "The duration in the `{{.used_name}}`-flag is too large.")

	durationMissingUnitErrorArg = i18n.NewFallbackConfig(
		"arg.type.duration.error.missing_unit.arg", "The duration in argument {{.position}} is missing a unit.")
	durationMissingUnitErrorFlag = i18n.NewFallbackConfig(
		"arg.type.duration.error.missing_unit.flag",
		"The duration in the `{{.used_name}}`-flag is missing a unit.")

	durationInvalidUnitError = i18n.NewFallbackConfig(
		"arg.type.duration.error.invalid_unit",
		"`{{.unit}}` is not a valid unit of time. "+
			"Valid units are `ms` for milliseconds, `s` for seconds, `min` for minutes, `h` for hours, `d` for days, "+
			"`w` for weeks, `m` for months (30 days), and `y` for years.")

	durationBelowMinErrorArg = i18n.NewFallbackConfig(
		"arg.type.duration.error.below_min.arg",
		"Argument {{.position}} may not be smaller than `{{.min}}`.")
	durationBelowMinErrorFlag = i18n.NewFallbackConfig(
		"arg.type.duration.error.below_min.flag",
		"The `{{.used_name}}`-flag may not be smaller than `{{.min}}`.")

	durationAboveMaxErrorArg = i18n.NewFallbackConfig(
		"arg.type.duration.error.above_max.arg",
		"Argument {{.position}} may not be larger than `{{.max}}`.")
	durationAboveMaxErrorFlag = i18n.NewFallbackConfig(
		"arg.type.duration.error.above_max.flag",
		"The `{{.used_name}}`-flag may not be larger than `{{.max}}`.")
)

// =============================================================================
// Time
// =====================================================================================

// ================================ Meta Data ================================

var (
	timeName = i18n.NewFallbackConfig("arg.type.time.name", "Time")

	timeDescriptionOptionalUTC = i18n.NewFallbackConfig(
		"arg.type.time.description.optional_utc",
		"A 24-hour formatted time, e.g. `13:01`. Optionally, you can add the offset from UTC behind, "+
			"e.g. `13:01 -0500` to use Panama's time zone.")
	timeDescriptionMustUTC = i18n.NewFallbackConfig(
		"arg.type.time.description.must_utc",
		"A 24-hour formatted time with UTC offset, e.g. `13:01 -0500` to use Panama's time zone.")
)

// ================================ Errors ================================

var (
	timeInvalidErrorOptionalUTCArg = i18n.NewFallbackConfig(
		"arg.type.time.error.invalid.optional_utc.arg",
		"The time in argument {{.position}} is invalid. Please use a time like `13:01` or `13:01 -0500`.")
	timeInvalidErrorOptionalUTCFlag = i18n.NewFallbackConfig(
		"arg.type.time.error.invalid.optional_utc.flag",
		"The time you used as `{{.used_name}}`-flag is invalid. Please use a time like `13:01` or `13:01 -0500`.")

	timeInvalidErrorMustUTCArg = i18n.NewFallbackConfig(
		"arg.type.time.error.invalid.must_utc.arg",
		"The time in argument {{.position}} is invalid. Please use a time like `13:01 -0500`.")
	timeInvalidErrorMustUTCFlag = i18n.NewFallbackConfig(
		"arg.type.time.error.invalid.must_utc.flag",
		"The time you used as `{{.used_name}}`-flag is invalid. Please use a time like `13:01 -0500`.")

	timeRequireUTCOffsetErrorArg = i18n.NewFallbackConfig(
		"arg.type.time.error.require_utc_offset.arg",
		"You need to add a UTC offset to the time in argument {{.position}}, "+
			"e.g. `13:01 -0500` to use Panama's time zone.")
	timeRequireUTCOffsetErrorFlag = i18n.NewFallbackConfig(
		"arg.type.time.error.require_utc_offset.flag",
		"You need to add a UTC offset to the time used as `{{.used_name}}`-flag, "+
			"e.g. `13:01 -0500` to use Panama's time zone.")

	timeBeforeMinErrorArg = i18n.NewFallbackConfig(
		"arg.type.time.error.before_min.arg", "The time in argument {{.position}} may not be before {{.min}}.")
	timeBeforeMinErrorFlag = i18n.NewFallbackConfig(
		"arg.type.time.error.before_min.flag",
		"The time you used as the `{{.used_name}}`-flag may not be before {{.min}}.")

	timeAfterMaxErrorArg = i18n.NewFallbackConfig(
		"arg.type.time.error.after_max.arg", "The time in argument {{.position}} may not be after {{.max}}.")
	timeAfterMaxErrorFlag = i18n.NewFallbackConfig(
		"arg.type.time.error.after_max.flag",
		"The time you used as the `{{.used_name}}`-flag may not be after {{.max}}.")
)

// =============================================================================
// Date
// =====================================================================================

// ================================ Meta Data ================================

var (
	dateName = i18n.NewFallbackConfig("arg.type.date.name", "Date")

	dateDescriptionOptionalUTC = i18n.NewFallbackConfig(
		"arg.type.date.description.optional_utc",
		"A date, e.g. `2020-10-31`. Optionally, you can add the offset from UTC behind, "+
			"e.g. `2020-10-31 -0600` to use Costa Rica's time zone.")
	dateDescriptionMustUTC = i18n.NewFallbackConfig(
		"arg.type.date.description.must_utc",
		"A date with UTC offset, e.g. `2020-10-31 -0600` to use Costa Rica's time zone.")
)

var (
	dateInvalidErrorNoUTCArg = i18n.NewFallbackConfig(
		"arg.type.date.error.invalid.no_utc.arg",
		"The date in argument {{.position}} is invalid. Please use a date like `2020-10-31`.")
	dateInvalidErrorNoUTCFlag = i18n.NewFallbackConfig(
		"arg.type.date.error.invalid.no_utc.flag",
		"The date you used as `{{.used_name}}`-flag is invalid. Please use a date like `2020-10-31`.")

	dateInvalidErrorOptionalUTCArg = i18n.NewFallbackConfig(
		"arg.type.date.error.invalid.optional_utc.arg",
		"The date in argument {{.position}} is invalid. Please use a date like `2020-10-31` or `2020-10-31 -0600`.")
	dateInvalidErrorOptionalUTCFlag = i18n.NewFallbackConfig(
		"arg.type.date.error.invalid.optional_utc.flag",
		"The date you used as `{{.used_name}}`-flag is invalid. "+
			"Please use a date like `2020-10-31` or `2020-10-31 -0600`.")

	dateInvalidErrorMustUTCArg = i18n.NewFallbackConfig(
		"arg.type.date.error.invalid.must_utc.arg",
		"The date in argument {{.position}} is invalid. Please use a date like `2020-10-31 -0600`.")
	dateInvalidErrorMustUTCFlag = i18n.NewFallbackConfig(
		"arg.type.date.error.invalid.must_utc.flag",
		"The date you used as `{{.used_name}}`-flag is invalid. Please use a date like `2020-10-31 -0600`.")

	dateRequireUTCOffsetErrorArg = i18n.NewFallbackConfig(
		"arg.type.date.error.require_utc_offset.arg",
		"You need to add a UTC offset to the date in argument {{.position}}, "+
			"e.g. `2020-10-31 -0600` to use Costa Rica's time zone.")
	dateRequireUTCOffsetErrorFlag = i18n.NewFallbackConfig(
		"arg.type.date.error.require_utc_offset.flag",
		"You need to add a UTC offset to the date used as `{{.used_name}}`-flag, "+
			"e.g. `2020-10-31 -0600` to use Costa Rica's time zone.")

	dateBeforeMinErrorArg = i18n.NewFallbackConfig(
		"arg.type.date.error.before_min.arg", "The date in argument {{.position}} may not be before {{.min}}.")
	dateBeforeMinErrorFlag = i18n.NewFallbackConfig(
		"arg.type.date.error.before_min.flag",
		"The date you used as `{{.used_name}}`-flag may not be before {{.min}}.")

	dateAfterMaxErrorArg = i18n.NewFallbackConfig(
		"arg.type.date.error.after_max.arg", "The date in argument {{.position}} may not be after {{.max}}.")
	dateAfterMaxErrorFlag = i18n.NewFallbackConfig(
		"arg.type.date.error.after_max.flag",
		"The date you used as `{{.used_name}}`-flag may not be after {{.max}}.")
)

// =============================================================================
// DateTime
// =====================================================================================

// ================================ Meta Data ================================

var (
	dateTimeName = i18n.NewFallbackConfig("arg.type.date_time.name", "Date and Time")

	dateTimeDescriptionOptionalUTC = i18n.NewFallbackConfig(
		"arg.type.date_time.description.optional_utc",
		"A date with time, e.g. `2020-10-31 13:01`. Optionally, you can add the offset from UTC behind, "+
			"e.g. `2020-10-31 13:01 +0200` to use South Africa's time zone.")
	dateTimeDescriptionMustUTC = i18n.NewFallbackConfig(
		"arg.type.date_time.description.must_utc",
		"A date with time, e.g. `2020-10-31 13:01 +0200` to use South Africa's time zone.")
)

var (
	dateTimeInvalidErrorOptionalUTCArg = i18n.NewFallbackConfig(
		"arg.type.date_time.error.invalid.optional_utc.arg",
		"The date/time combination in argument {{.position}} is invalid. "+
			"Please use a date like `2020-10-31 13:01` or `2020-10-31 13:01 +0200`.")
	dateTimeInvalidErrorOptionalUTCFlag = i18n.NewFallbackConfig(
		"arg.type.date_time.error.invalid.optional_utc.flag",
		"The date/time combination you used as `{{.used_name}}`-flag is invalid. "+
			"Please use a date like `2020-10-31 13:01` or `2020-10-31 13:01 +0200`.")

	dateTimeInvalidErrorMustUTCArg = i18n.NewFallbackConfig(
		"arg.type.date_time.error.invalid.must_utc.arg",
		"The date/time combination in argument {{.position}} is invalid. "+
			"Please use a date like `2020-10-31 13:01 +0200`.")
	dateTimeInvalidErrorMustUTCFlag = i18n.NewFallbackConfig(
		"arg.type.date_time.error.invalid.must_utc.flag",
		"The date/time combination you used as `{{.used_name}}`-flag is invalid. "+
			"Please use a date like `2020-10-31 13:01 +0200`.")
)

// =============================================================================
// TimeZone
// =====================================================================================

// ================================ Meta Data ================================

var (
	timeZoneName        = i18n.NewFallbackConfig("arg.type.time_zone.name", "Time Zone")
	timeZoneDescription = i18n.NewFallbackConfig(
		"arg.type.time_zone.description",
		"The name of an IANA time zone, e.g. `America/New_York`.")
)

// ================================ Errors ================================

var timeZoneInvalidError = i18n.NewFallbackConfig(
	"arg.type.time_zone.error.invalid", "`{{.raw}}` is not a valid IANA time zone name.")

// =============================================================================
// Text
// =====================================================================================

// ================================ Meta Data ================================

var (
	textName        = i18n.NewFallbackConfig("arg.type.text.name", "Text")
	textDescription = i18n.NewFallbackConfig("arg.type.text.description", "A text.")
)

// ================================ Errors ================================

var (
	textBelowMinLengthErrorArg = i18n.NewFallbackConfig(
		"arg.type.text.error.below_min_length.arg",
		"Argument {{.position}} must be at least {{.min}} characters long.")
	textBelowMinLengthErrorFlag = i18n.NewFallbackConfig(
		"arg.type.text.error.below_min_length.flag",
		"The `{{.used_name}}`-flag must be at least {{.min}} characters long.")

	textAboveMaxLengthErrorArg = i18n.NewFallbackConfig(
		"arg.type.text.error.above_max_length.arg",
		"Argument {{.position}} may not be longer than {{.max}} characters.")
	textAboveMaxLengthErrorFlag = i18n.NewFallbackConfig(
		"arg.type.text.error.above_max_length.flag",
		"The `{{.used_name}}`-flag may not be longer than {{.max}} characters.")
)

// =============================================================================
// Code
// =====================================================================================

// ================================ Meta Data ================================

var (
	codeName        = i18n.NewFallbackConfig("arg.type.code.name", "Code")
	codeDescription = i18n.NewFallbackConfig("arg.type.code.description", "A code block.")
)

// ================================ Errors ================================

var (
	codeInvalidErrorArg = i18n.NewFallbackConfig(
		"arg.type.code.error.invalid.arg", "Argument {{.position}} is not a valid code block.")
	codeInvalidErrorFlag = i18n.NewFallbackConfig(
		"arg.type.code.error.invalid.flag", "The `{{.used_name}}`-flag doesn't contain a valid code block.")
)

// =============================================================================
// Link
// =====================================================================================

// ================================ Meta Data ================================

var (
	linkName        = i18n.NewFallbackConfig("arg.type.link.name", "Link")
	linkDescription = i18n.NewFallbackConfig("arg.type.link.description", "A link to something.")
)

// ================================ Errors ================================

var (
	linkInvalidErrorArg = i18n.NewFallbackConfig(
		"arg.type.link.error.invalid.arg", "The link in argument {{.position}} is not valid.")
	linkInvalidErrorFlag = i18n.NewFallbackConfig(
		"arg.type.link.error.invalid.flag", "The link you used as `{{.used_name}}`-flag is not valid.")
)

// =============================================================================
// ID
// =====================================================================================

// ================================ Meta Data ================================

var (
	idName        = i18n.NewFallbackConfig("arg.type.id.name", "ID")
	idDescription = i18n.NewFallbackConfig("arg.type.id.description", "An id.")
)

// ================================ Errors ================================

var (
	idBelowMinLengthErrorArg = i18n.NewFallbackConfig(
		"arg.type.id.error.below_min_length.arg",
		"Argument {{.position}} must be at least {{.min}} characters long.")
	idBelowMinLengthErrorFlag = i18n.NewFallbackConfig(
		"arg.type.id.error.below_min_length.flag",
		"The `{{.used_name}}`-flag must be at least {{.min}} characters long.")

	idAboveMaxLengthErrorArg = i18n.NewFallbackConfig(
		"arg.type.id.error.above_max_length.arg",
		"Argument {{.position}} may not be longer than {{.max}} characters.")
	idAboveMaxLengthErrorFlag = i18n.NewFallbackConfig(
		"arg.type.id.error.above_max_length.flag",
		"The `{{.used_name}}`-flag may not be longer than {{.max}} characters.")

	idInvalidErrorArg = i18n.NewFallbackConfig(
		"arg.type.id.error.invalid.arg",
		"Argument {{.position}} is not a valid id.")
	idInvalidErrorFlag = i18n.NewFallbackConfig(
		"arg.type.id.error.invalid.flag",
		"The `{{.used_name}}`-flag is not a valid id.")
)

// =============================================================================
// Emoji
// =====================================================================================

// ================================ Meta Data ================================

var (
	emojiName = i18n.NewFallbackConfig("arg.type.emoji.name", "Emoji")

	emojiDescriptionNoID   = i18n.NewFallbackConfig("arg.type.emoji.description.no_id", "An emoji. "+emojiutil.Ghost)
	emojiDescriptionWithID = i18n.NewFallbackConfig(
		"arg.type.emoji.description.with_id", "An emoji or the id of an emoji. "+emojiutil.Ghost)
)

// ================================ Errors ================================

var (
	emojiInvalidError = i18n.NewFallbackConfig("arg.type.emoji.error.invalid", "`{{.raw}}` is not an emoji.")

	emojiCustomEmojiInDMError = i18n.NewFallbackConfig(
		"arg.type.emoji.error.custom_emoji_in_dm", "You can't use custom emojis in DMs.")

	emojiCustomEmojiErrorArg = i18n.NewFallbackConfig(
		"arg.type.emoji.error.custom_emoji.arg",
		"You can't use a custom emoji as argument {{.position}}.")
	emojiCustomEmojiErrorFlag = i18n.NewFallbackConfig(
		"arg.type.emoji.error.custom_emoji.flag",
		"You can't use a custom emoji as `{{.used_name}}`-flag.")

	emojiNoAccessError = i18n.NewFallbackConfig(
		"arg.type.emoji.error.no_access",
		"`{{.raw}}` is either not an emoji or I'm unable to access it. "+
			"Make sure to only use emojis from this server.")

	emojiIDNoAccessError = i18n.NewFallbackConfig(
		"arg.type.emoji.error.id_no_access",
		"`{{.raw}}` is either not a valid emoji id or I'm unable to access the emoji it belongs to. "+
			"Make sure to only use emojis from this server.")
)

// =============================================================================
// Member
// =====================================================================================

// ================================ Meta Data ================================

var (
	memberName = i18n.NewFallbackConfig("arg.type.member.name", "Member")

	memberDescriptionNoIDs = i18n.NewFallbackConfig(
		"arg.type.member.description.no_id", "A mention of a user in a server. For example @Wumpus.")
	memberDescriptionWithIDs = i18n.NewFallbackConfig(
		"arg.type.member.description.with_id",
		"A user mention or their id. For example @Wumpus or 123456789098765432.")
)

// =============================================================================
// User
// =====================================================================================

// ================================ Meta Data ================================

var (
	userName        = i18n.NewFallbackConfig("arg.type.user.name", "User")
	userDescription = i18n.NewFallbackConfig(
		"arg.type.user.description",
		"A user mention or their id. The command doesn't need to be invoked on the server the user is on. "+
			"For example: @Wumpus or 123456789098765432")
)

// ================================ Errors ================================

var (
	userInvalidError = i18n.NewFallbackConfig("arg.type.user.error.invalid", "`{{.raw}}` is not a user.")

	userIDInvalidError = i18n.NewFallbackConfig(
		"arg.type.user.error.id_invalid", "`{{.raw}}` is not a valid user id.")

	userInvalidMentionWithRawError = i18n.NewFallbackConfig(
		"arg.type.user.error.invalid_mention_with_raw", "{{.raw}} is not a valid user mention.")

	userInvalidMentionErrorArg = i18n.NewFallbackConfig(
		"arg.type.user.error.invalid_mention.arg",
		"The mention in argument {{.position}} is invalid. Make sure the user is still on the server.")
	userInvalidMentionErrorFlag = i18n.NewFallbackConfig(
		"arg.type.user.error.invalid_mention.flag",
		"The mention in the `{{.used_name}}`-flag is invalid. Make sure the user is still on the server.")
)

// =============================================================================
// Role
// =====================================================================================

// ================================ Meta Data ================================

var (
	roleName = i18n.NewFallbackConfig("arg.type.role.name", "Role")

	roleDescriptionNoID = i18n.NewFallbackConfig(
		"arg.type.role.description.no_id", "A role mention. For example @WumpusGang.")
	roleDescriptionWithID = i18n.NewFallbackConfig(
		"arg.type.role.description.with_id",
		"A role mention or an id of a role. For example @WumpusGang or 123456789098765432.")
)

// ================================ Errors ================================

var (
	roleInvalidError = i18n.NewFallbackConfig("arg.type.role.error.invalid", "`{{.raw}}` is not a role.")

	roleIDInvalidError = i18n.NewFallbackConfig(
		"arg.type.role.error.id_invalid", "`{{.raw}}` is not a valid role id.")

	roleInvalidMentionWithRawError = i18n.NewFallbackConfig(
		"arg.type.role.error.invalid_mention_with_raw", "{{.raw}} is not a valid role mention.")

	roleInvalidMentionErrorArg = i18n.NewFallbackConfig(
		"arg.type.role.error.invalid_mention.arg",
		"The role mention in argument {{.position}} is invalid. Make sure the still role exists.")
	roleInvalidMentionErrorFlag = i18n.NewFallbackConfig(
		"arg.type.role.error.invalid_mention.flag",
		"The role mention you used as `{{.used_name}}`-flag is invalid. Make sure the still role exists.")
)

// =============================================================================
// Channels
// =====================================================================================

// ================================ Errors ================================

var channelIDInvalidError = i18n.NewFallbackConfig(
	"arg.type.channel.error.id_invalid",
	"`{{.raw}}` is not a valid channel id.")

// =============================================================================
// TextChannel
// =====================================================================================

// ================================ Meta Data ================================

var (
	textChannelName = i18n.NewFallbackConfig("arg.type.text_channel.name", "Text Channel")

	textChannelDescriptionNoID = i18n.NewFallbackConfig(
		"arg.type.text_channel.description.no_id",
		"A mention of a text or announcement channel.")
	textChannelDescriptionWithID = i18n.NewFallbackConfig(
		"arg.type.text_channel.description.with_id",
		"A mention of a text or a announcement channel or an id of such.")
)

// ================================ Errors ================================

var (
	textChannelInvalidError = i18n.NewFallbackConfig(
		"arg.type.text_channel.error.invalid", "`{{.raw}}` is not a valid id channel.")

	textChannelInvalidMentionWithRawError = i18n.NewFallbackConfig(
		"arg.type.text_channel.error.invalid_mention_with_raw",
		"`{{.raw}}` is not a valid mention or id of a channel.")

	textChannelInvalidMentionErrorArg = i18n.NewFallbackConfig(
		"arg.type.text_channel.error.invalid_mention.arg",
		"The mention in argument {{.position}} does not belong to channel on this server.")
	textChannelInvalidMentionErrorFlag = i18n.NewFallbackConfig(
		"arg.type.text_channel.error.invalid_mention.flag",
		"The mention you used as the `{{.used_name}}`-flag does not belong to channel on this server.")

	textChannelIDGuildNotMatchingError = i18n.NewFallbackConfig(
		"arg.type.text_channel.error.id_guild_not_matching",
		"The id `{{.raw}}` belongs to a channel from another server.")

	textChannelIDInvalidTypeError = i18n.NewFallbackConfig(
		"arg.type.text_channel.error.id_invalid_type",
		"The id `{{.raw}}` doesn't belong to a text channel.")
)

// =============================================================================
// Category
// =====================================================================================

// ================================ Meta Data ================================

var (
	categoryName        = i18n.NewFallbackConfig("arg.type.category.name", "Category")
	categoryDescription = i18n.NewFallbackConfig(
		"arg.type.category.description",
		"The name of a category or its id.")
)

// ================================ Chooser ================================

var (
	categoryChooserTitle = i18n.NewFallbackConfig("arg.type.category.chooser.title", "Multiple Matches")

	categoryChooserDescription = i18n.NewFallbackConfig(
		"arg.type.category.chooser.description",
		"There are multiple categories that match the name you gave me. "+
			"Please choose the correct one by reacting with the corresponding emoji, "+
			"or react with {{.cancel_emoji}} to cancel.")

	categoryChooserMatch = i18n.NewFallbackConfig(
		"arg.type.category.chooser.match",
		"{{.emoji}} **{{.category_name}}** (position: {{.position}})")

	categoryChooserFullMatchesName = i18n.NewFallbackConfig(
		"arg.type.category.chooser.full_matches.name",
		"Full Matches")

	categoryChooserPartialMatchesName = i18n.NewFallbackConfig(
		"arg.type.category.chooser.partial_matches.name",
		"Partial Matches")

	categoryChooserTooManyPartialMatches = i18n.NewFallbackConfig(
		"arg.type.category.chooser.too_many_partial_matches",
		"There are {{.num_partial_matches}} additional partial matches. "+
			"Use the full name of the category or their id, to match any of these.")
)

type (
	categoryChooserDescriptionPlaceholders struct {
		CancelEmoji discord.APIEmoji
	}

	categoryChooserMatchPlaceholders struct {
		Emoji        discord.APIEmoji
		CategoryName string
		Position     int
	}

	categoryChooserTooManyPartialMatchesPlaceholders struct {
		NumPartialMatches int
	}
)

// ================================ Errors ================================

var (
	categoryNotFoundError = i18n.NewFallbackConfig(
		"arg.type.category.error.not_found",
		"I couldn't find a category with the name or id `{{.raw}}`. Make sure you spelled it correctly.")

	categoryIDInvalidErrorArg = i18n.NewFallbackConfig(
		"arg.type.category.error.id_invalid.arg",
		"Argument {{.position}} is not a valid category id.")
	categoryIDInvalidErrorFlag = i18n.NewFallbackConfig(
		"arg.type.category.error.id_invalid.flag",
		"The `{{.used_name}}`-flag doesn't contain a valid category id.")

	categoryIDInvalidTypeError = i18n.NewFallbackConfig(
		"arg.type.category.error.id_invalid_type",
		"The id `{{.raw}}` doesn't belong to a category.")

	categoryTooManyMatchesError = i18n.NewFallbackConfig(
		"arg.type.category.error.too_many_full_matches",
		"There are too many categories that match `{{.raw}}`. "+
			"You can either (temporarily) rename the category and try again, or use the id of category instead.")

	categoryTooManyPartialMatchesError = i18n.NewFallbackConfig(
		"arg.type.category.error.too_many_partial_matches",
		"There are too many categories that match `{{.raw}}`. "+
			"You can either try to find the category by using their full name, "+
			"or you can use the id of category instead.")
)

// =============================================================================
// VoiceChannel
// =====================================================================================

// ================================ Meta Data ================================

var (
	voiceChannelName        = i18n.NewFallbackConfig("arg.type.voice_channel.name", "Voice Channel")
	voiceChannelDescription = i18n.NewFallbackConfig(
		"arg.type.voice_channel.description",
		"The name of a voice channel or its id.")
)

// ================================ Chooser ================================

var (
	voiceChannelChooserTitle = i18n.NewFallbackConfig("arg.type.voice_channel.chooser.title", "Multiple Matches")

	voiceChannelChooserDescription = i18n.NewFallbackConfig(
		"arg.type.voice_channel.chooser.description",
		"There are multiple voice channels that match the name you gave me. "+
			"Please choose the correct one by reacting with the corresponding emoji, "+
			"or react with {{.cancel_emoji}} to cancel.")

	voiceChannelChooserRootMatch = i18n.NewFallbackConfig(
		"arg.type.category.chooser.match.root",
		"{{.emoji}} **{{.channel_name}}** (position: {{.position}})")

	voiceChannelChooserNestedMatch = i18n.NewFallbackConfig(
		"arg.type.voice_channel.chooser.match.nested",
		"{{.emoji}} **{{.channel_name}}** ({{.category_name}}, position: {{.position}})")

	voiceChannelChooserFullMatchesName = i18n.NewFallbackConfig(
		"arg.type.voice_channel.chooser.full_matches.name",
		"Full Matches")

	voiceChannelChooserPartialMatchesName = i18n.NewFallbackConfig(
		"arg.type.voice_channel.chooser.partial_matches.name",
		"Partial Matches")

	voiceChannelChooserTooManyPartialMatches = i18n.NewFallbackConfig(
		"arg.type.voice_channel.chooser.too_many_partial_matches",
		"There are {{.num_partial_matches}} additional partial matches. "+
			"Use the full name of the voice channel or their id, to match any of these.")
)

type (
	voiceChannelChooserDescriptionPlaceholders struct {
		CancelEmoji discord.APIEmoji
	}

	// voiceChannelChooserMatchPlaceholders is the placeholder struct used for
	// both voiceChannelChooserRootMatch and voiceChannelChooserNestedMatch.
	voiceChannelChooserMatchPlaceholders struct {
		Emoji        discord.APIEmoji
		CategoryName string
		ChannelName  string
		Position     int
	}

	voiceChannelChooserTooManyPartialMatchesPlaceholders struct {
		NumPartialMatches int
	}
)

// ================================ Errors ================================

var (
	voiceChannelNotFoundError = i18n.NewFallbackConfig(
		"arg.type.voice_channel.error.not_found",
		"I couldn't find a voice channel with the name or id `{{.raw}}`. Make sure you spelled it correctly.")

	voiceChannelIDInvalidErrorArg = i18n.NewFallbackConfig(
		"arg.type.voice_channel.error.id_invalid.arg",
		"Argument {{.position}} is not a valid voice channel id.")
	voiceChannelIDInvalidErrorFlag = i18n.NewFallbackConfig(
		"arg.type.voice_channel.error.id_invalid.flag",
		"The `{{.used_name}}`-flag doesn't contain a valid voice channel id.")

	voiceChannelIDInvalidTypeError = i18n.NewFallbackConfig(
		"arg.type.voice_channel.error.id_invalid_type",
		"The id `{{.raw}}` doesn't belong to a voice channel.")

	voiceChannelTooManyMatchesError = i18n.NewFallbackConfig(
		"arg.type.voice_channel.error.too_many_full_matches",
		"There are too many voice channels that match `{{.raw}}`. "+
			"You can either (temporarily) rename the voice channel and try again,"+
			" or use the id of the voice channel instead.")

	voiceChannelTooManyPartialMatchesError = i18n.NewFallbackConfig(
		"arg.type.voice_channel.error.too_many_partial_matches",
		"There are too many voice channels that match `{{.raw}}`. "+
			"You can either try to find the voice channel by using their full name, "+
			"or you can use the id of voice channel instead.")
)

// =============================================================================
// Command
// =====================================================================================

// ================================ Meta Data ================================

var (
	commandName        = i18n.NewFallbackConfig("arg.type.command.name", "Command")
	commandDescription = i18n.NewFallbackConfig(
		"arg.type.command.description", "The name of a command without the command's prefix.")
)

// ================================ Errors ================================

var (
	commandNotFoundError = i18n.NewFallbackConfig(
		"arg.type.command.error.not_found",
		"I don't know any commands by the name of `{{.raw}}`. Make sure you spelled it right.")

	commandNotFoundErrorProvidersUnavailable = i18n.NewFallbackConfig(
		"arg.type.command.error.not_found.providers_unavailable",
		"I couldn't find any commands by the name of `{{.raw}}`, "+
			"but I don't have access to some commands right now. Try again later or check your spelling.")
)

// =============================================================================
// Module
// =====================================================================================

// ================================ Meta Data ================================

var (
	moduleName        = i18n.NewFallbackConfig("arg.type.module.name", "Module")
	moduleDescription = i18n.NewFallbackConfig(
		"arg.type.module.description",
		"The name of a module, without the bot's prefix.")
)

// ================================ Errors ================================

var (
	moduleNotFoundError = i18n.NewFallbackConfig(
		"arg.type.module.error.not_found",
		"I don't know any modules by the name of `{{.raw}}`. Make sure you spelled it right.")

	moduleNotFoundErrorProvidersUnavailable = i18n.NewFallbackConfig(
		"arg.type.module.error.not_found.providers_unavailable",
		"I couldn't find any modules by the name of `{{.raw}}`, "+
			"but I don't have access to some modules right now. Try again later or check your spelling.")
)

// =============================================================================
// Plugin
// =====================================================================================

// ================================ Meta Data ================================

var (
	pluginName        = i18n.NewFallbackConfig("arg.type.plugin.name", "Command or Module")
	pluginDescription = i18n.NewFallbackConfig(
		"arg.type.plugin.description",
		"The name of a command or module, without the bot's prefix.")
)

// ================================ Errors ================================

var (
	pluginNotFoundError = i18n.NewFallbackConfig(
		"arg.type.plugin.error.not_found",
		"I don't know any commands or modules with the name `{{.raw}}`. Make sure you spelled it right.")

	pluginNotFoundErrorProvidersUnavailable = i18n.NewFallbackConfig(
		"arg.type.plugin.error.not_found.providers_unavailable",
		"I couldn't find any commands or modules with the name `{{.raw}}`, "+
			"however, I'm having trouble accessing some of my commands, so this may be why. "+
			"Try again later or check your spelling.")
)

// =============================================================================
// RegularExpression
// =====================================================================================

// ================================ Meta Data ================================

var (
	regexpName        = i18n.NewFallbackConfig("arg.type.regular_expression.name", "Regular Expression")
	regexpDescription = i18n.NewFallbackConfig(
		"arg.type.regular_expression.description",
		"A regular expression following the RE2/Go-flavor.")
)

// ================================ Errors ================================

var (
	regexpInvalidErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid.arg",
		"The regular expression in argument {{.position}} is invalid.")
	regexpInvalidErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid.flag",
		"The regular expression you used in the `{{.used_name}}`-flag is invalid.")

	regexpInvalidCharClassErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_character_class.arg",
		"The regular expression in argument {{.position}} uses an invalid character class:\n```\n{{.expression}}```")
	regexpInvalidCharClassErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_character_class.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses an invalid character class:\n"+
			"```\n{{.expression}}```")

	regexpInvalidCharRangeErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_character_range.arg",
		"The regular expression in argument {{.position}} uses an invalid character class range:\n"+
			"```\n{{.expression}}```")
	regexpInvalidCharRangeErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_character_range.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses an invalid character class range:\n"+
			"```\n{{.expression}}```")

	regexpInvalidEscapeErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_escape.arg",
		"The regular expression in argument {{.position}} uses an invalid escape sequence:\n```\n{{.expression}}```")
	regexpInvalidEscapeErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_escape.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses an invalid escape sequence:\n"+
			"```\n{{.expression}}```")

	regexpInvalidNamedCaptureErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_named_capture.arg",
		"The regular expression in argument {{.position}} uses an invalid named capture:\n```\n{{.expression}}```")
	regexpInvalidNamedCaptureErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_named_capture.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses an invalid named capture:\n"+
			"```\n{{.expression}}```")

	regexpInvalidPerlOpErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_perl_operation.arg",
		"The regular expression in argument {{.position}} uses invalid or unsupported Perl syntax:\n"+
			"```\n{{.expression}}```")
	regexpInvalidPerlOpErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_perl_operation.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses invalid or unsupported Perl syntax:\n"+
			"```\n{{.expression}}```")

	regexpInvalidRepeatOpErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_repeat_operation.arg",
		"The regular expression in argument {{.position}} has two consecutive `+`.")
	regexpInvalidRepeatOpErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_repeat_operation.flag",
		"The regular expression you used in the `{{.used_name}}`-flag has two consecutive `+`.")

	regexpInvalidRepeatSizeErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_repeat_size.arg",
		"The regular expression in argument {{.position}} uses an invalid invalid repeat count:\n"+
			"```\n{{.expression}}```")
	regexpInvalidRepeatSizeErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_repeat_size.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses an invalid invalid repeat count:\n"+
			"```\n{{.expression}}```")

	regexpInvalidUTF8ErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_utf8.arg",
		"The regular expression in argument {{.position}} uses invalid UTF-8:\n```\n{{.expression}}```")
	regexpInvalidUTF8ErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.invalid_utf8.flag",
		"The regular expression you used in the `{{.used_name}}`-flag uses invalid UTF-8:\n```\n{{.expression}}```")

	regexpMissingBracketErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.missing_bracket.arg",
		"The regular expression in argument {{.position}} is missing a closing `]`:\n```\n{{.expression}}```")
	regexpMissingBracketErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.missing_bracket.flag",
		"The regular expression you used in the `{{.used_name}}`-flag is missing a closing `]`:\n"+
			"```\n{{.expression}}```")

	regexpMissingParenErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.missing_parentheses.arg",
		"The regular expression in argument {{.position}} is missing a closing `)`:\n```\n{{.expression}}```")
	regexpMissingParenErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.missing_parentheses.flag",
		"The regular expression you used in the `{{.used_name}}`-flag is missing a closing `)`:\n"+
			"```\n{{.expression}}```")

	regexpMissingRepeatArgErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.missing_repeat_argument.arg",
		"The regular expression in argument {{.position}} is missing an argument to the repetition operator:\n"+
			"```\n{{.expression}}```")
	regexpMissingRepeatArgErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.missing_repeat_argument.flag",
		"The regular expression you used in the `{{.used_name}}`-flag is missing an argument to the repetition "+
			"operator:\n```\n{{.expression}}```")

	regexpTrailingBackslashErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.trailing_backslash.arg",
		"The regular expression in argument {{.position}} has a trailing backlash at the end of the expression.")
	regexpTrailingBackslashErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.trailing_backslash.flag",
		"The regular expression you used in the `{{.used_name}}`-flag has a trailing backlash at the end of the "+
			"expression.")

	regexpUnexpectedParenErrorArg = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.unexpected_parentheses.arg",
		"The regular expression in argument {{.position}} has an unexpected `)`:\n```\n{{.expression}}```")
	regexpUnexpectedParenErrorFlag = i18n.NewFallbackConfig(
		"arg.type.regular_expression.error.unexpected_parentheses.flag",
		"The regular expression you used in the `{{.used_name}}`-flag has an unexpected `)`:\n"+
			"```\n{{.expression}}```")
)
