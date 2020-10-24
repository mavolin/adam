package arg

import (
	"fmt"
	"math"
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/diamondburned/arikawa/utils/httputil"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestTextChannel_Parse(t *testing.T) {
	successCases := []struct {
		name string

		raw             string
		allowChannelIDs bool

		expect *discord.Channel
	}{
		{
			name: "mention",
			raw:  "<#123>",
			expect: &discord.Channel{
				ID:      123,
				GuildID: 456,
				Type:    discord.GuildText,
			},
		},
		{
			name:            "id",
			raw:             "123",
			allowChannelIDs: true,
			expect: &discord.Channel{
				ID:      123,
				GuildID: 456,
				Type:    discord.GuildText,
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				TextChannelAllowIDs = c.allowChannelIDs

				m, s := state.NewMocker(t)

				ctx := &Context{
					Context: &plugin.Context{
						MessageCreateEvent: &state.MessageCreateEvent{
							MessageCreateEvent: &gateway.MessageCreateEvent{
								Message: discord.Message{
									GuildID: c.expect.GuildID,
								},
							},
						},
					},
					Raw: c.raw,
				}

				m.Channel(*c.expect)

				actual, err := TextChannel.Parse(s, ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)

				m.Eval()
			})
		}
	})

	failureCases := []struct {
		name string

		raw      string
		allowIDs bool

		expectArg, expectFlag *i18n.Config
	}{
		{
			name:       "mention id range",
			raw:        fmt.Sprintf("<#%d9>", uint64(math.MaxUint64)),
			expectArg:  textChannelInvalidMentionErrorArg,
			expectFlag: textChannelInvalidMentionErrorFlag,
		},
		{
			name:       "invalid - ids not allowed",
			raw:        "abc",
			allowIDs:   false,
			expectArg:  textChannelInvalidMentionWithRawError,
			expectFlag: textChannelInvalidMentionWithRawError,
		},
		{
			name:       "invalid - ids allowed",
			raw:        "abc",
			allowIDs:   true,
			expectArg:  textChannelInvalidError,
			expectFlag: textChannelInvalidError,
		},
	}

	apiErrorCases := []struct {
		name string

		raw      string
		allowIDs bool

		expectArg, expectFlag *i18n.Config
	}{
		{
			name:       "mention - channel not found",
			raw:        "<#123>",
			expectArg:  textChannelInvalidMentionErrorArg,
			expectFlag: textChannelInvalidMentionErrorFlag,
		},
		{
			name:       "id - channel not found",
			raw:        "123",
			allowIDs:   true,
			expectArg:  channelIDInvalidError,
			expectFlag: channelIDInvalidError,
		},
	}

	apiFailureCases := []struct {
		name string

		raw      string
		allowIDs bool
		channel  discord.Channel

		expectArg, expectFlag *i18n.Config
	}{
		{
			name: "mention in dm",
			raw:  "<#123>",
			channel: discord.Channel{
				ID:      123,
				GuildID: 0,
				Type:    discord.GuildText,
			},
			expectArg:  textChannelGuildNotMatchingError,
			expectFlag: textChannelGuildNotMatchingError,
		},
		{
			name: "mention - invalid type",
			raw:  "<#123>",
			channel: discord.Channel{
				ID:      123,
				GuildID: 456,
				Type:    discord.DirectMessage,
			},
			expectArg:  textChannelInvalidTypeError,
			expectFlag: textChannelInvalidTypeError,
		},
		{
			name:     "id in dm",
			raw:      "123",
			allowIDs: true,
			channel: discord.Channel{
				ID:      123,
				GuildID: 0,
				Type:    discord.GuildText,
			},
			expectArg:  textChannelIDGuildNotMatchingError,
			expectFlag: textChannelIDGuildNotMatchingError,
		},
		{
			name:     "id - invalid type",
			raw:      "123",
			allowIDs: true,
			channel: discord.Channel{
				ID:      123,
				GuildID: 456,
				Type:    discord.DirectMessage,
			},
			expectArg:  textChannelIDInvalidTypeError,
			expectFlag: textChannelIDInvalidTypeError,
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				TextChannelAllowIDs = c.allowIDs

				ctx := &Context{
					Context: &plugin.Context{
						MessageCreateEvent: &state.MessageCreateEvent{
							MessageCreateEvent: &gateway.MessageCreateEvent{
								Message: discord.Message{
									GuildID: 456,
								},
							},
						},
					},
					Raw:  c.raw,
					Kind: KindArg,
				}

				expect := c.expectArg
				expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

				_, actual := TextChannel.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

				ctx.Kind = KindFlag

				expect = c.expectFlag
				expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

				_, actual = TextChannel.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
			})
		}

		for _, c := range apiErrorCases {
			t.Run(c.name, func(t *testing.T) {
				TextChannelAllowIDs = c.allowIDs

				srcMocker, _ := state.NewMocker(t)
				srcMocker.Error(http.MethodGet, "/channels/123", httputil.HTTPError{
					Status:  http.StatusNotFound,
					Code:    10003,
					Message: "Unknown channel",
				})

				ctx := &Context{
					Context: &plugin.Context{
						MessageCreateEvent: &state.MessageCreateEvent{
							MessageCreateEvent: &gateway.MessageCreateEvent{
								Message: discord.Message{
									GuildID: 456,
								},
							},
						},
					},
					Raw:  c.raw,
					Kind: KindArg,
				}

				expect := c.expectArg
				expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

				m, s := state.CloneMocker(srcMocker, t)

				_, actual := TextChannel.Parse(s, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

				m.Eval()

				ctx.Kind = KindFlag

				expect = c.expectFlag
				expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

				m, s = state.CloneMocker(srcMocker, t)

				_, actual = TextChannel.Parse(s, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

				m.Eval()
			})
		}

		for _, c := range apiFailureCases {
			t.Run(c.name, func(t *testing.T) {
				TextChannelAllowIDs = c.allowIDs

				srcMocker, _ := state.NewMocker(t)
				srcMocker.Channel(c.channel)

				ctx := &Context{
					Context: &plugin.Context{
						MessageCreateEvent: &state.MessageCreateEvent{
							MessageCreateEvent: &gateway.MessageCreateEvent{
								Message: discord.Message{
									GuildID: 456,
								},
							},
						},
					},
					Raw:  c.raw,
					Kind: KindArg,
				}

				expect := c.expectArg
				expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

				m, s := state.CloneMocker(srcMocker, t)

				_, actual := TextChannel.Parse(s, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

				m.Eval()

				ctx.Kind = KindFlag

				expect = c.expectFlag
				expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

				m, s = state.CloneMocker(srcMocker, t)

				_, actual = TextChannel.Parse(s, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

				m.Eval()
			})
		}
	})
}
