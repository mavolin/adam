package arg

import (
	"github.com/mavolin/disstate/v2/pkg/state"

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
	info, err := genArgsInfo(l, c.Required, c.Optional, c.Flags, c.Variadic)
	if err != nil {
		return nil
	}

	return []plugin.ArgsInfo{info}
}
