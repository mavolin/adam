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

func (d *DiscordDataProvider) Channel() (*discord.Channel, error) {
	return d.ChannelReturn, d.ChannelError
}

func (d *DiscordDataProvider) Guild() (*discord.Guild, error) {
	return d.GuildReturn, d.GuildError
}

func (d *DiscordDataProvider) Self() (*discord.Member, error) {
	return d.SelfReturn, d.SelfError
}

type PluginProvider struct {
	CommandsReturn []plugin.Command
	ModulesReturn  []plugin.Module

	RuntimeCommandsReturn                     [][]plugin.Command
	RuntimeCommandsError, RuntimeCommandError error

	RuntimeModulesReturn                    [][]plugin.Module
	RuntimeModulesError, RuntimeModuleError error
}

func (p *PluginProvider) Commands() (cp []plugin.Command) {
	cp = make([]plugin.Command, len(p.CommandsReturn))
	copy(cp, p.CommandsReturn)

	return
}

func (p *PluginProvider) Command(id plugin.Identifier) plugin.Command {
	all := id.All()

	if len(all) <= 1 { // just root or invalid
		return nil
	}

	if len(all) == 2 { // top-level command
		return findCommand(p.CommandsReturn, all[1].Name())
	}

	mod := findModule(p.ModulesReturn, all[1].Name())
	if mod == nil {
		return nil
	}

	for _, id := range all[2 : len(all)-1] {
		mod = findModule(mod.Modules(), id.Name())
	}

	return findCommand(mod.Commands(), all[len(all)-1].Name())
}

func (p *PluginProvider) Modules() []plugin.Module {
	cp := make([]plugin.Module, len(p.ModulesReturn))
	copy(cp, p.ModulesReturn)

	return cp
}

func (p *PluginProvider) Module(id plugin.Identifier) plugin.Module {
	all := id.All()

	if len(all) <= 1 { // just root or invalid
		return nil
	}

	mod := findModule(p.ModulesReturn, all[1].Name())
	if mod == nil {
		return nil
	}

	for _, id := range all[2:] {
		mod = findModule(mod.Modules(), id.Name())
	}

	return mod
}

func (p *PluginProvider) RuntimeCommands() ([][]plugin.Command, error) {
	var cp [][]plugin.Command

	if p.RuntimeCommandsReturn != nil {
		cp = make([][]plugin.Command, len(p.CommandsReturn))
		copy(cp, p.RuntimeCommandsReturn)
	}

	return cp, p.RuntimeCommandsError
}

func (p *PluginProvider) RuntimeCommand(id plugin.Identifier) (plugin.Command, error) {
	if p.RuntimeCommandError != nil {
		return nil, p.RuntimeCommandError
	}

	all := id.All()

	if len(all) <= 1 { // just root or invalid
		return nil, p.RuntimeCommandError
	}

	if len(all) == 2 { // top-level command
		for _, cmds := range p.RuntimeCommandsReturn {
			cmd := findCommand(cmds, all[1].Name())
			if cmd != nil {
				return cmd, p.RuntimeCommandError
			}
		}

		return nil, p.RuntimeCommandError
	}

	for _, mods := range p.RuntimeModulesReturn {
		mod := findModule(mods, all[1].Name())
		if mod == nil {
			continue
		}

		for _, id := range all[2 : len(all)-1] {
			mod = findModule(mod.Modules(), id.Name())
			if mod == nil {
				continue
			}
		}

		cmd := findCommand(mod.Commands(), all[len(all)-1].Name())
		if cmd != nil {
			return cmd, p.RuntimeCommandError
		}
	}

	return nil, p.RuntimeCommandError
}

func (p *PluginProvider) RuntimeModules() ([][]plugin.Module, error) {
	var cp [][]plugin.Module

	if p.RuntimeModulesReturn != nil {
		cp = make([][]plugin.Module, len(p.CommandsReturn))
		copy(cp, p.RuntimeModulesReturn)
	}

	return cp, p.RuntimeModulesError
}

func (p *PluginProvider) RuntimeModule(id plugin.Identifier) (plugin.Module, error) {
	all := id.All()

	if len(all) <= 1 { // just root or invalid
		return nil, nil
	}

	for _, mods := range p.RuntimeModulesReturn {
		mod := findModule(mods, all[1].Name())
		if mod == nil {
			continue
		}

		for _, id := range all[2:] {
			mod = findModule(mod.Modules(), id.Name())
			if mod == nil {
				continue
			}
		}

		if mod != nil {
			return mod, p.RuntimeModuleError
		}
	}

	return nil, nil
}

func findCommand(cmds []plugin.Command, name string) plugin.Command {
	for _, cmd := range cmds {
		if cmd.Meta().GetName() == name {
			return cmd
		}

		for _, alias := range cmd.Meta().GetAliases() {
			if alias == name {
				return cmd
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

func (h *ErrorHandler) HandleError(err interface{}) {
	h.mut.Lock()
	defer h.mut.Unlock()

	for i, expect := range h.expectErr {
		if reflect.DeepEqual(expect, err) {
			h.expectErr = append(h.expectErr[:i], h.expectErr[i+1:]...)
			return
		}
	}

	panic("unexpected call to plugin.ErrorHandler.HandleError")
}

func (h *ErrorHandler) HandleErrorSilent(err interface{}) {
	h.mut.Lock()
	defer h.mut.Unlock()

	for i, expect := range h.expectSilent {
		if reflect.DeepEqual(expect, err) {
			h.expectSilent = append(h.expectSilent[:i], h.expectSilent[i+1:]...)
			return
		}
	}

	panic("unexpected call to plugin.ErrorHandler.HandleErrorSilent")
}
