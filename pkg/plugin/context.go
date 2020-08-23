package plugin

import (
	"time"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/getsentry/sentry-go"
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/utils/discordutil"
)

// NewContext creates a new Context using the passed state.State.
// All other fields must be set manually.
func NewContext(s *state.State) *Context {
	return &Context{
		s: s,
	}
}

// Context contains context information about a command.
type Context struct {
	// MessageCreateEvent contains the event data about the invoking message.
	*state.MessageCreateEvent

	// Hub is the sentry.Hub of the command.
	Hub *sentry.Hub

	// Localizer is the localizer set to the guilds language.
	*localization.Localizer

	// Args contains the arguments supplied to the bot.
	// They are guaranteed to be valid and parsed according to the type spec.
	Args Args
	// Flags contains the flags supplied to the bot.
	// They are guaranteed to be valid and parsed according to the type spec.
	Flags Flags

	// CommandIdentifier is the Identifier of the command.
	CommandIdentifier Identifier

	// DiscordDataProvider is an embedded interface that provides additional
	// data fetched from Discord's API.
	DiscordDataProvider

	// Prefix is the prefix of the bot in the guild.
	Prefix string
	// Location is the timezone of the guild.
	Location *time.Location

	// HelpCommandIdentifier is the identifier of the help command.
	HelpCommandIdentifier Identifier

	// BotOwnerIDs contains the ids of the bot owners.
	BotOwnerIDs []discord.UserID

	// Provider is an embedded interface that provides access to the Commands
	// and Modules of the Bot, as well as the runtime commands and modules
	// for the guild.
	Provider

	s *state.State
}

// IsBotOwner checks if the invoking user is a bot owner.
func (c *Context) IsBotOwner() bool {
	for _, owner := range c.BotOwnerIDs {
		if c.Author.ID == owner {
			return true
		}
	}

	return false
}

// Reply replies with the passed message in the channel the command was
// originally sent in.
func (c *Context) Reply(content string) (*discord.Message, error) {
	return c.s.SendText(c.ChannelID, content)
}

// Replyl replies with the message translated from the passed
// localization.Config in the channel the command was originally sent in.
func (c *Context) Replyl(cfg localization.Config) (*discord.Message, error) {
	s, err := c.Localizer.Localize(cfg)
	if err != nil {
		return nil, err
	}

	return c.Reply(s)
}

// Replylt replies with the message translated from the passed term in the
// channel the command was originally sent in.
func (c *Context) Replylt(term localization.Term) (*discord.Message, error) {
	return c.Replyl(term.AsConfig())
}

// ReplyEmbed replies with the passed discord.Embed in the channel the command
// was originally sent in.
func (c *Context) ReplyEmbed(e discord.Embed) (*discord.Message, error) {
	return c.s.SendEmbed(c.ChannelID, e)
}

// ReplyEmbedBuilder builds the discord.Embed from the passed
// discordutil.EmbedBuilder and sends it in the channel the command was sent
// in.
func (c *Context) ReplyEmbedBuilder(e *discordutil.EmbedBuilder) (*discord.Message, error) {
	embed, err := e.Build(c.Localizer)
	if err != nil {
		return nil, err
	}

	return c.ReplyEmbed(embed)
}

// ReplyMessage sends the passed api.SendMessageData to the channel the command
// was originally sent in.
func (c *Context) ReplyMessage(data api.SendMessageData) (*discord.Message, error) {
	return c.s.SendMessageComplex(c.ChannelID, data)
}

// Reply replies with the passed message in the channel the command was
// originally sent in.
func (c *Context) ReplyDM(content string) (*discord.Message, error) {
	return c.ReplyMessageDM(api.SendMessageData{Content: content})
}

// Replyl replies with the message translated from the passed
// localization.Config in the channel the command was originally sent in.
func (c *Context) ReplyDMl(cfg localization.Config) (*discord.Message, error) {
	s, err := c.Localizer.Localize(cfg)
	if err != nil {
		return nil, err
	}

	return c.ReplyDM(s)
}

// Replylt replies with the message translated from the passed term in the
// channel the command was originally sent in.
func (c *Context) ReplyDMlt(term localization.Term) (*discord.Message, error) {
	return c.ReplyDMl(term.AsConfig())
}

// ReplyEmbed replies with the passed discord.Embed in the channel the command
// was originally sent in.
func (c *Context) ReplyEmbedDM(e discord.Embed) (*discord.Message, error) {
	return c.ReplyMessageDM(api.SendMessageData{Embed: &e})
}

// ReplyEmbedBuilder builds the discord.Embed from the passed
// discordutil.EmbedBuilder and sends it in the channel the command was sent
// in.
func (c *Context) ReplyEmbedBuilderDM(e *discordutil.EmbedBuilder) (*discord.Message, error) {
	embed, err := e.Build(c.Localizer)
	if err != nil {
		return nil, err
	}

	return c.ReplyEmbedDM(embed)
}

// ReplyMessage sends the passed api.SendMessageData to the channel the command
// was originally sent in.
func (c *Context) ReplyMessageDM(data api.SendMessageData) (*discord.Message, error) {
	channel, err := c.s.CreatePrivateChannel(c.Author.ID)
	if err != nil {
		return nil, err
	}

	return c.s.SendMessageComplex(channel.ID, data)
}

type (
	// DiscordDataProvider is an embeddable interface used to extend a Context
	// with additional information.
	DiscordDataProvider interface {
		// Channel returns the channel the
		Channel() (*discord.Channel, error)
		// Guild returns the guild the message was sent in.
		// If this happened in a private channel, Guild will return nil, nil.
		Guild() (*discord.Guild, error)
		// Self returns the bot member, if this happened guild.
		// If this happened in a private channel, Self will return nil, nil.
		Self() (*discord.Member, error)
	}

	// Provider provides copies if the plugins of the bot in the Context.
	// The returned slices can therefore be freely modified.
	//
	// Copies are only created on call of one of the methods.
	Provider interface {
		// Commands returns a copy of the bot's commands.
		// Calling Commands will only trigger a copy once, and will return
		// the same copy for all subsequent calls.
		Commands() []Command
		// Modules returns a copy of the bot's modules.
		// Calling Modules will only trigger a copy once, and will return
		// the same copy for all subsequent calls.
		Modules() []Module

		// RuntimeCommands returns a copy of the runtime commands in this
		// guild.
		// The outer slice represents the individual runtime command providers.
		// Calling RuntimeCommands will only trigger a copy once, and will
		// return the same copy for all subsequent calls.
		RuntimeCommands() ([][]Command, error)
		// RuntimeModules returns a copy of the runtime modules in this guild.
		// The outer slice represents the individual runtime module providers.
		// Calling RuntimeModules will only trigger a copy once, and will
		// returns the same copy for all subsequent calls.
		RuntimeModules() ([][]Module, error)
	}
)
