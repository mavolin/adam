// Package say provides the say command.
package say

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/plugin"
)

// Say is the say command.
type Say struct {
	command.Meta
}

var _ plugin.Command = new(Say) // compile-time check.

// New creates a new Say command.
func New() *Say {
	return &Say{
		Meta: command.Meta{
			Name:             "say",
			Aliases:          []string{"repeat"},
			ShortDescription: "Repeats what you say.",
			ExampleArgs:      plugin.ExampleArgs{{Args: []string{"Hello"}}},
			ArgParser:        arg.RawParser,
			Args: &arg.Config{
				RequiredArgs: []arg.RequiredArg{
					{
						Name:        "Text",
						Description: "The text you want me to say.",
						Type:        arg.SimpleText,
					},
				},
			},
			ChannelTypes:   plugin.AllChannels,
			BotPermissions: discord.PermissionSendMessages,
		},
	}
}

func (s *Say) Invoke(_ *state.State, ctx *plugin.Context) (interface{}, error) {
	return ctx.Args.String(0), nil
}
