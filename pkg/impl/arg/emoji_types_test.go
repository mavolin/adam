package arg

import (
	"fmt"
	"math"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	emojiutil "github.com/mavolin/adam/pkg/utils/emoji"
)

func TestEmoji_Parse(t *testing.T) {
	apiSuccessCases := []struct {
		name string

		raw           string
		allowEmojiIDs bool

		expect *discord.Emoji
	}{
		{
			name: "custom emoji",
			raw:  "<:thonk:456>",
			expect: &discord.Emoji{
				Name: "thonk",
				ID:   456,
				User: discord.User{ID: 1},
			},
		},
		{
			name:          "id",
			raw:           "789",
			allowEmojiIDs: true,
			expect: &discord.Emoji{
				Name: "pepe",
				ID:   789,
				User: discord.User{ID: 1},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range apiSuccessCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := state.NewMocker(t)
				defer m.Eval()

				ctx := &Context{
					Context: &plugin.Context{
						Message: discord.Message{GuildID: 123},
					},
					Raw: c.raw,
				}

				m.Emojis(ctx.GuildID, []discord.Emoji{*c.expect})

				EmojiAllowIDs = c.allowEmojiIDs

				actual, err := Emoji.Parse(s, ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}

		t.Run("unicode emoji", func(t *testing.T) {
			expect := &discord.Emoji{Name: emojiutil.Cloud}

			ctx := &Context{Raw: expect.Name}

			actual, err := Emoji.Parse(nil, ctx)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})
	})

	failureCases := []struct {
		name string

		raw           string
		guild         bool
		allowEmojiIDs bool
		customEmojis  bool

		expectArg, expectFlag *i18n.Config
	}{
		{
			name:         "custom emoji - no custom emojis allowed",
			raw:          "<:abc:123>",
			guild:        true,
			customEmojis: false,
			expectArg:    emojiCustomEmojiErrorArg,
			expectFlag:   emojiCustomEmojiErrorFlag,
		},
		{
			name:         "custom emoji in dm",
			raw:          "<:abc:123>",
			guild:        false,
			customEmojis: true,
			expectArg:    emojiCustomEmojiInDMError,
			expectFlag:   emojiCustomEmojiInDMError,
		},
		{
			name:         "custom emoji id range error",
			raw:          fmt.Sprintf("<:abc:%v9>", uint64(math.MaxUint64)),
			guild:        true,
			customEmojis: true,
			expectArg:    emojiInvalidError,
			expectFlag:   emojiInvalidError,
		},
		{
			name:          "id - no id allowed",
			raw:           "123",
			guild:         true,
			allowEmojiIDs: false,
			customEmojis:  true,
			expectArg:     emojiInvalidError,
			expectFlag:    emojiInvalidError,
		},
		{
			name:          "id - no custom emojis allowed",
			raw:           "123",
			guild:         true,
			allowEmojiIDs: true,
			customEmojis:  false,
			expectArg:     emojiCustomEmojiErrorArg,
			expectFlag:    emojiCustomEmojiErrorFlag,
		},
		{
			name:          "id in dm",
			raw:           "123",
			guild:         false,
			allowEmojiIDs: true,
			customEmojis:  true,
			expectArg:     emojiCustomEmojiInDMError,
			expectFlag:    emojiCustomEmojiInDMError,
		},
		{
			name:          "id invalid",
			raw:           "abc",
			guild:         true,
			allowEmojiIDs: true,
			customEmojis:  true,
			expectArg:     emojiInvalidError,
			expectFlag:    emojiInvalidError,
		},
	}

	apiFailureCases := []struct {
		name string

		raw           string
		allowEmojiIDs bool

		expectArg, expectFlag *i18n.Config
	}{
		{
			name:       "custom emoji not found",
			raw:        "<:abc:123>",
			expectArg:  emojiNoAccessError,
			expectFlag: emojiNoAccessError,
		},
		{
			name:          "emoji id not found",
			raw:           "123",
			allowEmojiIDs: true,
			expectArg:     emojiIDNoAccessError,
			expectFlag:    emojiIDNoAccessError,
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				EmojiAllowIDs = c.allowEmojiIDs

				ctx := &Context{
					Raw:     c.raw,
					Context: new(plugin.Context),
					Kind:    KindArg,
				}

				if c.guild {
					ctx.GuildID = 456
				}

				emoji := new(emoji)
				emoji.customEmojis = c.customEmojis

				expect := newArgParsingErr(c.expectArg, ctx, nil)

				_, actual := emoji.Parse(nil, ctx)
				assert.Equal(t, expect, actual)

				ctx.Kind = KindFlag

				expect = newArgParsingErr(c.expectFlag, ctx, nil)

				_, actual = emoji.Parse(nil, ctx)
				assert.Equal(t, expect, actual)
			})
		}

		for _, c := range apiFailureCases {
			t.Run(c.name, func(t *testing.T) {
				srcMocker, _ := state.NewMocker(t)

				EmojiAllowIDs = c.allowEmojiIDs

				ctx := &Context{
					Raw:     c.raw,
					Context: &plugin.Context{Message: discord.Message{GuildID: 456}},
					Kind:    KindArg,
				}

				srcMocker.Emojis(ctx.GuildID, []discord.Emoji{})

				expect := newArgParsingErr(c.expectArg, ctx, nil)

				m, s := state.CloneMocker(srcMocker, t)

				_, actual := Emoji.Parse(s, ctx)
				assert.Equal(t, expect, actual)

				m.Eval()

				ctx.Kind = KindFlag

				expect = newArgParsingErr(c.expectFlag, ctx, nil)

				m, s = state.CloneMocker(srcMocker, t)

				_, actual = Emoji.Parse(s, ctx)
				assert.Equal(t, expect, actual)

				m.Eval()
			})
		}
	})
}

func TestRawEmoji_Parse(t *testing.T) {
	successCases := []struct {
		name string

		raw string

		expect discord.APIEmoji
	}{
		{
			name:   "unicode",
			raw:    emojiutil.SmilingFaceWithHeartEyes,
			expect: discord.APIEmoji(emojiutil.SmilingFaceWithHeartEyes),
		},
		{
			name:   "custom emoji",
			raw:    "<:abc:123>",
			expect: "abc:123",
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				ctx := &Context{Raw: c.raw}

				actual, err := RawEmoji.Parse(nil, ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		ctx := &Context{Raw: "abc"}

		expect := newArgParsingErr(emojiInvalidError, ctx, nil)

		_, actual := RawEmoji.Parse(nil, ctx)
		assert.Equal(t, expect, actual)
	})
}
