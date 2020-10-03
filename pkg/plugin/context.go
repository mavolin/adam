package plugin

import (
	"time"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/internal/constant"
	"github.com/mavolin/adam/pkg/localization"
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
	*localization.Localizer

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
		return nil, withStack(err)
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
// localization.Config in a direct message to the invoking user.
func (c *Context) ReplyDMl(cfg localization.Config) (*discord.Message, error) {
	s, err := c.Localizer.Localize(cfg)
	if err != nil {
		return nil, err
	}

	return c.ReplyDM(s)
}

// Replylt replies with the message translated from the passed term in a direct
// message to the invoking user.
func (c *Context) ReplyDMlt(term localization.Term) (*discord.Message, error) {
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
		return nil, withStack(err)
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
		return withStack(err)
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
		return withStack(err)
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
func (c *Context) DeleteInvoke() error { return withStack(c.s.DeleteMessage(c.ChannelID, c.ID)) }

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
		// AllCommands returns a copy of all top-level commands.
		// Commands[0] contains the built-in commands of the bot, and is named
		// 'built_in'.
		// If there are no built-in top-level commands, Commands will be nil.
		//
		// AllCommands will always return valid data, even if error != nil.
		// If a bot.RuntimePluginProvider returns an error, it will be wrapped
		// in a bot.RuntimePluginProviderError, that contains the name of the
		// runtime plugin provider and the original error.
		//
		// If multiple errors occur, the error will be of type
		// errors.MultiError.
		AllCommands() ([]CommandRepository, error)
		// AllModules returns a copy of all top-level modules.
		// Modules[0] contains the built-in commands of the bot, and is named
		// 'built_in'.
		// If there are no built-in top-level modules, Modules will be nil.
		//
		// AllModules will always return valid data, even if error != nil.
		// If a bot.RuntimePluginProvider returns an error, it will be wrapped
		// in a bot.RuntimePluginProviderError that contains the name of the
		// runtime plugin provider and the original error.
		//
		// If multiple errors occur, the error will be of type
		// errors.MultiError.
		AllModules() ([]ModuleRepository, error)

		// Commands returns merged version of AllCommands for simpler
		// iteration.
		//
		// Commands will always return valid data, even if error != nil.
		// However, all runtime plugin providers that returned an error won't
		// be included, and their error will be returned wrapped in a
		// bot.RuntimePluginProviderError.
		// If multiple errors occur, a errors.MultiError filled with
		// bot.RuntimePluginProviderErrors will be returned.
		Commands() ([]RegisteredCommand, error)
		// Modules returns a merged version of AllModules, as the command
		// router uses it.
		//
		// Modules will always return valid data, even if error != nil.
		// However, all runtime plugin providers that returned an error won't
		// be included, and their error will be returned wrapped in a
		// bot.RuntimePluginProviderError.
		// If multiple errors occur, a errors.MultiError filled with
		// bot.RuntimePluginProviderErrors will be returned.
		Modules() ([]RegisteredModule, error)

		// Command returns the RegisteredCommand with the passed Identifier.
		//
		// Note that Identifiers may only consist of the command's name, not
		// their alias.
		//
		// It will return nil, nil if no command matching the identifier was
		// found.
		// An error will only be returned if one of the
		// bot.RuntimePluginProviders returns an error, and should therefore
		// be seen as an indicator that the command may exist, but is
		// unavailable right now.
		// If so, it will be wrapped in a bot.RuntimePluginProviderError.
		// If multiple errors occur, the returned error will be of type
		// errors.MultiError.
		Command(Identifier) (RegisteredCommand, error)
		// Module returns the RegisteredModule with the passed Identifier.
		//
		// It will return nil, nil if no module matching the identifier was
		// found.
		// An error will only be returned if one of the runtime plugin
		// providers returns an error, and should therefore be seen as an
		// indicator that the module may exist, but is unavailable right now.
		// If so, it will be wrapped in a bot.RuntimePluginProviderError.
		// If multiple errors occur, the returned error will be of type
		// errors.MultiError.
		Module(Identifier) (RegisteredModule, error)

		// FindCommand returns the RegisteredCommand with the passed invoke.
		//
		// Note that Identifiers may only consist of the command's name, not
		// their alias.
		//
		// It will return nil, nil if no command matching the passed invoke was
		// found.
		// An error will only be returned if one of the runtime plugin
		// providers returns an error, and should therefore be seen as an
		// indicator that the command may exist, but is unavailable right now.
		// If so, it will be wrapped in a bot.RuntimePluginProviderError.
		// If multiple errors occur, the returned error will be of type
		// errors.MultiError.
		FindCommand(invoke string) (RegisteredCommand, error)
		// FindModule returns the RegisteredModule with the passed invoke.
		//
		// It will return nil, nil if no module matching the passed invoke was
		// found.
		// An error will only be returned if one of the runtime plugin
		// providers returns an error, and should therefore be seen as an
		// indicator that the module may exist, but is unavailable right now.
		// If so, it will be wrapped in a bot.RuntimePluginProviderError.
		// If multiple errors occur, the returned error will be of type
		// errors.MultiError.
		FindModule(invoke string) (RegisteredModule, error)
	}

	// ModuleRepository is the struct returned by Provider.AllModules.
	ModuleRepository struct {
		// Name is the name of the bot.RuntimePluginProvider that provides
		// these modules.
		Name string
		// Modules are the modules were returned by the
		// bot.RuntimePluginProvider.
		Modules []Module
	}

	// CommandRepository is the struct returned by Provider.AllCommands.
	CommandRepository struct {
		// Name is the name of the bot.RuntimePluginProvider that provides
		// these commands.
		Name string
		// Commands are the commands that were returned by the
		// bot.RuntimePluginProvider.
		Commands []Command
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
