package throttler

import (
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"
)

func Test_snowflakeThrottler_expire(t *testing.T) {
	testCases := []struct {
		name      string
		throttled []time.Time
		before    time.Time
		expect    []time.Time
	}{
		{
			name: "none",
			throttled: []time.Time{
				time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
			},
			before: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
			expect: []time.Time{
				time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "none - exact match",
			throttled: []time.Time{
				time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
			},
			before: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
			expect: []time.Time{
				time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "before",
			throttled: []time.Time{
				time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
			},
			before: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
			expect: []time.Time{
				time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			var s discord.Snowflake = 123

			st := newSnowflakeThrottler(10, 10)

			st.throttled[s] = c.throttled

			st.expire(s, c.before)

			assert.Equal(t, c.expect, st.throttled[s])
		})
	}
}

func Test_snowflakeThrottler_check(t *testing.T) {
	t.Run("blocked", func(t *testing.T) {
		var s discord.Snowflake = 123

		now = func() time.Time {
			return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
		}

		st := newSnowflakeThrottler(2, 30*time.Second)
		st.throttled[s] = []time.Time{
			time.Date(2020, 1, 1, 11, 59, 40, 0, time.UTC),
			time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
		}

		cancelFunc, actualDuration := st.check(s)
		assert.Nil(t, cancelFunc)
		assert.Equal(t, 10*time.Second, actualDuration)
	})

	t.Run("pass", func(t *testing.T) {
		var s discord.Snowflake = 123

		now = func() time.Time {
			return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
		}

		st := newSnowflakeThrottler(2, 30*time.Second)
		st.throttled[s] = []time.Time{
			time.Date(2020, 1, 1, 11, 59, 29, 0, time.UTC),
			time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
		}

		cancelFunc, _ := st.check(s)
		assert.NotNil(t, cancelFunc)
	})

	t.Run("cancel", func(t *testing.T) {
		var s discord.Snowflake = 123

		now = func() time.Time {
			return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
		}

		st := newSnowflakeThrottler(2, 30*time.Second)

		cancelFunc, _ := st.check(s)
		assert.NotNil(t, cancelFunc)

		assert.Len(t, st.throttled[s], 1)

		cancelFunc()

		assert.Len(t, st.throttled[s], 0)
	})
}
