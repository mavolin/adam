package arg

import (
	"regexp"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	emojiutil "github.com/mavolin/adam/pkg/utils/emoji"
)

// EmojiAllowIDs is a global flag that allows you to specify whether Emojis
// may be noted as Snowflakes.
//
// Defaults to false.
var EmojiAllowIDs = false

// =============================================================================
// Emoji
// =====================================================================================

var (
	// Emoji is the Type used for unicode and custom emojis.
	// Due to Discord-API limitations the type currently only supports custom
	// emojis from the invoking guild.
	// However, if such features become available, this will be modified.
	// Use RawEmoji to use the raw emoji, which is not bound to such
	// limitations.
	// Alternatively, you could obtain the emoji through a reaction.
	//
	// Go type: *discord.Emoji
	Emoji = &emoji{customEmojis: true}
	// UnicodeEmoji is the type used for unicode emojis.
	//
	// Go type: *discord.Emoji
	UnicodeEmoji = &emoji{customEmojis: false}
)

type emoji struct {
	customEmojis bool
}

func (e emoji) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(emojiName) // we have a fallback
	return name
}

func (e emoji) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(emojiDescription) // we have a fallback
	return desc
}

var customEmojiRegexp = regexp.MustCompile(`^<a?:.+:(?P<id>\d+)>$`)

func (e emoji) Parse(s *state.State, ctx *Context) (interface{}, error) {
	if emojiutil.IsValid(ctx.Raw) {
		return &discord.Emoji{Name: ctx.Raw}, nil
	}

	if !e.customEmojis {
		return nil, newArgParsingErr(emojiOnlyUnicodeErrorArg, emojiOnlyUnicodeErrorFlag, ctx, nil)
	} else if ctx.GuildID == 0 {
		return nil, newArgParsingErr(emojiCustomEmojiInDMError, emojiCustomEmojiInDMError, ctx, nil)
	}

	if matches := customEmojiRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		rawID := matches[1]

		id, err := discord.ParseSnowflake(rawID)
		if err != nil { // range err
			return nil, newArgParsingErr(emojiInvalidError, emojiInvalidError, ctx, nil)
		}

		emoji, err := s.Emoji(ctx.GuildID, discord.EmojiID(id))
		if err != nil {
			return nil, newArgParsingErr(emojiNoAccessError, emojiNoAccessError, ctx, nil)
		}

		return emoji, nil
	}

	if !EmojiAllowIDs {
		return nil, newArgParsingErr(emojiInvalidError, emojiInvalidError, ctx, nil)
	}

	id, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgParsingErr(emojiInvalidError, emojiInvalidError, ctx, nil)
	}

	emoji, err := s.Emoji(ctx.GuildID, discord.EmojiID(id))
	if err != nil {
		return nil, newArgParsingErr(emojiIDNoAccessErrorArg, emojiIDNoAccessErrorFlag, ctx, nil)
	}

	return emoji, nil
}

func (e emoji) Default() interface{} {
	return (*discord.Emoji)(nil)
}
