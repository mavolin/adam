// Package messageutil provides utilities for awaiting replies and reactions.
package messageutil

import (
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// AwaitReply awaits a reply using the default waiter.
func AwaitReply(
	s *state.State, ctx *plugin.Context, initialTimeout, typingTimeout time.Duration,
) (*discord.Message, error) {
	return NewReplyWaiterFromDefault(s, ctx).
		Await(initialTimeout, typingTimeout)
}
