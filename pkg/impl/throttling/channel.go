package throttling

import (
	"time"

	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/plugin"
)

// channel is a plugin.Throttler that works on a per-channel basis.
type channel struct {
	throttler *snowflakeThrottler
}

// PerUser returns a new plugin.Throttler that works on a per-channel basis.
// It allows at maximum the passed number of invokes in the passed duration.
func PerChannel(maxInvokes uint, duration time.Duration) plugin.Throttler {
	return &channel{
		throttler: newSnowflakeThrottler(maxInvokes, duration),
	}
}

func (g *channel) Check(ctx *plugin.Context) (func(), error) {
	cancelFunc, available := g.throttler.check(discord.Snowflake(ctx.ChannelID))

	if cancelFunc == nil {
		return nil, genError(available, channelErrorSecond, channelErrorMinute)
	}

	return cancelFunc, nil
}
