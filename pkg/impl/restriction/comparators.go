package restriction

import (
	"regexp"
	"strings"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

const (
	// indentMultiplier defines the amount of whitespaces per indent level.
	indentMultiplier = 2
	// entryPrefix is the prefix used in front of every entry.
	entryPrefix = "• "
)

// EmbeddableError is an error with different messages for embedding in a any
// or all error than when returned directly.
type EmbeddableError struct {
	// EmbeddableVersion is the version used when embedded in an any or all
	// error.
	EmbeddableVersion *plugin.RestrictionError
	// DefaultVersion is the version returned if the error won't get embedded.
	DefaultVersion error
}

func (e *EmbeddableError) Wrap(*state.State, *plugin.Context) error { return e.DefaultVersion }
func (e *EmbeddableError) Error() string                            { return e.DefaultVersion.Error() }

// All asserts that all of the passed plugin.RestrictionFuncs return nil.
// If not, it will create an error containing a list of all missing
// requirements using the returned errors.
//
// This list can only be created, if the returned errors are of type
// *plugin.RestrictionError or *EmbeddableError, or are error lists returned by
// nested Alls or Anys.
//
// If at least one of the passed plugin.RestrictionFuncs returns a fatal
// plugin.RestrictionError, the error produced by the returned function will be
// fatal as well.
func All(funcs ...plugin.RestrictionFunc) plugin.RestrictionFunc { //nolint:gocognit
	return func(s *state.State, ctx *plugin.Context) error {
		if len(funcs) == 0 {
			return nil
		} else if len(funcs) == 1 {
			return funcs[0](s, ctx)
		}

		missing := new(allError)

		var embeddable *EmbeddableError

		for _, f := range funcs {
			err := f(s, ctx)
			if err == nil {
				continue
			}

			if rerr := new(plugin.RestrictionError); errors.As(err, &rerr) {
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				if errors.Is(err, plugin.DefaultRestrictionError) {
					return err
				} else if errors.Is(err, plugin.DefaultFatalRestrictionError) {
					return err
				}

				missing.restrictions = append(missing.restrictions, rerr)
			} else if eerr := new(EmbeddableError); errors.As(err, &eerr) {
				embeddable = eerr

				if errors.Is(eerr.EmbeddableVersion, plugin.DefaultRestrictionError) {
					return err
				} else if errors.Is(err, plugin.DefaultFatalRestrictionError) {
					return err
				}

				missing.restrictions = append(missing.restrictions, eerr.EmbeddableVersion)
			} else if aerr := new(allError); errors.As(err, &aerr) { // we can just merge
				missing.restrictions = append(missing.restrictions, aerr.restrictions...)
				missing.anys = append(missing.anys, aerr.anys...)
			} else if aerr := new(anyError); errors.As(err, &aerr) {
				missing.anys = append(missing.anys, aerr)
			} else {
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				return err
			}
		}

		// check if we have collected only a single error, and return it
		// directly if so
		switch {
		case len(missing.restrictions) == 1 && len(missing.anys) == 0:
			if embeddable != nil { // if it is embeddable, it will be stored here
				return embeddable
			}

			return missing.restrictions[0]
		case len(missing.restrictions) == 0 && len(missing.anys) == 1:
			return missing.anys[0]
		}

		// check if we have an error at all
		if len(missing.restrictions) != 0 || len(missing.anys) != 0 {
			return missing
		}

		return nil
	}
}

// Allf works like All, but returns the passed returnError if one of the
// plugin.RestrictionFuncs returns an error.
func Allf(returnError error, funcs ...plugin.RestrictionFunc) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		for _, f := range funcs {
			err := f(s, ctx)
			if err != nil {
				return returnError
			}
		}

		return nil
	}
}

// Any asserts that at least one of the passed plugin.RestrictionFuncs returns
// no error.
// If not it will return a list of the needed requirements using the returned
// errors.
//
// This list can only be created, if the returned errors are of type
// *plugin.RestrictionError or *EmbeddableError, or are error lists returned by
// nested Alls or Anys.
//
// If at least one of the passed plugin.RestrictionFuncs returns a non-fatal
// plugin.RestrictionError, the error produced by the returned function will
// not be fatal as well.
func Any(funcs ...plugin.RestrictionFunc) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if len(funcs) == 0 {
			return nil
		} else if len(funcs) == 1 {
			return funcs[0](s, ctx)
		}

		missing := new(anyError)

		for _, f := range funcs {
			err := f(s, ctx)
			if err == nil {
				return nil
			}

			if rerr := new(plugin.RestrictionError); errors.As(err, &rerr) {
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				if errors.Is(err, plugin.DefaultRestrictionError) {
					return err
				} else if errors.Is(err, plugin.DefaultFatalRestrictionError) {
					return err
				}

				missing.restrictions = append(missing.restrictions, rerr)
			} else if eerr := new(EmbeddableError); errors.As(err, &eerr) {
				if errors.Is(eerr.EmbeddableVersion, plugin.DefaultRestrictionError) {
					return err
				} else if errors.Is(err, plugin.DefaultFatalRestrictionError) {
					return err
				}

				missing.restrictions = append(missing.restrictions, eerr.EmbeddableVersion)
			} else if aerr := new(anyError); errors.As(err, &aerr) { // we can just merge
				missing.restrictions = append(missing.restrictions, aerr.restrictions...)
				missing.alls = append(missing.alls, aerr.alls...)
			} else if aerr := new(allError); errors.As(err, &aerr) {
				missing.alls = append(missing.alls, aerr)
			} else {
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				return err
			}
		}

		// missing contains at least two errors, otherwise we would have
		// returned already
		return missing
	}
}

