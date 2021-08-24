package throttler

import (
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// guild is a plugin.Throttler that works on a per-guild basis.
type guild struct {
	guildThrottler *snowflakeThrottler
	userThrottler  plugin.Throttler
}

var _ plugin.Throttler = new(guild)

// PerGuild returns a new plugin.Throttler that works on a per-guild basis.
// It allows at maximum the passed number of invokes in the passed duration.
//
// All commands invoked in direct messages will be throttled on a per-user
// basis.
func PerGuild(maxInvokes uint, duration time.Duration) plugin.Throttler {
	return &guild{
		guildThrottler: newSnowflakeThrottler(maxInvokes, duration),
		userThrottler:  PerUser(maxInvokes, duration),
	}
}

func (g *guild) Check(_ *state.State, ctx *plugin.Context) (func(), error) {
	if ctx.GuildID == 0 {
		return g.userThrottler.Check(nil, ctx)
	}

	cancelFunc, available := g.guildThrottler.check(discord.Snowflake(ctx.GuildID))
	if cancelFunc == nil {
		return nil, genError(available, guildThrottledErrorSecond, guildThrottledErrorMinute)
	}

	return cancelFunc, nil
}
