package bot

import (
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// newMessageCreateEvent creates a new state.MessageCreateEvent from the passed
// plugin.Context.
func newMessageCreateEvent(ctx *plugin.Context) *state.MessageCreateEvent {
	return &state.MessageCreateEvent{
		MessageCreateEvent: &gateway.MessageCreateEvent{
			Message: ctx.Message,
			Member:  ctx.Member,
		},
		Base: ctx.Base,
	}
}

// newMessageUpdateEvent creates a new state.MessageUpdateEvent from the passed
// plugin.Context.
func newMessageUpdateEvent(ctx *plugin.Context) *state.MessageUpdateEvent {
	return &state.MessageUpdateEvent{
		MessageUpdateEvent: &gateway.MessageUpdateEvent{
			Message: ctx.Message,
			Member:  ctx.Member,
		},
		Base: ctx.Base,
	}
}
