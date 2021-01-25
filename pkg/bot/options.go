package bot

import (
	"fmt"
	"log"
	"runtime/debug"
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

	// DefaultChannelTypes are the default plugin.ChannelTypes, used if neither
	// the parent modules of a command nor the command itself define channel
	// types.
	//
	// Default: plugin.AllChannels
	DefaultChannelTypes plugin.ChannelTypes
	// DefaultRestrictions are the default restrictions, used if neither the
	// parent modules of a command nor the command itself define any.
	//
	// Default: nil
	DefaultRestrictions plugin.RestrictionFunc
	// DefaultThrottler is the default throttler, used if neither the parent
	// modules of a command nor the command itself define a throttler.
	// Note that the throttler will not work on a per-command basis, but
	// globally for all commands that use it.
	//
	// Default: nil
	DefaultThrottler plugin.Throttler

	// ThrottlerCancelChecker is the function run every time a command returns
	// with a non-nil error.
	// If the function returns true, the command's throttler will not count the
	// invoke.
	//
	// Settings this field has no effect, if ManualChecks is set to true.
	//
	// Default: DefaultThrottlerErrorCheck
	ThrottlerCancelChecker func(error) bool

	// AsyncPluginProviders specifies whether the plugins providers should be
	// fetched asynchronously or one by one.
	//
	// Default: false
	AsyncPluginProviders bool

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
	// Default: func(err error) { log.Println(err) }
	StateErrorHandler func(error)
	// StatePanicHandler is the panic handler of the *state.State, called if an
	// event handler panics.
	//
	// Default:
	// 	func(recovered interface{}) {
	// 		log.Printf("recovered from panic: %+v\n%s", recovered, debug.Stack())
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

	// ManualChecks specifies whether the checks are performed by the user and
	// should not be added automatically.
	//
	// By default, the following checks are performed in the following order.
	// Steps enclosed in parentheses are just steps executed in-between checks.
	//
	//	(command invoke received and routed) ->
	//	CheckChannelTypes ->
	//	CheckBotPermissions ->
	// 	NewThrottlerChecker(b.ThrottlerCancelChecker) ->
	//	(run custom middlewares) ->
	//	CheckRestrictions ->
	//	ParseArgs ->
	//	(invoke command)
	//
	// If you set this to true, you are responsible for ensuring that all
	// (desired) checks are performed, by adding them as middlewares.
	// All default checks can be found in package bot.
	//
	// Note that leaving out these checks may lead to unexpected or undesired
	// behavior, as default and third-party plugins will assume that these
	// checks are run.
	// Therefore this is only recommended, to change order or use custom
	// checks.
	ManualChecks bool
}

// SetDefaults fills the defaults for all options, that weren't manually set.
func (o *Options) SetDefaults() (err error) {
	if o.SettingsProvider == nil {
		o.SettingsProvider = NewStaticSettingsProvider(",")
	}

	if len(o.Status) == 0 {
		o.Status = gateway.OnlineStatus
	}

	if o.DefaultChannelTypes == 0 {
		o.DefaultChannelTypes = plugin.AllChannels
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
		o.StatePanicHandler = func(recovered interface{}) {
			log.Printf("recovered from panic: %+v\n%s", recovered, debug.Stack())
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

// SettingsProvider is the function used to retrieve the settings in the guild.
//
// The passed *state.Base is the base of the event triggering settings check.
// This will either stem from either message create event, or a message update
// event, if Options.EditAge is greater than 0.
//
// If the returned *18n.Localizer is nil, a fallback localizer will be used,
// that always uses fallback messages.
//
// Note that messages sent in a direct message don't require a prefix.
// However, all prefixes for a direct messages will still be stripped, if the
// message starts with one.
type SettingsProvider func(b *state.Base, m *discord.Message) (prefixes []string, localizer *i18n.Localizer)

// NewStaticSettingsProvider creates a new SettingsProvider that returns the
// same prefixes for all guilds and users.
// The returned localizer will always be a fallback localizer.
func NewStaticSettingsProvider(prefixes ...string) SettingsProvider {
	return func(*state.Base, *discord.Message) ([]string, *i18n.Localizer) {
		return prefixes, nil
	}
}

// =============================================================================
// Defaults
// =====================================================================================

func DefaultThrottlerErrorCheck(err error) bool {
	ierr := new(errors.InformationalError)
	return !errors.As(err, ierr)
}

func DefaultGatewayErrorHandler(err error) {
	// ignore error used on reconnect

	var cerr *websocket.CloseError
	if errors.As(err, &cerr) && websocket.IsCloseError(cerr, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		return
	} else if errors.Is(err, syscall.ECONNRESET) {
		return
	}

	log.Println(err)
}

func DefaultErrorHandler(err error, s *state.State, ctx *plugin.Context) {
	for i := 0; i < 4 && err != nil; i++ { // prevent error cycle
		var Err errors.Error
		if !errors.As(err, &Err) {
			Err = errors.WithStack(err).(errors.Error) //nolint:errorlint
		}

		err = Err.Handle(s, ctx)
	}
}

func DefaultPanicHandler(recovered interface{}, s *state.State, ctx *plugin.Context) {
	err, ok := recovered.(error)
	if ok {
		err = errors.Wrap(err, "panic")
	} else {
		err = fmt.Errorf("panic: %+v", recovered)
	}

	for i := 0; i < 4 && err != nil; i++ { // prevent error cycle
		var Err errors.Error
		if !errors.As(err, &Err) {
			Err = errors.WithStack(err).(errors.Error) //nolint:errorlint
		}

		err = Err.Handle(s, ctx)
		if i == 0 {
			fmt.Println(string(debug.Stack()))
		}
	}
}
