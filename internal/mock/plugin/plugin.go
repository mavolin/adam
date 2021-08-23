package plugin

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// =============================================================================
// Command
// =====================================================================================

type Command struct {
	Name    string
	Aliases []string

	ShortDescription string
	LongDescription  string

	Args        plugin.ArgConfig
	ArgParser   plugin.ArgParser
	ExampleArgs plugin.ExampleArgs

	Hidden bool

	ChannelTypes   plugin.ChannelTypes
	BotPermissions discord.Permissions
	Restrictions   plugin.RestrictionFunc
	Throttler      plugin.Throttler

	Meta       plugin.CommandMeta
	InvokeFunc func(*state.State, *plugin.Context) (interface{}, error)
}

var _ plugin.Command = Command{}

func (c Command) GetName() string                                   { return c.Name }
func (c Command) GetAliases() []string                              { return c.Aliases }
func (c Command) GetShortDescription(*i18n.Localizer) string        { return c.ShortDescription }
func (c Command) GetLongDescription(*i18n.Localizer) string         { return c.LongDescription }
func (c Command) GetArgs() plugin.ArgConfig                         { return c.Args }
func (c Command) GetArgParser() plugin.ArgParser                    { return c.ArgParser }
func (c Command) GetExampleArgs(*i18n.Localizer) plugin.ExampleArgs { return c.ExampleArgs }
func (c Command) IsHidden() bool                                    { return c.Hidden }
func (c Command) GetChannelTypes() plugin.ChannelTypes              { return c.ChannelTypes }
func (c Command) GetBotPermissions() discord.Permissions            { return c.BotPermissions }

func (c Command) IsRestricted(s *state.State, ctx *plugin.Context) error {
	if c.Restrictions == nil {
		return nil
	}

	return c.Restrictions(s, ctx)
}

func (c Command) GetThrottler() plugin.Throttler { return c.Throttler }

func (c Command) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	return c.InvokeFunc(s, ctx)
}

// =============================================================================
// Module
// =====================================================================================

type Module struct {
	Name             string
	ShortDescription string
	LongDescription  string

	Commands []plugin.Command
	Modules  []plugin.Module
}

var _ plugin.Module = Module{}

func (m Module) GetName() string                            { return m.Name }
func (m Module) GetShortDescription(*i18n.Localizer) string { return m.ShortDescription }
func (m Module) GetLongDescription(*i18n.Localizer) string  { return m.LongDescription }
func (m Module) GetCommands() []plugin.Command              { return m.Commands }
func (m Module) GetModules() []plugin.Module                { return m.Modules }

// =============================================================================
// Throttler
// =====================================================================================

// Throttler is the mocked version of a plugin.Throttler.
type Throttler struct {
	checkReturn error
	// Canceled is set to true, if the function returned by the Throttler's
	// Check method is called.
	Canceled bool
}

var _ plugin.Throttler = new(Throttler)

// NewThrottler creates a new mocked Throttler with the given return value
// for check.
func NewThrottler(checkReturn error) *Throttler {
	return &Throttler{checkReturn: checkReturn}
}

func (t *Throttler) Check(*state.State, *plugin.Context) (func(), error) {
	return func() {
		t.Canceled = true
	}, t.checkReturn
}

// =============================================================================
// Restriction
// =====================================================================================

// RestrictionFunc creates a restriction func that returns the passed error.
func RestrictionFunc(ret error) plugin.RestrictionFunc {
	return func(*state.State, *plugin.Context) error { return ret }
}

// =============================================================================
// RestrictionErrorWrapper
// =====================================================================================

type RestrictionErrorWrapper struct {
	Return error
}

func (m *RestrictionErrorWrapper) Wrap(*state.State, *plugin.Context) error {
	return m.Return
}

func (m *RestrictionErrorWrapper) Error() string {
	return "RestrictionErrorWrapper"
}
