// Package ping provides the ping command.
package ping

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/plugin"
)

// Ping is the ping command.
type Ping struct {
	command.Meta
}

var _ plugin.Command = new(Ping) // compile-time check

// New creates a new Ping command.
func New() *Ping {
	return &Ping{
		Meta: command.Meta{
			Name:             "ping",
			ShortDescription: "Tells you the ping to Discord.",
			// we can leave this empty, if we want to use the ShortDescription
			// as LongDescription as well
			LongDescription: "",
			ChannelTypes:    plugin.AllChannels,
			BotPermissions:  discord.PermissionSendMessages,
		},
	}
}

func (p *Ping) Invoke(*state.State, *plugin.Context) (interface{}, error) {
	return "Pong!", nil
}
