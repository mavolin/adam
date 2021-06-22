package bot

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

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

func (h ctxErrorHandler) HandleErrorSilently(err error) {
	h.HandleError(errors.Silent(err))
}

// =============================================================================
// plugin.DiscordDataProvider
// =====================================================================================

type discordDataProvider struct {
	s *state.State

	mut sync.Mutex

	guild    *discord.Guild
	guildErr error
	guildWG  *sync.WaitGroup

	channel    *discord.Channel
	channelErr error
	channelWG  *sync.WaitGroup

	self    *discord.Member
	selfErr error
	selfWG  *sync.WaitGroup

	guildID   discord.GuildID
	channelID discord.ChannelID
	selfID    discord.UserID
}

func (d *discordDataProvider) GuildAsync() func() (*discord.Guild, error) { //nolint:dupl
	if d.guild != nil || d.guildErr != nil {
		return func() (*discord.Guild, error) { return d.guild, d.guildErr }
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
			return d.guild, d.guildErr
		}
	}

	d.guildWG = new(sync.WaitGroup)
	d.guildWG.Add(1)

	go func() {
		d.guild, err = d.s.Guild(d.guildID)
		d.guildErr = errors.WithStack(err)

		d.guildWG.Done()
	}()

	return func() (*discord.Guild, error) {
		d.guildWG.Wait()
		return d.guild, d.guildErr
	}
}

func (d *discordDataProvider) ChannelAsync() func() (*discord.Channel, error) { //nolint:dupl
	if d.channel != nil || d.channelErr != nil {
		return func() (*discord.Channel, error) { return d.channel, d.channelErr }
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
			return d.channel, d.channelErr
		}
	}

	d.channelWG = new(sync.WaitGroup)
	d.channelWG.Add(1)

	go func() {
		d.channel, err = d.s.Channel(d.channelID)
		d.channelErr = errors.WithStack(err)

		d.channelWG.Done()
	}()

	return func() (*discord.Channel, error) {
		d.channelWG.Wait()
		return d.channel, d.channelErr
	}
}

func (d *discordDataProvider) SelfAsync() func() (*discord.Member, error) { //nolint:dupl
	if d.self != nil || d.selfErr != nil {
		return func() (*discord.Member, error) { return d.self, d.selfErr }
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
			return d.self, d.selfErr
		}
	}

	d.selfWG = new(sync.WaitGroup)
	d.selfWG.Add(1)

	go func() {
		d.self, err = d.s.Member(d.guildID, d.selfID)
		d.selfErr = errors.WithStack(err)

		d.selfWG.Done()
	}()

	return func() (*discord.Member, error) {
		d.selfWG.Wait()
		return d.self, d.selfErr
	}
}
