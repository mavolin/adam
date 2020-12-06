package throttler

import (
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// user is a plugin.Throttler that works on a per-user basis.
type user struct {
	throttler *snowflakeThrottler
}

var _ plugin.Throttler = new(user)

// PerUser returns a new plugin.Throttler that works on a per-user basis.
// It allows at maximum the passed number of invokes in the passed duration.
func PerUser(maxInvokes uint, duration time.Duration) plugin.Throttler {
	return &user{throttler: newSnowflakeThrottler(maxInvokes, duration)}
}

func (g *user) Check(_ *state.State, ctx *plugin.Context) (func(), error) {
	cancelFunc, available := g.throttler.check(discord.Snowflake(ctx.Author.ID))

	if cancelFunc == nil {
		return nil, genError(available, userThrottledErrorSecond, userThrottledErrorMinute)
	}

	return cancelFunc, nil
}
