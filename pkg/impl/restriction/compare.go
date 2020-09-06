package restriction

import (
	"regexp"
	"strings"

	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

const (
	// indentMultiplicator defines the amount of whitespaces per indent level.
	indentMultiplicator = 2
	// entryPrefix is the prefix used in front of every entry.
	entryPrefix = "• "
)

// EmbeddableError is an error with different messages for embedding in a any
// or all error than when returned directly.
type EmbeddableError struct {
	// EmbeddableVersion is the version used when embedded in an any or all
	// error.
	// This MUST be of type *errors.RestrictionError or
	// *errors.FatalRestrictionError.
	EmbeddableVersion error
	// DefaultVersion is the version returned if the error won't get embedded.
	DefaultVersion error
}

func (e *EmbeddableError) Wrap(*state.State, *plugin.Context) error { return e.DefaultVersion }
func (e *EmbeddableError) Error() string                            { return e.DefaultVersion.Error() }

// All asserts that all of the passed plugin.RestrictionFuncs return nil.
// If not, it will create an error containing a list of all missing
// requirements using the returned errors.
// For this list to be created, the error must either be of type
// *errors.RestrictionError or *EmbeddableError, or must be a nested All or
// Any.
func All(funcs ...plugin.RestrictionFunc) plugin.RestrictionFunc {
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

			switch err := err.(type) {
			case *errors.RestrictionError:
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				if errors.Is(err, errors.DefaultRestrictionError) {
					return err
				}

				missing.restrictions = append(missing.restrictions, err)
			case *errors.FatalRestrictionError:
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				if errors.Is(err, errors.DefaultFatalRestrictionError) {
					return err
				}

				missing.fatalRestrictions = append(missing.fatalRestrictions, err)
			case *EmbeddableError:
				embeddable = err

				switch emver := err.EmbeddableVersion.(type) {
				case *errors.RestrictionError:
					if errors.Is(emver, errors.DefaultRestrictionError) {
						return err
					}

					missing.restrictions = append(missing.restrictions, emver)
				case *errors.FatalRestrictionError:
					if errors.Is(emver, errors.DefaultFatalRestrictionError) {
						return err
					}

					missing.fatalRestrictions = append(missing.fatalRestrictions, emver)
				default:
					return errors.DefaultFatalRestrictionError
				}

			// we can just merge
			case *allError:
				missing.restrictions = append(missing.restrictions, err.restrictions...)
				missing.anys = append(missing.anys, err.anys...)
			// there will always be 2 or more func in the error
			case *anyError:
				missing.anys = append(missing.anys, err)
			default:
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				return err
			}
		}

		// check if we have collected only a single error, and return it
		// directly if so
		switch {
		case len(missing.restrictions) == 1 && len(missing.fatalRestrictions) == 0 && len(missing.anys) == 0:
			if embeddable != nil { // if it is embeddable, it will be stored here
				return embeddable
			}

			return missing.restrictions[0]
		case len(missing.restrictions) == 0 && len(missing.fatalRestrictions) == 1 && len(missing.anys) == 0:
			if embeddable != nil {
				return embeddable
			}

			return missing.fatalRestrictions[0]
		case len(missing.restrictions) == 0 && len(missing.fatalRestrictions) == 0 && len(missing.anys) == 1:
			return missing.anys[0]
		}

		// check if we have an error at all
		if len(missing.restrictions) != 0 || len(missing.fatalRestrictions) != 0 || len(missing.anys) != 0 {
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
// For this list to be created, the error must either be of type
// *errors.RestrictionError or *EmbeddableError, or must be a nested All or
// Any.
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

			switch err := err.(type) {
			case *errors.RestrictionError:
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				if errors.Is(err, errors.DefaultRestrictionError) {
					return err
				}

				missing.restrictions = append(missing.restrictions, err)
			case *errors.FatalRestrictionError:
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				if errors.Is(err, errors.DefaultFatalRestrictionError) {
					return err
				}

				missing.fatalRestrictions = append(missing.fatalRestrictions, err)
			case *EmbeddableError:
				switch emver := err.EmbeddableVersion.(type) {
				case *errors.RestrictionError:
					if errors.Is(emver, errors.DefaultRestrictionError) {
						return err
					}

					missing.restrictions = append(missing.restrictions, emver)
				case *errors.FatalRestrictionError:
					if errors.Is(emver, errors.DefaultFatalRestrictionError) {
						return err
					}

					missing.fatalRestrictions = append(missing.fatalRestrictions, emver)
				default:
					return errors.DefaultFatalRestrictionError
				}
			// we can just merge
			case *anyError:
				missing.restrictions = append(missing.restrictions, err.restrictions...)
				missing.alls = append(missing.alls, err.alls...)
			// there will always be 2 or more funcs in the error
			case *allError:
				missing.alls = append(missing.alls, err)
			default:
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
	restrictions      []*errors.RestrictionError
	fatalRestrictions []*errors.FatalRestrictionError
	anys              []*anyError
}

func (e *allError) format(indentLvl int, l *localization.Localizer) (s string, fatal bool, err error) {
	indent, nlIndent := genIndent(indentLvl)

	fatal = len(e.fatalRestrictions) > 0

	for _, m := range e.fatalRestrictions {
		desc, err := m.Description(l)
		if err != nil {
			return "", false, err
		}

		s += "\n" + indent + entryPrefix + newlineRegexp.ReplaceAllStringFunc(desc, func(s string) string {
			return "\n" + nlIndent + s[1:]
		})
	}

	for _, m := range e.restrictions {
		desc, err := m.Description(l)
		if err != nil {
			return "", false, err
		}

		s += "\n" + indent + entryPrefix + newlineRegexp.ReplaceAllStringFunc(desc, func(s string) string {
			return "\n" + nlIndent + s[1:]
		})
	}

	// we can ignore the error, as there is a fallback
	anyMessage, _ := l.Localize(anyMessageInline)

	for _, m := range e.anys {
		s += "\n" + indent + entryPrefix + anyMessage + "\n"

		msg, subFatal, err := m.format(indentLvl+1, l)
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
		return errors.NewFatalRestrictionError(header + "\n\n" + missing)
	}

	return errors.NewRestrictionError(header + "\n\n" + missing)
}

func (e *allError) Error() string {
	return "allError"
}

type anyError struct {
	restrictions      []*errors.RestrictionError
	fatalRestrictions []*errors.FatalRestrictionError
	alls              []*allError
}

func (e *anyError) format(indentLvl int, l *localization.Localizer) (s string, fatal bool, err error) {
	indent, nlIndent := genIndent(indentLvl)

	fatal = len(e.restrictions) == 0

	for _, m := range e.fatalRestrictions {
		desc, err := m.Description(l)
		if err != nil {
			return "", false, err
		}

		s += "\n" + indent + entryPrefix + newlineRegexp.ReplaceAllStringFunc(desc, func(s string) string {
			return "\n" + nlIndent + s[1:]
		})
	}

	for _, m := range e.restrictions {
		desc, err := m.Description(l)
		if err != nil {
			return "", false, err
		}

		s += "\n" + indent + entryPrefix + newlineRegexp.ReplaceAllStringFunc(desc, func(s string) string {
			return "\n" + nlIndent + s[1:]
		})
	}

	// we can ignore the error, as there is a fallback
	allMessage, _ := l.Localize(allMessageInline)

	for _, m := range e.alls {
		s += "\n" + indent + entryPrefix + allMessage + "\n"

		msg, subFatal, err := m.format(indentLvl+1, l)
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
		return errors.NewFatalRestrictionError(header + "\n\n" + missing)
	}

	return errors.NewRestrictionError(header + "\n\n" + missing)
}

func (e *anyError) Error() string {
	return "anyError"
}

// genIndent generates the necessary amount of indent for the passed indent
// level.
// The first return value is the normal indent, and the second is the indent
// needed for newlines within an entry.
func genIndent(indentLvl int) (indent, nlIndent string) {
	// use an "ideographic space" for indenting, as discord strips whitespace
	// on new lines in embeds.
	indent = strings.Repeat("\u3000", indentLvl*indentMultiplicator)
	nlIndent = strings.Repeat("\u3000", indentLvl*indentMultiplicator)

	if nlIndent == "" { // prefix a zero-width whitespace if the indentLvl is 0
		nlIndent += "\u200b"
	}

	nlIndent += strings.Repeat(" ", len(entryPrefix))

	return
}
