package plugin

import (
	"time"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
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

	// Localizer is the localizer set to the guilds language.
	*localization.Localizer

	// Args contains the arguments supplied to the bot.
	// They are guaranteed to be valid and parsed according to the type spec.
	Args Args
	// Flags contains the flags supplied to the bot.
	// They are guaranteed to be valid and parsed according to the type spec.
	Flags Flags

	// Command is the Command that is being invoked.
	Command Command

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

	// ErrorHandler is an embedded interface that provides error handling
	// capabilities to the command.
	ErrorHandler

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
		Commands() []Command
		// Command returns the Command with the passed Identifier, or nil if no
		// such command exists.
		Command(Identifier) Command
		// Module returns the Module with the passed Identifier, or nil if no
		// such module exists.
		Module(Identifier) Module
		// Modules returns a copy of the bot's modules.
		Modules() []Module

		// RuntimeCommands returns a copy of the runtime commands in this
		// guild.
		// The outer slice represents the individual runtime command providers.
		RuntimeCommands() ([][]Command, error)
		// RuntimeCommand returns the first runtime command with the passed
		// Identifier, or (nil, nil) if no such command exists.
		RuntimeCommand(Identifier) (Command, error)
		// RuntimeModule returns the first runtime module with the passed
		// Identifier, or (nil, nil) if no such module exists.
		RuntimeModule(Identifier) (Module, error)
		// RuntimeModules returns a copy of the runtime modules in this guild.
		// The outer slice represents the individual runtime module providers.
		RuntimeModules() ([][]Module, error)
	}

	// ErrorHandler is an embeddable interface used to provide direct error
	// handling capabilities to a command.
	// This is useful if an error is encountered, that should be captured
	// through the bot's error handler, but execution can remain uninterrupted.
	ErrorHandler interface {
		// HandleError hands the error to the bot's error handler.
		HandleError(err interface{})
		// HandleErroSilent wraps the error using errors.Silent and hands it to
		// the bot's error handler.
		HandleErrorSilent(err interface{})
	}
)
