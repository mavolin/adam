package plugin

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/mavolin/disstate/v4/pkg/event"

	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/utils/permutil"
)

// Context contains context information about a command.
type Context struct {
	// Message is the invoking message.
	discord.Message
	// Member is the invoking member, or nil if the command was invoked in a
	// direct message.
	*discord.Member

	// Base is the *event.Base of the MessageCreateEvent or MessageUpdateEvent
	// that triggered the invoke.
	*event.Base

	// Localizer is the localizer set to the guild's or user's language.
	*i18n.Localizer

	// InvokeIndex is the starting index of the invoke as found in Content.
	// The invoke ends at ArgsIndex-1 and is trailed by whitespace (' ', '\n').
	InvokeIndex int
	// ArgsIndex is the starting index of the argument as found in Content.
	ArgsIndex int

	// Args contains the arguments supplied to the bot.
	// They are guaranteed to be valid and parsed according to the type spec.
	Args Args
	// Flags contains the flags supplied to the bot.
	// They are guaranteed to be valid and parsed according to the type spec.
	Flags Flags

	// InvokedCommand is the ResolvedCommand that is being invoked.
	InvokedCommand ResolvedCommand

	// Prefixes contains the prefixes of the bot as defined for the invoking
	// guild or user.
	// It does not include the bot's mention, which is always a valid
	// prefix.
	// It may be empty, in which case the command was invoked using the bot's
	// mention.
	//
	// Note that direct messages do not require prefixes.
	// However, the bot's mention or any other prefix returned by the bot's
	// bot.SettingsProvider (as stored in this variable), will be stripped if
	// the message starts with such.
	Prefixes []string

	// BotOwnerIDs contains the ids of the bot owners, as defined in the bot's
	// bot.Options.
	BotOwnerIDs []discord.UserID

	// Replier is the interface used to send replies to a command.
	//
	// Defaults to replier.WrapState, as found in impl/replier.
	Replier Replier

	// Provider is an embedded interface that provides access to the commands
	// and Modules of the Bot, as well as the runtime commands and modules
	// for the guild.
	Provider

	// ErrorHandler is an embedded interface that provides error handling
	// capabilities to the command.
	ErrorHandler

	// DiscordDataProvider is an embedded interface that gives direct access to
	// common data types needed during execution.
	// Its asynchronous methods are supplemented by blocking methods provided
	// by the context.
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

// UsedPrefix returns the prefix used to invoke the command.
func (ctx *Context) UsedPrefix() string {
	return strings.TrimRight(ctx.Content[:ctx.InvokeIndex], shared.Whitespace)
}

// RawInvoke returns the raw invoke stripped of prefix and args as the user
// typed it.
func (ctx *Context) RawInvoke() string {
	return strings.TrimRight(ctx.Content[ctx.InvokeIndex:ctx.ArgsIndex], shared.Whitespace)
}

// RawArgs returns the raw arguments, as the user typed them.
func (ctx *Context) RawArgs() string {
	return ctx.Content[ctx.ArgsIndex:]
}

// Reply replies with the passed message in the channel the command was
// originally sent in.
// The message will be formatted as fmt.Sprint(content...).
func (ctx *Context) Reply(content ...interface{}) (*discord.Message, error) {
	return ctx.ReplyMessage(api.SendMessageData{Content: fmt.Sprint(content...)})
}

// Replyf replies with the passed message in the channel the command was
// originally sent in.
// The message will be formatted as fmt.Sprintf(format, a...).
func (ctx *Context) Replyf(format string, a ...interface{}) (*discord.Message, error) {
	return ctx.ReplyMessage(api.SendMessageData{Content: fmt.Sprintf(format, a...)})
}

// Replyl replies with the message generated from the passed *i18n.Config in
// the channel the command was originally sent in.
func (ctx *Context) Replyl(c *i18n.Config) (*discord.Message, error) {
	s, err := ctx.Localizer.Localize(c)
	if err != nil {
		return nil, err
	}

	return ctx.Reply(s)
}

// ReplyEmbeds replies with the passed discord.Embeds in the channel the
// command was originally sent in.
func (ctx *Context) ReplyEmbeds(embeds ...discord.Embed) (*discord.Message, error) {
	return ctx.ReplyMessage(api.SendMessageData{Embeds: embeds})
}

// ReplyMessage sends the passed api.SendMessageData to the channel the command
// was originally sent in.
func (ctx *Context) ReplyMessage(data api.SendMessageData) (*discord.Message, error) {
	return ctx.Replier.Reply(ctx, data)
}

// ReplyDM replies with the passed message in a direct message to the invoking
// user.
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

// ReplylDM replies with the message translated from the passed *i18n.Config in
// a direct message to the invoking user.
func (ctx *Context) ReplylDM(c *i18n.Config) (*discord.Message, error) {
	s, err := ctx.Localizer.Localize(c)
	if err != nil {
		return nil, err
	}

	return ctx.ReplyDM(s)
}

// ReplyEmbedsDM replies with the passed discord.Embeds in a direct message
// to the invoking user.
func (ctx *Context) ReplyEmbedsDM(embeds ...discord.Embed) (*discord.Message, error) {
	return ctx.ReplyMessageDM(api.SendMessageData{Embeds: embeds})
}

// ReplyMessageDM sends the passed api.SendMessageData in a direct message to
// the invoking user.
func (ctx *Context) ReplyMessageDM(data api.SendMessageData) (msg *discord.Message, err error) {
	return ctx.Replier.ReplyDM(ctx, data)
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

// EditEmbeds replaces the embeds of the message with the passed id in the
// invoking channel.
func (ctx *Context) EditEmbeds(messageID discord.MessageID, embeds ...discord.Embed) (*discord.Message, error) {
	return ctx.EditMessage(messageID, api.EditMessageData{Embeds: &embeds})
}

// EditMessage sends the passed api.EditMessageData to the channel the command
// was originally sent in.
func (ctx *Context) EditMessage(messageID discord.MessageID, data api.EditMessageData) (*discord.Message, error) {
	return ctx.Replier.Edit(ctx, messageID, data)
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

// EditEmbedsDM replaces the embeds of the message with the passed id in the
// invoking channel.
func (ctx *Context) EditEmbedsDM(messageID discord.MessageID, embeds ...discord.Embed) (*discord.Message, error) {
	return ctx.EditMessageDM(messageID, api.EditMessageData{Embeds: &embeds})
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

// ParentChannel returns the parent *discord.Channel the command was invoked
// in.
func (ctx *Context) ParentChannel() (*discord.Channel, error) {
	return ctx.ParentChannelAsync()()
}

// Self returns the *discord.Member that belongs to the bot.
// It will return (nil, nil) if the command was not invoked in a guild.
func (ctx *Context) Self() (*discord.Member, error) {
	return ctx.SelfAsync()()
}

// SelfPermissions returns the discord.Permissions the bot has in the invoking
// channel.
// If the command is executed in a direct message, permutil.DMPermissions will
// be returned instead.
func (ctx *Context) SelfPermissions() (discord.Permissions, error) {
	if ctx.GuildID == 0 {
		return permutil.DMPermissions, nil
	}

	gf := ctx.GuildAsync()
	sf := ctx.SelfAsync()

	ch, err := ctx.Channel()
	if err != nil {
		return 0, err
	}

	if ch.Type == discord.GuildNewsThread || ch.Type == discord.GuildPublicThread ||
		ch.Type == discord.GuildPrivateThread {
		ch, err = ctx.ParentChannel()
		if err != nil {
			return 0, nil
		}
	}

	s, err := sf()
	if err != nil {
		return 0, err
	}

	g, err := gf()
	if err != nil {
		return 0, err
	}

	return discord.CalcOverwrites(*g, *ch, *s), nil
}

// UserPermissions returns the permissions of the invoking user in the
// channel.
// If this command is executed in a direct message, permutil.DMPermissions will
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
		// Reply sends a message in the invoking channel.
		Reply(ctx *Context, data api.SendMessageData) (*discord.Message, error)
		// ReplyDM sends the passed message in a direct message to the user.
		ReplyDM(ctx *Context, data api.SendMessageData) (*discord.Message, error)

		// Edit edits a message in the invoking channel
		Edit(ctx *Context, messageID discord.MessageID, data api.EditMessageData) (*discord.Message, error)
		// EditDM edits a message sent in the direct message channel with the
		// invoking user.
		EditDM(ctx *Context, messageID discord.MessageID, data api.EditMessageData) (*discord.Message, error)
	}

	// DiscordDataProvider is an embeddable interface used to extend a Context
	// with additional information.
	DiscordDataProvider interface {
		// GuildAsync returns a callback returning guild the message was sent
		// in.
		// If the command was invoked in a private channel, Guild will return
		// (nil, nil).
		GuildAsync() func() (*discord.Guild, error)
		// ChannelAsync returns a callback returning the Channel the message
		// was sent in.
		ChannelAsync() func() (*discord.Channel, error)
		// ParentChannelAsync returns a callback returning the parent of the
		// Channel the message was sent in.
		ParentChannelAsync() func() (*discord.Channel, error)
		// SelfAsync returns a callback returning the member object of the bot
		// in the calling guild.
		// If the command was used in a private channel, SelfAsync will return
		// (nil, nil).
		SelfAsync() func() (*discord.Member, error)
	}

	// Provider provides copies of the plugins of the bot.
	// The returned slices can therefore be freely modified.
	//
	// Copies are only created on call of one of the methods.
	Provider interface {
		// PluginSources returns a slice of Sources containing all commands and
		// modules of the bot.
		// PluginSources()[0] contains the built-in plugins of the bot, and is
		// named BuiltInSource.
		//
		// To check if any runtime plugin sources returned an error, call
		// UnavailablePluginSources.
		PluginSources() []Source

		// Commands returns all top-level commands sorted in ascending order by
		// name.
		//
		// To check if any of the plugin sources returned an error, call
		// UnavailablePluginProviders.
		// If that is the case, the data returned might be incomplete.
		Commands() []ResolvedCommand
		// Modules returns all top-level modules sorted in ascending order by
		// name.
		//
		// To check if any of the plugin sources returned an error, call
		// UnavailablePluginProviders.
		// If that is the case, the data returned might be incomplete.
		Modules() []ResolvedModule

		// Command returns the ResolvedCommand with the passed ID.
		//
		// Note that Identifiers may only consist of the command's name, not
		// their alias.
		//
		// It will return nil if no command matching the identifier was found.
		//
		// To check if any of the runtime plugin sources returned an error,
		// call UnavailablePluginProviders.
		Command(ID) ResolvedCommand
		// Module returns the ResolvedModule with the passed ID.
		//
		// It will return nil if no module matching the identifier was found.
		//
		// To check if any of the plugin sources returned an error, call
		// UnavailablePluginProviders.
		// If that is the case, the module's description might not be available
		// or differ from the description that is used if all plugin sources
		// function properly.
		Module(ID) ResolvedModule

		// FindCommand returns the ResolvedCommand with the passed invoke.
		//
		// It will return nil if no command matching the passed invoke was
		// found.
		//
		// To check if any of the plugin sources returned an error, call
		// UnavailablePluginProviders.
		FindCommand(invoke string) ResolvedCommand
		// FindCommandWithArgs is the same as FindCommand, but allows invoke
		// to contain trailing arguments.
		//
		// If a command is found, it is returned alongside the arguments.
		// Otherwise, (nil, "") will be returned.
		FindCommandWithArgs(invoke string) (cmd ResolvedCommand, args string)
		// FindModule returns the ResolvedModule with the passed invoke.
		//
		// It will return nil if no module matching the passed invoke was
		// found.
		//
		// To check if any of the plugin sources returned an error, call
		// UnavailablePluginProviders.
		// If that is the case, the module's description might not be available
		// or differ from the description that is used if all plugin sources
		// function properly.
		FindModule(invoke string) ResolvedModule

		// UnavailablePluginSources returns a list of all unavailable plugin
		// sources.
		// If no runtime plugins were requested yet, it will request them and
		// return the list of unavailable ones.
		//
		// If the length of the returned slice is 0, all plugin sources are
		// available.
		UnavailablePluginSources() []UnavailableSource
	}

	// Source is the struct returned by Provider.PluginSources.
	// It contains the top-level plugins of a single repository.
	Source struct {
		// Name is the name of the source that provides these plugins.
		Name string
		// Commands are the top-level commands of the repository.
		Commands []Command
		// Modules are the top-level modules of the repository.
		Modules []Module
	}

	// UnavailableSource contains information about an unavailable plugin
	// source.
	UnavailableSource struct {
		// Name is the name of the plugin source.
		Name string
		// Error is the error returned by the plugin source.
		Error error
	}

	// ErrorHandler is an embedded interface used to provide error handling
	// capabilities through a Context.
	ErrorHandler interface {
		// HandleError hands the error to the bot's error handler.
		HandleError(err error)
		// HandleErrorSilently wraps the error using errors.Silent and hands it
		// to the bot's error handler.
		HandleErrorSilently(err error)
	}
)
