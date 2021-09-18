package embed

import (
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/examples/localized_bot/internal/term"
	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/msgbuilder"
)

type Embed struct {
	command.LocalizedMeta
}

var _ plugin.Command = new(Embed)

func New() *Embed {
	return &Embed{
		LocalizedMeta: command.LocalizedMeta{
			Name:             "embed",
			Aliases:          nil,
			ShortDescription: shortDescription,
			LongDescription:  longDescription,
			BotPermissions:   discord.PermissionSendMessages,
		},
	}
}

func (e *Embed) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	var cancelled bool
	var title, desc discord.Message

	_, err := msgbuilder.New(s, ctx).
		WithContentl(titleQuestion).
		WithAwaitedResponse(&title, 10*time.Second, 5*time.Second).
		WithComponent(msgbuilder.NewActionRow(&cancelled).
			With(msgbuilder.NewButtonl(discord.DangerButton, term.Cancel, true))).
		ReplyAndAwait(20 * time.Second)
	if err != nil {
		return nil, err
	}

	_, err = msgbuilder.New(s, ctx).
		WithContentl(descriptionQuestion).
		WithAwaitedResponse(&desc, 10*time.Second, 5*time.Second).
		WithComponent(msgbuilder.NewActionRow(&cancelled).
			With(msgbuilder.NewButtonl(discord.DangerButton, term.Cancel, true))).
		ReplyAndAwait(20 * time.Second)
	if err != nil {
		return nil, err
	}

	return discord.Embed{Title: title.Content, Description: desc.Content}, err
}
