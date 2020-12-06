package throttler

import (
	"testing"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/plugin"
)

func Test_guild_Check(t *testing.T) {
	t.Run("dm", func(t *testing.T) {
		t.Run("blocked", func(t *testing.T) {
			now = func() time.Time {
				return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
			}

			ctx := &plugin.Context{
				Message: discord.Message{
					GuildID: 0,
					Author:  discord.User{ID: 123},
				},
			}

			guild := PerGuild(2, 30*time.Second).(*guild)
			guild.userThrottler.(*user).throttler.throttled[discord.Snowflake(ctx.Author.ID)] = []time.Time{
				time.Date(2020, 1, 1, 11, 59, 40, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
			}

			cancelFunc, err := guild.Check(nil, ctx)
			assert.Nil(t, cancelFunc)
			assert.Equal(t, genError(10*time.Second, userThrottledErrorSecond, userThrottledErrorMinute), err)
		})

		t.Run("pass", func(t *testing.T) {
			now = func() time.Time {
				return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
			}

			ctx := &plugin.Context{
				Message: discord.Message{
					GuildID: 0,
					Author:  discord.User{ID: 123},
				},
			}

			guild := PerGuild(2, 30*time.Second).(*guild)
			guild.userThrottler.(*user).throttler.throttled[discord.Snowflake(ctx.Author.ID)] = []time.Time{
				time.Date(2020, 1, 1, 11, 59, 29, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
			}

			cancelFunc, _ := guild.Check(nil, ctx)
			assert.NotNil(t, cancelFunc)
		})
	})

	t.Run("guild", func(t *testing.T) {
		t.Run("blocked", func(t *testing.T) {
			now = func() time.Time {
				return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
			}

			ctx := &plugin.Context{Message: discord.Message{GuildID: 123}}

			guild := PerGuild(2, 30*time.Second).(*guild)
			guild.guildThrottler.throttled[discord.Snowflake(ctx.GuildID)] = []time.Time{
				time.Date(2020, 1, 1, 11, 59, 40, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
			}

			cancelFunc, err := guild.Check(nil, ctx)
			assert.Nil(t, cancelFunc)
			assert.Equal(t, genError(10*time.Second, guildThrottledErrorSecond, guildThrottledErrorMinute), err)
		})

		t.Run("pass", func(t *testing.T) {
			var s discord.Snowflake = 123

			now = func() time.Time {
				return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
			}

			guild := PerGuild(2, 30*time.Second).(*guild)
			guild.guildThrottler.throttled[s] = []time.Time{
				time.Date(2020, 1, 1, 11, 59, 29, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
			}

			ctx := &plugin.Context{Message: discord.Message{GuildID: 132}}

			cancelFunc, _ := guild.Check(nil, ctx)
			assert.NotNil(t, cancelFunc)
		})
	})
}
