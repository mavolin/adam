package bot

import (
	"log"
	"runtime/debug"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state/store"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/gorilla/websocket"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/adam/pkg/plugin"
)

// Options contains different configurations for a Bot.
//nolint:maligned // only one-time use anyway, ordered by importance, we can take the (temporary) few bytes
type Options struct {
	// Token is the bot token without the 'Bot ' prefix.
	//
	// This field is required.
	Token string

	// SettingsProvider is the settings provider for the bot.
	// If left nil, only the mention prefix will be usable.
	//
	// Default: StaticSettings()
	SettingsProvider SettingsProvider
	// Owners are the ids of the bot owners.
	// These are accessible through plugin.Context.BotOwnerIDs.
	//
	// Default: nil
	Owners []discord.UserID
	// EditAge is the oldest age an edit message may have, to trigger a
	// command.
	// If a message older than EditAge is edited, it will be ignored.
	// If this is set to 0 or less, edited messages won't be watched.
	//
	// Default: 0
	EditAge time.Duration

	// Status is the status of the bot.
	//
	// Default: gateway.OnlineStatus
	Status discord.Status
	// ActivityType is the type of activity.
	// ActivityName must be set for this to take effect.
	//
	// Default: discord.GameActivity
	ActivityType discord.ActivityType
	// ActivityName is the name of the activity the bot will display, if any.
	// If this left empty, the bot won't display any activity.
	//
	// Default: None
	ActivityName string
	// ActivityURL is the URL of the activity.
	// Currently, this is only used if the activity is set to Streaming.
	//
	// Default: None
	ActivityURL discord.URL

	// ArgParser is the plugin.ArgParser used to parse the arguments of all
	// commands that don't define a custom one.
	//
	// Default: &arg.DelimiterParser{Delimiter: ','}
	ArgParser plugin.ArgParser

	// AllowBot specifies whether bots may trigger commands.
	//
	// Settings this field has no effect if NoDefaultMiddlewares is set to
	// true.
	//
	// Default: false
	AllowBot bool

	// ThrottlerCancelChecker is the function run every time a command returns
	// with a non-nil error.
	// If the function returns true, the command's throttler will not count the
	// invoke.
	//
	// Settings this field has no effect if NoDefaultMiddlewares is set to
	// true.
	//
	// Default: DefaultThrottlerCancelCheck
	ThrottlerCancelChecker func(error) bool

	// Cabinet is the store.Cabinet used for caching.
	// Use store.NoopCabinet to deactivate caching.
	//
	// Default: defaultstore.New()
	Cabinet *store.Cabinet

	// GatewayErrorHandler is the error handler of the gateway.
	//
	// Default: DefaultGatewayErrorHandler
	GatewayErrorHandler func(error)

	// StateErrorHandler is the error handler of the *state.State, called if an
	// event handler returns with an error.
	//
	// Default: func(err error) { log.Println("event handler:", err.String()) }
	StateErrorHandler func(error)
	// StatePanicHandler is the panic handler of the *state.State, called if an
	// event handler panics.
	//
	// Default:
	// 	func(rec interface{}) {
	// 		log.Printf("event handler: panic: %+v\n%s\n", rec)
	//	}
	StatePanicHandler func(recovered interface{})

	// ErrorHandler is the handler called if a command returns with a non-nil
	// error.
	//
	// Default: DefaultErrorHandler
	ErrorHandler func(error, *state.State, *plugin.Context)
	// PanicHandler is the handler called if a command panics.
	//
	// Default: DefaultPanicHandler
	PanicHandler func(recovered interface{}, s *state.State, ctx *plugin.Context)

	// HTTPClient is the http client that will be used to make requests.
	//
	// Default: httputil.NewClient()
	HTTPClient *httputil.Client

	// MessageCreateMiddlewares are the middlewares invoked before routing the
	// command, if the command was received through a message create event.
	//
	// The signature of the middleware funcs must satisfy the requirements
	// of state middlewares.
	MessageCreateMiddlewares []interface{}
	// MessageCreateMiddlewares are the middlewares invoked before routing the
	// command, if the command was received through a message update event.
	//
	// The signature of the middleware funcs must satisfy the requirements
	// of state middlewares.
	MessageUpdateMiddlewares []interface{}

	// TotalShards is the total number of shards.
	// If it is <= 0, the recommended number of shards will be used.
	TotalShards int
	// ShardIDs are the custom shard ids this Bot instance will use.
	//
	// If setting this, you also need to set TotalShards.
	//
	// Default: 0..TotalShards
	ShardIDs []int

	// Gateways are the initial gateways to use.
	// It is an alternative to TotalShards and ShardIDs, and you shouldn't set
	// both.
	Gateways []*gateway.Gateway

	// Rescale is the function called if Discord closes any of the gateways
	// with a 4011 close code aka. 'Sharding Required'.
	//
	// Usage
	//
	// To update the state's shard manager, you must call update.
	// All zero-value options in the Options you give to update, will be set to
	// the options you used when initially creating the state.
	// This does not apply to TotalShards, ShardIDs, and Gateways, which will
	// assume the defaults described in their respective documentation.
	// Furthermore, setting ErrorHandler or PanicHandler will have no effect.
	//
	// After calling update, you should reopen the state, by calling Open.
	// Alternatively, you can call open individually for State.Gateways().
	// Note, however, that you should call Sate.Handler.Open(State.Events) once
	// before calling Gateway.Open, should you choose to open individually.
	//
	// During update, the state's State field will be replaced, as well as the
	// gateways and the rescale function. The event handler will remain
	// untouched which is why you don't need to readd your handlers.
	//
	// Default
	//
	// If you set neither TotalShards nor Gateways, this will default to the
	// below unless you define a custom Rescale function.
	//
	// 	func(update func(Options) *State) {
	//		s, err := update(Options{})
	//		if err != nil {
	//			log.Println("could not update state during rescale:", err.Error())
	//			return
	//		}
	//
	//		err = s.Open(context.Background())
	//		if err != nil {
	//			log.Println("could not open state during rescale:", err.Error())
	//		}
	//	}
	//
	// Otherwise, you are required to set this function yourself.
	Rescale func(update func(state.Options) (*state.State, error))

	// NoDefaultMiddlewares, if true, does not add the default middlewares
	// upon creation of the bot.
	// These middlewares are responsible for validation as well as filling some
	// of the plugin.Context's fields.
	//
	// If setting this to true those checks, or equivalents of them, should be
	// manually added.
	// Although possible, it is highly discouraged to disable certain checks
	// unless the resulting behavior is explicitly desired, as default or
	// third-party plugins may rely on these checks in order to perform as
	// intended.
	//
	// When the first middleware is called not all fields of the context will
	// be set.
	// Instead, they are set by the default middlewares.
	// If you want to swap out a default middleware for a custom
	// implementation, refer to its doc to see which fields it sets.
	//
	// Also keep in mind that middlewares added after field-setting
	// middlewares, may rely on those field being set.
	// Fully removing a default middleware that sets some fields without
	// setting  them otherwise is highly discouraged, as plugins will assume
	// all fields to be set.
	//
	// To see which fields are always set, i.e. which fields are available
	// to all middlewares, refer to the doc of Bot.Route.
	//
	// By default, the following middlewares are added upon creation of the
	// bot.
	//
	//  Bot.AddMiddleware(CheckMessageType)
	//  Bot.AddMiddleware(CheckHuman) // if Options.AllowBot is true
	//	Bot.AddMiddleware(NewSettingsRetriever(Options.SettingsProvider))
	//  Bot.AddMiddleware(CheckPrefix)
	//	Bot.AddMiddleware(FindCommand)
	//	Bot.AddMiddleware(CheckChannelTypes)
	//	Bot.AddMiddleware(CheckBotPermissions)
	//	Bot.AddMiddleware(NewThrottlerChecker(Options.ThrottlerCancelChecker))
	//
	//	Bot.AddPostMiddleware(CheckRestrictions)
	//	Bot.AddPostMiddleware(ParseArgs)
	//	Bot.AddPostMiddleware(InvokeCommand)
	NoDefaultMiddlewares bool
}

