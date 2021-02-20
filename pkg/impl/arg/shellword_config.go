package arg //nolint:dupl

import (
	"strings"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// ShellwordConfig is a plugin.ArgConfig that roughly follows the parsing rules
// of the Bourne shell.
//
// Flags
//
// Flags can be placed both before and after arguments.
// For simplicity, flags always start with a single minus, double minuses are
// not permitted.
//
// Arguments
//
// Arguments are space separated.
// To use arguments with whitespace quotes, both single and double, can be
// used.
// Additionally, lines of code as well as code blocks will be parsed as a
// single argument.
//
// Escapes
//
// Escapes are only permitted if using double quotes.
// Valid escapes are '\\' and '\"', all other combinations will be parsed
// literally to make usage easier for users unaware of escapes.
type ShellwordConfig struct {
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
	_ plugin.ArgConfig  = ShellwordConfig{}
	_ plugin.ArgsInfoer = ShellwordConfig{}
)

func (c ShellwordConfig) Parse(args string, s *state.State, ctx *plugin.Context) (plugin.Args, plugin.Flags, error) {
	parser := newShellwordParser(args, c, s, ctx)
	return parser.parse()
}

func (c ShellwordConfig) Info(l *i18n.Localizer) []plugin.ArgsInfo {
	info := genArgsInfo(l, c.Required, c.Optional, c.Flags, c.Variadic)
	info.ArgsFormatter = newShellwordFormatter(info)

	return []plugin.ArgsInfo{info}
}

// LocalizedShellwordConfig is a plugin.ArgConfig that roughly follows the
// parsing rules of the Bourne shell.
//
// Flags
//
// Flags can be placed both before and after arguments.
// For simplicity, flags always start with a single minus, double minuses are
// not permitted.
//
// Arguments
//
// Arguments are space separated.
// To use arguments with whitespace quotes, both single and double, can be
// used.
// Additionally, lines of code as well as code blocks will be parsed as a
// single argument.
//
// Escapes
//
// Escapes are only permitted if using double quotes.
// Valid escapes are '\\' and '\"', all other combinations will be parsed
// literally to make usage easier for users unaware of escapes.
type LocalizedShellwordConfig struct {
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
	_ plugin.ArgConfig  = LocalizedShellwordConfig{}
	_ plugin.ArgsInfoer = LocalizedShellwordConfig{}
)

func (c LocalizedShellwordConfig) Parse(
	args string, s *state.State, ctx *plugin.Context,
) (plugin.Args, plugin.Flags, error) {
	parser := newShellwordParserl(args, c, s, ctx)
	return parser.parse()
}

func (c LocalizedShellwordConfig) Info(l *i18n.Localizer) []plugin.ArgsInfo {
	info, ok := genArgsInfol(l, c.Required, c.Optional, c.Flags, c.Variadic)
	if !ok {
		return nil
	}

	info.ArgsFormatter = newShellwordFormatter(info)

	return []plugin.ArgsInfo{info}
}

func newShellwordFormatter(info plugin.ArgsInfo) func(f plugin.ArgFormatter) string {
	return func(f plugin.ArgFormatter) string {
		var b strings.Builder

		// provision 20 characters per arg
		b.Grow((len(info.Required) + len(info.Optional)) * 20)

		for i, ai := range info.Required {
			if i > 0 {
				b.WriteString(" ")
			}

			variadic := info.Variadic && i == len(info.Required)-1 && len(info.Optional) == 0
			b.WriteString(f(ai, false, variadic))
		}

		for i, ai := range info.Optional {
			if i > 0 || len(info.Required) > 0 {
				b.WriteString(" ")
			}

			variadic := info.Variadic && i == len(info.Optional)-1
			b.WriteString(f(ai, true, variadic))
		}

		return b.String()
	}
}
