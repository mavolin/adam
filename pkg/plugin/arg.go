package plugin

import (
	"regexp"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
)

type (
	// ArgConfig is the abstraction of the commands argument and flag
	// configuration.
	// It provides mean to parse arguments, and generate help messages.
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
	ArgsInfoer interface {
		// Info returns localized information about the arguments and flags of
		// a command.
		Info(l *localization.Localizer) ([]ArgsInfo, error)
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

	// TypeInfo contains information about a type.
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

// LookupBool returns the argument with the passed index as bool.
func (a Args) LookupBool(i int) (v bool, ok bool) {
	v, ok = a[i].(bool)
	return
}

// LookupInt returns the argument with the passed index as int.
// ok is false if the value at the passed index is not of type int.
func (a Args) LookupInt(i int) (v int, ok bool) {
	v, ok = a[i].(int)
	return
}

// LookupInt64 returns the argument with the passed index as int64.
// ok is false if the value at the passed index is not of type int64.
func (a Args) LookupInt64(i int) (v int64, ok bool) {
	v, ok = a[i].(int64)
	return
}

// LookupUint returns the argument with the passed index as uint.
// ok is false if the value at the passed index is not of type uint.
func (a Args) LookupUint(i int) (v uint, ok bool) {
	v, ok = a[i].(uint)
	return
}

// LookupUint64 returns the argument with the passed index as uint64.
// ok is false if the value at the passed index is not of type uint64.
func (a Args) LookupUint64(i int) (v uint64, ok bool) {
	v, ok = a[i].(uint64)
	return
}

// LookupFloat32 returns the argument with the passed index as float32.
// ok is false if the value at the passed index is not of type float32.
func (a Args) LookupFloat32(i int) (v float32, ok bool) {
	v, ok = a[i].(float32)
	return
}

// LookupFloat64 returns the argument with the passed index as float64.
// ok is false if the value at the passed index is not of type float64.
func (a Args) LookupFloat64(i int) (v float64, ok bool) {
	v, ok = a[i].(float64)
	return
}

// LookupString returns the argument with the passed index as string.
// ok is false if the value at the passed index is not of type string.
func (a Args) LookupString(i int) (v string, ok bool) {
	v, ok = a[i].(string)
	return
}

// LookupMember returns the argument with the passed index as *discord.Member.
// ok is false if the value at the passed index is not of type *discord.Member.
func (a Args) LookupMember(i int) (v *discord.Member, ok bool) {
	v, ok = a[i].(*discord.Member)
	return
}

// LookupChannel returns the argument with the passed index as
// *discord.Channel.
// ok is false if the value at the passed index is not of type
// *discord.Channel.
func (a Args) LookupChannel(i int) (v *discord.Channel, ok bool) {
	v, ok = a[i].(*discord.Channel)
	return
}

// LookupRole returns the argument with the passed index as *discord.Role.
// ok is false if the value at the passed index is not of type *discord.Role.
func (a Args) LookupRole(i int) (v *discord.Role, ok bool) {
	v, ok = a[i].(*discord.Role)
	return
}

// LookupEmoji returns the argument with the passed index as *discord.Emoji.
// ok is false if the value at the passed index is not of type *discord.Emoji.
func (a Args) LookupEmoji(i int) (v *discord.Emoji, ok bool) {
	v, ok = a[i].(*discord.Emoji)
	return
}

// LookupDuration returns the argument with the passed index as time.Duration.
// ok is false if the value at the passed index is not of type time.Duration.
func (a Args) LookupDuration(i int) (v time.Duration, ok bool) {
	v, ok = a[i].(time.Duration)
	return
}

// LookupTime returns the argument with the passed index as time.Time.
// ok is false if the value at the passed index is not of type time.Time.
func (a Args) LookupTime(i int) (v time.Time, ok bool) {
	v, ok = a[i].(time.Time)
	return
}

// LookupRegexp returns the flag with the passed index as *regexp.Regexp.
// ok is false if the value at the passed index is not of type *regexp.Regexp.
func (a Args) LookupRegexp(i int) (v *regexp.Regexp, ok bool) {
	v, ok = a[i].(*regexp.Regexp)
	return
}

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

// MultiBool returns the flag with the passed name as []bool.
func (f Flags) MultiBool(name string) []bool { return f[name].([]bool) }

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

// LookupBool returns the flag with the passed name as bool.
// ok is false if there is either no flag with the passed key, or the value is
// not of type bool.
func (f Flags) LookupBool(name string) (v bool, ok bool) {
	v, ok = f[name].(bool)
	return
}

// LookupInt returns the flag with the passed name as int.
// ok is false if there is either no flag with the passed key, or the value is
// not of type int.
func (f Flags) LookupInt(name string) (v int, ok bool) {
	v, ok = f[name].(int)
	return
}

// LookupInt64 returns the flag with the passed name as int64.
// ok is false if there is either no flag with the passed key, or the value is
// not of type int64.
func (f Flags) LookupInt64(name string) (v int64, ok bool) {
	v, ok = f[name].(int64)
	return
}

// LookupUint returns the flag with the passed name as uint.
// ok is false if there is either no flag with the passed key, or the value is
// not of type uint.
func (f Flags) LookupUint(name string) (v uint, ok bool) {
	v, ok = f[name].(uint)
	return
}

// LookupUint64 returns the flag with the passed name as uint64.
// ok is false if there is either no flag with the passed key, or the value is
// not of type uint64.
func (f Flags) LookupUint64(name string) (v uint64, ok bool) {
	v, ok = f[name].(uint64)
	return
}

// LookupFloat32 returns the flag with the passed name as float32.
// ok is false if there is either no flag with the passed key, or the value is
// not of type float32.
func (f Flags) LookupFloat32(name string) (v float32, ok bool) {
	v, ok = f[name].(float32)
	return
}

// LookupFloat64 returns the flag with the passed name as float64.
// ok is false if there is either no flag with the passed key, or the value is
// not of type float64.
func (f Flags) LookupFloat64(name string) (v float64, ok bool) {
	v, ok = f[name].(float64)
	return
}

// LookupString returns the flag with the passed name as string.
// ok is false if there is either no flag with the passed key, or the value is
// not of type string.
func (f Flags) LookupString(name string) (v string, ok bool) {
	v, ok = f[name].(string)
	return
}

// LookupMember returns the flag with the passed name as *discord.Member.
// ok is false if there is either no flag with the passed key, or the value is
// not of type *discord.Member.
func (f Flags) LookupMember(name string) (v *discord.Member, ok bool) {
	v, ok = f[name].(*discord.Member)
	return
}

// LookupChannel returns the flag with the passed name as *discord.Channel.
// ok is false if there is either no flag with the passed key, or the value is
// not of type *discord.Channel.
func (f Flags) LookupChannel(name string) (v *discord.Channel, ok bool) {
	v, ok = f[name].(*discord.Channel)
	return
}

// LookupRole returns the flag with the passed name as *discord.Role.
// ok is false if there is either no flag with the passed key, or the value is
// not of type *discord.Role.
func (f Flags) LookupRole(name string) (v *discord.Role, ok bool) {
	v, ok = f[name].(*discord.Role)
	return
}

// LookupEmoji returns the flag with the passed name as *discord.Emoji.
// ok is false if there is either no flag with the passed key, or the value is
// not of type *discord.Emoji.
func (f Flags) LookupEmoji(name string) (v *discord.Emoji, ok bool) {
	v, ok = f[name].(*discord.Emoji)
	return
}

// LookupDuration returns the flag with the passed name as time.Duration.
// ok is false if there is either no flag with the passed key, or the value is
// not of type time.Duration.
func (f Flags) LookupDuration(name string) (v time.Duration, ok bool) {
	v, ok = f[name].(time.Duration)
	return
}

// LookupTime returns the flag with the passed name as time.Time.
// ok is false if there is either no flag with the passed key, or the value is
// not of type time.Time.
func (f Flags) LookupTime(name string) (v time.Time, ok bool) {
	v, ok = f[name].(time.Time)
	return
}

// LookupRegexp returns the flag with the passed name as *regexp.Regexp.
// ok is false if there is either no flag with the passed key, or the value is
// not of type *regexp.*regexp.Regexp.
func (f Flags) LookupRegexp(name string) (v *regexp.Regexp, ok bool) {
	v, ok = f[name].(*regexp.Regexp)
	return
}

// LookupMultiInt returns the flag with the passed name as []int.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []int.
func (f Flags) LookupMultiInt(name string) (v []int, ok bool) {
	v, ok = f[name].([]int)
	return
}

// LookupMultiInt64 returns the flag with the passed name as []int64.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []int64.
func (f Flags) LookupMultiInt64(name string) (v []int64, ok bool) {
	v, ok = f[name].([]int64)
	return
}

// LookupMultiUint returns the flag with the passed name as []uint.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []uint.
func (f Flags) LookupMultiUint(name string) (v []uint, ok bool) {
	v, ok = f[name].([]uint)
	return
}

// LookupMultiUint64 returns the flag with the passed name as []uint64.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []uint64.
func (f Flags) LookupMultiUint64(name string) (v []uint64, ok bool) {
	v, ok = f[name].([]uint64)
	return
}

// LookupMultiFloat32 returns the flag with the passed name as []float32.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []float32.
func (f Flags) LookupMultiFloat32(name string) (v []float32, ok bool) {
	v, ok = f[name].([]float32)
	return
}

// LookupMultiFloat64 returns the flag with the passed name as []Float64.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []Float64.
func (f Flags) LookupMultiFloat64(name string) (v []float64, ok bool) {
	v, ok = f[name].([]float64)
	return
}

// LookupMultiString returns the flag with the passed name as []string.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []string.
func (f Flags) LookupMultiString(name string) (v []string, ok bool) {
	v, ok = f[name].([]string)
	return
}

// LookupMultiBool returns the flag with the passed name as []discord.Member.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []discord.Member.
func (f Flags) LookupMultiMember(name string) (v []discord.Member, ok bool) {
	v, ok = f[name].([]discord.Member)
	return
}

// LookupMultiChannel returns the flag with the passed name as []discord.Channel.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []discord.Channel.
func (f Flags) LookupMultiChannel(name string) (v []discord.Channel, ok bool) {
	v, ok = f[name].([]discord.Channel)
	return
}

// LookupMultiRole returns the flag with the passed name as []discord.Role.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []discord.Role.
func (f Flags) LookupMultiRole(name string) (v []discord.Role, ok bool) {
	v, ok = f[name].([]discord.Role)
	return
}

// LookupMultiEmoji returns the flag with the passed name as []discord.Emoji.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []discord.Emoji.
func (f Flags) LookupMultiEmoji(name string) (v []discord.Emoji, ok bool) {
	v, ok = f[name].([]discord.Emoji)
	return
}

// LookupMultiDuration returns the flag with the passed name as []time.Duration.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []time.Duration.
func (f Flags) LookupMultiDuration(name string) (v []time.Duration, ok bool) {
	v, ok = f[name].([]time.Duration)
	return
}

// LookupMultiTime returns the flag with the passed name as []time.Time.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []time.Time.
func (f Flags) LookupMultiTime(name string) (v []time.Time, ok bool) {
	v, ok = f[name].([]time.Time)
	return
}

// LookupMultiRegexp returns the flag with the passed name as []*regexp.Regexp.
// ok is false if there is either no flag with the passed key, or the value is
// not of type []*regexp.Regexp.
func (f Flags) LookupMultiRegexp(name string) (v []*regexp.Regexp, ok bool) {
	v, ok = f[name].([]*regexp.Regexp)
	return
}