// SetDefaults fills the defaults for all options, that weren't manually set.
func (o *Options) SetDefaults() (err error) {
	if o.SettingsProvider == nil {
		o.SettingsProvider = StaticSettings()
	}

	if o.ArgParser == nil {
		o.ArgParser = &arg.DelimiterParser{Delimiter: ','}
	}

	if o.ThrottlerCancelChecker == nil {
		o.ThrottlerCancelChecker = DefaultThrottlerCancelCheck
	}

	if o.GatewayErrorHandler == nil {
		o.GatewayErrorHandler = DefaultGatewayErrorHandler
	}

	if o.StateErrorHandler == nil {
		o.StateErrorHandler = func(err error) { log.Println("event handler:", err.Error()) }
	}

	if o.StatePanicHandler == nil {
		o.StatePanicHandler = func(rec interface{}) {
			log.Printf("event handler: panic: %+v\n%s\n", rec, debug.Stack())
		}
	}

	if o.ErrorHandler == nil {
		o.ErrorHandler = DefaultErrorHandler
	}

	if o.PanicHandler == nil {
		o.PanicHandler = DefaultPanicHandler
	}

	return nil
}

// SettingsProvider is the function used to retrieve the settings for the guild
// or direct message.
//
// The passed *event.Base is the base of the event triggering the settings
// check.
// This will either stem from a message create event, or a message update
// event, if Options.EditAge is greater than 0.
//
// First Return Value
//
// The first return value contains the prefixes used by the guild or user.
// In a guild, the message must start with one of the prefixes or with a bot
// mention.
// Direct Messages are not subject to this limitation.
// However, if prefixes are returned for a direct message invoke and the user
// prefixes their message with one, the prefix will still be stripped.
// Similarly, if the user prefixes a direct message with the bot's mention, it
// will also get stripped.
// This behavior can be changed by replacing the CheckPrefix default
// middleware.
//
// All spaces and newlines between prefix and the rest of the message will be
// removed before given to the router.
// If prefixes is empty, the only valid prefix will be a mention of the bot.
//
// Prefix matching is non-exhaustive, meaning the first matching prefix will be
// used and the bot will not look for longer prefix that matches as well.
//
// Second Return Value
//
// The second return value is the *1i8n.Localizer, used to generate
// translations.
// New localizers can be created using i18n.NewLocalizer.
//
// Third Return Value
//
// The last return value is an ok-type bool.
// If false, the message will be discarded, regardless of whether there is a
// matching prefix.
// This intended for use in error scenarios, when no prefix or localizer can be
// obtained, and fallbacks are undesired.
//
// Error Handling
//
// SettingsProvider intentionally uses a bool instead of an error return value
// to signal unsuccessful execution.
// This has two reasons.
//
// Mainly, because it is not ensured if the message we are checking is even a
// command invoke.
// Suppose you are unable to retrieve your bots settings.
// This would lead to the bot responding to every message on every server it is
// with some sort of error.
//
// Secondly, errors in adam are represented through errors.Error.
// However, to invoke errors.Error.Handle a plugin.Context is required,
// which at this point is not generated yet, because the message has not been
// identified as a command.
//
// For those reasons, error handling is left implementation-specific, and you
// are responsible for ensuring that errors are properly captured.
type SettingsProvider func(b *event.Base, m *discord.Message) (prefixes []string, localizer *i18n.Localizer, ok bool)

