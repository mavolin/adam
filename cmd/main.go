package main

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/disstate/pkg/state"
	"github.com/mavolin/logstract/pkg/logstract"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

func main() {
	logstract.Logger = func(lvl logstract.Lvl, msg string, fields logstract.Fields) {
		if fields != nil {
			msg += " (" + joinFields(fields, ": ", ", ") + ")"
		}

		fmt.Println(msg)
	}

	s, _ := state.New("Bot Njc2NTMxNDQ3ODc1NjMzMTY3.XkHDAg.USGIx_exATqp3aS1Q-oXU-C9aWA")

	ctx := plugin.NewContext(s)
	ctx.MessageCreateEvent = &state.MessageCreateEvent{
		MessageCreateEvent: &gateway.MessageCreateEvent{
			Message: discord.Message{
				GuildID:   691672327154303056,
				ChannelID: 691672327154303056,
				Author: discord.User{
					ID: 256827968133791744,
				},
			},
			Member: &discord.Member{
				RoleIDs: []discord.RoleID{279690191076327436},
			},
		},
	}
	ctx.DiscordDataProvider = provider{
		guild: discord.Guild{
			ID: 691672327154303056,
			Roles: []discord.Role{
				{
					ID:   415268185911328768,
					Name: "18+",
				},
				{
					ID:   279690191076327436,
					Name: "Kaka",
				},
			},
			OwnerID: 123,
		},
		channel: discord.Channel{
			NSFW: false,
		},
	}
	ctx.CommandIdentifier = ".mod.ban"
	ctx.BotOwnerIDs = []discord.UserID{}
	ctx.Localizer = localization.NewManager(nil).Localizer("de_DE")

	f := restriction.Any(
		restriction.UserPermissions(discord.PermissionAdministrator),
		restriction.All(
			restriction.AllRoles(415268185911328768),
			restriction.BotOwner,
		),
	)

	err := f(s, ctx)

	// err = err.(*restriction.EmbeddableError).Wrap(s, ctx)

	err.(plugin.RestrictionErrorWrapper).Wrap(s, ctx).(errors.Handler).Handle(s, ctx)
}

func joinFields(f logstract.Fields, keyValSep, entrySep string) string {
	n := (len(f) - 1) * (len(keyValSep) + len(entrySep))

	for k, v := range f {
		n += len(k) + len(fmt.Sprintf("%+v", v))
	}

	var b strings.Builder

	b.Grow(n)

	first := true

	for k, v := range f {
		if !first {
			b.WriteString(entrySep)
		}

		b.WriteString(k)
		b.WriteString(keyValSep)
		b.WriteString(fmt.Sprintf("%+v", v))

		first = false
	}

	return b.String()
}

type provider struct {
	guild   discord.Guild
	channel discord.Channel
}

func (p provider) Channel() (*discord.Channel, error) {
	return &p.channel, nil
}

func (p provider) Guild() (*discord.Guild, error) {
	return &p.guild, nil
}

func (p provider) Self() (*discord.Member, error) {
	panic("implement me")
}
