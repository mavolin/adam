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

			expect := newArgParsingErr(userInvalidMentionErrorArg, ctx, nil)

			_, actual := User.Parse(nil, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = KindFlag
			expect = newArgParsingErr(userInvalidMentionErrorFlag, ctx, nil)

			_, actual = User.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
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
				Raw:  userID.Mention(),
				Kind: KindArg,
			}

			srcMocker.Error(http.MethodGet, "/users/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := newArgParsingErr(userInvalidMentionErrorArg, ctx, nil)

			m, s := state.CloneMocker(srcMocker, t)

			_, actual := User.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			m.Eval()

			ctx.Kind = KindFlag
			expect = newArgParsingErr(userInvalidMentionErrorFlag, ctx, nil)

			m, s = state.CloneMocker(srcMocker, t)

			_, actual = User.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			m.Eval()
		})

		t.Run("not id", func(t *testing.T) {
			ctx := &Context{Raw: "abc"}

			expect := newArgParsingErr(userInvalidError, ctx, nil)

			_, actual := User.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("id user not found", func(t *testing.T) {
			m, s := state.NewMocker(t)

			var userID discord.UserID = 123

			ctx := &Context{
				Raw: "123",
			}

			m.Error(http.MethodGet, "/users/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := newArgParsingErr(userIDInvalidError, ctx, nil)

			_, actual := User.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			m.Eval()
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
				Context: &plugin.Context{
					MessageCreateEvent: &state.MessageCreateEvent{
						MessageCreateEvent: &gateway.MessageCreateEvent{
							Message: discord.Message{
								GuildID: 123,
							},
						},
					},
				},
				Raw: "<@456>",
			},
			expect: &discord.Member{
				User: discord.User{ID: 456},
			},
		},
		{
			name: "id",
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
			expect: &discord.Member{
				User: discord.User{ID: 456},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := state.NewMocker(t)

				m.Member(c.ctx.GuildID, *c.expect)

				actual, err := Member.Parse(s, c.ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)

				m.Eval()
			})
		}

		t.Run("mention", func(t *testing.T) {
			expect := &discord.Member{
				User: discord.User{ID: 456},
				Deaf: true,
			}

			ctx := &Context{
				Context: &plugin.Context{
					MessageCreateEvent: &state.MessageCreateEvent{
						MessageCreateEvent: &gateway.MessageCreateEvent{
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
				Context: &plugin.Context{
					MessageCreateEvent: &state.MessageCreateEvent{
						MessageCreateEvent: &gateway.MessageCreateEvent{
							Message: discord.Message{
								GuildID: 123,
							},
						},
					},
				},
				Raw:  fmt.Sprintf("<@%d9>", uint64(math.MaxUint64)),
				Kind: KindArg,
			}

			expect := newArgParsingErr(userInvalidMentionErrorArg, ctx, nil)

			_, actual := Member.Parse(nil, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = KindFlag
			expect = newArgParsingErr(userInvalidMentionErrorFlag, ctx, nil)

			_, actual = Member.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("mention member not found", func(t *testing.T) {
			srcMocker, _ := state.NewMocker(t)

			var userID discord.UserID = 456

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
				Raw:  userID.Mention(),
				Kind: KindArg,
			}

			srcMocker.Error(http.MethodGet, "/guilds/"+ctx.GuildID.String()+"/members/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := newArgParsingErr(userInvalidMentionErrorArg, ctx, nil)

			m, s := state.CloneMocker(srcMocker, t)

			_, actual := Member.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			m.Eval()

			ctx.Kind = KindFlag
			expect = newArgParsingErr(userInvalidMentionErrorFlag, ctx, nil)

			m, s = state.CloneMocker(srcMocker, t)

			_, actual = Member.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			m.Eval()
		})

		t.Run("not id", func(t *testing.T) {
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
				Raw: "abc",
			}

			expect := newArgParsingErr(userInvalidError, ctx, nil)

			_, actual := Member.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("id user not found", func(t *testing.T) {
			m, s := state.NewMocker(t)

			var userID discord.UserID = 456

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
				Raw: userID.String(),
			}

			m.Error(http.MethodGet, "/guilds/"+ctx.GuildID.String()+"/members/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := newArgParsingErr(userIDInvalidError, ctx, nil)

			_, actual := Member.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			m.Eval()
		})
	})
}
