package arg //nolint:dupl

import (
	"strings"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// CommaConfig is an unlocalized plugin.ArgConfig, that uses a comma to
// separate flags and arguments.
// Literal commas can be escaped using a double comma (',,').
//
// Flags can be placed both at the beginning and the end of the arguments.
// Their names mustn't contain double minuses, commas or whitespace.
//
// Additionally, in order to distinguish flags from arguments argument that
// start with a minus must be escaped using a double minus.
// To ease usability for users unaware of this escapes are not needed if one
// of the following conditions is fulfilled:
//
// 1. The minus is used in any of the required arguments except the first.
//
// 2. There are no flags.
//
// Even if one of those exceptions applies, double minuses will still be parsed
// as a single minus to preserve predictability.
type CommaConfig struct {
	// Required contains the required arguments.
	Required []RequiredArg
	// Optional contains the optional arguments.
	Optional []OptionalArg
	// Variadic specifies whether the last possibly specifiable argument is
	// variadic.
	Variadic bool

	// Flags contains the flags.
	Flags []Flag
}

var (
	_ plugin.ArgConfig  = CommaConfig{}
	_ plugin.ArgsInfoer = CommaConfig{}
)

func (c CommaConfig) Parse(args string, s *state.State, ctx *plugin.Context) (plugin.Args, plugin.Flags, error) {
	parser := newCommaParser(args, c, s, ctx)
	return parser.parse()
}

func (c CommaConfig) Info(l *i18n.Localizer) []plugin.ArgsInfo {
	// still use a Localizer, so replacements for the type information will
	// be used, if there are any
	info := genArgsInfo(l, c.Required, c.Optional, c.Flags, c.Variadic)
	info.ArgsFormatter = newCommaFormatter(info)

	return []plugin.ArgsInfo{info}
}

// LocalizedCommaConfig is a localized plugin.ArgConfig, that uses a comma to
// separate flags and arguments.
// Literal commas can be escaped using a double comma (',,').
//
// Flags can be placed both at the beginning and the end of the arguments.
// Their names mustn't contain double minuses, commas or whitespace.
//
// Additionally, in order to distinguish flags from arguments argument that
// start with a minus must be escaped using a double minus.
// To ease usability for users unaware of this escapes are not needed if one
// of the following conditions is fulfilled:
//
// 1. The minus is used in any of the required arguments except the first.
//
// 2. There are no flags.
//
// Even if one of those exceptions applies, double minuses will still be parsed
// as a single minus to preserve predictability.
type LocalizedCommaConfig struct {
	// Required contains the required arguments.
	Required []LocalizedRequiredArg
	// Optional contains the optional arguments.
	Optional []LocalizedOptionalArg
	// Variadic specifies whether the last possibly specifiable argument is
	// variadic.
	Variadic bool

	// Flags contains the flags.
	Flags []LocalizedFlag
}

var (
	_ plugin.ArgConfig  = LocalizedCommaConfig{}
	_ plugin.ArgsInfoer = LocalizedCommaConfig{}
)

func (c LocalizedCommaConfig) Parse(
	args string, s *state.State, ctx *plugin.Context,
) (plugin.Args, plugin.Flags, error) {
	parser := newCommaParserl(args, c, s, ctx)
	return parser.parse()
}

func (c LocalizedCommaConfig) Info(l *i18n.Localizer) []plugin.ArgsInfo {
	info, ok := genArgsInfol(l, c.Required, c.Optional, c.Flags, c.Variadic)
	if !ok {
		return nil
	}

	info.ArgsFormatter = newCommaFormatter(info)

	return []plugin.ArgsInfo{info}
}

func newCommaFormatter(info plugin.ArgsInfo) func(f plugin.ArgFormatter) string {
	return func(f plugin.ArgFormatter) string {
		var b strings.Builder

		// provision 20 characters per arg
		b.Grow((len(info.Required) + len(info.Optional)) * 20)

		for i, ai := range info.Required {
			if i > 0 {
				b.WriteString(", ")
			}

			variadic := info.Variadic && i == len(info.Required)-1 && len(info.Optional) == 0
			b.WriteString(f(ai, false, variadic))
		}

		for i, ai := range info.Optional {
			if i > 0 || len(info.Required) > 0 {
				b.WriteString(", ")
			}

			variadic := info.Variadic && i == len(info.Optional)-1
			b.WriteString(f(ai, true, variadic))
		}

		return b.String()
	}
}
