package plugin

import (
	"regexp"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

type (
	// ArgConfig is the abstraction of the commands argument and flag
	// configuration.
	//
	// Default implementations can be found in impl/arg.
	ArgConfig interface {
		// GetRequiredArgs returns the required arguments of the command.
		GetRequiredArgs() []RequiredArg
		// GetOptionalArgs returns the optional arguments of the command.
		GetOptionalArgs() []OptionalArg
		// IsVariadic returns whether the last argument is variadic, i.e. that
		// it may be specified multiple times.
		IsVariadic() bool

		// GetFlags returns the flags of the command.
		GetFlags() []Flag
	}

	// ArgParser is the abstraction of a parser that uses the information
	// provided by the ArgConfig to parse the arguments and flags supplied to
	// a command.
	//
	// Every bot instance defines a global ArgParser that can be overridden by
	// the individual commands.
	ArgParser interface {
		// Parse parses the passed arguments and stores them in the passed
		// *plugin.Context.
		// args is the trimmed message, with prefix and command stripped.
		Parse(args string, argConfig ArgConfig, s *state.State, ctx *Context) error

		// FormatArgs formats the passed arguments and flags so that they would
		// present a valid input for the passed ArgConfig by properly escaping
		// and delimiting the individual arguments and flags.
		//
		// It is guaranteed that the passed args and flags are valid in
		// themselves, i.e. that they fulfil the requirements defined by their
		// config.
		//
		// See package arg for example implementations.
		FormatArgs(argConfig ArgConfig, args []string, flags map[string]string) string
		// FormatUsage formats the passed arguments and flags so that they are
		// properly delimited.
		// It should ignore the need for escapes, as the produced output is
		// solely intended to be used for usage illustrations such as
		//	<Required Argument 1>, <Required Argument 2>, [Optional Argument 1]
		// The above output would be produced if the args slice contained
		// {"<Required Argument 1>", "<Required Argument 2>",
		// "[Optional Argument 1"} and the ArgParser uses a "," as delimiter.
		FormatUsage(argConfig ArgConfig, args []string) string
		// FormatFlag formats the passed name of a flag as it would be required
		// if using that flag.
		// For example "my-flag" could become "-my-flag" if using a
		// shellword-like flag notation.
		FormatFlag(name string) string
	}

	// RequiredArg is the interface used to access information about a single
	// required argument.
	RequiredArg interface {
		// GetName returns the name of the argument.
		GetName(*i18n.Localizer) string
		// GetType returns the ArgType of the argument.
		GetType() ArgType
		// GetDescription returns the optional description of the argument.
		GetDescription(*i18n.Localizer) string
	}

	// OptionalArg is the interface used to access information about a single
	// optional argument.
	OptionalArg interface {
		// GetName returns the name of the argument.
		GetName(*i18n.Localizer) string
		// GetType returns the ArgType of the argument.
		GetType() ArgType
		// GetDefault is the default value of the argument.
		//
		// If Default is (interface{})(nil), ArgType.Default() will be used.
		GetDefault() interface{}
		// GetDescription returns the optional description of the argument.
		GetDescription(*i18n.Localizer) string
	}

	// Flag contains information about a flag.
	Flag interface {
		// GetName returns the name of the flag.
		GetName() string
		// GetAliases returns the optional aliases of the flag.
		GetAliases() []string
		// GetType returns information about the type of the flag.
		GetType() ArgType
		// GetDefault is the default value of the flag.
		//
		// If Default is (interface{})(nil), ArgType.Default() will be used.
		GetDefault() interface{}
		// GetDescription returns the optional description of the flag.
		GetDescription(*i18n.Localizer) string
		// IsMulti returns whether the flag may be used multiple times.
		IsMulti() bool
	}

	// ArgType contains information about a flag or arg type.
	// The returned name and description must be the same for all arguments of
	// the guild.
	ArgType interface {
		// GetName returns the name of the type.
		// The name should be a noun.
		GetName(*i18n.Localizer) string
		// GetDescription returns the description of the type.
		GetDescription(*i18n.Localizer) string
		// Parse parses the argument or flag using the passed Context.
		//
		// The first return value must always be of the same type.
		Parse(s *state.State, ctx *ParseContext) (interface{}, error)
		// GetDefault returns the default value for the type.
		// See Flag.Default or OptionalArg.Default for more info.
		//
		// It must return a value that is of the type returned by Parse.
		GetDefault() interface{}
	}
)

// ArgKind specifies whether a flag or an argument is being parsed.
type ArgKind string

const (
	// KindArg is the Kind used for argument.
	KindArg = "arg"
	// KindFlag is the Kind used for flags.
	KindFlag = "flag"
)

// ParseContext is the context passed to ArgType.Parse.
type ParseContext struct {
	*Context

	// Raw is the raw argument or flag.
	Raw string
	// Name is the name of the argument or flag.
	// It includes possible prefixes such as minuses.
	Name string
	// UsedName is the alias of the flag the Context represents.
	// If the name of the flag was used, or the context represents an
	// argument, UsedName will be equal to Name.
	// It includes possible prefixes such as minuses.
	UsedName string
	// Index contains the index of the argument, if the context represents
	// an argument.
	Index int
	// Kind specifies whether a flag or an argument is being parsed.
	Kind ArgKind
}

// Args are the parsed arguments of a command.
type Args []interface{}

// Bool returns the argument with the passed index as bool.
func (a Args) Bool(i int) bool { return a[i].(bool) }

// Int returns the argument with the passed index as int.
func (a Args) Int(i int) int { return a[i].(int) }

// Int64 returns the argument with the passed index as int64.
func (a Args) Int64(i int) int64 { return a[i].(int64) }

// Uint returns the argument with the passed index as uint.
func (a Args) Uint(i int) uint { return a[i].(uint) }

// Uint64 returns the argument with the passed index as uint64.
func (a Args) Uint64(i int) uint64 { return a[i].(uint64) }

// Float32 returns the argument with the passed index as float32.
func (a Args) Float32(i int) float32 { return a[i].(float32) }

// Float64 returns the argument with the passed index as float64.
func (a Args) Float64(i int) float64 { return a[i].(float64) }

// String returns the argument with the passed index as string.
func (a Args) String(i int) string { return a[i].(string) }

// User returns the argument with the passed index as *discord.User.
func (a Args) User(i int) *discord.User { return a[i].(*discord.User) }

// Member returns the argument with the passed index as *discord.Member.
func (a Args) Member(i int) *discord.Member { return a[i].(*discord.Member) }

// Channel returns the argument with the passed index as *discord.Channel.
func (a Args) Channel(i int) *discord.Channel { return a[i].(*discord.Channel) }

// Role returns the argument with the passed index as *discord.Role.
func (a Args) Role(i int) *discord.Role { return a[i].(*discord.Role) }

// APIEmoji returns the argument with the passed index as discord.APIEmoji.
func (a Args) APIEmoji(i int) discord.APIEmoji { return a[i].(discord.APIEmoji) }

// Emoji returns the argument with the passed index as *discord.Emoji.
func (a Args) Emoji(i int) *discord.Emoji { return a[i].(*discord.Emoji) }

// Duration returns the argument with the passed index as time.Duration.
func (a Args) Duration(i int) time.Duration { return a[i].(time.Duration) }

// Time returns the argument with the passed index as time.Time.
func (a Args) Time(i int) time.Time { return a[i].(time.Time) }

// Location returns the argument with the passed index as *time.Location.
func (a Args) Location(i int) *time.Location { return a[i].(*time.Location) }

// Regexp returns the flag with the passed index as *regexp.Regexp.
func (a Args) Regexp(i int) *regexp.Regexp { return a[i].(*regexp.Regexp) }

// Command returns the argument with the passed index as *ResolvedCommand.
func (a Args) Command(i int) *ResolvedCommand { return a[i].(*ResolvedCommand) }

// Module returns the argument with the passed index as *RegisteredRegexp.
func (a Args) Module(i int) *ResolvedModule { return a[i].(*ResolvedModule) }

// Ints returns the argument with the passed index as []int.
func (a Args) Ints(i int) []int { return a[i].([]int) }

// Int64s returns the argument with the passed index as []int64.
func (a Args) Int64s(i int) []int64 { return a[i].([]int64) }

// Uints returns the argument with the passed index as []uint.
func (a Args) Uints(i int) []uint { return a[i].([]uint) }

// Uint64s returns the argument with the passed index as []uint64.
func (a Args) Uint64s(i int) []uint64 { return a[i].([]uint64) }

// Float32s returns the argument with the passed index as []float32.
func (a Args) Float32s(i int) []float32 { return a[i].([]float32) }

// Float64s returns the argument with the passed index as []Float64.
func (a Args) Float64s(i int) []float64 { return a[i].([]float64) }

// Strings returns the argument with the passed index as []string.
func (a Args) Strings(i int) []string { return a[i].([]string) }

// Members returns the argument with the passed index as []*discord.Member.
func (a Args) Members(i int) []*discord.Member { return a[i].([]*discord.Member) }

// Users returns the argument with the passed index as []*discord.Member.
func (a Args) Users(i int) []*discord.User { return a[i].([]*discord.User) }

// Channels returns the argument with the passed index as []discord.Channel.
func (a Args) Channels(i int) []*discord.Channel { return a[i].([]*discord.Channel) }

// Roles returns the argument with the passed index as []*discord.Role.
func (a Args) Roles(i int) []*discord.Role { return a[i].([]*discord.Role) }

// APIEmojis returns the argument with the passed index as []discord.APIEmoji.
func (a Args) APIEmojis(i int) []discord.APIEmoji { return a[i].([]discord.APIEmoji) }

// Emojis returns the argument with the passed index as []*discord.Emoji.
func (a Args) Emojis(i int) []*discord.Emoji { return a[i].([]*discord.Emoji) }

// Durations returns the argument with the passed index as []time.Duration.
func (a Args) Durations(i int) []time.Duration { return a[i].([]time.Duration) }

// Times returns the argument with the passed index as []time.Time.
func (a Args) Times(i int) []time.Time { return a[i].([]time.Time) }

// Locations returns the argument with the passed index as []*time.Location.
func (a Args) Locations(i int) []*time.Location { return a[i].([]*time.Location) }

// Regexps returns the argument with the passed index as []*regexp.Regexp.
func (a Args) Regexps(i int) []*regexp.Regexp { return a[i].([]*regexp.Regexp) }

// Commands returns the argument with the passed index as []*ResolvedCommand.
func (a Args) Commands(i int) []*ResolvedCommand { return a[i].([]*ResolvedCommand) }

// Modules returns the argument with the passed index as []*ResolvedModule.
func (a Args) Modules(i int) []*ResolvedModule { return a[i].([]*ResolvedModule) }

// Flags are the parsed flags of a command.
type Flags map[string]interface{}

// Bool returns the flag with the passed name as bool.
func (f Flags) Bool(name string) bool { return f[name].(bool) }

// Int returns the flag with the passed name as int.
func (f Flags) Int(name string) int { return f[name].(int) }

// Int64 returns the flag with the passed name as int64.
func (f Flags) Int64(name string) int64 { return f[name].(int64) }

// Uint returns the flag with the passed name as uint.
func (f Flags) Uint(name string) uint { return f[name].(uint) }

// Uint64 returns the flag with the passed name as uint64.
func (f Flags) Uint64(name string) uint64 { return f[name].(uint64) }

// Float32 returns the flag with the passed name as float31.
func (f Flags) Float32(name string) float32 { return f[name].(float32) }

// Float64 returns the flag with the passed name as float64.
func (f Flags) Float64(name string) float64 { return f[name].(float64) }

// String returns the flag with the passed name as string.
func (f Flags) String(name string) string { return f[name].(string) }

// User returns the flag with the passed name as *discord.User.
func (f Flags) User(name string) *discord.User { return f[name].(*discord.User) }

// Member returns the flag with the passed name as *discord.Member.
func (f Flags) Member(name string) *discord.Member { return f[name].(*discord.Member) }

// Channel returns the flag with the passed name as *discord.Channel.
func (f Flags) Channel(name string) *discord.Channel { return f[name].(*discord.Channel) }

// Role returns the flag with the passed name as *discord.Role.
func (f Flags) Role(name string) *discord.Role { return f[name].(*discord.Role) }

// APIEmoji returns the flag with the passed name as discord.APIEmoji.
func (f Flags) APIEmoji(name string) discord.APIEmoji { return f[name].(discord.APIEmoji) }

// Emoji returns the flag with the passed name as *discord.Emoji.
func (f Flags) Emoji(name string) *discord.Emoji { return f[name].(*discord.Emoji) }

// Duration returns the flag with the passed name as time.Duration.
func (f Flags) Duration(name string) time.Duration { return f[name].(time.Duration) }

// Time returns the flag with the passed name as time.Time.
func (f Flags) Time(name string) time.Time { return f[name].(time.Time) }

// Location returns the flag with the passed name as *time.Location.
func (f Flags) Location(name string) *time.Location { return f[name].(*time.Location) }

// Regexp returns the flag with the passed name as *regexp.Regexp.
func (f Flags) Regexp(name string) *regexp.Regexp { return f[name].(*regexp.Regexp) }

// Command returns the flag with the passed name as *ResolvedCommand.
func (f Flags) Command(name string) *ResolvedCommand { return f[name].(*ResolvedCommand) }

// Module returns the flag with the passed name as *ResolvedModule.
func (f Flags) Module(name string) *ResolvedModule { return f[name].(*ResolvedModule) }

// Ints returns the flag with the passed name as []int.
func (f Flags) Ints(name string) []int { return f[name].([]int) }

// Int64s returns the flag with the passed name as []int64.
func (f Flags) Int64s(name string) []int64 { return f[name].([]int64) }

// Uints returns the flag with the passed name as []uint.
func (f Flags) Uints(name string) []uint { return f[name].([]uint) }

// Uint64s returns the flag with the passed name as []uint64.
func (f Flags) Uint64s(name string) []uint64 { return f[name].([]uint64) }

// Float32s returns the flag with the passed name as []float32.
func (f Flags) Float32s(name string) []float32 { return f[name].([]float32) }

// Float64s returns the flag with the passed name as []Float64.
func (f Flags) Float64s(name string) []float64 { return f[name].([]float64) }

// Strings returns the flag with the passed name as []string.
func (f Flags) Strings(name string) []string { return f[name].([]string) }

// Users returns the flag with the passed name as []*discord.User.
func (f Flags) Users(name string) []*discord.User { return f[name].([]*discord.User) }

// Members returns the flag with the passed name as []*discord.Member.
func (f Flags) Members(name string) []*discord.Member { return f[name].([]*discord.Member) }

// Channels returns the flag with the passed name as []*discord.Channel.
func (f Flags) Channels(name string) []*discord.Channel { return f[name].([]*discord.Channel) }

// Roles returns the flag with the passed name as []*discord.Role.
func (f Flags) Roles(name string) []*discord.Role { return f[name].([]*discord.Role) }

// APIEmojis returns the flag with the passed name as []discord.APIEmoji.
func (f Flags) APIEmojis(name string) []discord.APIEmoji { return f[name].([]discord.APIEmoji) }

// Emojis returns the flag with the passed name as []*discord.Emoji.
func (f Flags) Emojis(name string) []*discord.Emoji { return f[name].([]*discord.Emoji) }

// Durations returns the flag with the passed name as []time.Duration.
func (f Flags) Durations(name string) []time.Duration { return f[name].([]time.Duration) }

// Times returns the flag with the passed name as []time.Time.
func (f Flags) Times(name string) []time.Time { return f[name].([]time.Time) }

// Locations returns the flag with the passed name as []*time.Location.
func (f Flags) Locations(name string) []*time.Location { return f[name].([]*time.Location) }

// Regexps returns the flag with the passed name as []*regexp.Regexp.
func (f Flags) Regexps(name string) []*regexp.Regexp { return f[name].([]*regexp.Regexp) }

// Commands returns the flag with the passed name as []*ResolvedCommand.
func (f Flags) Commands(name string) []*ResolvedCommand { return f[name].([]*ResolvedCommand) }

// Modules returns the flag with the passed name as []*ResolvedModule.
func (f Flags) Modules(name string) []*ResolvedModule { return f[name].([]*ResolvedModule) }
