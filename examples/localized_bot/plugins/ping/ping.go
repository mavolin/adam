// Package ping provides the ping command.
package ping

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/plugin"
)

// Ping is the ping command.
type Ping struct {
	command.LocalizedMeta
}

var _ plugin.Command = new(Ping) // compile-time check

// New creates a new Ping command.
func New() *Ping {
	return &Ping{
		LocalizedMeta: command.LocalizedMeta{
			Name:             "ping",
			ShortDescription: shortDescription,
			// we can leave this empty, if we want to use the ShortDescription
			// as LongDescription as well
			LongDescription: nil,
			BotPermissions:  discord.PermissionSendMessages,
		},
	}
}

func (p *Ping) Invoke(*state.State, *plugin.Context) (interface{}, error) {
	// "Pong" is interlingual, so we don't need to localize it
	return "Pong!", nil
}
