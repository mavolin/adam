package plugin

import (
	"time"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/internal/constant"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/utils/embedutil"
)

// ErrInsufficientSendPermissions is an informational error that signals
// that a message wasn't sent, because the bot lacks permissions.
// This error should not be handled.
var ErrInsufficientSendPermissions = &noHandlingError{
	s: "insufficient permissions to send message",
}

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
	*i18n.Localizer

	// Args contains the arguments supplied to the bot.
	// They are guaranteed to be valid and parsed according to the type spec.
	Args Args
	// Flags contains the flags supplied to the bot.
	// They are guaranteed to be valid and parsed according to the type spec.
	Flags Flags

	// InvokedCommand is the RegisteredCommand that is being invoked.
	InvokedCommand *RegisteredCommand

	// DiscordDataProvider is an embedded interface that provides additional
	// data fetched from Discord's API.
	DiscordDataProvider

	// Prefix is the prefix of the bot in the guild.
	// If the guild has prefixes disabled, Prefix will be empty.
	Prefix string
	// Location is the timezone of the guild.
	Location *time.Location

	// BotOwnerIDs contains the ids of the bot owners.
	BotOwnerIDs []discord.UserID

	// ResponseMiddlewares contains the middlewares that should be used when
	// awaiting a response.
	// These following types are permitted:
	//		• func(*state.State, interface{})
	//		• func(*state.State, interface{}) error
	//		• func(*state.State, *state.Base)
	//		• func(*state.State, *state.Base) error
	//		• func(*state.State, *state.MessageCreateEvent)
	//		• func(*state.State, *state.MessageCreateEvent) error
	ResponseMiddlewares []interface{}

	// Provider is an embedded interface that provides access to the Commands
	// and Modules of the Bot, as well as the runtime commands and modules
	// for the guild.
	Provider

	// ErrorHandler is an embedded interface that provides error handling
	// capabilities to the command.
	ErrorHandler

	s *state.State

	dmID discord.ChannelID

	guildReplies []discord.MessageID
	dmReplies    []discord.MessageID
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
	return c.ReplyMessage(api.SendMessageData{
		Content: content,
	})

}

// Replyl replies with the message translated from the passed
// i18n.Config in the channel the command was originally sent in.
func (c *Context) Replyl(cfg i18n.Config) (*discord.Message, error) {
	s, err := c.Localizer.Localize(cfg)
	if err != nil {
		return nil, err
	}

	return c.Reply(s)
}

// Replylt replies with the message translated from the passed term in the
// channel the command was originally sent in.
func (c *Context) Replylt(term i18n.Term) (*discord.Message, error) {
	return c.Replyl(term.AsConfig())
}

// ReplyEmbed replies with the passed discord.Embed in the channel the command
// was originally sent in.
func (c *Context) ReplyEmbed(e discord.Embed) (*discord.Message, error) {
	return c.ReplyMessage(api.SendMessageData{
		Embed: &e,
	})
}

// ReplyEmbedBuilder builds the discord.Embed from the passed
// discordutil.EmbedBuilder and sends it in the channel the command was sent
// in.
func (c *Context) ReplyEmbedBuilder(e *embedutil.Builder) (*discord.Message, error) {
	embed, err := e.Build(c.Localizer)
	if err != nil {
		return nil, err
	}

	return c.ReplyEmbed(embed)
}

// ReplyMessage sends the passed api.SendMessageData to the channel the command
// was originally sent in.
func (c *Context) ReplyMessage(data api.SendMessageData) (*discord.Message, error) {
	perms, err := c.SelfPermissions()
	if err != nil {
		return nil, err
	}

	if !perms.Has(discord.PermissionSendMessages) {
		return nil, ErrInsufficientSendPermissions
	}

	msg, err := c.s.SendMessageComplex(c.ChannelID, data)
	if err != nil {
		return nil, errWithStack(err)
	}

	if c.GuildID == 0 {
		c.dmReplies = append(c.dmReplies, msg.ID)
	} else {
		c.guildReplies = append(c.guildReplies, msg.ID)
	}

	return msg, nil
}

// ReplyDM replies with the passed message in tin a direct message to the
// invoking user.
func (c *Context) ReplyDM(content string) (*discord.Message, error) {
	return c.ReplyMessageDM(api.SendMessageData{
		Content: content,
	})
}

// ReplyDMl replies with the message translated from the passed
// i18n.Config in a direct message to the invoking user.
func (c *Context) ReplyDMl(cfg i18n.Config) (*discord.Message, error) {
	s, err := c.Localizer.Localize(cfg)
	if err != nil {
		return nil, err
	}

	return c.ReplyDM(s)
}

