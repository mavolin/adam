// Package kick provides the kick command.
package kick

import (
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/impl/throttler"
	"github.com/mavolin/adam/pkg/plugin"
)

// Kick is the kick command.
type Kick struct {
	command.Meta
}

// New creates a new Kick command.
func New() *Kick {
	return &Kick{
		Meta: command.Meta{
			Name:             "kick",
			Aliases:          nil,
			ShortDescription: "Kicks a user.",
			ExampleArgs: plugin.ExampleArgs{
				{Args: []string{"@Clyde", "@Clyde, self-botting"}},
			},
			Args: &arg.Config{
				RequiredArgs: []arg.RequiredArg{
					{
						Name:        "Member",
						Type:        arg.Member,
						Description: "The member you want to kick.",
					},
				},
			},
			Hidden:         false,
			ChannelTypes:   plugin.GuildChannels,
			BotPermissions: discord.PermissionSendMessages | discord.PermissionKickMembers,
			Restrictions:   restriction.UserPermissions(discord.PermissionManageGuild),
			// up to 25 people every four minutes
			Throttler: throttler.PerGuild(25, 4*time.Minute),
		},
	}
}

func (k *Kick) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	m := ctx.Args.Member(0)
	if m.User.ID == ctx.Author.ID {
		return nil, errors.NewUserError("You can't kick yourself.")
	}

	if err := s.Kick(ctx.GuildID, m.User.ID); err != nil {
		return nil, err
	}

	return fmt.Sprintf("👮 %s has been kicked!", m.User.Username), nil
}
