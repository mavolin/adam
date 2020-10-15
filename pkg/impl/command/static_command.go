package command

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// StaticCommand is a simple plugin.Command, that responds with the same
// content every time it is called.
type StaticCommand struct {
	// Meta contains the meta information about the command.
	Meta plugin.CommandMeta

	bot.MiddlewareManager

	// Reply is the reply that will be sent when called.
	// All types accepted by plugin.Command.Invoke as 1st return type are
	// allowed.
	Reply interface{}
}

var _ plugin.Command = new(StaticCommand)

func (c *StaticCommand) GetName() string      { return c.Meta.GetName() }
func (c *StaticCommand) GetAliases() []string { return c.Meta.GetAliases() }

func (c *StaticCommand) GetShortDescription(l *i18n.Localizer) string {
	return c.Meta.GetShortDescription(l)
}

func (c *StaticCommand) GetLongDescription(l *i18n.Localizer) string {
	return c.Meta.GetLongDescription(l)
}

func (c *StaticCommand) GetExamples(l *i18n.Localizer) []string  { return c.Meta.GetExamples(l) }
func (c *StaticCommand) GetArgs() plugin.ArgConfig               { return c.Meta.GetArgs() }
func (c *StaticCommand) IsHidden() bool                          { return c.Meta.IsHidden() }
func (c *StaticCommand) GetChannelTypes() plugin.ChannelTypes    { return c.Meta.GetChannelTypes() }
func (c *StaticCommand) GetBotPermissions() *discord.Permissions { return c.Meta.GetBotPermissions() }

func (c *StaticCommand) GetRestrictionFunc() plugin.RestrictionFunc {
	return c.Meta.GetRestrictionFunc()
}
func (c *StaticCommand) GetThrottler() plugin.Throttler { return c.Meta.GetThrottler() }

func (c *StaticCommand) Invoke(*state.State, *plugin.Context) (interface{}, error) {
	return c.Reply, nil
}