// Replylt replies with the message translated from the passed term in a direct
// message to the invoking user.
func (c *Context) ReplyDMlt(term i18n.Term) (*discord.Message, error) {
	return c.ReplyDMl(term.AsConfig())
}

// ReplyEmbedDM replies with the passed discord.Embed in a direct message
// to the invoking user.
func (c *Context) ReplyEmbedDM(e discord.Embed) (*discord.Message, error) {
	return c.ReplyMessageDM(api.SendMessageData{
		Embed: &e,
	})
}

// ReplyEmbedBuilder builds the discord.Embed from the passed embedutil.Builder
// and sends it in a direct message to the invoking user.
func (c *Context) ReplyEmbedBuilderDM(e *embedutil.Builder) (*discord.Message, error) {
	embed, err := e.Build(c.Localizer)
	if err != nil {
		return nil, err
	}

	return c.ReplyEmbedDM(embed)
}

// ReplyMessageDM sends the passed api.SendMessageData in a direct message to
// the invoking user.
func (c *Context) ReplyMessageDM(data api.SendMessageData) (*discord.Message, error) {
	if !c.dmID.IsValid() {
		ch, err := c.s.CreatePrivateChannel(c.Author.ID)
		if err != nil {
			return nil, err
		}

		c.dmID = ch.ID
	}

	msg, err := c.s.SendMessageComplex(c.dmID, data)
	if err != nil {
		return nil, errWithStack(err)
	}

	c.dmReplies = append(c.dmReplies, msg.ID)

	return msg, nil
}

// DeleteDMReplies deletes all replies sent to the invoking user in a private
// channel during the execution of the command.
//
// Note that only those messages sent via the Context will be deleted.
func (c *Context) DeleteDMReplies() error {
	if len(c.dmReplies) == 0 {
		return nil
	}

	err := c.s.DeleteMessages(c.dmID, c.dmReplies)
	if err != nil {
		return errWithStack(err)
	}

	c.dmReplies = nil

	return nil
}

// DeleteGuildReplies deletes all replies sent to the invoking user in a guild.
// during the execution of the command.
//
// Note that only those messages sent via the Context will be deleted.
func (c *Context) DeleteGuildReplies() error {
	if len(c.guildReplies) == 0 {
		return nil
	}

	err := c.s.DeleteMessages(c.ChannelID, c.guildReplies)
	if err != nil {
		return errWithStack(err)
	}

	c.guildReplies = nil

	return nil
}

// DeleteAllReplies deletes all replies sent to the invoking user, during the
// execution of the command.
//
// Note that only those messages sent via the Context will be deleted.
func (c *Context) DeleteAllReplies() error {
	err := c.DeleteGuildReplies()
	if err != nil {
		return err
	}

	return c.DeleteDMReplies()
}

// DeleteInvoke deletes the invoking message.
func (c *Context) DeleteInvoke() error {
	return errWithStack(c.s.DeleteMessage(c.ChannelID, c.ID))
}

// DeleteInvokeInBackground deletes the invoking message in a separate
// goroutine.
// If it encounters an error, it will pass it to Context.HandleErrorSilent.
func (c *Context) DeleteInvokeInBackground() {
	go func() {
		err := c.DeleteInvoke()
		if err != nil {
			c.HandleErrorSilent(err)
		}
	}()
}

// SelfPermissions checks if the bot has the passed permissions.
// If this command is executed in a direct message, constant.DMPermissions will
// be returned instead.
func (c *Context) SelfPermissions() (discord.Permissions, error) {
	if c.GuildID == 0 {
		return constant.DMPermissions, nil
	}

	g, err := c.Guild()
	if err != nil {
		return 0, err
	}

	ch, err := c.Channel()
	if err != nil {
		return 0, err
	}

	s, err := c.Self()
	if err != nil {
		return 0, err
	}

	return discord.CalcOverwrites(*g, *ch, *s), nil
}

// UserPermissions returns the permissions of the invoking user in this
// channel.
// If this command is executed in a direct message, constant.DMPermissions will
// be returned instead.
func (c *Context) UserPermissions() (discord.Permissions, error) {
	if c.GuildID == 0 {
		return constant.DMPermissions, nil
	}

	g, err := c.Guild()
	if err != nil {
		return 0, err
	}

	ch, err := c.Channel()
	if err != nil {
		return 0, err
	}

	return discord.CalcOverwrites(*g, *ch, *c.Member), nil
}

