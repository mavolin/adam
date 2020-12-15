// Package bot provides the Bot handling all commands.
package bot

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/session"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

type Bot struct {
	State *state.State

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
	*MiddlewareManager
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
func (b *Bot) Open() error {
	if b.State.Gateway.Identifier.Intents == 0 {
		// todo: derive intents
		b.AddIntents(gateway.IntentGuilds)
	}

	if err := b.State.Open(); err != nil {
		return err
	}

	return nil
}

// AddIntents adds the passed gateway.Intents to the bot.
func (b *Bot) AddIntents(i gateway.Intents) {
	b.State.Gateway.AddIntents(i)
}
