package mock

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

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
	c := plugin.NewRegisteredCommandWithParent(nil, cmd.GetRestrictionFunc())

	c.Source = cmd
	// c.SourceParents = nil
	c.ProviderName = providerName
	c.Identifier = plugin.Identifier("." + cmd.GetName())
	c.Name = cmd.GetName()
	c.Aliases = cmd.GetAliases()
	c.Args = cmd.GetArgs()
	c.Hidden = cmd.IsHidden()
	c.ChannelTypes = cmd.GetChannelTypes()

	if perms := cmd.GetBotPermissions(); perms != nil {
		c.BotPermissions = *perms
	}

	c.Throttler = cmd.GetThrottler()

	return c
}

// GenerateRegisteredCommandWithParents creates a new RegisteredCommand from
// the passed module.
// It then returns the command with the given identifier found in the module.
//
// The passed module must be the root module.
func GenerateRegisteredCommandWithParents(
	providerName string, smod plugin.Module, cmdID plugin.Identifier,
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

	Args ArgConfig

	Examples       []string
	Hidden         bool
	ChannelTypes   plugin.ChannelTypes
	BotPermissions *discord.Permissions
	Restrictions   plugin.RestrictionFunc
	Throttler      plugin.Throttler
}

var _ plugin.CommandMeta = CommandMeta{}

func (c CommandMeta) GetName() string                            { return c.Name }
func (c CommandMeta) GetAliases() []string                       { return c.Aliases }
func (c CommandMeta) GetShortDescription(*i18n.Localizer) string { return c.ShortDescription }
func (c CommandMeta) GetLongDescription(*i18n.Localizer) string  { return c.LongDescription }
func (c CommandMeta) GetArgs() plugin.ArgConfig                  { return c.Args }
func (c CommandMeta) GetExamples(*i18n.Localizer) []string       { return c.Examples }
func (c CommandMeta) IsHidden() bool                             { return c.Hidden }
func (c CommandMeta) GetChannelTypes() plugin.ChannelTypes       { return c.ChannelTypes }
func (c CommandMeta) GetBotPermissions() *discord.Permissions    { return c.BotPermissions }
func (c CommandMeta) GetRestrictionFunc() plugin.RestrictionFunc { return c.Restrictions }
func (c CommandMeta) GetThrottler() plugin.Throttler             { return c.Throttler }

type ArgConfig struct {
	Expect string

	ArgsReturn  plugin.Args
	FlagsReturn plugin.Flags
	ErrorReturn error
}

var _ plugin.ArgConfig = ArgConfig{}

func (a ArgConfig) Parse(_ string, _ *state.State, _ *plugin.Context) (plugin.Args, plugin.Flags, error) {
	return a.ArgsReturn, a.FlagsReturn, a.ErrorReturn
}
