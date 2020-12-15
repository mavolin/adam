package throttler

import (
	"sort"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
)

type snowflakeThrottler struct {
	throttled  map[discord.Snowflake][]time.Time
	maxInvokes uint
	duration   time.Duration

	mutex sync.Mutex
}

// newSnowflakeThrottler creates a new throttler with the passed inclusive
// maxInvokes and the passed expiration duration that works on a per-snowflake
// basis.
func newSnowflakeThrottler(maxInvokes uint, duration time.Duration) *snowflakeThrottler {
	return &snowflakeThrottler{
		throttled:  make(map[discord.Snowflake][]time.Time),
		maxInvokes: maxInvokes,
		duration:   duration,
	}
}

// expire removes all invokes for the passed discord.Snowflake before the
// passed time.Time.
//
// It will not lock the mutex.
func (t *snowflakeThrottler) expire(s discord.Snowflake, before time.Time) {
	throttled := t.throttled[s]

	i := sort.Search(len(throttled), func(i int) bool {
		return !throttled[i].Before(before)
	})

	t.throttled[s] = throttled[i:]
}

// now is a function that returns the current time.
// Made replaceable for testing.
var now = time.Now

// check checks if the entity with the passed snowflake should be throttled.
// If so check will return nil and the duration until the command can be
// invoked again.
// Otherwise, check will return a cancelFunc and add the invoke.
func (t *snowflakeThrottler) check(s discord.Snowflake) (func(), time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	now := now()
	t.expire(s, now.Add(-t.duration))

	throttled := t.throttled[s]

	// maxInvokes reached, prevent invoke
	if len(throttled) >= int(t.maxInvokes) {
		return nil, t.duration - now.Sub(throttled[0])
	}

	t.throttled[s] = append(throttled, now)

	return func() {
		t.mutex.Lock()
		defer t.mutex.Unlock()

		throttled := t.throttled[s]
		for i, cmp := range throttled {
			if cmp == now {
				t.throttled[s] = append(throttled[:i], throttled[i+1:]...)
				return
			}
		}
	}, 0
}