// Anyf works like Any, but returns the passed returnError if all of the passed
// plugin.RestrictionFuncs return an error.
func Anyf(returnError error, funcs ...plugin.RestrictionFunc) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		for _, f := range funcs {
			err := f(s, ctx)
			if err == nil {
				return nil
			}
		}

		return returnError
	}
}

var newlineRegexp = regexp.MustCompile(`\n[^\n]`)

type allError struct {
	restrictions []*plugin.RestrictionError
	anys         []*anyError
}

func (e *allError) format(indentLvl int, l *i18n.Localizer) (s string, fatal bool, err error) {
	indent, nlIndent := genIndent(indentLvl)

	fatal = false

	for _, r := range e.restrictions {
		if r.Fatal {
			fatal = true
		}

		desc, err := r.Description(l)
		if err != nil {
			return "", false, err
		}

		s += "\n" + indent + entryPrefix + newlineRegexp.ReplaceAllStringFunc(desc, func(s string) string {
			return "\n" + nlIndent + s[1:]
		})
	}

	// we can ignore the error, as there is a fallback
	anyMessage, _ := l.Localize(anyMessageInline)

	for _, a := range e.anys {
		s += "\n" + indent + entryPrefix + anyMessage + "\n"

		msg, subFatal, err := a.format(indentLvl+1, l)
		if err != nil {
			return "", false, err
		}

		if subFatal {
			fatal = true
		}

		s += msg
	}

	// strip the first newline
	return s[1:], fatal, nil
}

func (e *allError) Wrap(_ *state.State, ctx *plugin.Context) error {
	missing, fatal, err := e.format(0, ctx.Localizer)
	if err != nil {
		return err
	}

	header, _ := ctx.Localize(allMessageHeader)

	if fatal {
		return plugin.NewFatalRestrictionError(header + "\n\n" + missing)
	}

	return plugin.NewRestrictionError(header + "\n\n" + missing)
}

func (e *allError) Error() string {
	return "allError"
}

type anyError struct {
	restrictions []*plugin.RestrictionError
	alls         []*allError
}

func (e *anyError) format(indentLvl int, l *i18n.Localizer) (s string, fatal bool, err error) {
	indent, nlIndent := genIndent(indentLvl)

	fatal = true

	for _, r := range e.restrictions {
		if !r.Fatal {
			fatal = false
		}

		desc, err := r.Description(l)
		if err != nil {
			return "", false, err
		}

		s += "\n" + indent + entryPrefix + newlineRegexp.ReplaceAllStringFunc(desc, func(s string) string {
			return "\n" + nlIndent + s[1:]
		})
	}

	// we can ignore the error, as there is a fallback
	allMessage, _ := l.Localize(allMessageInline)

	for _, a := range e.alls {
		s += "\n" + indent + entryPrefix + allMessage + "\n"

		msg, subFatal, err := a.format(indentLvl+1, l)
		if err != nil {
			return "", false, err
		}

		if !subFatal {
			fatal = false
		}

		s += msg
	}

	// strip the first newline
	return s[1:], fatal, nil
}

func (e *anyError) Wrap(_ *state.State, ctx *plugin.Context) error {
	missing, fatal, err := e.format(0, ctx.Localizer)
	if err != nil {
		return err
	}

	header, _ := ctx.Localize(anyMessageHeader)

	if fatal {
		return plugin.NewFatalRestrictionError(header + "\n\n" + missing)
	}

	return plugin.NewRestrictionError(header + "\n\n" + missing)
}

func (e *anyError) Error() string {
	return "anyError"
}

// genIndent generates the necessary amount of indent for the passed indent
// level.
// The first return value is the normal indent, and the second is the indent
// needed for newlines within an entry.
func genIndent(indentLvl int) (indent, nlIndent string) {
	// use an "ideographic space" for indenting, as Discord strips whitespace
	// on new lines in embeds.
	indent = strings.Repeat("\u3000", indentLvl*indentMultiplier)
	nlIndent = strings.Repeat("\u3000", indentLvl*indentMultiplier)

	if nlIndent == "" { // prefix a zero-width whitespace if the indentLvl is 0
		nlIndent += "\u200b"
	}

	nlIndent += strings.Repeat(" ", len(entryPrefix))

	return
}