// StaticSettings creates a new SettingsProvider that returns the same prefixes
// for all guilds and users.
// The returned localizer will always be a fallback localizer.
func StaticSettings(prefixes ...string) SettingsProvider {
	return func(*event.Base, *discord.Message) ([]string, *i18n.Localizer, bool) {
		return prefixes, i18n.NewFallbackLocalizer(), true
	}
}

// =============================================================================
// Defaults
// =====================================================================================

// DefaultThrottlerCancelCheck checks if the error is an *InformationalError.
// If so, it returns false otherwise it returns true.
func DefaultThrottlerCancelCheck(err error) bool {
	var ierr *errors.InformationalError
	return !errors.As(err, &ierr)
}

func DefaultGatewayErrorHandler(err error) {
	if !FilterGatewayError(err) {
		log.Println(err)
	}
}

// FilterGatewayError filters out informational reconnect errors.
func FilterGatewayError(err error) bool {
	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-close-event-codes
	const errShardingRequired = 4011

	var cerr *websocket.CloseError
	return errors.As(err, &cerr) && cerr.Code == errShardingRequired
}

func DefaultErrorHandler(err error, s *state.State, ctx *plugin.Context) {
	errors.Handle(s, ctx, err)
}

func DefaultPanicHandler(recovered interface{}, s *state.State, ctx *plugin.Context) {
	err, ok := recovered.(error)
	if ok {
		err = errors.Wrap(err, "panic")
	} else {
		err = errors.NewWithStackf("panic: %+v", recovered)
	}

	errors.Handle(s, ctx, err)
}
