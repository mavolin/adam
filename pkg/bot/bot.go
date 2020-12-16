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

	commands []plugin.Command
	modules  []plugin.Module

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
	gw.Identifier.Presence.Status = o.Status
	gw.Identifier.Shard = &o.Shard

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

func (b *Bot) AddModule(mod plugin.Module) {
	b.modules = append(b.modules, mod)
}
