package throttling

import (
	"testing"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/plugin"
)

func Test_user_Check(t *testing.T) {
	t.Run("blocked", func(t *testing.T) {
		var s discord.Snowflake = 123

		now = func() time.Time {
			return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
		}

		user := PerUser(2, 30*time.Second).(*user)
		user.throttler.throttled[s] = []time.Time{
			time.Date(2020, 1, 1, 11, 59, 40, 0, time.UTC),
			time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
		}

		ctx := &plugin.Context{
			MessageCreateEvent: &state.MessageCreateEvent{
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						Author: discord.User{
							ID: discord.UserID(s),
						},
					},
				},
			},
		}

		cancelFunc, err := user.Check(ctx)
		assert.Nil(t, cancelFunc)
		assert.Equal(t, genError(10*time.Second, userErrorSecond, userErrorMinute), err)
	})

	t.Run("pass", func(t *testing.T) {
		var s discord.Snowflake = 123

		now = func() time.Time {
			return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
		}

		user := PerUser(2, 30*time.Second).(*user)
		user.throttler.throttled[s] = []time.Time{
			time.Date(2020, 1, 1, 11, 59, 29, 0, time.UTC),
			time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
		}

		ctx := &plugin.Context{
			MessageCreateEvent: &state.MessageCreateEvent{
				MessageCreateEvent: &gateway.MessageCreateEvent{
					Message: discord.Message{
						Author: discord.User{
							ID: discord.UserID(s),
						},
					},
				},
			},
		}

		cancelFunc, _ := user.Check(ctx)
		assert.NotNil(t, cancelFunc)
	})
}
