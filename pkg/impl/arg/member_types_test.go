package arg

import (
	"net/http"
	"strconv"
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

func TestMember_Parse(t *testing.T) {
	successCases := []struct {
		name string

		raw      string
		allowIDs bool

		expect *discord.Member
	}{
		{
			name: "mention",
			raw:  "<@456>",
			expect: &discord.Member{
				User: discord.User{ID: 456},
			},
		},
		{
			name:     "id",
			raw:      "456",
			allowIDs: true,
			expect: &discord.Member{
				User: discord.User{ID: 456},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := state.NewMocker(t)

				MemberAllowIDs = c.allowIDs

				ctx := &Context{
					Context: &plugin.Context{
						MessageCreateEvent: &state.MessageCreateEvent{
							MessageCreateEvent: &gateway.MessageCreateEvent{
								Message: discord.Message{
									GuildID: 123,
								},
							},
						},
					},
					Raw: c.raw,
				}

				m.Member(ctx.GuildID, *c.expect)

				actual, err := Member.Parse(s, ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)

				m.Eval()
			})
		}
	})

	var failureCasesWithoutMember = []struct {
		name string

		ctx      *Context
		allowIDs bool

		expectArg, expectFlag *i18n.Config
	}{
		{
			name: "snowflake out of range",
			ctx: &Context{
				Context: &plugin.Context{
					MessageCreateEvent: &state.MessageCreateEvent{
						MessageCreateEvent: &gateway.MessageCreateEvent{
							Message: discord.Message{
								GuildID: 123,
							},
						},
					},
				},
				Raw: "<@99999999999999999999999999999999999999999999999999>",
			},
			expectArg:  userInvalidMentionArg,
			expectFlag: userInvalidMentionFlag,
		},
		{
			name: "id - no ids allowed",
			ctx: &Context{
				Context: &plugin.Context{
					MessageCreateEvent: &state.MessageCreateEvent{
						MessageCreateEvent: &gateway.MessageCreateEvent{
							Message: discord.Message{
								GuildID: 123,
							},
						},
					},
				},
				Raw: "456",
			},
			allowIDs:   false,
			expectArg:  userInvalidMentionWithRaw,
			expectFlag: userInvalidMentionWithRaw,
		},
		{
			name: "invalid id",
			ctx: &Context{
				Context: &plugin.Context{
					MessageCreateEvent: &state.MessageCreateEvent{
						MessageCreateEvent: &gateway.MessageCreateEvent{
							Message: discord.Message{
								GuildID: 123,
							},
						},
					},
				},
				Raw: "abc",
			},
			allowIDs:   true,
			expectArg:  userInvalidIDWithRaw,
			expectFlag: userInvalidIDWithRaw,
		},
	}

	failureCasesWithMemberCall := []struct {
		name string

		raw         string
		memberError httputil.HTTPError
		allowIDs    bool

		expectArg, expectFlag *i18n.Config
	}{
		{
			name: "mention - member not found",
			raw:  "<@456>",
			memberError: httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			},
			expectArg:  userInvalidMentionArg,
			expectFlag: userInvalidMentionFlag,
		},
		{
			name: "id - member not found",
			raw:  "456",
			memberError: httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			},
			allowIDs:   true,
			expectArg:  userInvalidIDArg,
			expectFlag: userInvalidIDFlag,
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCasesWithoutMember {
			MemberAllowIDs = c.allowIDs

			c.expectArg.Placeholders = attachDefaultPlaceholders(c.expectArg.Placeholders, c.ctx)

			c.ctx.Kind = KindArg

			_, actual := Member.Parse(nil, c.ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectArg), actual)

			c.ctx.Kind = KindFlag

			c.expectFlag.Placeholders = attachDefaultPlaceholders(c.expectFlag.Placeholders, c.ctx)

			_, actual = Member.Parse(nil, c.ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectFlag), actual)
		}

		for _, c := range failureCasesWithMemberCall {
			srcMocker, _ := state.NewMocker(t)

			MemberAllowIDs = c.allowIDs

			ctx := &Context{
				Context: &plugin.Context{
					MessageCreateEvent: &state.MessageCreateEvent{
						MessageCreateEvent: &gateway.MessageCreateEvent{
							Message: discord.Message{
								GuildID: 123,
							},
						},
					},
				},
				Raw:  c.raw,
				Kind: KindArg,
			}

			srcMocker.Error(http.MethodGet, "/guilds/"+ctx.GuildID.String()+"/members/"+strconv.Itoa(456), c.memberError)

			c.expectArg.Placeholders = attachDefaultPlaceholders(c.expectArg.Placeholders, ctx)

			_, s := state.CloneMocker(srcMocker, t)

			_, actual := Member.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectArg), actual)

			ctx = &Context{
				Context: &plugin.Context{
					MessageCreateEvent: &state.MessageCreateEvent{
						MessageCreateEvent: &gateway.MessageCreateEvent{
							Message: discord.Message{
								GuildID: 123,
							},
						},
					},
				},
				Raw:  c.raw,
				Kind: KindFlag,
			}

			c.expectFlag.Placeholders = attachDefaultPlaceholders(c.expectFlag.Placeholders, ctx)

			_, s = state.CloneMocker(srcMocker, t)

			_, actual = Member.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectFlag), actual)
		}
	})
}
