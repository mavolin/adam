package throttler

import (
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// member is a plugin.Throttler that works on a per-member basis.
type member struct {
	memberThrottler map[discord.GuildID]*snowflakeThrottler
	userThrottler   plugin.Throttler

	memberMutex sync.Mutex

	maxInvokes uint
	duration   time.Duration
}

var _ plugin.Throttler = new(member)

// PerMember returns a new plugin.Throttler that works on a per-member basis.
// It allows at maximum the passed number of invokes in the passed duration.
// Effectively, this is the same as PerUser but filtered additionally by guild.
//
// All commands invoked in direct messages will be throttled on a per-user
// basis.
func PerMember(maxInvokes uint, duration time.Duration) plugin.Throttler {
	return &member{
		memberThrottler: make(map[discord.GuildID]*snowflakeThrottler),
		userThrottler:   PerUser(maxInvokes, duration),
		maxInvokes:      maxInvokes,
		duration:        duration,
	}
}

func (g *member) Check(s *state.State, ctx *plugin.Context) (func(), error) {
	if ctx.GuildID == 0 {
		return g.userThrottler.Check(s, ctx)
	}

	g.memberMutex.Lock()
	defer g.memberMutex.Unlock()

	gt := g.memberThrottler[ctx.GuildID]
	if gt == nil {
		gt = newSnowflakeThrottler(g.maxInvokes, g.duration)
	}

	cancelFunc, available := gt.check(discord.Snowflake(ctx.Author.ID))
	g.memberThrottler[ctx.GuildID] = gt

	if cancelFunc == nil {
		return nil, genError(available, memberThrottledErrorSecond, memberThrottledErrorMinute)
	}

	g.memberThrottler[ctx.GuildID] = gt

	return cancelFunc, nil
}
