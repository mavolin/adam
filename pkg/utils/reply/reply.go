// Package reply provides utilities for handling replies
package reply

import (
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// Await awaits a reply using the default waiter.
func Await(
	s *state.State, ctx *plugin.Context, initialTimeout, typingTimeout time.Duration,
) (*discord.Message, error) {
	return NewDefaultWaiter(s, ctx).
		Await(initialTimeout, typingTimeout)
}
