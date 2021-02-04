// Package kick provides the kick command.
package kick

import (
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
	command.LocalizedMeta
}

// New creates a new Kick command.
func New() *Kick {
	return &Kick{
		LocalizedMeta: command.LocalizedMeta{
			Name:             "kick",
			Aliases:          nil,
			ShortDescription: shortDescription,
			ExampleArgs:      examples,
			Args: arg.LocalizedCommaConfig{
				Required: []arg.LocalizedRequiredArg{
					{
						Name:        argMemberName,
						Type:        arg.Member,
						Description: argMemberDescription,
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
		return nil, errors.NewUserErrorl(selfKickError)
	}

	if err := s.Kick(ctx.GuildID, m.User.ID); err != nil {
		return nil, err
	}

	return success.
		WithPlaceholders(successPlaceholders{
			Username: m.User.Username,
		}), nil
}
