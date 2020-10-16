// Package arg provides implementations for the argument abstractions found in
// Package plugins.
package arg

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
)

type Kind string

const (
	KindArgument = "argument"
	KindFlag     = "flag"
)

type (
	// RequiredArgument is an unlocalized required argument.
	RequiredArgument struct {
		// Name is the name of the argument.
		Name i18nutil.Text
		// Type is the type of the argument.
		Type Type
		// Description is an optional short description of the argument.
		Description i18nutil.Text
	}
	// OptionalArgument is an unlocalized optional argument.
	OptionalArgument struct {
		// Name is the name of the argument.
		Name i18nutil.Text
		// Type is the type of the argument.
		Type Type
		// Default is the default value of the argument.
		//
		// If Default is (interface{})(nil), the default of Type will be used,
		// as returned by Type.Default() will be used.
		Default interface{}
		// Description is an optional short description of the argument.
		Description i18nutil.Text
	}
	// Flag is an unlocalized flag.
	Flag struct {
		// Name is the name of the flag.
		Name string
		// Aliases contains the optional aliases of the flag.
		Aliases []string
		// Type is the type of the flag.
		Type Type
		// Default is the default value of the flag, and is used if the flag
		// doesn't get set.
		//
		// If Default is (interface{})(nil), the default of Type will be used,
		// as returned by Type.Default() will be used.
		Default interface{}
		// Description is an optional short description of the flag.
		Description i18nutil.Text
		// Multi specifies whether this flag can be used multiple times.
		Multi bool
	}

	// Type is the abstraction of a type.
	Type interface {
		// Name returns the name of the type.
		// The name should be a noun.
		Name(l *i18n.Localizer) string
		// Description is an optional short description of the type.
		Description(l *i18n.Localizer) string
		// Parse parses the argument or flag using the passed Context.
		//
		// The first return value must always be of the same type.
		Parse(s *state.State, ctx *Context) (interface{}, error)
		// Default returns the default value for the type.
		// See Flag.Default or OptionalArgument.Default for more info.
		//
		// It must return a value that is of the type returned by Parse.
		Default() interface{}
	}

	// Context is the context passed to Type.Parse.
	Context struct {
		*plugin.Context

		// Raw is the raw argument or flag.
		Raw string
		// Name is the name of the argument or flag.
		Name string
		// UsedName is the alias of the flag the Context represents.
		// If the name of the flag was used, or the context represents an
		// argument, UsedName will be equal to Name.
		UsedName string
		// Index contains the index of the argument, if the context represents
		// an argument.
		Index int
		// Kind specifies whether the parsed argument is a flag or an argument.
		Kind Kind
	}
)
