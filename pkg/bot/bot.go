// Package bot provides the Bot handling all commands.
package bot

import (
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/session"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// Bot is the bot executing all commands.
type Bot struct {
	State *state.State
	*MiddlewareManager

	commands        []plugin.Command
	modules         []plugin.Module
	pluginProviders []*pluginProvider

	// ----- Settings -----

	PrefixProvider PrefixProvider
	Owners         []discord.UserID
	EditThreshold  uint

	AllowBot   bool
	SendTyping bool

	PluginDefaults plugin.Defaults

	ThrottlerErrorCheck func(error) bool

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

	gw := gateway.NewCustomGateway(o.GatewayURL, o.Token)

	gw.WSTimeout = o.GatewayTimeout
	gw.WS.Timeout = o.GatewayTimeout
	gw.ErrorLog = o.GatewayErrorHandler
	gw.Identifier = &gateway.Identifier{
		IdentifyData: gateway.IdentifyData{
			Shard:    &o.Shard,
			Presence: &gateway.UpdateStatusData{Status: o.Status},
		},
	}

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

	b.PrefixProvider = o.PrefixProvider
	b.Owners = o.Owners
	b.EditThreshold = o.EditThreshold
	b.AllowBot = o.AllowBot
	b.SendTyping = o.SendTyping
	b.PluginDefaults = plugin.Defaults{
		ChannelTypes:   o.DefaultChannelTypes,
		BotPermissions: o.DefaultBotPermissions,
		Restrictions:   o.DefaultRestrictions,
		Throttler:      o.DefaultThrottler,
	}
	b.ThrottlerErrorCheck = o.ThrottlerErrorCheck
	b.ErrorHandler = o.ErrorHandler
	b.PanicHandler = o.PanicHandler

	return b, nil
}

// Open opens a connection to the gateway and starts the bot.
//
// If no gateway.Intents were added to the State before opening, Open will
// derive intents from the registered handlers.
// Additionally, gateway.IntentGuilds will be added, to ensure caching of guild
// data.
func (b *Bot) Open() error {
	if b.State.Gateway.Identifier.Intents == 0 {
		b.AddIntents(b.State.DeriveIntents())
		b.AddIntents(gateway.IntentGuilds)
	}

	return b.State.Open()
}

// Wait blockingly waits for SIGINT and returns, when it receives it.
func (b *Bot) Wait() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
}

// AddIntents adds the passed gateway.Intents to the bot.
func (b *Bot) AddIntents(i gateway.Intents) {
	b.State.Gateway.AddIntents(i)
}

// AddCommand adds the passed command to the bot.
func (b *Bot) AddCommand(cmd plugin.Command) {
	b.commands = append(b.commands, cmd)
}

// AddModule adds the passed module to the Bot.
func (b *Bot) AddModule(mod plugin.Module) {
	b.modules = append(b.modules, mod)
}

// AddPluginProvider adds the passed PluginProvider under the passed name.
// The is similar to a key and can be used later on to distinguish between
// different plugin providers.
// It is typically snake_case.
//
// 'built_in' is not allowed as name, and AddPluginProvider will panic if
// attempting to use it.
//
// If there is another plugin provider with the passed name, it will be removed
// first.
//
// The plugin providers will be used in the order they are added in.
func (b *Bot) AddPluginProvider(name string, p PluginProvider, defaults plugin.Defaults) {
	if name == plugin.BuiltInProvider {
		panic("you cannot use " + name + " as plugin provider")
	}

	if p == nil {
		return
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
