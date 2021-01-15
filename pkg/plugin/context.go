package plugin

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/internal/errorutil"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/utils/embedutil"
	"github.com/mavolin/adam/pkg/utils/permutil"
)

// Context contains context information about a command.
type Context struct {
	// Message is the invoking message.
	discord.Message
	// Member is the invoking member, if this happened in a guild.
	*discord.Member

	// Base is the Base MessageCreateEvent or MessageUpdateEvent that triggered
	// the invoke.
	*state.Base

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

	// Prefixes contains the prefixes of the bot in the guild.
	// Length may be 0, if the guild allows the use of mentions only.
	Prefixes []string

	// BotOwnerIDs contains the ids of the bot owners.
	BotOwnerIDs []discord.UserID

	// ReplyMiddlewares contains the middlewares that should be used when
	// awaiting a reply.
	//
	// The following types are permitted:
	//		• func(*state.State, interface{})
	//		• func(*state.State, interface{}) error
	//		• func(*state.State, *state.Base)
	//		• func(*state.State, *state.Base) error
	//		• func(*state.State, *state.MessageCreateEvent)
	//		• func(*state.State, *state.MessageCreateEvent) error
	ReplyMiddlewares []interface{}

	// Replier is the interface used to send replies to a command.
	// Defaults to replier.WrapState, found in impl/replier
	Replier Replier

	// Provider is an embedded interface that provides access to the Commands
	// and Modules of the Bot, as well as the runtime commands and modules
	// for the guild.
	Provider

	// ErrorHandler is an embedded interface that provides error handling
	// capabilities to the command.
	ErrorHandler

	// DiscordDataProvider is an embedded interface that provides additional
	// data fetched from Discord's API.
	DiscordDataProvider
}

// IsBotOwner checks if the invoking user is a bot owner.
func (ctx *Context) IsBotOwner() bool {
	for _, owner := range ctx.BotOwnerIDs {
		if ctx.Author.ID == owner {
			return true
		}
	}

	return false
}

// Reply replies with the passed message in the channel the command was
// originally sent in.
// The message will be formatted as fmt.Sprint(content...).
func (ctx *Context) Reply(content ...interface{}) (*discord.Message, error) {
	return ctx.ReplyMessage(api.SendMessageData{Content: fmt.Sprint(content...)})
}

// Reply replies with the passed message in the channel the command was
// originally sent in.
// The message will be formatted as fmt.Sprintf(format, a...).
func (ctx *Context) Replyf(format string, a ...interface{}) (*discord.Message, error) {
	return ctx.ReplyMessage(api.SendMessageData{Content: fmt.Sprintf(format, a...)})
}

// Replyl replies with the message generated from the passed i18n.Config in the
// channel the command was originally sent in.
func (ctx *Context) Replyl(c *i18n.Config) (*discord.Message, error) {
	s, err := ctx.Localizer.Localize(c)
	if err != nil {
		return nil, err
	}

	return ctx.Reply(s)
}

// Replylt replies with the message translated from the passed term in the
// channel the command was originally sent in.
func (ctx *Context) Replylt(term i18n.Term) (*discord.Message, error) {
	return ctx.Replyl(term.AsConfig())
}

// ReplyEmbed replies with the passed discord.Embed in the channel the command
// was originally sent in.
func (ctx *Context) ReplyEmbed(e discord.Embed) (*discord.Message, error) {
	return ctx.ReplyMessage(api.SendMessageData{Embed: &e})
}

// ReplyEmbedBuilder builds the discord.Embed from the passed
// embedutil.Builder and sends it in the channel the command was sent
// in.
func (ctx *Context) ReplyEmbedBuilder(e *embedutil.Builder) (*discord.Message, error) {
	embed, err := e.Build(ctx.Localizer)
	if err != nil {
		return nil, err
	}

	return ctx.ReplyEmbed(embed)
}

// ReplyMessage sends the passed api.SendMessageData to the channel the command
// was originally sent in.
func (ctx *Context) ReplyMessage(data api.SendMessageData) (*discord.Message, error) {
	return ctx.Replier.Reply(ctx, data)
}

