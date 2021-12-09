// Package bot provides the Bot handling all commands.
package bot

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state/store"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/internal/resolved"
	"github.com/mavolin/adam/pkg/plugin"
)

// Bot is the bot executing all commands.
type Bot struct {
	State *state.State

	MiddlewareManager
	postMiddlewares MiddlewareManager

	pluginResolver *resolved.PluginResolver

	selfID discord.UserID

	// ----- Settings -----

	Owners []discord.UserID

	EditAge time.Duration

	ErrorHandler func(error, *state.State, *plugin.Context)

	PanicHandler             func(recovered interface{}, s *state.State, ctx *plugin.Context)
	MessageCreateMiddlewares []interface{}
	MessageUpdateMiddlewares []interface{}
}

// A PluginSourceFunc is the function used to retrieve additional plugins from
// other sources only available at runtime.
// A typical example for this would be custom commands or tags.
//
// PluginProviders will be called in the order they were added to a Bot, until
// one of the returns a matching plugin.
//
// If there are no plugins to return, all return values should be nil.
// If there is an error the returned plugins will be discarded, and the error
// will be noted in the Context of the command, available via
// Context.UnavailablePluginSources().
type PluginSourceFunc = resolved.PluginSourceFunc

// New creates a new Bot from the passed options.
// The Options.Token field must be set.
func New(o Options) (b *Bot, err error) {
	b = new(Bot)

	if err = o.SetDefaults(); err != nil {
		return nil, err
	}

	var activity *discord.Activity

	if o.ActivityName != "" {
		activity = &discord.Activity{
			Name: o.ActivityName,
			Type: o.ActivityType,
			URL:  o.ActivityURL,
		}
	}

	b.State, err = state.New(state.Options{
		Token:        o.Token,
		Status:       o.Status,
		Activity:     activity,
		Cabinet:      o.Cabinet,
		TotalShards:  o.TotalShards,
		ShardIDs:     o.ShardIDs,
		Gateways:     o.Gateways,
		HTTPClient:   o.HTTPClient,
		Rescale:      o.Rescale,
		ErrorHandler: o.StateErrorHandler,
		PanicHandler: o.StatePanicHandler,
	})
	if err != nil {
		return nil, err
	}

	self, err := b.State.Me()
	if err != nil {
		return nil, err
	}

	b.selfID = self.ID

	b.Owners = o.Owners
	b.EditAge = o.EditAge
	b.ErrorHandler = o.ErrorHandler
	b.PanicHandler = o.PanicHandler

	b.pluginResolver = resolved.NewPluginResolver(o.ArgParser)

	if !o.NoDefaultMiddlewares {
		b.AddMiddleware(CheckMessageType)

		if o.AllowBot {
			b.AddMiddleware(CheckHuman)
		}

		b.AddMiddleware(NewSettingsRetriever(o.SettingsProvider))
		b.AddMiddleware(CheckPrefix)
		b.AddMiddleware(FindCommand)
		b.AddMiddleware(CheckChannelTypes)
		b.AddMiddleware(CheckBotPermissions)
		b.AddMiddleware(NewThrottlerChecker(o.ThrottlerCancelChecker))

		b.AddPostMiddleware(CheckRestrictions)
		b.AddPostMiddleware(ParseArgs)
		b.AddPostMiddleware(InvokeCommand)
	}

	return b, nil
}

