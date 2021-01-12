// Package command provides implementations for the command abstractions found
// in package plugin.
package command

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/plugin"
)

// StaticCommand is a simple plugin.Command, that responds with the same
// content every time it is called.
type StaticCommand struct {
	plugin.CommandMeta
	*bot.MiddlewareManager

	reply interface{}
}

// NewStaticCommand creates a new *StaticCommand using the passed reply and
// plugin.CommandMeta.
//
// Reply may be of any type supported as first return value by
// plugin.Command.Invoke.
func NewStaticCommand(reply interface{}, meta plugin.CommandMeta) *StaticCommand {
	return &StaticCommand{
		CommandMeta:       meta,
		MiddlewareManager: new(bot.MiddlewareManager),
		reply:             reply,
	}
}

var _ plugin.Command = new(StaticCommand)

func (c *StaticCommand) Invoke(*state.State, *plugin.Context) (interface{}, error) {
	return c.reply, nil
}
