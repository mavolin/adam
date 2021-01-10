package arg

import (
	"fmt"
	"math"
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/httputil"
	"github.com/mavolin/disstate/v3/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
				Context: new(plugin.Context),
				Raw:     "<@123>",
			},
			expect: &discord.User{ID: 123},
		},
		{
			name:   "id",
			ctx:    &Context{Raw: "123"},
			expect: &discord.User{ID: 123},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := state.NewMocker(t)
				defer m.Eval()

				m.User(*c.expect)

				actual, err := User.Parse(s, c.ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}

		t.Run("mention", func(t *testing.T) {
			expect := &discord.User{ID: 123}

			ctx := &Context{
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
		t.Run("mention id range", func(t *testing.T) {
			ctx := &Context{
				Raw:  fmt.Sprintf("<@%d9>", uint64(math.MaxUint64)),
				Kind: KindArg,
			}

			expect := newArgumentError(userInvalidMentionErrorArg, ctx, nil)

			_, actual := User.Parse(nil, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = KindFlag
			expect = newArgumentError(userInvalidMentionErrorFlag, ctx, nil)

			_, actual = User.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("mention user not found", func(t *testing.T) {
			srcMocker, _ := state.NewMocker(t)

			var userID discord.UserID = 123

			ctx := &Context{
				Context: new(plugin.Context),
				Raw:     userID.Mention(),
				Kind:    KindArg,
			}

			srcMocker.Error(http.MethodGet, "/users/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := newArgumentError(userInvalidMentionErrorArg, ctx, nil)

			m, s := state.CloneMocker(srcMocker, t)

			_, actual := User.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			m.Eval()

			ctx.Kind = KindFlag
			expect = newArgumentError(userInvalidMentionErrorFlag, ctx, nil)

			m, s = state.CloneMocker(srcMocker, t)

			_, actual = User.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			m.Eval()
		})

		t.Run("not id", func(t *testing.T) {
			ctx := &Context{Raw: "abc"}

			expect := newArgumentError(userInvalidError, ctx, nil)

			_, actual := User.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("id user not found", func(t *testing.T) {
			m, s := state.NewMocker(t)
			defer m.Eval()

			var userID discord.UserID = 123

			ctx := &Context{Raw: "123"}

			m.Error(http.MethodGet, "/users/"+userID.String(), httputil.HTTPError{
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
	successCases := []struct {
		name string

		ctx *Context

		expect *discord.Member
	}{
		{
			name: "mention fallback",
			ctx: &Context{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     "<@456>",
			},
			expect: &discord.Member{
				User: discord.User{ID: 456},
			},
		},
		{
			name: "id",
			ctx: &Context{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     "456",
			},
			expect: &discord.Member{
				User: discord.User{ID: 456},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := state.NewMocker(t)
				defer m.Eval()

				m.Member(c.ctx.GuildID, *c.expect)

				actual, err := Member.Parse(s, c.ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}

		t.Run("mention", func(t *testing.T) {
			expect := &discord.Member{
				User: discord.User{ID: 456},
				Deaf: true,
			}

			ctx := &Context{
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
		t.Run("mention id range", func(t *testing.T) {
			ctx := &Context{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     fmt.Sprintf("<@%d9>", uint64(math.MaxUint64)),
				Kind:    KindArg,
			}

			expect := newArgumentError(userInvalidMentionErrorArg, ctx, nil)

			_, actual := Member.Parse(nil, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = KindFlag
			expect = newArgumentError(userInvalidMentionErrorFlag, ctx, nil)

			_, actual = Member.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("mention member not found", func(t *testing.T) {
			srcMocker, _ := state.NewMocker(t)

			var userID discord.UserID = 456

			ctx := &Context{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     userID.Mention(),
				Kind:    KindArg,
			}

			srcMocker.Error(http.MethodGet, "/guilds/"+ctx.GuildID.String()+"/members/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := newArgumentError(userInvalidMentionErrorArg, ctx, nil)

			m, s := state.CloneMocker(srcMocker, t)

			_, actual := Member.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			m.Eval()

			ctx.Kind = KindFlag
			expect = newArgumentError(userInvalidMentionErrorFlag, ctx, nil)

			m, s = state.CloneMocker(srcMocker, t)

			_, actual = Member.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			m.Eval()
		})

		t.Run("not id", func(t *testing.T) {
			ctx := &Context{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     "abc",
			}

			expect := newArgumentError(userInvalidError, ctx, nil)

			_, actual := Member.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("id user not found", func(t *testing.T) {
			m, s := state.NewMocker(t)
			defer m.Eval()

			var userID discord.UserID = 456

			ctx := &Context{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     userID.String(),
			}

			m.Error(http.MethodGet, "/guilds/"+ctx.GuildID.String()+"/members/"+userID.String(), httputil.HTTPError{
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
