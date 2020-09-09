package mock

import (
	"reflect"
	"sync"

	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/plugin"
)

// DiscordDataProvider is the mock implementation of
// plugin.DiscordDataProvider.
type DiscordDataProvider struct {
	ChannelReturn *discord.Channel
	ChannelError  error

	GuildReturn *discord.Guild
	GuildError  error

	SelfReturn *discord.Member
	SelfError  error
}

func (d DiscordDataProvider) Channel() (*discord.Channel, error) {
	return d.ChannelReturn, d.ChannelError
}

func (d DiscordDataProvider) Guild() (*discord.Guild, error) {
	return d.GuildReturn, d.GuildError
}

func (d DiscordDataProvider) Self() (*discord.Member, error) {
	return d.SelfReturn, d.SelfError
}

// PluginProvider is a mock for a plugin.Provider.
// For simplicity, this plugin.Provider won't merge modules, so you should make
// sure that only one Repository has a module with a given name.
type PluginProvider struct {
	AllCommandsReturn []plugin.CommandRepository
	AllModulesReturn  []plugin.ModuleRepository

	AllCommandsError, AllModulesError error
	CommandsError, ModulesError       error
	CommandError, ModuleError         error
	FindCommandError, FindModuleError error
}

func (p PluginProvider) AllCommands() ([]plugin.CommandRepository, error) {
	if p.AllCommandsReturn == nil {
		return nil, p.AllCommandsError
	}

	cp := make([]plugin.CommandRepository, len(p.AllCommandsReturn))
	copy(cp, p.AllCommandsReturn)

	return cp, p.AllCommandsError
}

func (p PluginProvider) AllModules() ([]plugin.ModuleRepository, error) {
	if p.AllModulesReturn == nil {
		return nil, p.AllModulesError
	}

	cp := make([]plugin.ModuleRepository, len(p.AllModulesReturn))
	copy(cp, p.AllModulesReturn)

	return cp, p.AllModulesError
}

func (p PluginProvider) Commands() ([]plugin.Command, error) {
	var qty int

	for _, r := range p.AllCommandsReturn {
		qty += len(r.Commands)
	}

	cmds := make([]plugin.Command, 0, qty)

	i := 0

	for _, r := range p.AllCommandsReturn {
		copy(cmds[i:i+len(r.Commands)], r.Commands)
		i = len(r.Commands)
	}

	return cmds, p.CommandsError
}

func (p PluginProvider) Modules() ([]plugin.Module, error) {
	var qty int

	for _, r := range p.AllModulesReturn {
		qty += len(r.Modules)
	}

	cmds := make([]plugin.Module, 0, qty)

	i := 0

	for _, r := range p.AllModulesReturn {
		copy(cmds[i:i+len(r.Modules)], r.Modules)
		i = len(r.Modules)
	}

	return cmds, p.ModulesError
}

func (p PluginProvider) Command(id plugin.Identifier) (plugin.Command, error) {
	all := id.All()

	if len(all) <= 1 { // just root or invalid
		return nil, p.CommandError
	}

	if len(all) == 2 { // top-level command
		for _, r := range p.AllCommandsReturn {
			cmd := findCommand(r.Commands, all[1].Name(), false)
			if cmd != nil {
				return cmd, p.CommandError
			}
		}

		return nil, p.CommandError
	}

	var mod plugin.Module

	for _, r := range p.AllModulesReturn {
		mod = findModule(r.Modules, all[1].Name())
		if mod != nil {
			goto ModFound
		}
	}

	return nil, p.CommandError

ModFound:
	for _, id := range all[2 : len(all)-1] {
		mod = findModule(mod.Modules(), id.Name())
	}

	return findCommand(mod.Commands(), all[len(all)-1].Name(), false), nil
}

