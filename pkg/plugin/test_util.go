package plugin

import (
	"errors"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/i18n"
)

// mockLocalizer is a copy of mock.Localizer, used to prevent import cycles.
type mockLocalizer struct {
	t *testing.T

	def      string
	onReturn map[i18n.Term]string
	errOn    map[i18n.Term]struct{}
}

func newMockedLocalizer(t *testing.T) *mockLocalizer {
	return &mockLocalizer{
		t:        t,
		onReturn: make(map[i18n.Term]string),
		errOn:    make(map[i18n.Term]struct{}),
	}
}

func (l *mockLocalizer) on(term i18n.Term, response string) *mockLocalizer {
	l.onReturn[term] = response
	return l
}

func (l *mockLocalizer) build() *i18n.Localizer {
	m := i18n.NewManager(func(lang string) i18n.LangFunc {
		return func(term i18n.Term, _ map[string]interface{}, _ interface{}) (string, error) {
			r, ok := l.onReturn[term]
			if ok {
				return r, nil
			}

			_, ok = l.errOn[term]
			if ok {
				return r, errors.New("error")
			}

			if l.def == "" {
				assert.Failf(l.t, "unexpected call to Localize", "unknown term %s", term)

				return string(term), errors.New("unknown term")
			}

			return l.def, nil
		}
	})

	return m.Localizer("")
}

type mockCommand struct {
	name           string
	aliases        []string
	args           ArgConfig
	shortDesc      string
	longDesc       string
	examples       []string
	hidden         bool
	channelTypes   ChannelTypes
	botPermissions *discord.Permissions
	restrictions   RestrictionFunc
	throttler      Throttler
	invokeFunc     func(*state.State, *Context) (interface{}, error)
}

func (c mockCommand) GetName() string                            { return c.name }
func (c mockCommand) GetAliases() []string                       { return c.aliases }
func (c mockCommand) GetArgs() ArgConfig                         { return c.args }
func (c mockCommand) GetShortDescription(*i18n.Localizer) string { return c.shortDesc }
func (c mockCommand) GetLongDescription(*i18n.Localizer) string  { return c.longDesc }
func (c mockCommand) GetExamples(*i18n.Localizer) []string       { return c.examples }
func (c mockCommand) IsHidden() bool                             { return c.hidden }
func (c mockCommand) GetChannelTypes() ChannelTypes              { return c.channelTypes }
func (c mockCommand) GetBotPermissions() *discord.Permissions    { return c.botPermissions }
func (c mockCommand) GetRestrictionFunc() RestrictionFunc        { return c.restrictions }
func (c mockCommand) GetThrottler() Throttler                    { return c.throttler }

func (c mockCommand) Invoke(s *state.State, ctx *Context) (interface{}, error) {
	return c.invokeFunc(s, ctx)
}

type mockModule struct {
	name                  string
	shortDesc             string
	longDesc              string
	Hidden                bool
	defaultChannelTypes   ChannelTypes
	defaultBotPermissions *discord.Permissions
	defaultRestrictions   RestrictionFunc
	defaultThrottler      Throttler
	commands              []Command
	modules               []Module
}

func (c mockModule) GetName() string                                { return c.name }
func (c mockModule) GetShortDescription(*i18n.Localizer) string     { return c.shortDesc }
func (c mockModule) GetLongDescription(*i18n.Localizer) string      { return c.longDesc }
func (c mockModule) IsHidden() bool                                 { return c.Hidden }
func (c mockModule) GetDefaultChannelTypes() ChannelTypes           { return c.defaultChannelTypes }
func (c mockModule) GetDefaultBotPermissions() *discord.Permissions { return c.defaultBotPermissions }
func (c mockModule) GetDefaultRestrictionFunc() RestrictionFunc     { return c.defaultRestrictions }
func (c mockModule) GetDefaultThrottler() Throttler                 { return c.defaultThrottler }
func (c mockModule) Commands() []Command                            { return c.commands }
func (c mockModule) Modules() []Module                              { return c.modules }

type mockThrottler struct {
	cmp string // used to make throttlers unique
}

func (m mockThrottler) Check(*Context) (func(), error) { return func() {}, nil }

// mockDiscordDataProvider is a copy of mock.DiscordDataProvider to prevent
// import cycles.
type mockDiscordDataProvider struct {
	ChannelReturn *discord.Channel
	ChannelError  error

	GuildReturn *discord.Guild
	GuildError  error

	SelfReturn *discord.Member
	SelfError  error
}

func (d mockDiscordDataProvider) Channel() (*discord.Channel, error) {
	return d.ChannelReturn, d.ChannelError
}

func (d mockDiscordDataProvider) Guild() (*discord.Guild, error) {
	return d.GuildReturn, d.GuildError
}

func (d mockDiscordDataProvider) Self() (*discord.Member, error) {
	return d.SelfReturn, d.SelfError
}

// removeRegisteredModuleFuncs sets all functions stored in the passed
// RegisteredModule to nil.
// Additionally, it does the same for all submodules and subcommands
// recursively.
func removeRegisteredModuleFuncs(mod *RegisteredModule) {
	for i := range mod.Commands {
		removeRegisteredCommandFuncs(mod.Commands[i])
	}

	for i := range mod.Modules {
		removeRegisteredModuleFuncs(mod.Modules[i])
	}
}

// removeRegisteredCommandFuncs sets all functions stored in the passed
// RegisteredCommand to nil.
func removeRegisteredCommandFuncs(cmd *RegisteredCommand) {
	cmd.restrictionFunc = nil
}
