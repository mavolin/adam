package mock

import (
	"reflect"
	"strings"
	"sync"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// AttachState is a utility function that attaches a state.State to a
// plugin.Context.
// Responses won't be copied.
func AttachState(s *state.State, ctx *plugin.Context) *plugin.Context {
	cp := plugin.NewContext(s)

	cp.MessageCreateEvent = ctx.MessageCreateEvent
	cp.Localizer = ctx.Localizer
	cp.Args = ctx.Args
	cp.Flags = ctx.Flags
	cp.InvokedCommand = ctx.InvokedCommand
	cp.DiscordDataProvider = ctx.DiscordDataProvider
	cp.Prefix = ctx.Prefix
	cp.Location = ctx.Location
	cp.BotOwnerIDs = ctx.BotOwnerIDs
	cp.ResponseMiddlewares = ctx.ResponseMiddlewares
	cp.Provider = ctx.Provider
	cp.ErrorHandler = ctx.ErrorHandler

	return cp
}

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
// Additionally, RegisteredCommands will always assume the settings of their
// direct command aquivalent, this means no checks on parents will be
// performed.
type PluginProvider struct {
	PluginRepositoriesReturn []plugin.Repository

	PluginRepositoriesError           error
	CommandsError, ModulesError       error
	CommandError, ModuleError         error
	FindCommandError, FindModuleError error
}

func (p PluginProvider) PluginRepositories() ([]plugin.Repository, error) {
	if p.PluginRepositoriesReturn == nil {
		return nil, p.PluginRepositoriesError
	}

	cp := make([]plugin.Repository, len(p.PluginRepositoriesReturn))
	copy(cp, p.PluginRepositoriesReturn)

	return cp, p.PluginRepositoriesError
}

func (p PluginProvider) Commands() ([]plugin.RegisteredCommand, error) {
	var qty int

	for _, r := range p.PluginRepositoriesReturn {
		qty += len(r.Commands)
	}

	cmds := make([]plugin.RegisteredCommand, 0, qty)

	for _, r := range p.PluginRepositoriesReturn {
		cmds = append(cmds, asRegisteredCommands(r.Commands, nil)...)
	}

	return cmds, p.CommandsError
}

func (p PluginProvider) Modules() ([]plugin.RegisteredModule, error) {
	var qty int

	for _, r := range p.PluginRepositoriesReturn {
		qty += len(r.Modules)
	}

	mods := make([]plugin.RegisteredModule, 0, qty)

	for _, r := range p.PluginRepositoriesReturn {
		mods = append(mods, asRegisteredModules(r.Modules)...)
	}

	return mods, p.ModulesError
}

func (p PluginProvider) Command(id plugin.Identifier) (plugin.RegisteredCommand, error) {
	cmd, err := p.FindCommand(id.AsInvoke())
	if err != nil {
		return cmd, p.CommandError
	}

	return cmd, p.CommandError
}

func (p PluginProvider) Module(id plugin.Identifier) (plugin.RegisteredModule, error) {
	mod, err := p.FindModule(id.AsInvoke())
	if err != nil {
		return mod, p.ModuleError
	}

	return mod, p.ModuleError
}

func (p PluginProvider) FindCommand(invoke string) (plugin.RegisteredCommand, error) {
	id := plugin.IdentifierFromInvoke(invoke)
	all := id.All()[1:]

	if len(all) <= 1 { // just root or invalid
		return nil, p.CommandError
	}

	if len(all) == 2 { // top-level command
		for _, r := range p.PluginRepositoriesReturn {
			cmd := findCommand(r.Commands, all[1].Name(), false)
			if cmd != nil {
				return asRegisteredCommand(cmd, nil), p.CommandError
			}
		}

		return nil, p.CommandError
	}

	var mod plugin.Module

	for _, r := range p.PluginRepositoriesReturn {
		mod = findModule(r.Modules, all[1].Name())
		if mod != nil {
			goto ModFound
		}
	}

	return nil, p.CommandError

ModFound:
	rmod := asRegisteredModule(mod)

	cmdID := all[len(all)-1][1:]
	cmdID = cmdID[strings.Index(string(cmdID), ".")+1:]

	return rmod.FindCommand(cmdID.AsInvoke()), nil
}

func (p PluginProvider) FindModule(invoke string) (plugin.RegisteredModule, error) {
	id := plugin.IdentifierFromInvoke(invoke)
	all := id.All()[1:]

	if len(all) <= 1 { // just root or invalid
		return nil, p.FindModuleError
	}

	var mod plugin.Module

	for _, r := range p.PluginRepositoriesReturn {
		mod = findModule(r.Modules, all[1].Name())
		if mod != nil {
			goto ModFound
		}
	}

	return nil, p.FindModuleError

ModFound:
	rmod := asRegisteredModule(mod)

	if len(all) == 2 { // top-level module
		return rmod, p.FindModuleError
	}

	modID := all[len(all)-1][1:]
	modID = modID[strings.Index(string(modID), ".")+1:]

	return rmod.FindModule(modID.AsInvoke()), p.FindModuleError
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
				h.expectSilent = append(h.expectSilent[:i], h.expectSilent[i+1:]...)
				return
			}
		}
	}

	panic("unexpected call to plugin.ErrorHandler.HandleErrorSilent")
}
