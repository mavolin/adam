package resolved

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type mockCommand struct {
	name            string
	aliases         []string
	shortDesc       string
	longDesc        string
	args            plugin.ArgConfig
	argParser       plugin.ArgParser
	exampleArgs     plugin.ExampleArgs
	hidden          bool
	channelTypes    plugin.ChannelTypes
	botPermissions  discord.Permissions
	restrictionFunc plugin.RestrictionFunc
	throttler       plugin.Throttler
	invokeFunc      func(*state.State, *plugin.Context) (interface{}, error)
}

func (c mockCommand) GetName() string                                   { return c.name }
func (c mockCommand) GetAliases() []string                              { return c.aliases }
func (c mockCommand) GetShortDescription(*i18n.Localizer) string        { return c.shortDesc }
func (c mockCommand) GetLongDescription(*i18n.Localizer) string         { return c.longDesc }
func (c mockCommand) GetArgs() plugin.ArgConfig                         { return c.args }
func (c mockCommand) GetArgParser() plugin.ArgParser                    { return c.argParser }
func (c mockCommand) GetExampleArgs(*i18n.Localizer) plugin.ExampleArgs { return c.exampleArgs }
func (c mockCommand) IsHidden() bool                                    { return c.hidden }
func (c mockCommand) GetChannelTypes() plugin.ChannelTypes              { return c.channelTypes }
func (c mockCommand) GetBotPermissions() discord.Permissions            { return c.botPermissions }

func (c mockCommand) IsRestricted(s *state.State, ctx *plugin.Context) error {
	if c.restrictionFunc == nil {
		return nil
	}

	return c.restrictionFunc(s, ctx)
}

func (c mockCommand) GetThrottler() plugin.Throttler {
	return c.throttler
}

func (c mockCommand) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	return c.invokeFunc(s, ctx)
}

type mockModule struct {
	name      string
	shortDesc string
	longDesc  string
	commands  []plugin.Command
	modules   []plugin.Module
}

var _ plugin.Module = mockModule{}

func (m mockModule) GetName() string                            { return m.name }
func (m mockModule) GetShortDescription(*i18n.Localizer) string { return m.shortDesc }
func (m mockModule) GetLongDescription(*i18n.Localizer) string  { return m.longDesc }
func (m mockModule) Commands() []plugin.Command                 { return m.commands }
func (m mockModule) Modules() []plugin.Module                   { return m.modules }

type mockRestrictionErrorWrapper struct {
	WrapReturn error
}

func (m *mockRestrictionErrorWrapper) Wrap(*state.State, *plugin.Context) error {
	return m.WrapReturn
}

func (m *mockRestrictionErrorWrapper) Error() string {
	return "mockRestrictionErrorWrapper"
}

type mockThrottler struct {
	cmp string // used to make throttlers unique
}

func (m mockThrottler) Check(*state.State, *plugin.Context) (func(), error) { return func() {}, nil }

func newProviderFromSources(sources []plugin.Source) *PluginProvider {
	r := NewPluginResolver(nil)

	for _, source := range sources {
		source := source
		if source.Name == plugin.BuiltInSource {
			for _, cmd := range source.Commands {
				r.AddBuiltInCommand(cmd)
			}

			for _, mod := range source.Modules {
				r.AddBuiltInModule(mod)
			}
		} else {
			r.AddSource(source.Name,
				func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error) {
					return source.Commands, source.Modules, nil
				})
		}
	}

	return r.NewProvider(nil, nil)
}
