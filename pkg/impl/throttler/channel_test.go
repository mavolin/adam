package throttler

import (
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/plugin"
)

func Test_channel_Check(t *testing.T) {
	t.Run("blocked", func(t *testing.T) {
		now = func() time.Time {
			return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
		}

		ctx := &plugin.Context{Message: discord.Message{ChannelID: 123}}

		channel := PerChannel(2, 30*time.Second).(*channel)
		channel.throttler.throttled[discord.Snowflake(ctx.ChannelID)] = []time.Time{
			time.Date(2020, 1, 1, 11, 59, 40, 0, time.UTC),
			time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
		}

		cancelFunc, err := channel.Check(nil, ctx)
		assert.Nil(t, cancelFunc)
		assert.Equal(t, genError(10*time.Second, channelThrottledErrorSecond, channelThrottledErrorMinute), err)
	})

	t.Run("pass", func(t *testing.T) {
		now = func() time.Time {
			return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
		}

		ctx := &plugin.Context{Message: discord.Message{Author: discord.User{ID: 123}}}

		channel := PerChannel(2, 30*time.Second).(*channel)
		channel.throttler.throttled[discord.Snowflake(ctx.Author.ID)] = []time.Time{
			time.Date(2020, 1, 1, 11, 59, 29, 0, time.UTC),
			time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
		}

		cancelFunc, _ := channel.Check(nil, ctx)
		assert.NotNil(t, cancelFunc)
	})
}
