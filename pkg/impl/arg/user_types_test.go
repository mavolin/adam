package arg

import (
	"fmt"
	"math"
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/plugin"
)

func TestUser_Parse(t *testing.T) {
	t.Parallel()

	successCases := []struct {
		name string

		ctx *plugin.ParseContext

		expect *discord.User
	}{
		{
			name: "mention fallback",
			ctx: &plugin.ParseContext{
				Context: new(plugin.Context),
				Raw:     "<@123>",
			},
			expect: &discord.User{ID: 123},
		},
		{
			name:   "id",
			ctx:    &plugin.ParseContext{Raw: "123"},
			expect: &discord.User{ID: 123},
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for _, c := range successCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				m, s := state.NewMocker(t)

				m.User(*c.expect)

				actual, err := User.Parse(s, c.ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}

		t.Run("mention", func(t *testing.T) {
			t.Parallel()

			expect := &discord.User{ID: 123}

			ctx := &plugin.ParseContext{
				Context: &plugin.Context{
					Message: discord.Message{
						Mentions: []discord.GuildUser{{User: *expect}},
					},
				},
				Raw: expect.Mention(),
			}

			actual, err := User.Parse(nil, ctx)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		t.Run("mention id range", func(t *testing.T) {
			t.Parallel()

			ctx := &plugin.ParseContext{
				Raw:  fmt.Sprintf("<@%d9>", uint64(math.MaxUint64)),
				Kind: plugin.KindArg,
			}

			expect := newArgumentError(userInvalidMentionErrorArg, ctx, nil)

			_, actual := User.Parse(nil, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = plugin.KindFlag
			expect = newArgumentError(userInvalidMentionErrorFlag, ctx, nil)

			_, actual = User.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("mention user not found", func(t *testing.T) {
			t.Parallel()

			srcMocker, _ := state.NewMocker(t)

			var userID discord.UserID = 123

			ctx := &plugin.ParseContext{
				Context: new(plugin.Context),
				Raw:     userID.Mention(),
				Kind:    plugin.KindArg,
			}

			srcMocker.Error(http.MethodGet, "users/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := newArgumentError(userInvalidMentionErrorArg, ctx, nil)

			_, s := state.CloneMocker(srcMocker, t)

			_, actual := User.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = plugin.KindFlag
			expect = newArgumentError(userInvalidMentionErrorFlag, ctx, nil)

			_, s = state.CloneMocker(srcMocker, t)

			_, actual = User.Parse(s, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("not id", func(t *testing.T) {
			t.Parallel()

			ctx := &plugin.ParseContext{Raw: "abc"}

			expect := newArgumentError(userInvalidError, ctx, nil)

			_, actual := User.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("id user not found", func(t *testing.T) {
			t.Parallel()

			m, s := state.NewMocker(t)

			var userID discord.UserID = 123

			ctx := &plugin.ParseContext{Raw: "123"}

			m.Error(http.MethodGet, "users/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := newArgumentError(userIDInvalidError, ctx, nil)

			_, actual := User.Parse(s, ctx)
			assert.Equal(t, expect, actual)
		})
	})
}

func TestMember_Parse(t *testing.T) {
	t.Parallel()

	successCases := []struct {
		name string

		ctx *plugin.ParseContext

		expect *discord.Member
	}{
		{
			name: "mention fallback",
			ctx: &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     "<@456>",
			},
			expect: &discord.Member{
				User: discord.User{ID: 456},
			},
		},
		{
			name: "id",
			ctx: &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     "456",
			},
			expect: &discord.Member{
				User: discord.User{ID: 456},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for _, c := range successCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				m, s := state.NewMocker(t)

				m.Member(c.ctx.GuildID, *c.expect)

				actual, err := Member.Parse(s, c.ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}

		t.Run("mention", func(t *testing.T) {
			t.Parallel()

			expect := &discord.Member{
				User: discord.User{ID: 456},
				Deaf: true,
			}

			ctx := &plugin.ParseContext{
				Context: &plugin.Context{
					Message: discord.Message{
						GuildID: 123,
						Mentions: []discord.GuildUser{
							{
								User:   discord.User{ID: 456},
								Member: &discord.Member{Deaf: true},
							},
						},
					},
				},
				Raw: expect.User.Mention(),
			}

			actual, err := Member.Parse(nil, ctx)
			require.NoError(t, err)
			assert.Equal(t, expect, actual)
		})
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		t.Run("mention id range", func(t *testing.T) {
			t.Parallel()

			ctx := &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     fmt.Sprintf("<@%d9>", uint64(math.MaxUint64)),
				Kind:    plugin.KindArg,
			}

			expect := newArgumentError(userInvalidMentionErrorArg, ctx, nil)

			_, actual := Member.Parse(nil, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = plugin.KindFlag
			expect = newArgumentError(userInvalidMentionErrorFlag, ctx, nil)

			_, actual = Member.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("mention member not found", func(t *testing.T) {
			t.Parallel()

			srcMocker, _ := state.NewMocker(t)

			var userID discord.UserID = 456

			ctx := &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     userID.Mention(),
				Kind:    plugin.KindArg,
			}

			srcMocker.Error(http.MethodGet, "guilds/"+ctx.GuildID.String()+"/members/"+userID.String(),
				httputil.HTTPError{
					Status:  http.StatusNotFound,
					Code:    10013, // unknown user
					Message: "Unknown user",
				})

			expect := newArgumentError(userInvalidMentionErrorArg, ctx, nil)

			_, s := state.CloneMocker(srcMocker, t)

			_, actual := Member.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = plugin.KindFlag
			expect = newArgumentError(userInvalidMentionErrorFlag, ctx, nil)

			_, s = state.CloneMocker(srcMocker, t)

			_, actual = Member.Parse(s, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("not id", func(t *testing.T) {
			t.Parallel()

			ctx := &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     "abc",
			}

			expect := newArgumentError(userInvalidError, ctx, nil)

			_, actual := Member.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("id user not found", func(t *testing.T) {
			t.Parallel()

			m, s := state.NewMocker(t)

			var userID discord.UserID = 456

			ctx := &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     userID.String(),
			}

			m.Error(http.MethodGet, "guilds/"+ctx.GuildID.String()+"/members/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := newArgumentError(userIDInvalidError, ctx, nil)

			_, actual := Member.Parse(s, ctx)
			assert.Equal(t, expect, actual)
		})
	})
}
