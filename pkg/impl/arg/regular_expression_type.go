package arg

import (
	"errors"
	"regexp"
	resyntax "regexp/syntax"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// RegularExpression is the Type used for regular expressions.
//
// Go type: *regexp.Regexp
var RegularExpression plugin.ArgType = new(regularExpression)

type regularExpression struct{}

func (r regularExpression) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(regexpName) // we have a fallback
	return name
}

func (r regularExpression) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(regexpDescription) // we have a fallback
	return desc
}

func (r regularExpression) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	compiled, err := regexp.Compile(ctx.Raw)
	if err == nil {
		return compiled, nil
	}

	var regerr *resyntax.Error
	if !errors.As(err, &regerr) {
		return nil, newArgumentError2(regexpInvalidErrorArg, regexpInvalidErrorFlag, ctx, nil)
	}

	placeholders := map[string]interface{}{
		"expression": regerr.Expr,
	}

	switch regerr.Code {
	case resyntax.ErrInvalidCharClass:
		return nil,
			newArgumentError2(regexpInvalidCharClassErrorArg, regexpInvalidCharClassErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidCharRange:
		return nil,
			newArgumentError2(regexpInvalidCharRangeErrorArg, regexpInvalidCharRangeErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidEscape:
		return nil, newArgumentError2(regexpInvalidEscapeErrorArg, regexpInvalidEscapeErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidNamedCapture:
		return nil,
			newArgumentError2(regexpInvalidNamedCaptureErrorArg, regexpInvalidNamedCaptureErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidPerlOp:
		return nil, newArgumentError2(regexpInvalidPerlOpErrorArg, regexpInvalidPerlOpErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidRepeatOp:
		return nil, newArgumentError2(regexpInvalidRepeatOpErrorArg, regexpInvalidRepeatOpErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidRepeatSize:
		return nil,
			newArgumentError2(regexpInvalidRepeatSizeErrorArg, regexpInvalidRepeatSizeErrorFlag, ctx, placeholders)
	case resyntax.ErrInvalidUTF8:
		return nil, newArgumentError2(regexpInvalidUTF8ErrorArg, regexpInvalidUTF8ErrorFlag, ctx, placeholders)
	case resyntax.ErrMissingBracket:
		return nil, newArgumentError2(regexpMissingBracketErrorArg, regexpMissingBracketErrorFlag, ctx, placeholders)
	case resyntax.ErrMissingParen:
		return nil, newArgumentError2(regexpMissingParenErrorArg, regexpMissingParenErrorFlag, ctx, placeholders)
	case resyntax.ErrMissingRepeatArgument:
		return nil,
			newArgumentError2(regexpMissingRepeatArgErrorArg, regexpMissingRepeatArgErrorFlag, ctx, placeholders)
	case resyntax.ErrTrailingBackslash:
		return nil,
			newArgumentError2(regexpTrailingBackslashErrorArg, regexpTrailingBackslashErrorFlag, ctx, placeholders)
	case resyntax.ErrUnexpectedParen:
		return nil, newArgumentError2(regexpUnexpectedParenErrorArg, regexpUnexpectedParenErrorFlag, ctx, placeholders)
	case resyntax.ErrInternalError:
		fallthrough
	default:
		return nil, newArgumentError2(regexpInvalidErrorArg, regexpInvalidErrorFlag, ctx, placeholders)
	}
}

func (r regularExpression) GetDefault() interface{} {
	return (*regexp.Regexp)(nil)
}
