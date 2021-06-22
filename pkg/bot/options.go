package bot

import (
	"log"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state/store"
	"github.com/diamondburned/arikawa/v2/state/store/defaultstore"
	"github.com/diamondburned/arikawa/v2/utils/wsutil"
	"github.com/gorilla/websocket"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// Options contains different configurations for a Bot.
type Options struct { //nolint:maligned // only one-time use anyway, ordered by importance, we can take the (temporary) few bytes
	// Token is the bot token without the 'Bot' prefix.
	//
	// This field is required.
	Token string

	// SettingsProvider is the settings provider for the bot.
	// If left nil, only the mention prefix will be usable.
	//
	// Default: NewStaticSettingsProvider()
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
	Status gateway.Status
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

	// AllowBot specifies whether bots may trigger commands.
	//
	// Default: false
	AllowBot bool

	// NoAutoOpen defines whether to call the Open and Close methods of plugins
	// automatically when bot.Open() and bot.Close() is called.
	// Both Open and Close may take in an optional *bot.Bot parameter, and may
	// return an error.
	//
	// The call to Open will be made before the gateway is opened.
	// It is therefore safe to add ReadyEvent handlers.
	//
	// The call to Close will be made after the event listener is Closed.
	//
	// Default: false
	NoAutoOpen bool
	// AutoAddHandlers specifies whether all methods of plugins that resemble
	// a handler func should be added automatically.
	// All methods that don't represent a handler will be discarded.
	//
	// Default: false
	AutoAddHandlers bool

	// ThrottlerCancelChecker is the function run every time a command returns
	// with a non-nil error.
	// If the function returns true, the command's throttler will not count the
	// invoke.
	//
	// Settings this field has no effect, if NoDefaultMiddlewares is set to true.
	//
	// Default: DefaultThrottlerErrorCheck
	ThrottlerCancelChecker func(error) bool

	// Cabinet is the store.Cabinet used for caching.
	// Use store.NoopCabinet to deactivate caching.
	//
	// Default: defaultstore.New()
	Cabinet store.Cabinet

	// Shard is the shard of the bot.
	//
	// Default: gateway.Shard[0, 1]
	Shard gateway.Shard
	// GatewayURL is the url of the gateway to use.
	//
	// Default: gateway.URL()
	GatewayURL string
	// GatewayTimeout is the timeout for connecting and writing to the gateway
	// before failing and exiting.
	//
	// Default: wsutil.WSTimeout
	GatewayTimeout time.Duration
	// GatewayErrorHandler is the error handler of the gateway.
	//
	// Default: DefaultGatewayErrorHandler
	GatewayErrorHandler func(error)

	// StateErrorHandler is the error handler of the *state.State, called if an
	// event handler returns with an error.
	//
	// Default: func(err error) { log.Println(err.String()) }
	StateErrorHandler func(error)
	// StatePanicHandler is the panic handler of the *state.State, called if an
	// event handler panics.
	//
	// Default:
	// 	func(rec interface{}) {
	// 		log.Printf("recovered from panic: %+v\n%s\n", rec)
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

	// NoDefaultMiddlewares, if true, prevents the default middlewares from
	// being added on creation.
	// These middlewares are responsible for validating the user is allowed
	// to run a command and the bot is able to.
	//
	// If you set this to true, you are responsible for adding those checks,
	// or equivalents of them.
	// Although possible, it is highly discouraged to disable certain checks,
	// unless the resulting behavior is explicitly desired.
	// Default or third-party plugins may rely on these checks to perform as
	// intended.
	//
	// By default, the following middlewares are added upon creation of the
	// bot.
	//
	//	Bot.AddMiddleware(NewPrefixChecker(selfID))
	//	Bot.AddMiddleware(RouteCommand)
	//	Bot.AddMiddleware(CheckChannelTypes)
	//	Bot.AddMiddleware(CheckBotPermissions)
	//	Bot.AddMiddleware(NewThrottlerChecker(Bot.ThrottlerCancelChecker))
	//
	//	Bot.AddPostMiddleware(CheckRestrictions)
	//	Bot.AddPostMiddleware(ParseArgs)
	//	Bot.AddPostMiddleware(InvokeCommand)
	NoDefaultMiddlewares bool
}

// SetDefaults fills the defaults for all options, that weren't manually set.
func (o *Options) SetDefaults() (err error) {
	if o.SettingsProvider == nil {
		o.SettingsProvider = NewStaticSettingsProvider()
	}

	if len(o.Status) == 0 {
		o.Status = gateway.OnlineStatus
	}

	if o.ThrottlerCancelChecker == nil {
		o.ThrottlerCancelChecker = DefaultThrottlerErrorCheck
	}

	o.setCabinetDefaults()

	if o.Shard[1] == 0 {
		o.Shard = gateway.Shard{0, 1}
	}

	if len(o.GatewayURL) == 0 {
		o.GatewayURL, err = gateway.URL()
		if err != nil {
			return err
		}
	}

	if o.GatewayTimeout <= 0 {
		o.GatewayTimeout = wsutil.WSTimeout
	}

	if o.GatewayErrorHandler == nil {
		o.GatewayErrorHandler = DefaultGatewayErrorHandler
	}

	if o.StateErrorHandler == nil {
		o.StateErrorHandler = func(err error) { log.Println(err) }
	}

	if o.StatePanicHandler == nil {
		o.StatePanicHandler = func(rec interface{}) {
			log.Printf("recovered from panic: %+v\n", rec)
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

func (o *Options) setCabinetDefaults() {
	if o.Cabinet.MeStore == nil {
		o.Cabinet.MeStore = defaultstore.NewMe()
	}

	if o.Cabinet.ChannelStore == nil {
		o.Cabinet.ChannelStore = defaultstore.NewChannel()
	}

	if o.Cabinet.EmojiStore == nil {
		o.Cabinet.EmojiStore = defaultstore.NewEmoji()
	}

	if o.Cabinet.GuildStore == nil {
		o.Cabinet.GuildStore = defaultstore.NewGuild()
	}

	if o.Cabinet.MemberStore == nil {
		o.Cabinet.MemberStore = defaultstore.NewMember()
	}

	if o.Cabinet.MessageStore == nil {
		o.Cabinet.MessageStore = defaultstore.NewMessage(100)
	}

	if o.Cabinet.PresenceStore == nil {
		o.Cabinet.PresenceStore = defaultstore.NewPresence()
	}

	if o.Cabinet.RoleStore == nil {
		o.Cabinet.RoleStore = defaultstore.NewRole()
	}

	if o.Cabinet.VoiceStateStore == nil {
		o.Cabinet.VoiceStateStore = defaultstore.NewVoiceState()
	}
}

// SettingsProvider is the function used to retrieve the settings for the guild
// or direct message.
//
// The passed *state.Base is the base of the event triggering settings check.
// This will either stem from either message create event, or a message update
// event, if Options.EditAge is greater than 0.
//
// First Return Value
//
// The first return value contains the prefixes used by the guild or user.
// In a guild, the message must start with one of the prefixes or with a bot
// mention.
// Direct Messages are not subject to this limitation.
// However, if prefixes are returned for a direct message invoke (msg.GuildID
// == 0), or the message starts with a mention, the prefix will still be
// stripped before being routed.
//
// All spaces, tabs, and newlines, between a prefix and the rest of the message
// will be removed before given to the router.
// If prefixes is empty, the only valid prefix will be a mention of the bot.
//
// Prefix matching is lazy, meaning the first matching prefix will be used and
// the bot will not look for longer prefix matching as well (first come, first
// served principle).
// For example if the prefixes are "a", and "ab" message will never match the
// second prefix, because all message matching "ab" also match "a", and "a"
// appears before "ab" in the prefix list.
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
// are responsible for ensuring that error are properly captured.
type SettingsProvider func(b *state.Base, m *discord.Message) (prefixes []string, localizer *i18n.Localizer, ok bool)

// NewStaticSettingsProvider creates a new SettingsProvider that returns the
// same prefixes for all guilds and users.
// The returned localizer will always be a fallback localizer.
func NewStaticSettingsProvider(prefixes ...string) SettingsProvider {
	return func(*state.Base, *discord.Message) ([]string, *i18n.Localizer, bool) {
		return prefixes, nil, true
	}
}

// =============================================================================
// Defaults
// =====================================================================================

func DefaultThrottlerErrorCheck(err error) bool {
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
	var cerr *websocket.CloseError
	return (errors.As(err, &cerr) &&
		(cerr.Code == websocket.CloseGoingAway || cerr.Code == websocket.CloseAbnormalClosure)) ||
		errors.Is(err, syscall.ECONNRESET)
}

func DefaultErrorHandler(err error, s *state.State, ctx *plugin.Context) {
	errors.Handle(err, s, ctx, 4)
}

func DefaultPanicHandler(recovered interface{}, s *state.State, ctx *plugin.Context) {
	err, ok := recovered.(error)
	if ok {
		err = errors.Wrap(err, "panic")
	} else {
		err = errors.NewWithStackf("panic: %+v", recovered)
	}

	errors.Handle(err, s, ctx, 4)
}
