package bot

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state/store"
	"github.com/diamondburned/arikawa/v2/state/store/defaultstore"
	"github.com/diamondburned/arikawa/v2/utils/wsutil"
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
	// LocalizationFunc is the function used to retrieve i18n.LangFuncs used
	// for creating localized data.
	// Leave this empty, if you don't want to use localization.
	//
	// Default: nil
	LocalizationFunc i18n.Func
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
	// SendTyping specifies whether to send a typing event if the command has
	// the SendMessagesPermission.
	// The event will be sent in 6 second intervals, until the Invoke method
	// of the command returns.
	//
	// Default: false
	SendTyping bool

	// NoAutoOpen defines whether to call the Open and Close methods of plugins
	// automatically when bot.Open() and bot.Close is called.
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
	// Default: DefaultThrottlerErrorCheck
	ThrottlerCancelChecker func(error) bool

	// ReplyMiddlewares contains the middlewares that should be used when
	// awaiting a reply.
	//
	// The following types are permitted:
	//  • func(*state.State, interface{})
	//  • func(*state.State, interface{}) error
	//  • func(*state.State, *state.Base)
	//  • func(*state.State, *state.Base) error
	//  • func(*state.State, *state.MessageCreateEvent)
	//  • func(*state.State, *state.MessageCreateEvent) error
	ReplyMiddlewares []interface{}

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
	// Default: func(err error) { log.Println(err) }
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

	for i, m := range o.ReplyMiddlewares {
		switch m.(type) {
		case func(*state.State, interface{}):
		case func(*state.State, interface{}) error:
		case func(*state.State, *state.Base):
		case func(*state.State, *state.Base) error:
		case func(*state.State, *state.MessageCreateEvent):
		case func(*state.State, *state.MessageCreateEvent) error:
		default:
			return errors.NewWithStackf("bot: Options.ReplyMiddlewares[%d] is of unsupported type %T", i, m)
		}
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
		o.GatewayErrorHandler = func(err error) { log.Println(err) }
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
type SettingsProvider func(b *state.Base, m *discord.Message) (prefixes []string, lang string)

// NewStaticSettingsProvider creates a new SettingsProvider that returns the
// same prefixes for all guilds and users.
// The returned language will always be an empty string.
func NewStaticSettingsProvider(prefixes ...string) SettingsProvider {
	return func(*state.Base, *discord.Message) ([]string, string) {
		return prefixes, ""
	}
}

// =============================================================================
// Defaults
// =====================================================================================

func DefaultThrottlerErrorCheck(err error) bool {
	return !errors.As(err, new(errors.InformationalError))
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
