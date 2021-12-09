package embed

import (
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/msgbuilder"
)

type Embed struct {
	command.Meta
}

var _ plugin.Command = new(Embed)

func New() *Embed {
	return &Embed{
		Meta: command.Meta{
			Name:             "embed",
			ShortDescription: "Create an embed.",
			LongDescription:  "Create a custom embed from the input you give me.",
			BotPermissions:   discord.PermissionSendMessages,
		},
	}
}

func (e *Embed) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	var cancelled bool
	var title, desc discord.Message

	_, err := msgbuilder.New(s, ctx).
		WithContent("What should the title of the embed be?").
		WithAwaitedResponse(&title, 10*time.Second, 5*time.Second).
		WithComponent(msgbuilder.NewActionRow(&cancelled).
			With(msgbuilder.NewButton(discord.DangerButton, "Cancel", true))).
		ReplyAndAwait(20 * time.Second)
	if err != nil || cancelled {
		return nil, err
	}

	_, err = msgbuilder.New(s, ctx).
		WithContent("What should the description of the embed be?").
		WithAwaitedResponse(&desc, 10*time.Second, 5*time.Second).
		WithComponent(msgbuilder.NewActionRow(&cancelled).
			With(msgbuilder.NewButton(discord.DangerButton, "Cancel", true))).
		ReplyAndAwait(20 * time.Second)
	if err != nil || cancelled {
		return nil, err
	}

	return discord.Embed{Title: title.Content, Description: desc.Content}, err
}
