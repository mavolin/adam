package arg

import (
	"math"
	"strconv"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestRole_Parse(t *testing.T) {
	successCases := []struct {
		name string

		ctx *Context

		expect *discord.Role
	}{
		{
			name: "mention",
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
			expect: &discord.Role{ID: 456},
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
			expect: &discord.Role{ID: 456},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				m, s := state.NewMocker(t)

				m.Roles(c.ctx.GuildID, []discord.Role{*c.expect})

				actual, err := Role.Parse(s, c.ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)

				m.Eval()
			})
		}
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

			expect := roleInvalidMentionArg
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual := Role.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			ctx.Kind = KindFlag

			expect = roleInvalidMentionFlag
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual = Role.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
		})

		t.Run("mention role not found", func(t *testing.T) {
			srcMocker, _ := state.NewMocker(t)

			var roleID discord.UserID = 123

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
				Raw:  roleID.Mention(),
				Kind: KindArg,
			}

			srcMocker.Roles(ctx.GuildID, []discord.Role{})

			expect := roleInvalidMentionArg
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s := state.CloneMocker(srcMocker, t)

			_, actual := Role.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			m.Eval()

			ctx.Kind = KindFlag

			expect = roleInvalidMentionFlag
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s = state.CloneMocker(srcMocker, t)

			_, actual = Role.Parse(s, ctx)
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

			expect := roleInvalidIDWithRaw
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual := Role.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			ctx.Kind = KindFlag

			expect = roleInvalidIDWithRaw
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			_, actual = Role.Parse(nil, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)
		})

		t.Run("role id not found", func(t *testing.T) {
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

			srcMocker.Roles(ctx.GuildID, []discord.Role{})

			expect := roleInvalidIDArg
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s := state.CloneMocker(srcMocker, t)

			_, actual := Role.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			m.Eval()

			ctx.Kind = KindFlag

			expect = roleInvalidIDFlag
			expect.Placeholders = attachDefaultPlaceholders(expect.Placeholders, ctx)

			m, s = state.CloneMocker(srcMocker, t)

			_, actual = Role.Parse(s, ctx)
			assert.Equal(t, errors.NewArgumentParsingErrorl(expect), actual)

			m.Eval()
		})
	})
}
