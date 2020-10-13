package plugin

import (
	"regexp"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

// ArgTypeRaw can be used to indicate that an arg is not parsed, but used
// literally.
// See the docs of ArgsInfoer for more information.
var ArgTypeRaw = TypeInfo{
	Name:        "__raw__",
	Description: "__raw__",
}

type (
	// ArgConfig is the abstraction of the commands argument and flag
	// configuration.
	//
	// Default implementations can be found in impl/arg.
	ArgConfig interface {
		// Parse parses the passed arguments and returns the retrieved Args
		// and Flags.
		// args is the trimmed message, with prefix and command stripped.
		Parse(args string, s *state.State, ctx *Context) (Args, Flags, error)
	}

	// ArgsInfoer is an interface that can be optionally implemented by an
	// ArgConfig.
	// It provides meta information about the arguments and flags of a command.
	//
	// Raw Args
	//
	// As a special case, there is also support for raw arguments, i.e.
	// arguments that are used literally.
	//
	// An argument config is considered as a raw argument, if Info returns a
	// slice of length 1 with a single non-variadic required argument, no
	// optionals and no flags.
	// The required argument's Type must be ArgTypeRaw.
	//
	// If those requirements are fulfilled, help commands should display the
	// description of the required argument instead of creating their normal
	// argument help message.
	ArgsInfoer interface {
		// Info returns localized information about the arguments and flags of
		// a command.
		// If the returned slice is nil, it will be assumed, that no
		// information is available.
		Info(l *i18n.Localizer) []ArgsInfo
	}

	// ArgsInfo contains localized information about a command's arguments.
	ArgsInfo struct {
		// Prefix contains the prefix, if there are multiple argument
		// combinations.
		// Otherwise it should be empty.
		Prefix string
		// Required contains information about required arguments.
		Required []ArgInfo
		// Optional contains information about optional arguments.
		Optional []ArgInfo
		// Variadic specifies whether the last possibly specifiable argument
		// is variadic.
		Variadic bool

		// Flags contains information about the command's flags.
		Flags []FlagInfo
	}

	// ArgInfo contains information about an argument.
	ArgInfo struct {
		// Name is the name of the argument.
		Name string
		// Type contains information about the type of the argument.
		Type TypeInfo
		// Description is the optional description of the argument.
		Description string
	}

	// FlagInfo contains information about a flag.
	FlagInfo struct {
		// Name is the name of the flag.
		Name string
		// Aliases contains the optional aliases of the flag.
		Aliases []string
		// Type contains information about the type of the flag.
		Type TypeInfo
		// Description is the optional description of the flag.
		Description string
		// Multi specifies whether the flag may be used multiple times.
		Multi bool
	}

	// TypeInfo contains information about a flag or arg type.
	//
	// Because of the special meaning of ArgTypeRaw, TypeInfos with the name
	// and description '__raw__' should not be used.
	TypeInfo struct {
		// Name is the name of the type.
		Name string
		// Description is the optional description of the type.
		Description string
	}
)

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

// Member returns the argument with the passed index as *discord.Member.
func (a Args) Member(i int) *discord.Member { return a[i].(*discord.Member) }

// Channel returns the argument with the passed index as *discord.Channel.
func (a Args) Channel(i int) *discord.Channel { return a[i].(*discord.Channel) }

// Role returns the argument with the passed index as *discord.Role.
func (a Args) Role(i int) *discord.Role { return a[i].(*discord.Role) }

// Emoji returns the argument with the passed index as *discord.Emoji.
func (a Args) Emoji(i int) *discord.Emoji { return a[i].(*discord.Emoji) }

// Duration returns the argument with the passed index as time.Duration.
func (a Args) Duration(i int) time.Duration { return a[i].(time.Duration) }

// Time returns the argument with the passed index as time.Time.
func (a Args) Time(i int) time.Time { return a[i].(time.Time) }

// Regexp returns the flag with the passed index as *regexp.Regexp.
func (a Args) Regexp(i int) *regexp.Regexp { return a[i].(*regexp.Regexp) }

// VariadicInt returns the last argument as []int.
func (a Args) VariadicInt() []int { return a[len(a)-1].([]int) }

// VariadicInt64 returns the last argument as []int64.
func (a Args) VariadicInt64() []int64 { return a[len(a)-1].([]int64) }

// VariadicUint returns the last argument as []uint.
func (a Args) VariadicUint() []uint { return a[len(a)-1].([]uint) }

// VariadicUint64 returns the last argument as []uint64.
func (a Args) VariadicUint64() []uint64 { return a[len(a)-1].([]uint64) }

// VariadicFloat32 returns the last argument as []float32.
func (a Args) VariadicFloat32() []float32 { return a[len(a)-1].([]float32) }

// VariadicFloat64 returns the last argument as []Float64.
func (a Args) VariadicFloat64() []float64 { return a[len(a)-1].([]float64) }

// VariadicString returns the last argument as []string.
func (a Args) VariadicString() []string { return a[len(a)-1].([]string) }

// VariadicBool returns the last argument as []discord.Member.
func (a Args) VariadicMember() []discord.Member { return a[len(a)-1].([]discord.Member) }

// VariadicChannel returns the last argument as []discord.Channel.
func (a Args) VariadicChannel() []discord.Channel { return a[len(a)-1].([]discord.Channel) }

// VariadicRole returns the last argument as []discord.Role.
func (a Args) VariadicRole() []discord.Role { return a[len(a)-1].([]discord.Role) }

// VariadicEmoji returns the last argument as []discord.Emoji.
func (a Args) VariadicEmoji() []discord.Emoji { return a[len(a)-1].([]discord.Emoji) }

// VariadicDuration returns the last argument as []time.Duration.
func (a Args) VariadicDuration() []time.Duration { return a[len(a)-1].([]time.Duration) }

// VariadicTime returns the last argument as []time.Time.
func (a Args) VariadicTime() []time.Time { return a[len(a)-1].([]time.Time) }

// VariadicRegexp returns the last argument as []*regexp.Regexp.
func (a Args) VariadicRegexp() []*regexp.Regexp { return a[len(a)-1].([]*regexp.Regexp) }

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

// Member returns the flag with the passed name as *discord.Member.
func (f Flags) Member(name string) *discord.Member { return f[name].(*discord.Member) }

// Channel returns the flag with the passed name as *discord.Channel.
func (f Flags) Channel(name string) *discord.Channel { return f[name].(*discord.Channel) }

// Role returns the flag with the passed name as *discord.Role.
func (f Flags) Role(name string) *discord.Role { return f[name].(*discord.Role) }

// Emoji returns the flag with the passed name as *discord.Emoji.
func (f Flags) Emoji(name string) *discord.Emoji { return f[name].(*discord.Emoji) }

// Duration returns the flag with the passed name as time.Duration.
func (f Flags) Duration(name string) time.Duration { return f[name].(time.Duration) }

// Time returns the flag with the passed name as time.Time.
func (f Flags) Time(name string) time.Time { return f[name].(time.Time) }

// Regexp returns the flag with the passed name as *regexp.Regexp.
func (f Flags) Regexp(name string) *regexp.Regexp { return f[name].(*regexp.Regexp) }

// MultiInt returns the flag with the passed name as []int.
func (f Flags) MultiInt(name string) []int { return f[name].([]int) }

// MultiInt64 returns the flag with the passed name as []int64.
func (f Flags) MultiInt64(name string) []int64 { return f[name].([]int64) }

// MultiUint returns the flag with the passed name as []uint.
func (f Flags) MultiUint(name string) []uint { return f[name].([]uint) }

// MultiUint64 returns the flag with the passed name as []uint64.
func (f Flags) MultiUint64(name string) []uint64 { return f[name].([]uint64) }

// MultiFloat32 returns the flag with the passed name as []float32.
func (f Flags) MultiFloat32(name string) []float32 { return f[name].([]float32) }

// MultiFloat64 returns the flag with the passed name as []Float64.
func (f Flags) MultiFloat64(name string) []float64 { return f[name].([]float64) }

// MultiString returns the flag with the passed name as []string.
func (f Flags) MultiString(name string) []string { return f[name].([]string) }

// MultiBool returns the flag with the passed name as []discord.Member.
func (f Flags) MultiMember(name string) []discord.Member { return f[name].([]discord.Member) }

// MultiChannel returns the flag with the passed name as []discord.Channel.
func (f Flags) MultiChannel(name string) []discord.Channel { return f[name].([]discord.Channel) }

// MultiRole returns the flag with the passed name as []discord.Role.
func (f Flags) MultiRole(name string) []discord.Role { return f[name].([]discord.Role) }

// MultiEmoji returns the flag with the passed name as []discord.Emoji.
func (f Flags) MultiEmoji(name string) []discord.Emoji { return f[name].([]discord.Emoji) }

// MultiDuration returns the flag with the passed name as []time.Duration.
func (f Flags) MultiDuration(name string) []time.Duration { return f[name].([]time.Duration) }

// MultiTime returns the flag with the passed name as []time.Time.
func (f Flags) MultiTime(name string) []time.Time { return f[name].([]time.Time) }

// MultiRegexp returns the flag with the passed name as []*regexp.Regexp.
func (f Flags) MultiRegexp(name string) []*regexp.Regexp { return f[name].([]*regexp.Regexp) }
