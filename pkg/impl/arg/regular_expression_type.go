package arg

import (
	"regexp"
	resyntax "regexp/syntax"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

// RegularExpression is the Type used for regular expressions.
//
// Go type: *regexp.Regexp
var RegularExpression = new(regularExpression)

type regularExpression struct{}

func (r regularExpression) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(regexpName) // we have a fallback
	return name
}

func (r regularExpression) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(regexpDescription) // we have a fallback
	return desc
}

func (r regularExpression) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	compiled, err := regexp.Compile(ctx.Raw)
	if err == nil {
		return compiled, nil
	}

	serr, ok := err.(*resyntax.Error)
	if !ok {
		return nil, newArgParsingErr(regexpInvalidErrorArg, regexpInvalidErrorFlag, ctx, nil)
	}

	placeholders := map[string]interface{}{
		"expression": serr.Expr,
	}

	switch serr.Code {
	case resyntax.ErrInvalidCharClass:
		return nil,
			newArgParsingErr(regexpInvalidCharClassErrorArg, regexpInvalidCharClassErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidCharRange:
		return nil,
			newArgParsingErr(regexpInvalidCharRangeErrorArg, regexpInvalidCharRangeErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidEscape:
		return nil, newArgParsingErr(regexpInvalidEscapeErrorArg, regexpInvalidEscapeErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidNamedCapture:
		return nil,
			newArgParsingErr(regexpInvalidNamedCaptureErrorArg, regexpInvalidNamedCaptureErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidPerlOp:
		return nil, newArgParsingErr(regexpInvalidPerlOpErrorArg, regexpInvalidPerlOpErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidRepeatOp:
		return nil, newArgParsingErr(regexpInvalidRepeatOpErrorArg, regexpInvalidRepeatOpErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidRepeatSize:
		return nil,
			newArgParsingErr(regexpInvalidRepeatSizeErrorArg, regexpInvalidRepeatSizeErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidUTF8:
		return nil, newArgParsingErr(regexpInvalidUTF8ErrorArg, regexpInvalidUTF8ErrorFlag, ctx, placeholders)
	case resyntax.ErrMissingBracket:
		return nil, newArgParsingErr(regexpMissingBracketErrorArg, regexpMissingBracketErrorFlag, ctx, placeholders)
	case resyntax.ErrMissingParen:
		return nil, newArgParsingErr(regexpMissingParenErrorArg, regexpMissingParenErrorFlag, ctx, placeholders)
	case resyntax.ErrMissingRepeatArgument:
		return nil,
			newArgParsingErr(regexpMissingRepeatArgErrorArg, regexpMissingRepeatArgErrorFlag, ctx, placeholders)
	case resyntax.ErrTrailingBackslash:
		return nil,
			newArgParsingErr(regexpTrailingBackslashErrorArg, regexpTrailingBackslashErrorFlag, ctx, placeholders)
	case resyntax.ErrUnexpectedParen:
		return nil, newArgParsingErr(regexpUnexpectedParenErrorArg, regexpUnexpectedParenErrorFlag, ctx, placeholders)
	default:
		return nil, newArgParsingErr(regexpInvalidErrorArg, regexpInvalidErrorFlag, ctx, placeholders)
	}
}

func (r regularExpression) Default() interface{} {
	return (*regexp.Regexp)(nil)
}
