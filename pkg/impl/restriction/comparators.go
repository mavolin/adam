package restriction

import (
	"regexp"
	"strings"

	"github.com/mavolin/disstate/v4/pkg/state"

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

// All asserts that all the passed plugin.RestrictionFuncs return nil.
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
//nolint:gocognit
func All(funcs ...plugin.RestrictionFunc) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if len(funcs) == 0 {
			return nil
		} else if len(funcs) == 1 {
			return funcs[0](s, ctx)
		}

		missing, err := newAllError(ctx.Localizer)
		if err != nil {
			return err
		}

		var embeddable *EmbeddableError

		for _, f := range funcs {
			err := f(s, ctx)
			if err == nil {
				continue
			}

			if eerr := new(EmbeddableError); errors.As(err, &eerr) {
				embeddable = eerr

				ok, addErr := missing.addRestrictions(ctx, eerr.EmbeddableVersion)
				if addErr != nil {
					return addErr
				}

				if !ok {
					return err
				}
			} else if aerr := new(allError); errors.As(err, &aerr) { // we can just merge
				missing.restrictions = append(missing.restrictions, aerr.restrictions...)

				if aerr.fatal {
					missing.fatal = true
				}

				if err := missing.addAnys(ctx, aerr.anys...); err != nil {
					return err
				}
			} else if aerr := new(anyError); errors.As(err, &aerr) {
				if err := missing.addAnys(ctx, aerr); err != nil {
					return err
				}
			} else if rerr := new(plugin.RestrictionError); errors.As(err, &rerr) {
				ok, addErr := missing.addRestrictions(ctx, rerr)
				if addErr != nil {
					return addErr
				}

				if !ok {
					return err
				}
			} else {
				// there is no need to create a full error message, if we don't have complete
				// information about what is missing
				return err
			}
		}

		switch {
		// check if we have an error at all
		case len(missing.restrictions) == 0 && len(missing.anys) == 0:
			return nil

		// check if we have collected only a single error, and return it
		// directly if so
		case len(missing.restrictions) == 1 && len(missing.anys) == 0:
			if embeddable != nil { // if it is embeddable, it will be stored here
				return embeddable
			}

			return plugin.NewRestrictionError(missing.restrictions[0])
		case len(missing.restrictions) == 0 && len(missing.anys) == 1:
			return missing.anys[0]
		}

		return missing
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
//nolint:gocognit
func Any(funcs ...plugin.RestrictionFunc) plugin.RestrictionFunc {
	return func(s *state.State, ctx *plugin.Context) error {
		if len(funcs) == 0 {
			return nil
		} else if len(funcs) == 1 {
			return funcs[0](s, ctx)
		}

		missing, err := newAnyError(ctx.Localizer)
		if err != nil {
			return err
		}

		for _, f := range funcs {
			err := f(s, ctx)
			if err == nil {
				return nil
			}

			if eerr := new(EmbeddableError); errors.As(err, &eerr) {
				ok, addErr := missing.addRestrictions(ctx, eerr.EmbeddableVersion)
				if addErr != nil {
					return addErr
				}

				if !ok {
					return err
				}
			} else if aerr := new(anyError); errors.As(err, &aerr) { // we can just merge
				missing.restrictions = append(missing.restrictions, aerr.restrictions...)

				if !aerr.fatal {
					missing.fatal = false
				}

				if err := missing.addAlls(ctx, aerr.alls...); err != nil {
					return err
				}
			} else if aerr := new(allError); errors.As(err, &aerr) {
				if err := missing.addAlls(ctx, aerr); err != nil {
					return err
				}
			} else if rerr := new(plugin.RestrictionError); errors.As(err, &rerr) {
				ok, addErr := missing.addRestrictions(ctx, rerr)
				if addErr != nil {
					return addErr
				}

				if !ok {
					return err
				}
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

// =============================================================================
// Errors
// =====================================================================================

var newlineRegexp = regexp.MustCompile(`\n[^\n]`)

// =============================================================================
// allError
// =====================================================================================

type allError struct {
	header string

	restrictions []string
	fatal        bool
	anyMessage   string
	anys         []*anyError
}

func newAllError(l *i18n.Localizer) (*allError, error) {
	header, err := l.Localize(allMessageHeader)
	if err != nil {
		return nil, err
	}

	return &allError{header: header}, nil
}

func (e *allError) addRestrictions(ctx *plugin.Context, rerrs ...*plugin.RestrictionError) (ok bool, err error) {
	for _, rerr := range rerrs {
		// there is no need to create a full error message, if we don't have complete
		// information about what is missing
		if errors.Is(rerr, plugin.DefaultRestrictionError) {
			return false, nil
		} else if errors.Is(rerr, plugin.DefaultFatalRestrictionError) {
			return false, nil
		}

		if rerr.Fatal {
			e.fatal = true
		}

		desc, err := rerr.Description(ctx.Localizer)
		if err != nil {
			return false, err
		}

		e.restrictions = append(e.restrictions, desc)
	}

	return true, nil
}

func (e *allError) addAnys(ctx *plugin.Context, aerr ...*anyError) error {
	if len(aerr) == 0 {
		return nil
	}

	e.anys = append(e.anys, aerr...)

	if e.anyMessage != "" {
		return nil
	}

	var err error
	e.anyMessage, err = ctx.Localize(anyMessageInline)
	return err
}

func (e *allError) format(indentLvl int) (s string) {
	var b strings.Builder
	b.Grow(2048) // 2048 is the max size of an embed description

	if indentLvl == 0 {
		b.WriteString(e.header)
		b.WriteString("\n") // second newline will be added in the for-loops
	}

	indent, nlIndent := genIndent(indentLvl)

	for _, r := range e.restrictions {
		if b.Len() > 0 {
			b.WriteRune('\n')
		}

		b.WriteString(indent)
		b.WriteString(entryPrefix)

		desc := newlineRegexp.ReplaceAllStringFunc(r, func(s string) string {
			return "\n" + nlIndent + s[1:]
		})
		b.WriteString(desc)
	}

	for _, a := range e.anys {
		if b.Len() > 0 {
			b.WriteRune('\n')
		}

		b.WriteString(indent)
		b.WriteString(entryPrefix)
		b.WriteString(e.anyMessage)
		b.WriteRune('\n')
		b.WriteString(a.format(indentLvl + 1))
	}

	return b.String()
}

func (e *allError) asRestrictionError() *plugin.RestrictionError {
	desc := e.format(0)

	if e.fatal {
		return plugin.NewFatalRestrictionError(desc)
	}

	return plugin.NewRestrictionError(desc)
}

func (e *allError) As(target interface{}) bool {
	switch err := target.(type) {
	case **plugin.RestrictionError:
		*err = e.asRestrictionError()
		return true
	case *errors.Error:
		*err = e.asRestrictionError()
		return true
	default:
		return false
	}
}

func (e *allError) Error() string {
	return "allError"
}

// =============================================================================
// anyError
// =====================================================================================

type anyError struct {
	header string

	restrictions []string
	fatal        bool
	allMessage   string
	alls         []*allError
}

func newAnyError(l *i18n.Localizer) (*anyError, error) {
	header, err := l.Localize(anyMessageHeader)
	if err != nil {
		return nil, err
	}

	return &anyError{header: header, fatal: true}, nil
}

func (e *anyError) addRestrictions(ctx *plugin.Context, rerrs ...*plugin.RestrictionError) (ok bool, err error) {
	for _, rerr := range rerrs {
		// there is no need to create a full error message, if we don't have complete
		// information about what is missing
		if errors.Is(rerr, plugin.DefaultRestrictionError) {
			return false, nil
		} else if errors.Is(rerr, plugin.DefaultFatalRestrictionError) {
			return false, nil
		}

		if !rerr.Fatal {
			e.fatal = false
		}

		desc, err := rerr.Description(ctx.Localizer)
		if err != nil {
			return false, err
		}

		e.restrictions = append(e.restrictions, desc)
	}

	return true, nil
}

func (e *anyError) addAlls(ctx *plugin.Context, aerr ...*allError) error {
	if len(aerr) == 0 {
		return nil
	}

	e.alls = append(e.alls, aerr...)

	if e.allMessage != "" {
		return nil
	}

	var err error
	e.allMessage, err = ctx.Localize(allMessageInline)
	return err
}

func (e *anyError) format(indentLvl int) string {
	var b strings.Builder
	b.Grow(2048) // 2048 is the max size of an embed description

	if indentLvl == 0 {
		b.WriteString(e.header)
		b.WriteString("\n") // second newline will be added in the for-loops
	}

	indent, nlIndent := genIndent(indentLvl)

	for _, r := range e.restrictions {
		if b.Len() > 0 {
			b.WriteRune('\n')
		}

		b.WriteString(indent)
		b.WriteString(entryPrefix)

		desc := newlineRegexp.ReplaceAllStringFunc(r, func(s string) string {
			return "\n" + nlIndent + s[1:]
		})
		b.WriteString(desc)
	}

	for _, a := range e.alls {
		if b.Len() > 0 {
			b.WriteRune('\n')
		}

		b.WriteString(indent)
		b.WriteString(entryPrefix)
		b.WriteString(e.allMessage)
		b.WriteRune('\n')
		b.WriteString(a.format(indentLvl + 1))
	}

	return b.String()
}

func (e *anyError) asRestrictionError() *plugin.RestrictionError {
	desc := e.format(0)

	if e.fatal {
		return plugin.NewFatalRestrictionError(desc)
	}

	return plugin.NewRestrictionError(desc)
}

func (e *anyError) As(target interface{}) bool {
	switch err := target.(type) {
	case **plugin.RestrictionError:
		*err = e.asRestrictionError()
		return true
	case *errors.Error:
		*err = e.asRestrictionError()
		return true
	default:
		return false
	}
}

func (e *anyError) Error() string {
	return "anyError"
}

// genIndent generates the necessary amount of indent for the passed indent
// level.
// The first return value is the normal indent, and the second is the indent
// needed for newlines within an entry.
func genIndent(indentLvl int) (indent, nlIndent string) {
	// use an "ideographic space" for indenting, as Discord strips normal
	// whitespace on new lines in embeds
	indent = strings.Repeat("\u3000", indentLvl*indentMultiplier)
	nlIndent = strings.Repeat("\u3000", indentLvl*indentMultiplier)

	if nlIndent == "" { // prefix a zero-width whitespace if the indentLvl is 0
		nlIndent += "\u200b"
	}

	nlIndent += strings.Repeat(" ", len(entryPrefix))

	return
}
