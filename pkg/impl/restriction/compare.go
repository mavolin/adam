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

func ALL(funcs ...plugin.RestrictionFunc) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if len(funcs) == 0 {
			return nil
		} else if len(funcs) == 1 {
			return funcs[0](s, ctx)
		}

		missing := new(allError)

		for _, f := range funcs {
			err := f(s, ctx)

			switch err := err.(type) {
			case *errors.RestrictionError:
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				if err == errors.DefaultRestrictionError {
					return err
				}

				missing.restrictions = append(missing.restrictions, err)
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
				if err != nil {
					return err
				}
			}
		}

		// check if we have collected only a single error, and return it
		// directly if so
		if len(missing.restrictions) == 1 && len(missing.anys) == 0 {
			return missing.restrictions[0]
			// check if we have an error at all
		} else if len(missing.restrictions) != 0 || len(missing.anys) != 0 {
			return missing
		}

		return nil
	}
}

// ALLf works like ALL, but returns the passed returnError, if one of the
// plugin.RestrictionFuncs errors.
func ALLf(returnError error, funcs ...plugin.RestrictionFunc) plugin.RestrictionFunc {
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

func ANY(funcs ...plugin.RestrictionFunc) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if len(funcs) == 0 {
			return nil
		} else if len(funcs) == 1 {
			return funcs[0](s, ctx)
		}

		missing := new(anyError)

		for _, f := range funcs {
			err := f(s, ctx)

			switch err := err.(type) {
			case *errors.RestrictionError:
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				if err == errors.DefaultRestrictionError {
					return err
				}

				missing.restrictions = append(missing.restrictions, err)
			// we can just merge
			case *anyError:
				missing.restrictions = append(missing.restrictions, err.restrictions...)
				missing.alls = append(missing.alls, err.alls...)
			// there will always be 2 or more func in the error
			case *allError:
				missing.alls = append(missing.alls, err)
			default:
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				if err != nil {
					return err
				} else {
					return nil
				}
			}
		}

		// missing contains at least two errors, otherwise we would have
		// returned already
		return missing
	}
}

func ANYf(returnError error, funcs ...plugin.RestrictionFunc) plugin.RestrictionFunc {
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
	restrictions []*errors.RestrictionError
	anys         []*anyError
}

func (e *allError) format(indentLvl int, l *localization.Localizer) (string, error) {
	indent, nlIndent := genIndent(indentLvl)

	var s string

	for _, m := range e.restrictions {
		desc, err := m.Description(l)
		if err != nil {
			return "", err
		}

		s += "\n" + indent + entryPrefix + newlineRegexp.ReplaceAllStringFunc(desc, func(s string) string {
			return "\n" + nlIndent + s[1:]
		})
	}

	// we can ignore the error, as there is a fallback
	anyMessage, _ := l.Localize(anyMessageInline)

	for _, m := range e.anys {
		s += "\n" + indent + entryPrefix + anyMessage + "\n"

		msg, err := m.format(indentLvl+1, l)
		if err != nil {
			return "", err
		}

		s += msg
	}

	return s[1:], nil // strip the first newline
}

func (e *allError) Error() string {
	return "allError"
}

type anyError struct {
	restrictions []*errors.RestrictionError
	alls         []*allError
}

func (e *anyError) format(indentLvl int, l *localization.Localizer) (string, error) {
	indent, nlIndent := genIndent(indentLvl)

	var s string

	for _, m := range e.restrictions {
		desc, err := m.Description(l)
		if err != nil {
			return "", err
		}

		s += "\n" + indent + entryPrefix + newlineRegexp.ReplaceAllStringFunc(desc, func(s string) string {
			return "\n" + nlIndent + s[1:]
		})
	}

	// we can ignore the error, as there is a fallback
	allMessage, _ := l.Localize(allMessageInline)

	for _, m := range e.alls {
		s += "\n" + indent + entryPrefix + allMessage + "\n"

		msg, err := m.format(indentLvl+1, l)
		if err != nil {
			return "", err
		}

		s += msg
	}

	return s[1:], nil // strip the first newline
}

func (e *anyError) Error() string {
	return "anyError"
}

// genIndent generates the necessary amount of indent for the passed indent
// level.
// The first return value is the normal indent, and the second is the indent
// needed for newlines within an entry.
func genIndent(indentLvl int) (indent, nlIndent string) {
	// use a zero width space to prevent trimming
	indent = strings.Repeat("　", indentLvl*indentMultiplicator)
	nlIndent = strings.Repeat("　", indentLvl*indentMultiplicator+len(entryPrefix))

	return
}
