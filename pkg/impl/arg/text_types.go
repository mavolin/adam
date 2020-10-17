package arg

import (
	"regexp"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

// Text is the Type for a string.
type Text struct {
	// MinLength is the inclusive minimum length the text may have.
	MinLength uint
	// MaxLength is the inclusive maximum length the text may have.
	// If MaxLength is 0, the text won't have a maximum.
	MaxLength uint

	// Regexp is the regular expression the text must match.
	// If Regexp is set to nil/zero, the text won't be matched.
	//
	// If matching fails, RegexpError will be returned.
	Regexp *regexp.Regexp
	// RegexpErrorArg is the error message used if an argument doesn't match
	// the regular expression defined.
	//
	// Available Placeholders are:
	//
	// 		• name - the name of the argument
	// 		• raw - the raw argument
	// 		• position - the position of the text (1-indexed)
	// 		• regexp - the regular expression that needs to be matched
	//
	// Defaults to: regexpNotMatchingErrorArg
	RegexpErrorArg i18n.Config
	// RegexpErrorFlag is the error message used if a flag doesn't match the
	// regular expression defined.
	//
	// Available Placeholders are:
	//
	// 		• name - the full name of the flag
	// 		• used_name - the name of the flag the invoking user used
	// 		• raw - the raw flag without the flags name
	// 		• regexp - the regular expression that needs to be matched
	//
	// Defaults to: regexpNotMatchingErrorFlag
	RegexpErrorFlag i18n.Config
}

// SimpleText is a Text with no length boundaries and no regular expression.
var SimpleText = Text{}

func (t Text) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(textName) // we have a fallback
	return name
}

func (t Text) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(textDescription) // we have a fallback
	return desc
}

func (t Text) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	if uint(len(ctx.Raw)) < t.MinLength {
		return nil, newArgParsingErr(
			textBelowMinLengthErrorArg, textBelowMinLengthErrorFlag, ctx, map[string]interface{}{
				"min": t.MinLength,
			})
	} else if t.MaxLength > 0 && uint(len(ctx.Raw)) > t.MaxLength {
		return nil, newArgParsingErr(
			textAboveMaxLengthErrorArg, textAboveMaxLengthErrorFlag, ctx, map[string]interface{}{
				"max": t.MaxLength,
			})
	}

	if t.Regexp != nil && !t.Regexp.MatchString(ctx.Raw) {
		if ctx.Kind == KindArg && !t.RegexpErrorArg.IsValid() {
			t.RegexpErrorArg = regexpNotMatchingErrorArg
		} else if ctx.Kind == KindFlag && !t.RegexpErrorFlag.IsValid() {
			t.RegexpErrorFlag = regexpNotMatchingErrorFlag
		}

		return nil, newArgParsingErr(t.RegexpErrorArg, t.RegexpErrorFlag, ctx, map[string]interface{}{
			"regexp": t.Regexp.String(),
		})
	}

	return ctx.Raw, nil
}

func (t Text) Default() interface{} {
	return ""
}
