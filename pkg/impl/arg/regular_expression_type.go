package arg

import (
	"errors"
	"regexp"
	resyntax "regexp/syntax"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

// RegularExpression is the Type used for regular expressions.
//
// Go type: *regexp.Regexp
var RegularExpression Type = new(regularExpression)

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

	var regerr *resyntax.Error
	if !errors.As(err, &regerr) {
		return nil, newArgParsingErr2(regexpInvalidErrorArg, regexpInvalidErrorFlag, ctx, nil)
	}

	placeholders := map[string]interface{}{
		"expression": regerr.Expr,
	}

	switch regerr.Code {
	case resyntax.ErrInvalidCharClass:
		return nil,
			newArgParsingErr2(regexpInvalidCharClassErrorArg, regexpInvalidCharClassErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidCharRange:
		return nil,
			newArgParsingErr2(regexpInvalidCharRangeErrorArg, regexpInvalidCharRangeErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidEscape:
		return nil, newArgParsingErr2(regexpInvalidEscapeErrorArg, regexpInvalidEscapeErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidNamedCapture:
		return nil,
			newArgParsingErr2(regexpInvalidNamedCaptureErrorArg, regexpInvalidNamedCaptureErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidPerlOp:
		return nil, newArgParsingErr2(regexpInvalidPerlOpErrorArg, regexpInvalidPerlOpErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidRepeatOp:
		return nil, newArgParsingErr2(regexpInvalidRepeatOpErrorArg, regexpInvalidRepeatOpErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidRepeatSize:
		return nil,
			newArgParsingErr2(regexpInvalidRepeatSizeErrorArg, regexpInvalidRepeatSizeErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidUTF8:
		return nil, newArgParsingErr2(regexpInvalidUTF8ErrorArg, regexpInvalidUTF8ErrorFlag, ctx, placeholders)
	case resyntax.ErrMissingBracket:
		return nil, newArgParsingErr2(regexpMissingBracketErrorArg, regexpMissingBracketErrorFlag, ctx, placeholders)
	case resyntax.ErrMissingParen:
		return nil, newArgParsingErr2(regexpMissingParenErrorArg, regexpMissingParenErrorFlag, ctx, placeholders)
	case resyntax.ErrMissingRepeatArgument:
		return nil,
			newArgParsingErr2(regexpMissingRepeatArgErrorArg, regexpMissingRepeatArgErrorFlag, ctx, placeholders)
	case resyntax.ErrTrailingBackslash:
		return nil,
			newArgParsingErr2(regexpTrailingBackslashErrorArg, regexpTrailingBackslashErrorFlag, ctx, placeholders)
	case resyntax.ErrUnexpectedParen:
		return nil, newArgParsingErr2(regexpUnexpectedParenErrorArg, regexpUnexpectedParenErrorFlag, ctx, placeholders)
	case resyntax.ErrInternalError:
		fallthrough
	default:
		return nil, newArgParsingErr2(regexpInvalidErrorArg, regexpInvalidErrorFlag, ctx, placeholders)
	}
}

func (r regularExpression) Default() interface{} {
	return (*regexp.Regexp)(nil)
}
