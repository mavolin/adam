// Package ban provides the ban command.
package ban

import (
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
	command.LocalizedMeta
}

var _ plugin.Command = new(Ban) // compile-time check.

// New creates a new Ban command.
func New() *Ban {
	return &Ban{
		LocalizedMeta: command.LocalizedMeta{
			Name:             "ban",
			Aliases:          []string{"banhammer"},
			ShortDescription: shortDescription,
			ExampleArgs:      examples,
			Args: &arg.LocalizedConfig{
				RequiredArgs: []arg.LocalizedRequiredArg{
					{
						Name:        argMemberName,
						Type:        arg.Member,
						Description: argMemberDescription,
					},
				},
				OptionalArgs: []arg.LocalizedOptionalArg{
					{
						Name:        argReasonName,
						Type:        arg.SimpleText,
						Description: argReasonDescription,
					},
				},
				Flags: []arg.LocalizedFlag{
					{
						Name:        "days",
						Aliases:     []string{"d"},
						Type:        arg.IntegerWithBounds(0, 7),
						Default:     1,
						Description: flagDaysDescription,
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
		return nil, errors.NewUserErrorl(selfBanError)
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

	return success.
		WithPlaceholders(successPlaceholders{
			Username: m.User.Username,
		}), nil
}