// ReplyDM replies with the passed message in in a direct message to the
// invoking user.
// The message will be formatted as fmt.Sprint(content...).
func (ctx *Context) ReplyDM(content ...interface{}) (*discord.Message, error) {
	return ctx.ReplyMessageDM(api.SendMessageData{Content: fmt.Sprint(content...)})
}

// ReplyfDM replies with the passed message in the channel the command was
// originally sent in.
// The message will be formatted as fmt.Sprintf(format, a...).
func (ctx *Context) ReplyfDM(format string, a ...interface{}) (*discord.Message, error) {
	return ctx.ReplyDM(fmt.Sprintf(format, a...))
}

// ReplylDM replies with the message translated from the passed i18n.Config in
// a direct message to the invoking user.
func (ctx *Context) ReplylDM(c *i18n.Config) (*discord.Message, error) {
	s, err := ctx.Localizer.Localize(c)
	if err != nil {
		return nil, err
	}

	return ctx.ReplyDM(s)
}

// ReplyltDM replies with the message generated from the passed term in a
// direct message to the invoking user.
func (ctx *Context) ReplyltDM(term i18n.Term) (*discord.Message, error) {
	return ctx.ReplylDM(term.AsConfig())
}

// ReplyEmbedDM replies with the passed discord.Embed in a direct message
// to the invoking user.
func (ctx *Context) ReplyEmbedDM(e discord.Embed) (*discord.Message, error) {
	return ctx.ReplyMessageDM(api.SendMessageData{Embed: &e})
}

// ReplyEmbedBuilderDM builds the discord.Embed from the passed
// embedutil.Builder and sends it in a direct message to the invoking user.
func (ctx *Context) ReplyEmbedBuilderDM(e *embedutil.Builder) (*discord.Message, error) {
	embed, err := e.Build(ctx.Localizer)
	if err != nil {
		return nil, err
	}

	return ctx.ReplyEmbedDM(embed)
}

// ReplyMessageDM sends the passed api.SendMessageData in a direct message to
// the invoking user.
func (ctx *Context) ReplyMessageDM(data api.SendMessageData) (msg *discord.Message, err error) {
	msg, err = ctx.Replier.ReplyDM(ctx, data)
	return msg, errorutil.WithStack(err)
}

// Edit edits the message with the passed id in the invoking channel.
// The message will be formatted as fmt.Sprint(content...).
func (ctx *Context) Edit(messageID discord.MessageID, content ...interface{}) (*discord.Message, error) {
	return ctx.EditMessage(messageID, api.EditMessageData{
		Content: option.NewNullableString(fmt.Sprint(content...)),
	})
}

// Editf edits the message with the passed id in the invoking channel.
// The message will be formatted as fmt.Sprintf(format, a...).
func (ctx *Context) Editf(messageID discord.MessageID, format string, a ...interface{}) (*discord.Message, error) {
	return ctx.Edit(messageID, fmt.Sprintf(format, a...))
}

// Editl edits the message with passed id in the invoking channel, by replacing
// it with the text generated from the passed *i18n.Config.
func (ctx *Context) Editl(messageID discord.MessageID, c *i18n.Config) (*discord.Message, error) {
	s, err := ctx.Localizer.Localize(c)
	if err != nil {
		return nil, err
	}

	return ctx.Edit(messageID, s)
}

// Editlt edits the message with the passed id in the invoking channel, by
// replacing it with the text generated from the passed i18n.Term.
func (ctx *Context) Editlt(messageID discord.MessageID, term i18n.Term) (*discord.Message, error) {
	return ctx.Editl(messageID, term.AsConfig())
}

// EditEmbed replaces the embed of the message with the passed id in the
// invoking channel.
func (ctx *Context) EditEmbed(messageID discord.MessageID, e discord.Embed) (*discord.Message, error) {
	return ctx.EditMessage(messageID, api.EditMessageData{Embed: &e})
}

// EditEmbedBuilder builds the discord.Embed from the passed
// *embedutil.Builder, and replaces the embed of the message with the passed
// id in the invoking channel.
func (ctx *Context) EditEmbedBuilder(messageID discord.MessageID, e *embedutil.Builder) (*discord.Message, error) {
	embed, err := e.Build(ctx.Localizer)
	if err != nil {
		return nil, err
	}

	return ctx.EditEmbed(messageID, embed)
}