func (p PluginProvider) Module(id plugin.Identifier) (plugin.Module, error) {
	all := id.All()

	if len(all) <= 1 { // just root or invalid
		return nil, p.ModuleError
	}

	var mod plugin.Module

	for _, r := range p.AllModulesReturn {
		mod = findModule(r.Modules, all[1].Name())
		if mod != nil {
			goto ModFound
		}
	}

	return nil, p.ModuleError

ModFound:

	for _, id := range all[2:] {
		mod = findModule(mod.Modules(), id.Name())
	}

	return mod, p.ModuleError
}

func (p PluginProvider) FindCommand(invoke string) (plugin.Command, error) {
	id := plugin.IdentifierFromInvoke(invoke)
	all := id.All()[1:]

	var mod plugin.Module

	for _, r := range p.AllModulesReturn {
		mod = findModule(r.Modules, all[0].Name())
		if mod != nil {
			goto ModFound
		}
	}

	return nil, p.FindCommandError

ModFound:

	for _, id := range all[1 : len(all)-1] {
		mod = findModule(mod.Modules(), id.Name())
	}

	return findCommand(mod.Commands(), all[len(all)-1].Name(), true), p.FindCommandError
}

func (p PluginProvider) FindModule(invoke string) (plugin.Module, error) {
	id := plugin.IdentifierFromInvoke(invoke)
	all := id.All()[1:]

	var mod plugin.Module

	for _, r := range p.AllModulesReturn {
		mod = findModule(r.Modules, all[0].Name())
		if mod != nil {
			goto ModFound
		}
	}

	return nil, p.FindModuleError

ModFound:

	for _, id := range all[1:] {
		mod = findModule(mod.Modules(), id.Name())
	}

	return mod, p.FindModuleError
}

func findCommand(cmds []plugin.Command, name string, checkAliases bool) plugin.Command {
	for _, cmd := range cmds {
		if cmd.Meta().GetName() == name {
			return cmd
		}

		if checkAliases {
			for _, alias := range cmd.Meta().GetAliases() {
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
		if mod.Meta().GetName() == name {
			return mod
		}
	}

	return nil
}

type ErrorHandler struct {
	mut          sync.Mutex
	expectErr    []error
	expectSilent []error
}

func NewErrorHandler() *ErrorHandler {
	return new(ErrorHandler)
}

func (h *ErrorHandler) ExpectError(err error) *ErrorHandler {
	h.expectErr = append(h.expectErr, err)
	return h
}

func (h *ErrorHandler) ExpectSilentError(err error) *ErrorHandler {
	h.expectSilent = append(h.expectSilent, err)
	return h
}

func (h *ErrorHandler) HandleError(err error) {
	h.mut.Lock()
	defer h.mut.Unlock()

	for i, expect := range h.expectErr {
		if reflect.DeepEqual(err, expect) {
			h.expectErr = append(h.expectErr[:i], h.expectErr[i+1:]...)
			return
		}

		err2 := err

		for uerr, ok := err2.(interface{ Unwrap() error }); ok; uerr, ok = err2.(interface{ Unwrap() error }) {
			err2 = uerr.Unwrap()

			if reflect.DeepEqual(err2, expect) {
				h.expectErr = append(h.expectErr[:i], h.expectErr[i+1:]...)
				return
			}
		}
	}

	panic("unexpected call to plugin.ErrorHandler.HandleError")
}

func (h *ErrorHandler) HandleErrorSilent(err error) {
	h.mut.Lock()
	defer h.mut.Unlock()

	for i, expect := range h.expectSilent {
		if reflect.DeepEqual(err, expect) {
			h.expectSilent = append(h.expectSilent[:i], h.expectSilent[i+1:]...)
			return
		}

		err2 := err

		for uerr, ok := err2.(interface{ Unwrap() error }); ok; uerr, ok = err2.(interface{ Unwrap() error }) {
			err2 = uerr.Unwrap()

			if reflect.DeepEqual(err2, expect) {
				h.expectErr = append(h.expectSilent[:i], h.expectSilent[i+1:]...)
				return
			}
		}
	}

	panic("unexpected call to plugin.ErrorHandler.HandleErrorSilent")
}
