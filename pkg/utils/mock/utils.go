package mock

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

type registeredCommand struct {
	parent          plugin.RegisteredModule
	ident           plugin.Identifier
	name            string
	aliases         []string
	args            plugin.ArgConfig
	shortDescFunc   func(l *localization.Localizer) string
	longDescFunc    func(l *localization.Localizer) string
	examplesFunc    func(l *localization.Localizer) []string
	isHidden        bool
	channelTypes    plugin.ChannelTypes
	botPerms        discord.Permissions
	restrictionFunc plugin.RestrictionFunc
	throttler       plugin.Throttler
	invokeFunc      func(s *state.State, ctx *plugin.Context) (interface{}, error)
}

func (r registeredCommand) Parent() (plugin.RegisteredModule, error) { return r.parent, nil }
func (r registeredCommand) Identifier() plugin.Identifier            { return r.ident }
func (r registeredCommand) Name() string                             { return r.name }
func (r registeredCommand) Aliases() []string                        { return r.aliases }
func (r registeredCommand) Args() plugin.ArgConfig                   { return r.args }

func (r registeredCommand) ShortDescription(l *localization.Localizer) string {
	return r.shortDescFunc(l)
}

func (r registeredCommand) LongDescription(l *localization.Localizer) string {
	return r.longDescFunc(l)
}

func (r registeredCommand) Examples(l *localization.Localizer) []string { return r.examplesFunc(l) }
func (r registeredCommand) IsHidden() bool                              { return r.isHidden }
func (r registeredCommand) ChannelTypes() plugin.ChannelTypes           { return r.channelTypes }
func (r registeredCommand) BotPermissions() discord.Permissions         { return r.botPerms }

func (r registeredCommand) IsRestricted(s *state.State, ctx *plugin.Context) error {
	return r.restrictionFunc(s, ctx)
}

func (r registeredCommand) Throttler() plugin.Throttler { return r.throttler }

func (r registeredCommand) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	return r.invokeFunc(s, ctx)
}

func findCommand(cmds []plugin.Command, name string, checkAliases bool) plugin.Command {
	for _, cmd := range cmds {
		if cmd.GetName() == name {
			return cmd
		}

		if checkAliases {
			for _, alias := range cmd.GetAliases() {
				if alias == name {
					return cmd
				}
			}
		}
	}

	return nil
}

func findModule(mods []plugin.Module, name string) plugin.Module {
	for _, mod := range mods {
		if mod.GetName() == name {
			return mod
		}
	}

	return nil
}

type registeredModule struct {
	parent        plugin.RegisteredModule
	ident         plugin.Identifier
	name          string
	shortDescFunc func(l *localization.Localizer) string
	longDescFunc  func(l *localization.Localizer) string
	isHidden      bool
	cmds          []plugin.RegisteredCommand
	mods          []plugin.RegisteredModule
}

func (r registeredModule) Parent() (plugin.RegisteredModule, error) { return r.parent, nil }
func (r registeredModule) Identifier() plugin.Identifier            { return r.ident }
func (r registeredModule) Name() string                             { return r.name }

func (r registeredModule) ShortDescription(l *localization.Localizer) string {
	return r.shortDescFunc(l)
}

func (r registeredModule) LongDescription(l *localization.Localizer) string {
	return r.longDescFunc(l)
}

func (r registeredModule) IsHidden() bool { return r.isHidden }

func (r registeredModule) Commands() []plugin.RegisteredCommand {
	cp := make([]plugin.RegisteredCommand, len(r.cmds))
	copy(cp, r.cmds)

	return cp
}

func (r registeredModule) Modules() []plugin.RegisteredModule {
	cp := make([]plugin.RegisteredModule, len(r.mods))
	copy(cp, r.mods)

	return cp
}

func (r registeredModule) FindCommand(invoke string) plugin.RegisteredCommand {
	id := plugin.IdentifierFromInvoke(invoke)
	all := id.All()[1:]

	var mod plugin.RegisteredModule = r

Identifiers:
	for i := 0; i < len(all)-1; i++ {
		name := all[i].Name()

		for _, searchMod := range mod.Modules() {
			if searchMod.Name() == name {
				mod = searchMod
				continue Identifiers
			}
		}

		return nil
	}

	name := all[len(all)-1].Name()

	for _, cmd := range r.cmds {
		if cmd.Name() == name {
			return cmd
		}

		for _, alias := range cmd.Aliases() {
			if alias == name {
				return cmd
			}
		}
	}

	return nil
}

func (r registeredModule) FindModule(invoke string) plugin.RegisteredModule {
	id := plugin.IdentifierFromInvoke(invoke)
	all := id.All()[1:]

	var mod plugin.RegisteredModule = r

Identifiers:
	for _, id := range all {
		name := id.Name()

		for _, searchMod := range mod.Modules() {
			if searchMod.Name() == name {
				mod = searchMod
				continue Identifiers
			}
		}

		return nil
	}

	return nil
}

func asRegisteredModule(src plugin.Module) plugin.RegisteredModule {
	rmod := &registeredModule{
		parent:        nil,
		ident:         "." + plugin.Identifier(src.GetName()),
		name:          src.GetName(),
		shortDescFunc: src.GetShortDescription,
		longDescFunc:  src.GetLongDescription,
		isHidden:      src.IsHidden(),
		mods:          asRegisteredModules(src.Modules()),
	}

	rmod.cmds = asRegisteredCommands(src.Commands(), rmod)

	return rmod
}

func asRegisteredModules(src []plugin.Module) []plugin.RegisteredModule {
	mods := make([]plugin.RegisteredModule, 0, len(src))

	for _, mod := range src {
		mods = append(mods, asRegisteredModule(mod))
	}

	return mods
}

func asRegisteredCommand(src plugin.Command, parent plugin.RegisteredModule) plugin.RegisteredCommand {
	return &registeredCommand{
		parent:          parent,
		ident:           "." + plugin.Identifier(src.GetName()),
		name:            src.GetName(),
		aliases:         src.GetAliases(),
		args:            src.GetArgs(),
		shortDescFunc:   src.GetShortDescription,
		longDescFunc:    src.GetLongDescription,
		examplesFunc:    src.GetExamples,
		isHidden:        src.IsHidden(),
		channelTypes:    src.GetChannelTypes(),
		botPerms:        *src.GetBotPermissions(),
		restrictionFunc: src.GetRestrictionFunc(),
		throttler:       src.GetThrottler(),
		invokeFunc:      src.Invoke,
	}
}

func asRegisteredCommands(src []plugin.Command, parent plugin.RegisteredModule) []plugin.RegisteredCommand {
	cmds := make([]plugin.RegisteredCommand, 0, len(src))

	for _, cmd := range src {
		cmds = append(cmds, asRegisteredCommand(cmd, parent))
	}

	return cmds
}
