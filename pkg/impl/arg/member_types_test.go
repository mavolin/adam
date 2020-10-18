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
				Raw:  "<@" + strconv.FormatUint(math.MaxUint64, 10) + "9>",
				Kind: KindArg,
			}

			expect := userInvalidMentionArg
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual := Member.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			ctx.Kind = KindFlag

			expect = userInvalidMentionFlag
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual = Member.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
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

			expect := userInvalidMentionArg
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s := state.CloneMocker(srcMocker, t)

			_, actual := Member.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			m.Eval()

			ctx.Kind = KindFlag

			expect = userInvalidMentionFlag
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s = state.CloneMocker(srcMocker, t)

			_, actual = Member.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			m.Eval()
		})

		t.Run("invalid id", func(t *testing.T) {
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
				Raw:  "abc",
				Kind: KindArg,
			}

			expect := userInvalidIDWithRaw
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual := Member.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			ctx.Kind = KindFlag

			expect = userInvalidIDWithRaw
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual = Member.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
		})

		t.Run("id user not found", func(t *testing.T) {
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
				Raw:  userID.String(),
				Kind: KindArg,
			}

			srcMocker.Error(http.MethodGet, "/guilds/"+ctx.GuildID.String()+"/members/"+userID.String(), httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			expect := userInvalidIDArg
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s := state.CloneMocker(srcMocker, t)

			_, actual := Member.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			m.Eval()

			ctx.Kind = KindFlag

			expect = userInvalidIDFlag
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s = state.CloneMocker(srcMocker, t)

			_, actual = Member.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			m.Eval()
		})
	})
}

func TestMemberID_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := state.NewMocker(t)

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
			Raw: "456",
		}

		expect := &discord.Member{
			User: discord.User{
				ID: 456,
			},
		}

		m.Member(ctx.GuildID, *expect)

		actual, err := MemberID.Parse(s, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)

		m.Eval()
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("invalid id", func(t *testing.T) {
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

			desc := userInvalidIDWithRaw
			desc.Placeholders = attachDefaultPlaceholders(userInvalidIDWithRaw.Placeholders, ctx)

			expect := errors.NewArgumentParsingErrorl(desc)

			_, actual := MemberID.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("member not found", func(t *testing.T) {
			srcMocker, _ := state.NewMocker(t)

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
				Raw:  "456",
				Kind: KindArg,
			}

			srcMocker.Error(http.MethodGet, "/guilds/"+ctx.GuildID.String()+"/members/456", httputil.HTTPError{
				Status:  http.StatusNotFound,
				Code:    10013, // unknown user
				Message: "Unknown user",
			})

			desc := userInvalidIDArg
			desc.Placeholders = attachDefaultPlaceholders(desc.Placeholders, ctx)

			expect := errors.NewArgumentParsingErrorl(desc)

			_, s := state.CloneMocker(srcMocker, t)

			_, actual := MemberID.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = KindFlag

			desc = userInvalidIDFlag
			desc.Placeholders = attachDefaultPlaceholders(desc.Placeholders, ctx)

			expect = errors.NewArgumentParsingErrorl(desc)

			_, s = state.CloneMocker(srcMocker, t)

			_, actual = MemberID.Parse(s, ctx)
			assert.Equal(t, expect, actual)
		})
	})
}
