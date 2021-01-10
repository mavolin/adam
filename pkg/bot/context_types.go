package bot

import (
	"sort"
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

// =============================================================================
// plugin.Provider
// =====================================================================================

type ctxPluginProvider struct {
	base *state.Base
	msg  *discord.Message

	// repos contains the already collected repositories
	repos []plugin.Repository
	// remProviders are the remaining, i.e. uncalled, plugin providers.
	remProviders []*pluginProvider

	async bool

	commands []*plugin.RegisteredCommand
	modules  []*plugin.RegisteredModule

	unavailableProviders []plugin.UnavailablePluginProvider
}

func (p *ctxPluginProvider) PluginRepositories() []plugin.Repository {
	p.lazyRepos()
	return p.repos
}

func (p *ctxPluginProvider) lazyRepos() {
	if len(p.remProviders) == 0 {
		return
	}

	if p.async {
		r, up := pluginProvidersAsync(p.base, p.msg, p.remProviders)
		p.repos = append(p.repos, r...)
		p.unavailableProviders = append(p.unavailableProviders, up...)

		p.remProviders = nil
		return
	}

	for _, remp := range p.remProviders {
		cmds, mods, err := remp.provider(p.base, p.msg)
		if err != nil {
			p.unavailableProviders = append(p.unavailableProviders, plugin.UnavailablePluginProvider{
				Name:  remp.name,
				Error: err,
			})
		} else {
			p.repos = append(p.repos, plugin.Repository{
				ProviderName: remp.name,
				Modules:      mods,
				Commands:     cmds,
				Defaults:     remp.defaults,
			})
		}
	}

	p.remProviders = nil
}

func (p *ctxPluginProvider) Commands() []*plugin.RegisteredCommand {
	p.lazyCommands()
	return p.commands
}

func (p *ctxPluginProvider) lazyCommands() {
	if p.commands == nil {
		p.lazyRepos()
		p.commands = plugin.GenerateRegisteredCommands(p.repos)
	}
}

func (p *ctxPluginProvider) Modules() []*plugin.RegisteredModule {
	p.lazyModules()
	return p.modules
}

func (p *ctxPluginProvider) lazyModules() {
	if p.modules == nil {
		p.lazyRepos()
		p.modules = plugin.GenerateRegisteredModules(p.repos)
	}
}

func (p *ctxPluginProvider) Command(id plugin.Identifier) *plugin.RegisteredCommand {
	if id.IsRoot() {
		return nil
	}

	if id.Parent().IsRoot() { // top-lvl command
		p.lazyCommands()

		name := id.Name()

		i := sort.Search(len(p.commands), func(i int) bool {
			return p.commands[i].Name >= name
		})

		if i == len(p.commands) || p.commands[i].Name != name { // nothing found
			return nil
		}

		return p.commands[i]
	}

	mod := p.Module(id.Parent())
	if mod == nil {
		return nil
	}

	return mod.FindCommand(id.Name())
}

func (p *ctxPluginProvider) Module(id plugin.Identifier) *plugin.RegisteredModule {
	p.lazyModules()

	all := id.All()
	if len(all) <= 1 { // invalid or just root
		return nil
	}

	all = all[1:]

	name := all[0].Name()

	i := sort.Search(len(p.modules), func(i int) bool {
		return p.modules[i].Name >= name
	})

	if i == len(p.modules) { // nothing found
		return nil
	}

	mod := p.modules[i]

	for _, id := range all[1:] {
		mod = mod.FindModule(id.Name())
		if mod == nil {
			return nil
		}
	}

	return mod
}

func (p *ctxPluginProvider) FindCommand(invoke string) *plugin.RegisteredCommand {
	id := plugin.NewIdentifierFromInvoke(invoke)

	name := id.Name()

	if id.Parent().IsRoot() {
		for _, cmd := range p.Commands() {
			if cmd.Name == name {
				return cmd
			}

			for _, alias := range cmd.Aliases {
				if alias == name {
					return cmd
				}
			}
		}

		return nil
	}

	mod := p.Module(id.Parent())
	for _, cmd := range mod.Commands {
		if cmd.Name == name {
			return cmd
		}

		for _, alias := range cmd.Aliases {
			if alias == name {
				return cmd
			}
		}
	}

	return nil
}

func (p *ctxPluginProvider) FindModule(invoke string) *plugin.RegisteredModule {
	id := plugin.NewIdentifierFromInvoke(invoke)
	return p.Module(id)
}

func (p *ctxPluginProvider) UnavailablePluginProviders() []plugin.UnavailablePluginProvider {
	p.lazyRepos()
	return p.unavailableProviders
}

// =============================================================================
// plugin.ErrorHandler
// =====================================================================================

type ctxErrorHandler func(error)

func newCtxErrorHandler(
	s *state.State, ctx *plugin.Context, f func(error, *state.State, *plugin.Context),
) ctxErrorHandler {
	return func(err error) { f(err, s, ctx) }
}

func (h ctxErrorHandler) HandleError(err error) {
	if err != nil {
		h(err)
	}
}

func (h ctxErrorHandler) HandleErrorSilent(err error) {
	err = errors.Silent(err)
	if err != nil {
		h(err)
	}
}

// =============================================================================
// plugin.DiscordDataProvider
// =====================================================================================

type discordDataProvider struct {
	s *state.State

	mut sync.Mutex

	guild   *discord.Guild
	guildWG *sync.WaitGroup

	channel   *discord.Channel
	channelWG *sync.WaitGroup

	self   *discord.Member
	selfWG *sync.WaitGroup

	guildID   discord.GuildID
	channelID discord.ChannelID
	selfID    discord.UserID
}

func (d *discordDataProvider) GuildAsync() func() (*discord.Guild, error) { //nolint:dupl
	if d.guild != nil {
		return func() (*discord.Guild, error) { return d.guild, nil }
	}

	d.mut.Lock()
	defer d.mut.Unlock()

	g, err := d.s.Cabinet.Guild(d.guildID)
	if err == nil {
		d.guild = g
		return func() (*discord.Guild, error) { return g, nil }
	}

	if d.guildWG != nil {
		return func() (*discord.Guild, error) {
			d.guildWG.Wait()
			return d.guild, err
		}
	}

	d.guildWG = new(sync.WaitGroup)
	d.guildWG.Add(1)

	go func() {
		d.guild, err = d.s.Guild(d.guildID)
		d.guildWG.Done()
	}()

	return func() (*discord.Guild, error) {
		d.guildWG.Wait()
		return d.guild, err
	}
}

func (d *discordDataProvider) ChannelAsync() func() (*discord.Channel, error) { //nolint:dupl
	if d.channel != nil {
		return func() (*discord.Channel, error) { return d.channel, nil }
	}

	d.mut.Lock()
	defer d.mut.Unlock()

	c, err := d.s.Cabinet.Channel(d.channelID)
	if err == nil {
		d.channel = c
		return func() (*discord.Channel, error) { return c, nil }
	}

	if d.channelWG != nil {
		return func() (*discord.Channel, error) {
			d.channelWG.Wait()
			return d.channel, err
		}
	}

	d.channelWG = new(sync.WaitGroup)
	d.channelWG.Add(1)

	go func() {
		d.channel, err = d.s.Channel(d.channelID)
		d.channelWG.Done()
	}()

	return func() (*discord.Channel, error) {
		d.channelWG.Wait()
		return d.channel, err
	}
}

func (d *discordDataProvider) SelfAsync() func() (*discord.Member, error) { //nolint:dupl
	if d.self != nil {
		return func() (*discord.Member, error) { return d.self, nil }
	}

	d.mut.Lock()
	defer d.mut.Unlock()

	m, err := d.s.Cabinet.Member(d.guildID, d.selfID)
	if err == nil {
		d.self = m
		return func() (*discord.Member, error) { return m, nil }
	}

	if d.selfWG != nil {
		return func() (*discord.Member, error) {
			d.selfWG.Wait()
			return d.self, err
		}
	}

	d.selfWG = new(sync.WaitGroup)
	d.selfWG.Add(1)

	go func() {
		d.self, err = d.s.Member(d.guildID, d.selfID)
		d.selfWG.Done()
	}()

	return func() (*discord.Member, error) {
		d.selfWG.Wait()
		return d.self, err
	}
}
