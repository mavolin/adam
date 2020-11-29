package mock

import (
	"reflect"
	"sort"
	"strings"
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

var _ plugin.DiscordDataProvider = DiscordDataProvider{}

func (d DiscordDataProvider) Channel() (*discord.Channel, error) {
	return d.ChannelReturn, d.ChannelError
}

func (d DiscordDataProvider) Guild() (*discord.Guild, error) {
	return d.GuildReturn, d.GuildError
}

func (d DiscordDataProvider) Self() (*discord.Member, error) {
	return d.SelfReturn, d.SelfError
}

// PluginProvider is a mock implementation of plugin.Provider.
// Calls to FindX must be made with trimmed, space-separated invokes.
type PluginProvider struct {
	// PluginRepositoriesReturn is the value returned by PluginRepositories.
	// The first element's ProviderName must be 'built_in'.
	PluginRepositoriesReturn []plugin.Repository

	commands []*plugin.RegisteredCommand
	modules  []*plugin.RegisteredModule

	UnavailablePluginProvidersReturn []plugin.UnavailablePluginProvider
}

var _ plugin.Provider = PluginProvider{}

func (p PluginProvider) PluginRepositories() []plugin.Repository {
	return p.PluginRepositoriesReturn
}

func (p PluginProvider) Commands() []*plugin.RegisteredCommand {
	if p.PluginRepositoriesReturn == nil {
		return nil
	} else if p.commands == nil {
		p.commands = plugin.GenerateRegisteredCommands(p.PluginRepositoriesReturn)
	}

	return p.commands
}

func (p PluginProvider) Modules() []*plugin.RegisteredModule {
	if p.PluginRepositoriesReturn == nil {
		return nil
	} else if p.modules == nil {
		p.modules = plugin.GenerateRegisteredModules(p.PluginRepositoriesReturn)
	}

	return p.modules
}

func (p PluginProvider) Command(id plugin.Identifier) *plugin.RegisteredCommand {
	all := id.All()
	if len(all) <= 1 { // invalid or just root
		return nil
	}

	all = all[1:]

	if len(all) == 1 { // top-level command
		cmds := p.Commands()

		i := sort.Search(len(cmds), func(i int) bool {
			return cmds[i].Name == all[0].Name()
		})

		if i == len(cmds) { // nothing found
			return nil
		}

		return cmds[i]
	}

	mod := p.FindModule(all[0].Name())

	for _, id := range all[1 : len(all)-1] {
		mod = mod.FindModule(id.Name())
		if mod == nil {
			return nil
		}
	}

	return mod.FindCommand(id.Name())
}

func (p PluginProvider) Module(id plugin.Identifier) *plugin.RegisteredModule {
	all := id.All()
	if len(all) <= 1 { // invalid or just root
		return nil
	}

	all = all[1:]

	mods := p.Modules()

	i := sort.Search(len(mods), func(i int) bool {
		return mods[i].Name == all[0].Name()
	})

	if i == len(mods) { // nothing found
		return nil
	}

	mod := mods[i]

	for _, id := range all[1:] {
		mod = mod.FindModule(id.Name())
		if mod == nil {
			return nil
		}
	}

	return mod
}

func (p PluginProvider) FindCommand(invoke string) *plugin.RegisteredCommand {
	id := "." + strings.ReplaceAll(invoke, " ", ".")
	return p.Command(plugin.Identifier(id))
}

func (p PluginProvider) FindModule(invoke string) *plugin.RegisteredModule {
	id := "." + strings.ReplaceAll(invoke, " ", ".")
	return p.Module(plugin.Identifier(id))
}

func (p PluginProvider) UnavailablePluginProviders() []plugin.UnavailablePluginProvider {
	return p.UnavailablePluginProvidersReturn
}

type ErrorHandler struct {
	mut          sync.Mutex
	expectErr    []error
	expectSilent []error
}

var _ plugin.ErrorHandler = new(ErrorHandler)

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

		type unwrapper interface{ Unwrap() error }

		for uerr, ok := err2.(unwrapper); ok; uerr, ok = err2.(unwrapper) { //nolint:errorlint
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

		type unwrapper interface{ Unwrap() error }

		for uerr, ok := err2.(unwrapper); ok; uerr, ok = err2.(unwrapper) { //nolint:errorlint
			err2 = uerr.Unwrap()

			if reflect.DeepEqual(err2, expect) {
				h.expectSilent = append(h.expectSilent[:i], h.expectSilent[i+1:]...)
				return
			}
		}
	}

	panic("unexpected call to plugin.ErrorHandler.HandleErrorSilent")
}