type (
	// DiscordDataProvider is an embeddable interface used to extend a Context
	// with additional information.
	DiscordDataProvider interface {
		// Channel returns the channel the message was sent in.
		Channel() (*discord.Channel, error)
		// Guild returns the guild the message was sent in.
		// If this happened in a private channel, Guild will return nil, nil.
		Guild() (*discord.Guild, error)
		// Self returns the bot as a member, if the command was invoked in a
		// guild.
		// If this happened in a private channel, Self will return nil, nil.
		Self() (*discord.Member, error)
	}

	// Provider provides copies if the plugins of the bot in the Context.
	// The returned slices can therefore be freely modified.
	//
	// Copies are only created on call of one of the methods.
	Provider interface {
		// PluginRepositories returns plugin repositories containing all
		// commands and modules of the bot.
		// Repositories[0] contains the built-in plugins of the bot, and is
		// named 'built_in'.
		//
		// To check if any runtime plugin providers returned an error, call
		// UnavailablePluginProviders.
		PluginRepositories() []Repository

		// Commands returns all top-level commands sorted in ascending order by
		// name.
		//
		// To check if any of the runtime plugin providers returned an error,
		// call UnavailablePluginProviders.
		// If that is the case, the data returned might be incomplete.
		Commands() []RegisteredCommand
		// Modules returns all top-level modules sorted in ascending order by
		// name.
		//
		// To check if any of the runtime plugin providers returned an error,
		// call UnavailablePluginProviders.
		// If that is the case, the data returned might be incomplete.
		Modules() []RegisteredModule

		// Command returns the RegisteredCommand with the passed Identifier.
		//
		// Note that Identifiers may only consist of the command's name, not
		// their alias.
		//
		// It will return nil if no command matching the identifier was found.
		//
		// To check if any of the runtime plugin providers returned an error,
		// call UnavailablePluginProviders.
		Command(Identifier) *RegisteredCommand
		// Module returns the RegisteredModule with the passed Identifier.
		//
		// It will return nil if no module matching the identifier was found.
		//
		// To check if any of the runtime plugin providers returned an error,
		// call UnavailablePluginProviders.
		// If that is the case, the module's description might not be available
		// or differ from the description that is used if all plugin providers
		// function properly.
		Module(Identifier) *RegisteredModule

		// FindCommand returns the RegisteredCommand with the passed invoke.
		//
		// It will return nil if no command matching the passed invoke was
		// found.
		//
		// To check if any of the runtime plugin providers returned an error,
		// call UnavailablePluginProviders.
		FindCommand(invoke string) *RegisteredCommand
		// FindModule returns the RegisteredModule with the passed invoke.
		//
		// It will return nil if no module matching the passed invoke was
		// found.
		//
		// To check if any of the runtime plugin providers returned an error,
		// call UnavailablePluginProviders.
		// If that is the case, the module's description might not be available
		// or differ from the description that is used if all plugin providers
		// function properly.
		FindModule(invoke string) *RegisteredModule

		// UnavailablePluginProviders returns a list of all unavailable runtime
		// plugin providers.
		// If no runtime plugins were requested yet, it will request them and
		// return the list of unavailable ones.
		//
		// If the length of the returned slice is 0, all plugin providers are
		// available.
		UnavailablePluginProviders() []UnavailablePluginProvider
	}

	// Repository is the struct returned by Provider.PluginRepositories.
	// It contains the top-level plugins of a single repository.
	Repository struct {
		// ProviderName is the name of the bot.RuntimePluginProvider that provides
		// these plugins.
		ProviderName string
		// Modules are the top-level modules of the repository.
		Modules []Module
		// Commands are the top-level commands of the repository.
		Commands []Command

		// CommandDefaults are the global defaults for command settings, the
		// provider uses.
		CommandDefaults CommandDefaults
	}

	// UnavailablePluginProvider contains information about an unavailable
	// runtime plugin provider.
	UnavailablePluginProvider struct {
		// Name is the name of the runtime plugin provider.
		Name string
		// Error is the error returned by the runtime plugin provider.
		Error error
	}

	// ErrorHandler is an embeddable interface used to provide direct error
	// handling capabilities to a command.
	// This is useful if an error is encountered, that should be captured
	// through the bot's error handler, but execution can remain uninterrupted.
	ErrorHandler interface {
		// HandleError hands the error to the bot's error handler.
		HandleError(err error)
		// HandleErrorSilent wraps the error using errors.Silent and hands it
		// to the bot's error handler.
		HandleErrorSilent(err error)
	}
)
