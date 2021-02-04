// Package say provides the say command.
package say

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/plugin"
)

// Say is the say command.
type Say struct {
	command.LocalizedMeta
}

var _ plugin.Command = new(Say) // compile-time check.

// New creates a new Say command
func New() *Say {
	return &Say{
		LocalizedMeta: command.LocalizedMeta{
			Name:             "say",
			Aliases:          []string{"repeat"},
			ShortDescription: shortDescription,
			ExampleArgs:      examples,
			Args: arg.LocalizedRaw{
				Name:        argTextName,
				Description: argTextDescription,
			},
			ChannelTypes:   plugin.AllChannels,
			BotPermissions: discord.PermissionSendMessages,
		},
	}
}

func (s *Say) Invoke(_ *state.State, ctx *plugin.Context) (interface{}, error) {
	return ctx.Args.String(0), nil
}
