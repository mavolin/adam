// Package bot provides the Bot handling all commands.
package bot

import (
	"regexp"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/session"
	"github.com/diamondburned/arikawa/v2/state/store"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// Bot is the bot executing all commands.
type Bot struct {
	State *state.State

	*MiddlewareManager
	postMiddlewares *MiddlewareManager

	commands        []plugin.Command
	modules         []plugin.Module
	pluginProviders []*pluginProvider

	selfID            discord.UserID
	selfMentionRegexp *regexp.Regexp

	// ----- Settings -----

	SettingsProvider SettingsProvider
	Owners           []discord.UserID
	EditAge          time.Duration

	AllowBot bool

	autoOpen        bool
	autoAddHandlers bool
	manualChecks    bool

	AsyncPluginProviders bool

	PluginDefaults plugin.Defaults

	ThrottlerCancelChecker func(error) bool

	ErrorHandler func(error, *state.State, *plugin.Context)
	PanicHandler func(recovered interface{}, s *state.State, ctx *plugin.Context)
}

type pluginProvider struct {
	name     string
	provider PluginProvider
	defaults plugin.Defaults
}

// Plugin provider is the function used by plugin providers.
// PluginProviders will be called in the order they were added to a Bot, until
// one of the returns a matching plugin.
//
// If there are no plugins that match the context of the message, the
// PluginProvider should return (nil, nil, nil).
// If there is an error the returned plugins will be discarded, and the error
// will be noted in the Context of the command, available via
// Context.UnavailablePluginProviders().
type PluginProvider func(*state.Base, *discord.Message) ([]plugin.Command, []plugin.Module, error)

// New creates a new Bot from the passed options.
// The Options.Token field must be set.
func New(o Options) (*Bot, error) {
	b := new(Bot)

	if err := o.SetDefaults(); err != nil {
		return nil, err
	}

	gw := gateway.NewCustomGateway(o.GatewayURL, "Bot "+o.Token)

	gw.WSTimeout = o.GatewayTimeout
	gw.WS.Timeout = o.GatewayTimeout
	gw.ErrorLog = o.GatewayErrorHandler
	gw.Identifier.IdentifyData.Shard = &o.Shard
	gw.Identifier.IdentifyData.Presence = &gateway.UpdateStatusData{Status: o.Status}

	if len(o.ActivityName) > 0 {
		gw.Identifier.Presence.Activities = &[]discord.Activity{
			{
				Name: o.ActivityName,
				Type: o.ActivityType,
				URL:  o.ActivityURL,
			},
		}
	}

	b.State = state.NewFromSession(session.NewWithGateway(gw), o.Cabinet)
	b.State.ErrorHandler = o.StateErrorHandler
	b.State.PanicHandler = o.StatePanicHandler
	b.MiddlewareManager = new(MiddlewareManager)

	b.SettingsProvider = o.SettingsProvider
	b.Owners = o.Owners
	b.EditAge = o.EditAge
	b.AllowBot = o.AllowBot
	b.manualChecks = o.ManualChecks
	b.autoOpen = !o.NoAutoOpen
	b.autoAddHandlers = o.AutoAddHandlers
	b.PluginDefaults = plugin.Defaults{
		ChannelTypes: o.DefaultChannelTypes,
		Restrictions: o.DefaultRestrictions,
		Throttler:    o.DefaultThrottler,
	}
	b.ThrottlerCancelChecker = o.ThrottlerCancelChecker
	b.AsyncPluginProviders = o.AsyncPluginProviders
	b.ErrorHandler = o.ErrorHandler
	b.PanicHandler = o.PanicHandler

	if !b.manualChecks {
		b.MustAddMiddleware(CheckChannelTypes)
		b.MustAddMiddleware(CheckBotPermissions)
		b.MustAddMiddleware(NewThrottlerChecker(b.ThrottlerCancelChecker))

		b.MustAddPostMiddleware(CheckRestrictions)
		b.MustAddPostMiddleware(ParseArgs)
	}

	return b, nil
}

// Open opens a connection to the gateway and starts the bot.
//
// If no gateway.Intents were added to the State before opening, Open will
// derive intents from the registered handlers.
// Additionally, gateway.IntentGuilds will be added, if guild caching is
// enabled.
func (b *Bot) Open() error {
	if b.State.Gateway.Identifier.Intents == 0 {
		b.AddIntents(b.State.DeriveIntents())
		b.AddIntents(gateway.IntentGuildMessages)
		b.AddIntents(gateway.IntentDirectMessages)

		if b.State.Cabinet.GuildStore != store.Noop {
			b.AddIntents(gateway.IntentGuilds)
		}
	}

	if b.autoOpen {
		for _, cmd := range b.commands {
			if err := b.callOpen(cmd); err != nil {
				return err
			}
		}

		for _, mod := range b.modules {
			if err := b.openModule(mod); err != nil {
				return err
			}
		}
	}

	done := make(chan struct{})

	rm := b.State.MustAddHandler(func(_ *state.State, r *state.ReadyEvent) {
		b.selfID = r.User.ID
		b.selfMentionRegexp = regexp.MustCompile("^<@!?" + r.User.ID.String() + ">")

		done <- struct{}{}
	})

	b.State.MustAddHandler(func(_ *state.State, e *state.MessageCreateEvent) {
		b.Route(e.Base, &e.Message, e.Member)
	})

	if b.EditAge > 0 {
		b.State.MustAddHandler(func(_ *state.State, e *state.MessageUpdateEvent) {
			if time.Since(e.Timestamp.Time()) <= b.EditAge {
				b.Route(e.Base, &e.Message, e.Member)
			}
		})
	}

	err := b.State.Open()
	if err != nil {
		return err
	}

	<-done
	rm()

	return nil
}

func (b *Bot) openModule(mod plugin.Module) error {
	for _, cmd := range mod.Commands() {
		if err := b.callOpen(cmd); err != nil {
			return err
		}
	}

	for _, mod := range mod.Modules() {
		if err := b.openModule(mod); err != nil {
			return err
		}
	}

	return nil
}

// callOpen tries to call i.Open.
// Open may have an optional *Bot argument and an optional error return.
//
// If none of i's methods match those parameters, the function is a no-op.
//
// An error will only be returned if i.Open returns it.
func (b *Bot) callOpen(i interface{}) error {
	switch opener := i.(type) {
	case interface{ Open() }:
		opener.Open()
	case interface{ Open(*Bot) }:
		opener.Open(b)
	case interface{ Open() error }:
		return opener.Open()
	case interface{ Open(*Bot) error }:
		return opener.Open(b)
	}

	return nil
}

// Close closes the websocket connection to Discord's gateway.
func (b *Bot) Close() error {
	if err := b.State.Close(); err != nil {
		return err
	}

	if !b.autoOpen {
		return nil
	}

	for _, cmd := range b.commands {
		if err := b.callClose(cmd); err != nil {
			return err
		}
	}

	for _, mod := range b.modules {
		if err := b.closeModule(mod); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) closeModule(mod plugin.Module) error {
	for _, cmd := range mod.Commands() {
		if err := b.callClose(cmd); err != nil {
			return err
		}
	}

	for _, mod := range mod.Modules() {
		if err := b.closeModule(mod); err != nil {
			return err
		}
	}

	return nil
}

// callClose tries to call i.Close.
// Close may have an optional *Bot argument and an optional error return.
//
// If none of i's methods match those parameters, the function is a no-op.
//
// An error will only be returned if i.Close returns it.
func (b *Bot) callClose(i interface{}) error {
	switch closer := i.(type) {
	case interface{ Close() }:
		closer.Close()
	case interface{ Close(*Bot) }:
		closer.Close(b)
	case interface{ Close() error }:
		return closer.Close()
	case interface{ Close(*Bot) error }:
		return closer.Close(b)
	}

	return nil
}

// AddIntents adds the passed gateway.Intents to the bot.
func (b *Bot) AddIntents(i gateway.Intents) {
	b.State.Gateway.AddIntents(i)
}

// AddCommand adds the passed top-level command to the bot.
//
// If automatic handler adding is enabled, all methods of the Command
// representing a handler func will be added to the State's event handler.
func (b *Bot) AddCommand(cmd plugin.Command) {
	b.commands = append(b.commands, cmd)

	if b.autoAddHandlers {
		b.State.AutoAddHandlers(b.commands)
	}
}

// AddModule adds the passed top-level module to the Bot.
//
// If automatic handler adding is enabled, all methods of the Module
// representing a handler func will be added to the State's event handler.
// The same goes for all sub-modules and sub-commands of the module.
func (b *Bot) AddModule(mod plugin.Module) {
	b.modules = append(b.modules, mod)

	if b.autoAddHandlers {
		b.autoAddModuleHandlers(mod)
	}
}

func (b *Bot) autoAddModuleHandlers(mod plugin.Module) {
	for _, cmd := range mod.Commands() {
		b.State.AutoAddHandlers(cmd)
	}

	for _, mod := range mod.Modules() {
		b.autoAddModuleHandlers(mod)
	}
}

// AddPostMiddleware adds a middleware to the Bot, that is invoked after all
// command and module middlewares were called.
// The order of invocation of post middlewares is the same as the order they
// were added in.
//
// If the middleware's type is invalid, AddMiddleware will return
// ErrMiddleware.
//
// Valid middleware types are:
//	• func(*state.State, interface{})
//	• func(*state.State, interface{}) error
//	• func(*state.State, *state.Base)
//	• func(*state.State, *state.Base) error
//	• func(*state.State, *state.MessageCreateEvent)
//	• func(*state.State, *state.MessageCreateEvent) error
//	• func(*state.State, *state.MessageUpdateEvent)
//	• func(*state.State, *state.MessageUpdateEvent) error
//	• func(next CommandFunc) CommandFunc
func (b *Bot) AddPostMiddleware(f interface{}) error {
	return b.postMiddlewares.AddMiddleware(f)
}

// MustAddPostMiddleware is the same as AddPostMiddleware, but panics if
// AddPostMiddleware returns an error.
func (b *Bot) MustAddPostMiddleware(f interface{}) {
	b.postMiddlewares.MustAddMiddleware(f)
}

// AddPluginProvider adds the passed PluginProvider under the passed name.
// The is similar name to a key and can be used later on to distinguish between
// different plugin providers.
// It is typically snake_case.
//
// 'built_in' is reserved for built-in plugins, and AddPluginProvider will
// panic if attempting to use it.
//
// If there is another plugin provider with the passed name, it will be removed
// first.
//
// If defaults.ChannelTypes is 0, it will be set to plugin.AllChannels.
//
// The plugin providers will be used in the order they are added in.
func (b *Bot) AddPluginProvider(name string, p PluginProvider, defaults plugin.Defaults) {
	if p == nil {
		return
	}

	if name == plugin.BuiltInProvider {
		panic("you cannot use " + plugin.BuiltInProvider + " as plugin provider")
	}

	if defaults.ChannelTypes == 0 {
		defaults.ChannelTypes = plugin.AllChannels
	}

	for i, rp := range b.pluginProviders {
		if rp.name == name {
			b.pluginProviders = append(b.pluginProviders[:i], b.pluginProviders[i+1:]...)
		}
	}

	b.pluginProviders = append(b.pluginProviders, &pluginProvider{
		name:     name,
		provider: p,
		defaults: defaults,
	})
}
