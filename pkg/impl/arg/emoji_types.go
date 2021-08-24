package arg

import (
	"regexp"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	emojiutil "github.com/mavolin/adam/pkg/utils/emoji"
)

// EmojiAllowIDs is a global flag that defines whether Emojis may also be noted
// as plain Snowflakes.
var EmojiAllowIDs = false

// =============================================================================
// Emoji
// =====================================================================================

var (
	// Emoji is the Type used for unicode and custom emojis.
	// Due to Discord-API limitations the type currently only supports custom
	// emojis from the invoking guild.
	// Use RawEmoji to use the raw emoji, which is not bound to such
	// limitations.
	//
	// Go type: *discord.Emoji
	Emoji plugin.ArgType = &emoji{customEmojis: true}
	// UnicodeEmoji is the type used for unicode emojis.
	//
	// Go type: *discord.Emoji
	UnicodeEmoji plugin.ArgType = &emoji{customEmojis: false}
)

type emoji struct {
	customEmojis bool
}

func (e emoji) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(emojiName) // we have a fallback
	return name
}

func (e emoji) GetDescription(l *i18n.Localizer) string {
	if EmojiAllowIDs {
		desc, err := l.Localize(emojiDescriptionWithID)
		if err == nil {
			return desc
		}
	}

	desc, _ := l.Localize(emojiDescriptionNoID) // we have a fallback
	return desc
}

var customEmojiRegexp = regexp.MustCompile(`^<a?:(?P<name>[^:]+):(?P<id>\d+)>$`)

func (e emoji) Parse(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	if emojiutil.IsValid(ctx.Raw) {
		return &discord.Emoji{Name: ctx.Raw}, nil
	}

	if matches := customEmojiRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 3 {
		if !e.customEmojis {
			return nil, newArgumentError2(emojiCustomEmojiErrorArg, emojiCustomEmojiErrorFlag, ctx, nil)
		} else if ctx.GuildID == 0 {
			return nil, newArgumentError2(emojiCustomEmojiInDMError, emojiCustomEmojiInDMError, ctx, nil)
		}

		rawID := matches[2]

		id, err := discord.ParseSnowflake(rawID)
		if err != nil { // range err
			return nil, newArgumentError(emojiInvalidError, ctx, nil)
		}

		emoji, err := s.Emoji(ctx.GuildID, discord.EmojiID(id))
		if err != nil {
			return nil, newArgumentError(emojiNoAccessError, ctx, nil)
		}

		return emoji, nil
	}

	if !EmojiAllowIDs {
		return nil, newArgumentError(emojiInvalidError, ctx, nil)
	}

	id, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgumentError(emojiInvalidError, ctx, nil)
	}

	if !e.customEmojis {
		return nil, newArgumentError2(emojiCustomEmojiErrorArg, emojiCustomEmojiErrorFlag, ctx, nil)
	} else if ctx.GuildID == 0 {
		return nil, newArgumentError(emojiCustomEmojiInDMError, ctx, nil)
	}

	emoji, err := s.Emoji(ctx.GuildID, discord.EmojiID(id))
	if err != nil {
		return nil, newArgumentError(emojiIDNoAccessError, ctx, nil)
	}

	return emoji, nil
}

func (e emoji) GetDefault() interface{} {
	return (*discord.Emoji)(nil)
}

// =============================================================================
// RawEmoji
// =====================================================================================

// RawEmoji is the Type for used for emojis that are either default emojis or
// custom ones from any guild.
// This means that an emoji is only guaranteed to be available to the bot, if
// it is unicode.
// If the emoji is custom, it is only guaranteed that it follows the pattern
// of an emoji.
// Unlike Emoji, this type only accepts actual emojis but no ids.
//
// Go type: discord.APIEmoji
var RawEmoji plugin.ArgType = new(rawEmoji)

type rawEmoji struct{}

func (r rawEmoji) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(emojiName) // we have a fallback
	return name
}

func (r rawEmoji) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(emojiDescriptionNoID) // we have a fallback
	return desc
}

func (r rawEmoji) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	if emojiutil.IsValid(ctx.Raw) {
		return discord.APIEmoji(ctx.Raw), nil
	} else if matches := customEmojiRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 3 {
		return discord.APIEmoji(matches[1] + ":" + matches[2]), nil
	}

	return nil, newArgumentError(emojiInvalidError, ctx, nil)
}

func (r rawEmoji) GetDefault() interface{} {
	return discord.APIEmoji("")
}
