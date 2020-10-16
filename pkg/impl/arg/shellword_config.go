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
// To use arguments with whitespace, quotes, both single and double can be
// used.
// Additionally, lines of code as well as code blocks will also be parsed as a
// single argument.
//
// Escapes are only permitted, if using double quotes.
// Valid escapes are '\\' and '\"', all other combinations will be interpreted
// literally, to make usage easier for users unaware of shell notation.
type ShellwordConfig struct {
	// RequiredArgs contains the required arguments.
	RequiredArgs []RequiredArgument
	// OptionalArgs contains the optional arguments.
	OptionalArgs []OptionalArgument
	// Variadic specifies whether the last possibly specifiable argument is
	// variadic.
	Variadic bool

	// Flags contains the flags.
	Flags []Flag
}

var _ plugin.ArgConfig = ShellwordConfig{}

func (c ShellwordConfig) Parse(args string, s *state.State, ctx *plugin.Context) (plugin.Args, plugin.Flags, error) {
	parser := newShellwordParser(args, c, s, ctx)
	return parser.parse()
}

func (c ShellwordConfig) Info(l *i18n.Localizer) []plugin.ArgsInfo {
	info, err := genArgsInfo(l, c.RequiredArgs, c.OptionalArgs, c.Flags, c.Variadic)
	if err != nil {
		return nil
	}

	return []plugin.ArgsInfo{info}
}
