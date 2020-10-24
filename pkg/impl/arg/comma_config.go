package arg

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// CommaConfig is a simple plugin.ArgConfig.
// It uses a comma to separate flags/arguments, literal commas can be escaped
// using a double comma (',,').
//
// Flags can be placed both at the beginning and the end of the arguments.
// Flag names mustn't contain double minuses, commas and whitespace (space, tab
// or newline).
//
// Additionally, in order to distinguish flags from arguments, if an argument
// beings with a minus, it must be escaped using a double minus.
// To ease usability for users unaware of this, escapes are not needed, if one
// of the following cases is fulfilled:
//
// 1. The minus is used in any of the required arguments except the first.
//
// 2. There are no flags.
//
// Even if one of those exceptions applies, it will still be escaped to
// preserve predictability.
type CommaConfig struct {
	// RequiredArgs contains the required arguments.
	RequiredArgs []RequiredArg
	// OptionalArgs contains the optional arguments.
	OptionalArgs []OptionalArg
	// Variadic specifies whether the last possibly specifiable argument is
	// variadic.
	Variadic bool

	// Flags contains the flags.
	Flags []Flag
}

var _ plugin.ArgConfig = CommaConfig{}

func (c CommaConfig) Parse(args string, s *state.State, ctx *plugin.Context) (plugin.Args, plugin.Flags, error) {
	parser := newCommaParser(args, c, s, ctx)
	return parser.parse()
}

func (c CommaConfig) Info(l *i18n.Localizer) []plugin.ArgsInfo {
	info, err := genArgsInfo(l, c.RequiredArgs, c.OptionalArgs, c.Flags, c.Variadic)
	if err != nil {
		return nil
	}

	return []plugin.ArgsInfo{info}
}
