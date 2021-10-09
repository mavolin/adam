package resolved

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type Command struct {
	parent   plugin.ResolvedModule
	provider *PluginProvider

	sourceName    string
	source        plugin.Command
	sourceParents []plugin.Module

	id      plugin.ID
	aliases []string
}

var _ plugin.ResolvedCommand = new(Command)

func newCommand(
	parent plugin.ResolvedModule, provider *PluginProvider, sourceName string, sourceParents []plugin.Module,
	scmd plugin.Command,
) *Command {
	parentInvoke := ""
	if parent != nil {
		parentInvoke = parent.ID().AsInvoke() + " "
	}

	invoke := parentInvoke + scmd.GetName()
	if _, ok := provider.usedNames[invoke]; ok {
		return nil
	}

	provider.usedNames[invoke] = struct{}{}

	var aliases []string

	if len(scmd.GetAliases()) > 0 {
		aliases = make([]string, len(scmd.GetAliases()))
		copy(aliases, scmd.GetAliases())

		for i, alias := range aliases {
			if _, ok := provider.usedNames[parentInvoke+alias]; ok {
				copy(aliases[i:], aliases[i+1:])
				aliases = aliases[:len(aliases)-1]
			}

			provider.usedNames[parentInvoke+alias] = struct{}{}
		}
	}

	id := plugin.ID("." + scmd.GetName())
	if parent != nil {
		id = parent.ID() + id
	}

	return &Command{
		parent:        parent,
		provider:      provider,
		sourceName:    sourceName,
		source:        scmd,
		sourceParents: sourceParents,
		id:            id,
		aliases:       aliases,
	}
}

func (cmd *Command) Parent() plugin.ResolvedModule {
	cmd.provider.Resolve()
	return cmd.parent
}

func (cmd *Command) SourceName() string             { return cmd.sourceName }
func (cmd *Command) Source() plugin.Command         { return cmd.source }
func (cmd *Command) SourceParents() []plugin.Module { return cmd.sourceParents }
func (cmd *Command) ID() plugin.ID                  { return cmd.id }
func (cmd *Command) Name() string                   { return cmd.source.GetName() }
func (cmd *Command) Aliases() []string              { return cmd.aliases }

func (cmd *Command) ShortDescription(l *i18n.Localizer) string {
	return cmd.source.GetShortDescription(l)
}

func (cmd *Command) LongDescription(l *i18n.Localizer) string {
	desc := cmd.source.GetLongDescription(l)
	if desc != "" {
		return desc
	}

	return cmd.ShortDescription(l)
}

func (cmd *Command) Args() plugin.ArgConfig { return cmd.source.GetArgs() }

func (cmd *Command) ArgParser() plugin.ArgParser {
	if p := cmd.source.GetArgParser(); p != nil {
		return p
	}

	return cmd.provider.resolver.argParser
}

func (cmd *Command) ExampleArgs(l *i18n.Localizer) plugin.ExampleArgs {
	return cmd.source.GetExampleArgs(l)
}

func (cmd *Command) Examples(l *i18n.Localizer) []string {
	exampleArgs := cmd.ExampleArgs(l)
	examples := make([]string, len(exampleArgs))

	for i, exampleArg := range exampleArgs {
		examples[i] = cmd.ID().AsInvoke()

		exampleArgString := cmd.ArgParser().FormatArgs(cmd.Args(), exampleArg.Args, exampleArg.Flags)
		if exampleArgString != "" {
			examples[i] += " " + exampleArgString
		}
	}

	return examples
}

func (cmd *Command) IsHidden() bool { return cmd.source.IsHidden() }

func (cmd *Command) ChannelTypes() plugin.ChannelTypes {
	if t := cmd.source.GetChannelTypes(); t > 0 {
		return t
	}

	return plugin.AllChannels
}

func (cmd *Command) BotPermissions() discord.Permissions { return cmd.source.GetBotPermissions() }

func (cmd *Command) IsRestricted(s *state.State, ctx *plugin.Context) error {
	return cmd.source.IsRestricted(s, ctx)
}

func (cmd *Command) Throttler() plugin.Throttler { return cmd.source.GetThrottler() }

func (cmd *Command) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	return cmd.source.Invoke(s, ctx)
}