// Open opens a connection to the gateway and starts the bot.
//
// If no gateway.Intents were added to the State before opening, Open will
// derive intents from the registered handlers.
//
// gateway.IntentGuildMessages and gateway.IntentDirectMessages will always be
// added.
// Additionally, gateway.IntentGuilds will be added, if guild caching is
// enabled.
//
// Refer to the doc of State.Open to understand how the timeout is applied.
func (b *Bot) Open(timeout time.Duration) error {
	if i := b.State.Gateway.Identifier.Intents; i == nil || i == option.ZeroUint {
		b.AddIntents(b.State.DeriveIntents())
	}

	b.AddIntents(gateway.IntentGuildMessages)
	b.AddIntents(gateway.IntentDirectMessages)

	if b.State.Cabinet.GuildStore != store.Noop {
		b.AddIntents(gateway.IntentGuilds)
	}

	b.State.AddHandler(func(_ *state.State, e *event.MessageCreate) {
		b.Route(e.Base, &e.Message, e.Member)
	}, b.MessageCreateMiddlewares...)

	if b.EditAge > 0 {
		b.State.AddHandler(func(_ *state.State, e *event.MessageUpdate) {
			if time.Since(e.Timestamp.Time()) <= b.EditAge {
				b.Route(e.Base, &e.Message, e.Member)
			}
		}, b.MessageUpdateMiddlewares...)
	}

	return b.State.Open(timeout)
}

// Close closes all gateways handled by the bot.
//
// If an error occurs, Close will attempt to close all remaining gateways
// first, before returning. If multiple errors occur during that process, a
// MultiError will be returned.
//
// The passed context will only be checked while waiting for all event handlers
// to finish.
// Even if the context expires, Close guarantees that all gateways are closed,
// except if errors occurred.
func (b *Bot) Close(ctx context.Context) error {
	return b.State.Close(ctx)
}

// AddIntents adds the passed gateway.Intents to the bot.
func (b *Bot) AddIntents(i gateway.Intents) {
	b.State.Gateway.AddIntents(i)
}

// AddCommand adds the passed top-level command to the bot.
func (b *Bot) AddCommand(cmd plugin.Command) {
	b.pluginResolver.AddBuiltInCommand(cmd)
}

// AddModule adds the passed top-level module to the Bot.
//
// If automatic handler adding is enabled, all methods of the Module
// representing a handler func will be added to the State's event handler.
// The same goes for all sub-modules and sub-commands of the module.
func (b *Bot) AddModule(mod plugin.Module) {
	b.pluginResolver.AddBuiltInModule(mod)
}

// TryAddPostMiddleware adds a middleware to the Bot that is invoked after all
// command and module middlewares were called.
// The order of invocation of post middlewares is the same as the order they
// were added in.
//
// If the middleware's type is invalid, TryAddMiddleware will return
// ErrMiddleware.
//
// Valid middleware types are:
//	• func(*state.State, interface{})
//	• func(*state.State, interface{}) error
//	• func(*state.State, *event.Base)
//	• func(*state.State, *event.Base) error
//	• func(*state.State, *state.MessageCreateEvent)
//	• func(*state.State, *state.MessageCreateEvent) error
//	• func(*state.State, *state.MessageUpdateEvent)
//	• func(*state.State, *state.MessageUpdateEvent) error
//	• func(next CommandFunc) CommandFunc
func (b *Bot) TryAddPostMiddleware(f interface{}) error {
	return b.postMiddlewares.TryAddMiddleware(f)
}

// AddPostMiddleware is the same as TryAddPostMiddleware, but panics if
// TryAddPostMiddleware returns an error.
func (b *Bot) AddPostMiddleware(f interface{}) {
	b.postMiddlewares.AddMiddleware(f)
}

// AddPluginSource adds the passed PluginSourceFunc under the passed unique
// name.
// The name is similar to a key and can be used later on to distinguish between
// different plugin sources.
// It is typically snake_case.
//
// 'built_in' is reserved for built-in plugins, and AddPluginSource will panic
// if attempting to use it.
//
// The plugin sources will be used in the order they are added in.
func (b *Bot) AddPluginSource(name string, f PluginSourceFunc) {
	if f == nil {
		return
	}

	if name == plugin.BuiltInSource {
		panic("you cannot use " + plugin.BuiltInSource + " as plugin provider")
	}

	b.pluginResolver.AddSource(name, f)
}
