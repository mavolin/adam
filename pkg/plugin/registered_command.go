package plugin

import (
	"sort"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

var BuiltInProvider = "built_in"

// RegisteredCommand is a resolved command as returned by a Provider.
// In contrast to the regular command abstraction, RegisteredCommand's fields
// take into account it's parents settings, as the router would see them.
type RegisteredCommand struct {
	// double pointer used to determine if the parent is just nil or not set
	parent   **RegisteredModule
	provider Provider

	// ProviderName is the name of the plugin provider that provides the
	// command.
	//
	// If the command is built-in, ProviderName will be set to 'built_in'.
	ProviderName string
	// Source is the original Command this command is based on.
	Source Command
	// SourceParents contains the original parent Modules in ascending order
	// from lowest order to the closest parent.
	//
	// If the command is top-level, SourceParents will be nil.
	SourceParents []Module

	// ID is the identifier of the command.
	ID ID
	// Name is the name of the command.
	Name string
	// Aliases contains the optional aliases of the command.
	Aliases []string

	// Args is the argument configuration of the command.
	//
	// If this is nil, the command accepts no arguments.
	Args ArgConfig

	// Hidden specifies whether to show this command in the help.
	Hidden bool
	// ChannelTypes are the ChannelTypes this command can be run in.
	//
	// If the command itself did not define some, ChannelTypes will be
	// AllChannels.
	ChannelTypes ChannelTypes
	// BotPermissions are the permissions this command needs to execute.
	BotPermissions discord.Permissions
	// Throttler is the Throttler of this command.
	Throttler Throttler
}

// NewRegisteredCommandWithParent creates a new RegisteredCommand from the
// passed parent module using the passed RestrictionFunc.
// The RestrictionFunc may be nil.
func NewRegisteredCommandWithParent(p *RegisteredModule) *RegisteredCommand {
	return &RegisteredCommand{parent: &p}
}

// NewRegisteredCommandWithProvider creates a new RegisteredCommand from the
// passed Provider using the passed RestrictionFunc.
// The RestrictionFunc may be nil.
func NewRegisteredCommandWithProvider(p Provider) *RegisteredCommand {
	return &RegisteredCommand{provider: p}
}

// GenerateRegisteredCommands generates top-level RegisteredCommands from the
// passed Repositories.
func GenerateRegisteredCommands(repos []Repository) []*RegisteredCommand { //nolint:gocognit
	var maxLen int

	for _, repo := range repos {
		maxLen += len(repo.Commands)
	}

	rcmds := make([]*RegisteredCommand, 0, maxLen)

	usedAliases := make(map[string]struct{})

	for _, repo := range repos {
		for _, scmd := range repo.Commands {
			i := sort.Search(len(rcmds), func(i int) bool {
				return rcmds[i].Name >= scmd.GetName()
			})

			if i < len(rcmds) && rcmds[i].Name == scmd.GetName() {
				continue // skip on duplicate name
			}

			var parent *RegisteredModule = nil

			rcmd := &RegisteredCommand{
				parent:         &parent,
				ProviderName:   repo.ProviderName,
				Source:         scmd,
				ID:             ID("." + scmd.GetName()),
				Name:           scmd.GetName(),
				Args:           scmd.GetArgs(),
				Hidden:         scmd.IsHidden(),
				ChannelTypes:   scmd.GetChannelTypes(),
				BotPermissions: scmd.GetBotPermissions(),
				Throttler:      scmd.GetThrottler(),
			}

			if rcmd.ChannelTypes == 0 {
				rcmd.ChannelTypes = AllChannels
			}

			if saliases := scmd.GetAliases(); len(saliases) > 0 {
				rcmd.Aliases = make([]string, 0, len(saliases))

				for _, a := range saliases { // check for duplicate aliases
					if _, ok := usedAliases[a]; !ok {
						usedAliases[a] = struct{}{}
						rcmd.Aliases = append(rcmd.Aliases, a)
					}
				}
			}

			if i == len(rcmds) {
				rcmds = append(rcmds, rcmd)
			} else {
				rcmds = append(rcmds, rcmd) // make space for a new element
				copy(rcmds[i+1:], rcmds[i:])

				rcmds[i] = rcmd
			}
		}
	}

	return rcmds
}

// Parent returns the parent of this command.
// The returned RegisteredModule may not consists of all modules that share the
// same namespace, if some plugin providers are unavailable.
// Check PluginProvider.UnavailableProviders() to check if that is the case.
//
// In any case the module will contain the built-in module and the module that
// provides the command.
func (c *RegisteredCommand) Parent() *RegisteredModule {
	if c.parent != nil {
		return *c.parent
	}

	parent := c.provider.Module(c.ID.Parent())
	c.parent = &parent

	return parent
}

// ShortDescription returns an optional brief description of the command.
func (c *RegisteredCommand) ShortDescription(l *i18n.Localizer) string {
	return c.Source.GetShortDescription(l)
}

// LongDescription returns an optional long description of the command.
//
// If the command only provides a short description, that will be used instead.
func (c *RegisteredCommand) LongDescription(l *i18n.Localizer) string {
	if desc := c.Source.GetLongDescription(l); len(desc) > 0 {
		return desc
	}

	return c.ShortDescription(l)
}

// ExampleArgs returns optional example arguments of the command.
func (c *RegisteredCommand) ExampleArgs(l *i18n.Localizer) []string {
	return c.Source.GetExampleArgs(l)
}

// Examples returns the command's example arguments prefixed with their invoke.
// Invoke and example arguments are separated by a space.
func (c *RegisteredCommand) Examples(l *i18n.Localizer) []string {
	args := c.ExampleArgs(l)

	for i, arg := range args {
		args[i] = c.ID.AsInvoke() + " " + arg
	}

	return args
}

// IsRestricted returns whether or not this command is restricted.
func (c *RegisteredCommand) IsRestricted(s *state.State, ctx *Context) error {
	return c.Source.IsRestricted(s, ctx)
}

// Invoke invokes the command.
// See Command.Invoke for more details.
func (c *RegisteredCommand) Invoke(s *state.State, ctx *Context) (interface{}, error) {
	return c.Source.Invoke(s, ctx)
}
