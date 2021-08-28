package arg

import (
	"fmt"
	"math"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/plugin"
)

func TestRole_Parse(t *testing.T) {
	t.Parallel()

	successCases := []struct {
		name string

		ctx *plugin.ParseContext

		expect *discord.Role
	}{
		{
			name: "mention",
			ctx: &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     "<@&456>",
			},
			expect: &discord.Role{ID: 456},
		},
		{
			name: "id",
			ctx: &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     "456",
			},
			expect: &discord.Role{ID: 456},
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for _, c := range successCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				m, s := state.NewMocker(t)

				m.Roles(c.ctx.GuildID, []discord.Role{*c.expect})

				actual, err := Role.Parse(s, c.ctx)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		t.Run("mention id range", func(t *testing.T) {
			t.Parallel()

			ctx := &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     fmt.Sprintf("<@&%d9>", uint64(math.MaxUint64)),
				Kind:    plugin.KindArg,
			}

			expect := newArgumentError(roleInvalidMentionErrorArg, ctx, nil)

			_, actual := Role.Parse(nil, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = plugin.KindFlag
			expect = newArgumentError(roleInvalidMentionErrorFlag, ctx, nil)

			_, actual = Role.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("mention role not found", func(t *testing.T) {
			t.Parallel()

			srcMocker, _ := state.NewMocker(t)

			var roleID discord.RoleID = 123

			ctx := &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 456}},
				Raw:     roleID.Mention(),
				Kind:    plugin.KindArg,
			}

			srcMocker.Roles(ctx.GuildID, []discord.Role{})

			expect := newArgumentError(roleInvalidMentionErrorArg, ctx, nil)

			_, s := state.CloneMocker(srcMocker, t)

			_, actual := Role.Parse(s, ctx)
			assert.Equal(t, expect, actual)

			ctx.Kind = plugin.KindFlag
			expect = newArgumentError(roleInvalidMentionErrorFlag, ctx, nil)

			_, s = state.CloneMocker(srcMocker, t)

			_, actual = Role.Parse(s, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("not id", func(t *testing.T) {
			t.Parallel()

			ctx := &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     "abc",
			}

			expect := newArgumentError(roleInvalidError, ctx, nil)

			_, actual := Role.Parse(nil, ctx)
			assert.Equal(t, expect, actual)
		})

		t.Run("role id not found", func(t *testing.T) {
			t.Parallel()

			m, s := state.NewMocker(t)

			ctx := &plugin.ParseContext{
				Context: &plugin.Context{Message: discord.Message{GuildID: 123}},
				Raw:     "456",
			}

			m.Roles(ctx.GuildID, []discord.Role{})

			expect := newArgumentError(roleIDInvalidError, ctx, nil)

			_, actual := Role.Parse(s, ctx)
			assert.Equal(t, expect, actual)
		})
	})
}
