package throttler

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/plugin"
)

type mockThrottler struct {
	num int
	ref *int
}

func (t *mockThrottler) Check(*state.State, *plugin.Context) (func(), error) {
	*t.ref = t.num
	return func() {}, nil
}

func TestConditional_Check(t *testing.T) {
	testCases := []struct {
		name        string
		conditional Conditional
		ctx         *plugin.Context
		expect      int
	}{
		{
			name: "condition",
			conditional: Conditional{
				Conditions: []Condition{
					{
						Restrictions: restriction.BotOwner,
						Throttler:    &mockThrottler{num: 1},
					},
					{
						Restrictions: restriction.BotOwner,
						Throttler:    &mockThrottler{num: 2},
					},
				},
			},
			ctx: &plugin.Context{
				Message:     discord.Message{Author: discord.User{ID: 123}},
				BotOwnerIDs: []discord.UserID{123},
			},
			expect: 1,
		},
		{
			name: "default",
			conditional: Conditional{
				Conditions: []Condition{
					{
						Restrictions: restriction.BotOwner,
						Throttler:    &mockThrottler{num: 1},
					},
				},
				Default: &mockThrottler{num: 2},
			},
			ctx:    new(plugin.Context),
			expect: 2,
		},
		{
			name: "condition - nil throttler",
			conditional: Conditional{
				Conditions: []Condition{
					{
						Restrictions: restriction.BotOwner,
						Throttler:    nil,
					},
				},
				Default: &mockThrottler{num: 1},
			},
			ctx: &plugin.Context{
				Message:     discord.Message{Author: discord.User{ID: 123}},
				BotOwnerIDs: []discord.UserID{123},
			},
			expect: 0,
		},
		{
			name: "default - nil throttler",
			conditional: Conditional{
				Conditions: []Condition{
					{
						Restrictions: restriction.BotOwner,
						Throttler:    &mockThrottler{num: 1},
					},
				},
				Default: nil,
			},
			ctx:    new(plugin.Context),
			expect: 0,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ref := 0

			for _, con := range c.conditional.Conditions {
				if mt, ok := con.Throttler.(*mockThrottler); ok {
					mt.ref = &ref
				}
			}

			if mt, ok := c.conditional.Default.(*mockThrottler); ok {
				mt.ref = &ref
			}

			_, err := c.conditional.Check(nil, c.ctx)
			require.NoError(t, err)
			assert.Equal(t, c.expect, ref)
		})
	}
}
