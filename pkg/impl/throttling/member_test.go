package throttling

import (
	"testing"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/plugin"
)

func Test_member_Check(t *testing.T) {
	t.Run("dm", func(t *testing.T) {
		t.Run("blocked", func(t *testing.T) {
			var s discord.Snowflake = 123

			now = func() time.Time {
				return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
			}

			guild := PerMember(2, 30*time.Second).(*member)
			guild.userThrottler.(*user).throttler.throttled[s] = []time.Time{
				time.Date(2020, 1, 1, 11, 59, 40, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
			}

			ctx := &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 0,
							Author: discord.User{
								ID: discord.UserID(s),
							},
						},
					},
				},
			}

			cancelFunc, err := guild.Check(ctx)
			assert.Nil(t, cancelFunc)
			assert.Equal(t, genError(10*time.Second, userErrorSecond, userErrorMinute), err)
		})

		t.Run("pass", func(t *testing.T) {
			var s discord.Snowflake = 123

			now = func() time.Time {
				return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
			}

			guild := PerMember(2, 30*time.Second).(*member)
			guild.userThrottler.(*user).throttler.throttled[s] = []time.Time{
				time.Date(2020, 1, 1, 11, 59, 29, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
			}

			ctx := &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							GuildID: 0,
							Author: discord.User{
								ID: discord.UserID(s),
							},
						},
					},
				},
			}

			cancelFunc, _ := guild.Check(ctx)
			assert.NotNil(t, cancelFunc)
		})
	})

	t.Run("guild", func(t *testing.T) {
		t.Run("blocked", func(t *testing.T) {
			var (
				guildID discord.GuildID   = 123
				userID  discord.Snowflake = 456
			)

			now = func() time.Time {
				return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
			}

			guild := PerMember(2, 30*time.Second).(*member)

			guild.memberThrottler[guildID] = newSnowflakeThrottler(2, 30*time.Second)
			guild.memberThrottler[guildID].throttled[userID] = []time.Time{
				time.Date(2020, 1, 1, 11, 59, 40, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
			}

			ctx := &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							Author: discord.User{
								ID: discord.UserID(userID),
							},
							GuildID: guildID,
						},
					},
				},
			}

			cancelFunc, err := guild.Check(ctx)
			assert.Nil(t, cancelFunc)
			assert.Equal(t, genError(10*time.Second, memberErrorSecond, memberErrorMinute), err)
		})

		t.Run("pass", func(t *testing.T) {
			var (
				guildID discord.GuildID   = 123
				userID  discord.Snowflake = 456
			)

			now = func() time.Time {
				return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
			}

			guild := PerMember(2, 30*time.Second).(*member)

			guild.memberThrottler[guildID] = newSnowflakeThrottler(2, 30*time.Second)
			guild.memberThrottler[guildID].throttled[userID] = []time.Time{
				time.Date(2020, 1, 1, 11, 59, 29, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 59, 50, 0, time.UTC),
			}

			ctx := &plugin.Context{
				MessageCreateEvent: &state.MessageCreateEvent{
					MessageCreateEvent: &gateway.MessageCreateEvent{
						Message: discord.Message{
							Author: discord.User{
								ID: discord.UserID(userID),
							},
							GuildID: guildID,
						},
					},
				},
			}

			cancelFunc, _ := guild.Check(ctx)
			assert.NotNil(t, cancelFunc)
		})
	})
}
