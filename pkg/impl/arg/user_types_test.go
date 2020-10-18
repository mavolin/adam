package arg

import (
	"math"
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
	"github.com/mavolin/adam/pkg/plugin"
)

func TestUser_Parse(t *testing.T) {
	successCases := []struct {
		name string

		ctx *Context

		expect *discord.User
	}{
		{
			name: "mention fallback",
			ctx: &Context{
				Context: &plugin.Context{
					MessageCreateEvent: &state.MessageCreateEvent{
						MessageCreateEvent: new(gateway.MessageCreateEvent),
					},
				},
				Raw: "<@123>",
			},
			expect: &discord.User{ID: 123},
		},
		{
			name: "id",
			ctx: &Context{
				Raw: "123",
			},
			expect: &discord.User{ID: 123},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := state.NewMocker(t)

				m.User(*c.expect)

				actual, err := User.Parse(s, c.ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)

				m.Eval()
			})
		}

		t.Run("mention", func(t *testing.T) {
			expect := &discord.User{ID: 123}

			ctx := &Context{
				Context: &plugin.Context{
					MessageCreateEvent: &state.MessageCreateEvent{
						MessageCreateEvent: &gateway.MessageCreateEvent{
							Message: discord.Message{
								Mentions: []discord.GuildUser{
									{
										User: *expect,
									},
								},
							},
						},
					},
				},
				Raw: "<@" + expect.ID.String() + ">",
			}

			actual, err := User.Parse(nil, ctx)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("mention id range", func(t *testing.T) {
			ctx := &Context{
				Raw:  "<@" + strconv.FormatUint(math.MaxUint64, 10) + "9>",
				Kind: KindArg,
			}

			expect := userInvalidMentionArg
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual := User.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			ctx.Kind = KindFlag

			expect = userInvalidMentionFlag
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual = User.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
		})

		t.Run("mention user not found", func(t *testing.T) {
			srcMocker, _ := state.NewMocker(t)

			var userID discord.UserID = 123

			ctx := &Context{
				Context: &plugin.Context{
					MessageCreateEvent: &state.MessageCreateEvent{
						MessageCreateEvent: new(gateway.MessageCreateEvent),
					},
				},
				Raw:  "<@" + userID.String() + ">",
				Kind: KindArg,
			}

			srcMocker.Error(http.MethodGet, "/users/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := userInvalidMentionArg
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s := state.CloneMocker(srcMocker, t)

			_, actual := User.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			m.Eval()

			ctx.Kind = KindFlag

			expect = userInvalidMentionFlag
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s = state.CloneMocker(srcMocker, t)

			_, actual = User.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			m.Eval()
		})

		t.Run("invalid id", func(t *testing.T) {
			ctx := &Context{
				Raw:  "abc",
				Kind: KindArg,
			}

			expect := userInvalidIDWithRaw
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual := User.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			ctx.Kind = KindFlag

			expect = userInvalidIDWithRaw
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual = User.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
		})

		t.Run("id user not found", func(t *testing.T) {
			srcMocker, _ := state.NewMocker(t)

			var userID discord.UserID = 123

			ctx := &Context{
				Raw:  "123",
				Kind: KindArg,
			}

			srcMocker.Error(http.MethodGet, "/users/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := userInvalidIDArg
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s := state.CloneMocker(srcMocker, t)

			_, actual := User.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			m.Eval()

			ctx.Kind = KindFlag

			expect = userInvalidIDFlag
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s = state.CloneMocker(srcMocker, t)

			_, actual = User.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			m.Eval()
		})
	})
}
