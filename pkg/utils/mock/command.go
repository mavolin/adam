package mock

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type Command struct {
	plugin.CommandMeta
	InvokeFunc func(*state.State, *plugin.Context) (interface{}, error)
}

var _ plugin.Command = Command{}

func (c Command) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	return c.InvokeFunc(s, ctx)
}

// GenerateRegisteredCommand creates a mocked RegisteredCommand from the passed
// Command using the passed provider name.
func GenerateRegisteredCommand(providerName string, cmd Command) *plugin.RegisteredCommand {
	c := plugin.NewRegisteredCommandWithParent(nil)

	c.Source = cmd
	// c.SourceParents = nil
	c.ProviderName = providerName
	c.ID = plugin.ID("." + cmd.GetName())
	c.Name = cmd.GetName()
	c.Aliases = cmd.GetAliases()
	c.Args = cmd.GetArgs()
	c.Hidden = cmd.IsHidden()
	c.ChannelTypes = cmd.GetChannelTypes()
	c.BotPermissions = cmd.GetBotPermissions()
	c.Throttler = cmd.GetThrottler()

	return c
}

// GenerateRegisteredCommandWithParents creates a new RegisteredCommand from
// the passed module.
// It then returns the command with the given identifier found in the module.
//
// The passed module must be the root module.
func GenerateRegisteredCommandWithParents(
	providerName string, smod plugin.Module, cmdID plugin.ID,
) *plugin.RegisteredCommand {
	rmod := GenerateRegisteredModule(providerName, smod)
	if rmod == nil {
		return nil
	}

	all := cmdID.All()
	if len(all) <= 1 {
		return nil
	}

	for _, id := range all[1 : len(all)-1] { // range from first module to last
		rmod = rmod.FindModule(id.Name())
		if rmod == nil {
			return nil
		}
	}

	return rmod.FindCommand(cmdID.Name())
}

type CommandMeta struct {
	Name             string
	Aliases          []string
	ShortDescription string
	LongDescription  string

	Args plugin.ArgConfig

	ExampleArgs    []string
	Hidden         bool
	ChannelTypes   plugin.ChannelTypes
	BotPermissions discord.Permissions
	Restrictions   plugin.RestrictionFunc
	Throttler      plugin.Throttler
}

var _ plugin.CommandMeta = CommandMeta{}

func (m CommandMeta) GetName() string                            { return m.Name }
func (m CommandMeta) GetAliases() []string                       { return m.Aliases }
func (m CommandMeta) GetShortDescription(*i18n.Localizer) string { return m.ShortDescription }
func (m CommandMeta) GetLongDescription(*i18n.Localizer) string  { return m.LongDescription }
func (m CommandMeta) GetArgs() plugin.ArgConfig                  { return m.Args }
func (m CommandMeta) GetExampleArgs(*i18n.Localizer) []string    { return m.ExampleArgs }
func (m CommandMeta) IsHidden() bool                             { return m.Hidden }
func (m CommandMeta) GetChannelTypes() plugin.ChannelTypes       { return m.ChannelTypes }
func (m CommandMeta) GetBotPermissions() discord.Permissions     { return m.BotPermissions }

func (m CommandMeta) IsRestricted(s *state.State, ctx *plugin.Context) error {
	if m.Restrictions == nil {
		return nil
	}

	return m.Restrictions(s, ctx)
}

func (m CommandMeta) GetThrottler() plugin.Throttler { return m.Throttler }

type ArgConfig struct {
	Expect string

	ArgCombinationID string
	Args             plugin.Args
	Flags            plugin.Flags
	ErrorReturn      error
}

var _ plugin.ArgConfig = ArgConfig{}

func (a ArgConfig) Parse(_ string, _ *state.State, ctx *plugin.Context) error {
	ctx.ArgCombinationID = a.ArgCombinationID
	ctx.Args = a.Args
	ctx.Flags = a.Flags

	return a.ErrorReturn
}
