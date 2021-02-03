package bot

import (
	"strings"
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

const whitespace = " \t\n"

// newMessageCreateEvent creates a new state.MessageCreateEvent from the passed
// *plugin.Context.
func newMessageCreateEvent(ctx *plugin.Context) *state.MessageCreateEvent {
	return &state.MessageCreateEvent{
		MessageCreateEvent: &gateway.MessageCreateEvent{
			Message: ctx.Message,
			Member:  ctx.Member,
		},
		Base: ctx.Base,
	}
}

// newMessageUpdateEvent creates a new state.MessageUpdateEvent from the passed
// *plugin.Context.
func newMessageUpdateEvent(ctx *plugin.Context) *state.MessageUpdateEvent {
	return &state.MessageUpdateEvent{
		MessageUpdateEvent: &gateway.MessageUpdateEvent{
			Message: ctx.Message,
			Member:  ctx.Member,
		},
		Base: ctx.Base,
	}
}

func findCommand(
	invoke string, p plugin.Provider, repo plugin.Repository,
) (rcmd *plugin.RegisteredCommand, args string) {
	cmds := repo.Commands
	mods := repo.Modules

	var parents []plugin.Module

	var id plugin.Identifier = "."

	for len(invoke) > 0 {
		name := firstWord(invoke)

		invoke = invoke[len(name):]
		invoke = strings.TrimLeft(invoke, whitespace)

		for _, cmd := range cmds {
			if cmd.GetName() == name {
				id += plugin.Identifier(cmd.GetName())
				return newRegisteredCommandWithProvider(p, id, cmd, parents, &repo), invoke
			}

			for _, alias := range cmd.GetAliases() {
				if alias == name {
					id += plugin.Identifier(cmd.GetName())
					return newRegisteredCommandWithProvider(p, id, cmd, parents, &repo), invoke
				}
			}
		}

		mod := findModule(name, mods)
		if mod == nil {
			return nil, ""
		}

		parents = append(parents, mod)

		cmds = mod.Commands()
		mods = mod.Modules()
	}

	return nil, ""
}

func newRegisteredCommandWithProvider(
	p plugin.Provider, id plugin.Identifier, cmd plugin.Command, parents []plugin.Module, repo *plugin.Repository,
) *plugin.RegisteredCommand {
	rcmd := plugin.NewRegisteredCommandWithProvider(p, cmd.GetRestrictionFunc())
	rcmd.ProviderName = repo.ProviderName
	rcmd.Source = cmd
	rcmd.SourceParents = parents
	rcmd.Identifier = id
	rcmd.Name = cmd.GetName()
	rcmd.Aliases = cmd.GetAliases()
	rcmd.Args = cmd.GetArgs()
	rcmd.Hidden = cmd.IsHidden()

	rcmd.ChannelTypes = cmd.GetChannelTypes()
	if rcmd.ChannelTypes == 0 {
		rcmd.ChannelTypes = plugin.AllChannels
	}

	rcmd.BotPermissions = cmd.GetBotPermissions()
	rcmd.Throttler = cmd.GetThrottler()

	return rcmd
}

func findModule(name string, mods []plugin.Module) plugin.Module {
	for _, mod := range mods {
		if mod.GetName() == name {
			return mod
		}
	}

	return nil
}

func pluginProvidersAsync(
	base *state.Base, msg *discord.Message, providers []*pluginProvider,
) ([]plugin.Repository, []plugin.UnavailablePluginProvider) {
	var wg sync.WaitGroup
	wg.Add(len(providers))

	repos := make([]plugin.Repository, len(providers))
	var unavailableProviders []plugin.UnavailablePluginProvider

	var mut sync.Mutex

	for i, p := range providers {
		go func(p *pluginProvider, i int) {
			cmds, mods, err := p.provider(base, msg)

			mut.Lock()
			defer mut.Unlock()

			if err != nil {
				unavailableProviders = append(unavailableProviders, plugin.UnavailablePluginProvider{
					Name:  p.name,
					Error: err,
				})
			} else {
				repos[i] = plugin.Repository{
					ProviderName: p.name,
					Modules:      mods,
					Commands:     cmds,
				}
			}

			wg.Done()
		}(p, i)
	}

	wg.Wait()

	for i := 0; i < len(repos); i++ {
		if len(repos[i].ProviderName) == 0 {
			copy(repos[i:], repos[i+1:])
			repos = repos[:len(repos)-1]
		}
	}

	return repos, unavailableProviders
}

func firstWord(s string) string {
	for i, r := range s {
		if strings.ContainsRune(whitespace, r) {
			return s[:i]
		}
	}

	return s
}
