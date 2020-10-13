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
// Additionally, in order to distinguish flags from arguments, a double minus
// ('--') must be used to use a literal minus.
// To ease usability for users unaware of this, escapes are not needed, if one
// of the following cases is fulfilled:
//
// 1. The minus is used inside of a flag or argument, i.e. anywhere but the
// beginning.
//
// 2. The minus is used in any of the required arguments except the first.
//
// 3. There are no flags.
//
// This still means, however, that a double minus will always be interpreted as
// a single one.
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

func (c CommaConfig) Parse(args string, s *state.State, ctx *plugin.Context) (plugin.Args, plugin.Flags, error) {
	parser := newCommaParser(args, c, s, ctx)
	return parser.parse()
}

func (c CommaConfig) Info(l *i18n.Localizer) []plugin.ArgsInfo {
	info := plugin.ArgsInfo{
		Required: make([]plugin.ArgInfo, len(c.RequiredArgs)),
		Optional: make([]plugin.ArgInfo, len(c.OptionalArgs)),
		Flags:    make([]plugin.FlagInfo, len(c.Flags)),
		Variadic: c.Variadic,
	}

	var err error

	for i, arg := range c.RequiredArgs {
		info.Required[i], err = requiredArgInfo(arg, l)
		if err != nil {
			return nil
		}
	}

	for i, arg := range c.OptionalArgs {
		info.Optional[i], err = optionalArgInfo(arg, l)
		if err != nil {
			return nil
		}
	}

	for i, flag := range c.Flags {
		info.Flags[i], err = flagInfo(flag, l)
		if err != nil {
			return nil
		}
	}

	return []plugin.ArgsInfo{info}
}

func (c CommaConfig) flag(name string) *Flag {
	for _, flag := range c.Flags {
		if flag.Name == name {
			return &flag
		}

		for _, alias := range flag.Aliases {
			if alias == name {
				return &flag
			}
		}
	}

	return nil
}