// EditMessage sends the passed api.EditMessageData to the channel the command
// was originally sent in.
func (ctx *Context) EditMessage(messageID discord.MessageID, data api.EditMessageData) (*discord.Message, error) {
	msg, err := ctx.Replier.Edit(ctx, messageID, data)
	return msg, errorutil.WithStack(err)
}

// EditDM edits the message with the passed id in the direct message channel
// with the invoking user.
// The message will be formatted as fmt.Sprint(content...).
func (ctx *Context) EditDM(messageID discord.MessageID, content ...interface{}) (*discord.Message, error) {
	return ctx.EditMessageDM(messageID, api.EditMessageData{
		Content: option.NewNullableString(fmt.Sprint(content...)),
	})
}

// EditfDM edits the message with the passed id in the direct message channel
// with the invoking user.
// The message will be formatted as fmt.Sprintf(format, a...).
func (ctx *Context) EditfDM(messageID discord.MessageID, format string, a ...interface{}) (*discord.Message, error) {
	return ctx.EditDM(messageID, fmt.Sprintf(format, a...))
}

// EditlDM edits the message with passed id in the direct message channel with
// the invoking user, by replacing it with the text generated from the passed
// *i18n.Config.
func (ctx *Context) EditlDM(messageID discord.MessageID, c *i18n.Config) (*discord.Message, error) {
	s, err := ctx.Localizer.Localize(c)
	if err != nil {
		return nil, err
	}

	return ctx.EditDM(messageID, s)
}

// EditltDM edits the message with the passed id in the direct message channel
// with the invoking user, by replacing it with the text generated from the
// passed i18n.Term.
func (ctx *Context) EditltDM(messageID discord.MessageID, term i18n.Term) (*discord.Message, error) {
	return ctx.EditlDM(messageID, term.AsConfig())
}

// EditEmbedDM replaces the embed of the message with the passed id in the
// invoking channel.
func (ctx *Context) EditEmbedDM(messageID discord.MessageID, e discord.Embed) (*discord.Message, error) {
	return ctx.EditMessageDM(messageID, api.EditMessageData{Embed: &e})
}

// EditEmbedBuilderDM builds the discord.Embed from the passed
// *embedutil.Builder, and replaces the embed of the message with the passed
// id in the direct message channel with the invoking user.
func (ctx *Context) EditEmbedBuilderDM(messageID discord.MessageID, e *embedutil.Builder) (*discord.Message, error) {
	embed, err := e.Build(ctx.Localizer)
	if err != nil {
		return nil, err
	}

	return ctx.EditEmbedDM(messageID, embed)
}

// EditMessageDM sends the passed api.EditMessageData to the direct message
// channel with the invoking user.
func (ctx *Context) EditMessageDM(messageID discord.MessageID, data api.EditMessageData) (*discord.Message, error) {
	return ctx.Replier.EditDM(ctx, messageID, data)
}

// Guild returns the guild the command was invoked in.
func (ctx *Context) Guild() (*discord.Guild, error) {
	return ctx.GuildAsync()()
}

// Channel returns the *discord.Channel the command was invoked in.
func (ctx *Context) Channel() (*discord.Channel, error) {
	return ctx.ChannelAsync()()
}

// Self returns the *discord.Member that belongs to the bot.
// It will return (nil, nil) if the command was not invoked in a guild.
func (ctx *Context) Self() (*discord.Member, error) {
	return ctx.SelfAsync()()
}

// SelfPermissions checks if the bot has the passed permissions.
// If this command is executed in a direct message, constant.DMPermissions will
// be returned instead.
func (ctx *Context) SelfPermissions() (discord.Permissions, error) {
	if ctx.GuildID == 0 {
		return permutil.DMPermissions, nil
	}

	gf := ctx.GuildAsync()
	cf := ctx.ChannelAsync()

	s, err := ctx.Self()
	if err != nil {
		return 0, err
	}

	g, err := gf()
	if err != nil {
		return 0, err
	}

	ch, err := cf()
	if err != nil {
		return 0, err
	}

	return discord.CalcOverwrites(*g, *ch, *s), nil
}

