package throttler

import (
	"testing"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/plugin"
)

func Test_member_Check(t *testing.T) {
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

			guild := PerMember(2, 30*time.Second).(*member)
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

			guild := PerMember(2, 30*time.Second).(*member)
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

			ctx := &plugin.Context{
				Message: discord.Message{
					Author:  discord.User{ID: 123},
					GuildID: 465,
				},
			}

			guild := PerMember(2, 30*time.Second).(*member)

			guild.memberThrottler[ctx.GuildID] = newSnowflakeThrottler(2, 30*time.Second)
			guild.memberThrottler[ctx.GuildID].throttled[discord.Snowflake(ctx.Author.ID)] = []time.Time{
				time.Date(2020, 1, 1, 11, 59, 40, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
			}

			cancelFunc, err := guild.Check(nil, ctx)
			assert.Nil(t, cancelFunc)
			assert.Equal(t, genError(10*time.Second, memberThrottledErrorSecond, memberThrottledErrorMinute), err)
		})

		t.Run("pass", func(t *testing.T) {
			now = func() time.Time {
				return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
			}

			ctx := &plugin.Context{
				Message: discord.Message{
					Author: discord.User{
						ID: 123,
					},
					GuildID: 456,
				},
			}

			guild := PerMember(2, 30*time.Second).(*member)

			guild.memberThrottler[ctx.GuildID] = newSnowflakeThrottler(2, 30*time.Second)
			guild.memberThrottler[ctx.GuildID].throttled[discord.Snowflake(ctx.Author.ID)] = []time.Time{
				time.Date(2020, 1, 1, 11, 59, 29, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
			}

			cancelFunc, _ := guild.Check(nil, ctx)
			assert.NotNil(t, cancelFunc)
		})
	})
}
