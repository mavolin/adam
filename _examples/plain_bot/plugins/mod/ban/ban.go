// Package ban provides the ban command.
package ban

import (
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/impl/throttler"
	"github.com/mavolin/adam/pkg/plugin"
)

// Ban is the ban command.
type Ban struct {
	command.Meta
}

var _ plugin.Command = new(Ban) // compile-time check.

// New creates a new Ban command.
func New() *Ban {
	return &Ban{
		Meta: command.Meta{
			Name:             "ban",
			Aliases:          []string{"banhammer"},
			ShortDescription: "Bans someone.",
			ExampleArgs: plugin.ExampleArgs{
				{Args: []string{"@Wumpus", "@Wumpus, using offensive language"}},
			},
			Args: &arg.Config{
				RequiredArgs: []arg.RequiredArg{
					{
						Name:        "Member",
						Type:        arg.Member,
						Description: "The member you want to ban.",
					},
				},
				OptionalArgs: []arg.OptionalArg{
					{
						Name:        "Reason",
						Type:        arg.SimpleText,
						Description: "The reason for the ban.",
					},
				},
				Flags: []arg.Flag{
					{
						Name:        "days",
						Aliases:     []string{"d"},
						Type:        arg.IntegerWithBounds(0, 7),
						Default:     1,
						Description: "The amount of days to delete messages for. You can delete 7 days at most.",
					},
				},
			},
			ChannelTypes:   plugin.GuildChannels,
			BotPermissions: discord.PermissionSendMessages | discord.PermissionBanMembers,
			// moderators typically have this permissions, so use this to
			// validate authorization
			Restrictions: restriction.UserPermissions(discord.PermissionManageGuild),
			// up to 25 people every four minutes
			Throttler: throttler.PerGuild(25, 4*time.Minute),
		},
	}
}

func (b *Ban) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	m := ctx.Args.Member(0)
	if m.User.ID == ctx.Author.ID {
		return nil, errors.NewUserError("Good try, but you can ban yourself.")
	}

	banData := api.BanData{
		DeleteDays: option.NewUint(uint(ctx.Flags.Int("days"))),
	}
	if reason := ctx.Args.String(1); len(reason) > 0 {
		banData.Reason = option.NewString(reason)
	}

	if err := s.Ban(ctx.GuildID, m.User.ID, banData); err != nil {
		return nil, err
	}

	return fmt.Sprintf("👮 The banhammer has been slayed, and %s is no more!", m.User.Username), nil
}
