package bot

import (
	"sync"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

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
	if err != nil {
		h(errors.Silent(err))
	}
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

	parentChannel    *discord.Channel
	parentChannelErr error
	parentChannelWG  *sync.WaitGroup

	self    *discord.Member
	selfErr error
	selfWG  *sync.WaitGroup

	guildID   discord.GuildID
	channelID discord.ChannelID
	selfID    discord.UserID
}

func (d *discordDataProvider) GuildAsync() func() (*discord.Guild, error) {
	d.mut.Lock()
	defer d.mut.Unlock()

	if d.guild != nil || d.guildErr != nil {
		return func() (*discord.Guild, error) { return d.guild, d.guildErr }
	}

	if d.guildWG != nil {
		return func() (*discord.Guild, error) {
			d.guildWG.Wait()
			return d.guild, d.guildErr
		}
	}

	g, err := d.s.Cabinet.Guild(d.guildID)
	if err == nil {
		d.guild = g
		return func() (*discord.Guild, error) { return g, nil }
	}

	d.guildWG = new(sync.WaitGroup)
	d.guildWG.Add(1)

	go func() {
		g, err = d.s.Guild(d.guildID)

		d.mut.Lock()
		d.guild, d.guildErr = g, errors.WithStack(err)
		d.mut.Unlock()

		d.guildWG.Done()
	}()

	return func() (*discord.Guild, error) {
		d.guildWG.Wait()
		return d.guild, d.guildErr
	}
}

func (d *discordDataProvider) ChannelAsync() func() (*discord.Channel, error) {
	d.mut.Lock()
	defer d.mut.Unlock()

	if d.channel != nil || d.channelErr != nil {
		return func() (*discord.Channel, error) { return d.channel, d.channelErr }
	}

	if d.channelWG != nil {
		return func() (*discord.Channel, error) {
			d.channelWG.Wait()
			return d.channel, d.channelErr
		}
	}

	c, err := d.s.Cabinet.Channel(d.channelID)
	if err == nil {
		d.channel = c
		return func() (*discord.Channel, error) { return c, nil }
	}

	d.channelWG = new(sync.WaitGroup)
	d.channelWG.Add(1)

	go func() {
		c, err = d.s.Channel(d.channelID)

		d.mut.Lock()
		d.channel, d.channelErr = c, errors.WithStack(err)
		d.mut.Unlock()

		d.channelWG.Done()
	}()

	return func() (*discord.Channel, error) {
		d.channelWG.Wait()
		return d.channel, d.channelErr
	}
}

func (d *discordDataProvider) ParentChannelAsync() func() (*discord.Channel, error) {
	d.mut.Lock()
	defer d.mut.Unlock()

	if d.parentChannel != nil || d.parentChannelErr != nil {
		return func() (*discord.Channel, error) { return d.parentChannel, d.parentChannelErr }
	}

	if d.channelErr != nil {
		d.parentChannelErr = errors.Wrap(d.channelErr, "could not get child channel to extract parent id from")
		return func() (*discord.Channel, error) { return nil, d.parentChannelErr }
	}

	if d.parentChannelWG != nil {
		return func() (*discord.Channel, error) {
			d.parentChannelWG.Wait()
			return d.parentChannel, d.parentChannelErr
		}
	}

	d.parentChannelWG = new(sync.WaitGroup)
	d.parentChannelWG.Add(1)

	go func() {
		c, err := d.ChannelAsync()()
		if err != nil {
			d.mut.Lock()
			d.parentChannelErr = errors.Wrap(d.channelErr, "could not get child channel to extract parent id from")
			d.mut.Unlock()

			return
		}

		if c.ParentID == 0 {
			d.mut.Lock()
			d.parentChannelErr = errors.Wrapf(d.channelErr, "%d has no parent channel", d.channel.ID)
			d.mut.Unlock()

			return
		}

		parent, err := d.s.Channel(c.ParentID)

		d.mut.Lock()
		d.parentChannel, d.parentChannelErr = parent, errors.WithStack(err)
		d.mut.Unlock()

		d.parentChannelWG.Done()
	}()

	return func() (*discord.Channel, error) {
		d.parentChannelWG.Wait()
		return d.parentChannel, d.parentChannelErr
	}
}

//nolint:dupl
func (d *discordDataProvider) SelfAsync() func() (*discord.Member, error) {
	d.mut.Lock()
	defer d.mut.Unlock()

	if d.self != nil || d.selfErr != nil {
		return func() (*discord.Member, error) { return d.self, d.selfErr }
	}

	if d.selfWG != nil {
		return func() (*discord.Member, error) {
			d.selfWG.Wait()
			return d.self, d.selfErr
		}
	}

	m, err := d.s.Cabinet.Member(d.guildID, d.selfID)
	if err == nil {
		d.self = m
		return func() (*discord.Member, error) { return m, nil }
	}

	d.selfWG = new(sync.WaitGroup)
	d.selfWG.Add(1)

	go func() {
		m, err = d.s.Member(d.guildID, d.selfID)

		d.mut.Lock()
		d.self, d.selfErr = m, errors.WithStack(err)
		d.mut.Unlock()

		d.selfWG.Done()
	}()

	return func() (*discord.Member, error) {
		d.selfWG.Wait()
		return d.self, d.selfErr
	}
}