// UserPermissions returns the permissions of the invoking user in the
// channel.
// If this command is executed in a direct message, constant.DMPermissions will
// be returned instead.
func (ctx *Context) UserPermissions() (discord.Permissions, error) {
	if ctx.GuildID == 0 {
		return permutil.DMPermissions, nil
	}

	gf := ctx.GuildAsync()

	ch, err := ctx.Channel()
	if err != nil {
		return 0, err
	}

	g, err := gf()
	if err != nil {
		return 0, err
	}

	return discord.CalcOverwrites(*g, *ch, *ctx.Member), nil
}

type (
	// Replier is the interface used to send replies to a command.
	//
	// This allows the user to define special behavior for commands, such as
	// the ability to delete answers after a set amount of time, after the bot
	// responds.
	Replier interface {
		// ReplyMessage sends a message in the invoking channel.
		Reply(ctx *Context, data api.SendMessageData) (*discord.Message, error)
		// ReplyDM sends the passed message in a direct message to the user.
		ReplyDM(ctx *Context, data api.SendMessageData) (*discord.Message, error)

		// Edit edits a message in the invoking channel
		Edit(ctx *Context, messageID discord.MessageID, data api.EditMessageData) (*discord.Message, error)
		// Edit edits a message sent in the direct message channel with the
		// invoking user.
		EditDM(ctx *Context, messageID discord.MessageID, data api.EditMessageData) (*discord.Message, error)
	}

	// DiscordDataProvider is an embeddable interface used to extend a Context
	// with additional information.
	DiscordDataProvider interface {
		// Guild returns a callback returning guild the message was sent in.
		// If this happened in a private channel, Guild will return nil, nil.
		GuildAsync() func() (*discord.Guild, error)
		// Self returns a callback returning the *discord.Member the bot
		// Channel returns a callback returning channel the message was sent
		// in.
		ChannelAsync() func() (*discord.Channel, error)
		// represents in the calling guild.
		// If this happened in a private channel, Self will return (nil, nil).
		SelfAsync() func() (*discord.Member, error)
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
		// To check if any of the plugin providers returned an error, call
		// UnavailablePluginProviders.
		// If that is the case, the data returned might be incomplete.
		Commands() []*RegisteredCommand
		// Modules returns all top-level modules sorted in ascending order by
		// name.
		//
		// To check if any of the plugin providers returned an error, call
		// UnavailablePluginProviders.
		// If that is the case, the data returned might be incomplete.
		Modules() []*RegisteredModule

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
		// To check if any of the plugin providers returned an error, call
		// UnavailablePluginProviders.
		// If that is the case, the module's description might not be available
		// or differ from the description that is used if all plugin providers
		// function properly.
		Module(Identifier) *RegisteredModule

		// FindCommand returns the RegisteredCommand with the passed invoke.
		//
		// It will return nil if no command matching the passed invoke was
		// found.
		//
		// To check if any of the plugin providers returned an error, call
		// UnavailablePluginProviders.
		FindCommand(invoke string) *RegisteredCommand
		// FindModule returns the RegisteredModule with the passed invoke.
		//
		// It will return nil if no module matching the passed invoke was
		// found.
		//
		// To check if any of the plugin providers returned an error, call
		// UnavailablePluginProviders.
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
		// Commands are the top-level commands of the repository.
		Commands []Command
		// Modules are the top-level modules of the repository.
		Modules []Module

		// Defaults are the global defaults for settings, the provider
		// uses.
		Defaults Defaults
	}

	// Defaults are the defaults used as fallback if a command does not define
	// a setting.
	Defaults struct {
		// ChannelTypes specifies the default channel types.
		ChannelTypes ChannelTypes
		// Restrictions is the default restriction func.
		Restrictions RestrictionFunc
		// Throttler is the default global throttler.
		// Note that the same Throttler will be shared across all commands that
		// don't define a custom one.
		Throttler Throttler
	}

	// UnavailablePluginProvider contains information about an unavailable
	// plugin provider.
	UnavailablePluginProvider struct {
		// Name is the name of the plugin provider.
		Name string
		// Error is the error returned by the plugin provider.
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
