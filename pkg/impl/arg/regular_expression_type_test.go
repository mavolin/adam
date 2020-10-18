package arg

import (
	"regexp"
	resyntax "regexp/syntax"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
)

func TestRegularExpression_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := regexp.MustCompile("abc")

		ctx := &Context{Raw: expect.String()}

		actual, err := RegularExpression.Parse(nil, ctx)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	failureCases := []struct {
		name resyntax.ErrorCode

		raw string

		expectArg, expectFlag *i18n.Config
		expression            string
	}{
		// resyntax.ErrInvalidCharClass - this error seems to be never returned
		{
			name:       resyntax.ErrInvalidCharRange,
			raw:        `[\x{ffff}-\x{aaaa}]`,
			expectArg:  regexpInvalidCharRangeErrorArg,
			expectFlag: regexpInvalidCharRangeErrorFlag,
			expression: `\x{ffff}-\x{aaaa}`,
		},
		{
			name:       resyntax.ErrInvalidEscape,
			raw:        `\x`,
			expectArg:  regexpInvalidEscapeErrorArg,
			expectFlag: regexpInvalidEscapeErrorFlag,
			expression: `\x`,
		},
		{
			name:       resyntax.ErrInvalidNamedCapture,
			raw:        `(?P<\>abc)`,
			expectArg:  regexpInvalidNamedCaptureErrorArg,
			expectFlag: regexpInvalidNamedCaptureErrorFlag,
			expression: `(?P<\>`,
		},
		{
			name:       resyntax.ErrInvalidPerlOp,
			raw:        `(?<=abc)`,
			expectArg:  regexpInvalidPerlOpErrorArg,
			expectFlag: regexpInvalidPerlOpErrorFlag,
			expression: `(?<`,
		},
		{
			name:       resyntax.ErrInvalidRepeatOp,
			raw:        `a++`,
			expectArg:  regexpInvalidRepeatOpErrorArg,
			expectFlag: regexpInvalidRepeatOpErrorFlag,
			expression: `++`,
		},
		{
			name:       resyntax.ErrInvalidRepeatSize,
			raw:        `a{4,3}`,
			expectArg:  regexpInvalidRepeatSizeErrorArg,
			expectFlag: regexpInvalidRepeatSizeErrorFlag,
			expression: `{4,3}`,
		},
		// resyntax.ErrInvalidUTF8 - no clue how to produce that
		{
			name:       resyntax.ErrMissingBracket,
			raw:        `[abc`,
			expectArg:  regexpMissingBracketErrorArg,
			expectFlag: regexpMissingBracketErrorFlag,
			expression: `[abc`,
		},
		{
			name:       resyntax.ErrMissingParen,
			raw:        `(a|b`,
			expectArg:  regexpMissingParenErrorArg,
			expectFlag: regexpMissingParenErrorFlag,
			expression: `(a|b`,
		},
		{
			name:       resyntax.ErrMissingRepeatArgument,
			raw:        `+`,
			expectArg:  regexpMissingRepeatArgErrorArg,
			expectFlag: regexpMissingRepeatArgErrorFlag,
			expression: `+`,
		},
		{
			name:       resyntax.ErrTrailingBackslash,
			raw:        `\`,
			expectArg:  regexpTrailingBackslashErrorArg,
			expectFlag: regexpTrailingBackslashErrorFlag,
			expression: "",
		},
		{
			name:       resyntax.ErrUnexpectedParen,
			raw:        `)`,
			expectArg:  regexpUnexpectedParenErrorArg,
			expectFlag: regexpUnexpectedParenErrorFlag,
			expression: `)`,
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(string(c.name), func(t *testing.T) {
				c.expectArg = c.expectArg.
					WithPlaceholders(map[string]interface{}{
						"expression": c.expression,
					})

				c.expectFlag = c.expectFlag.
					WithPlaceholders(map[string]interface{}{
						"expression": c.expression,
					})

				ctx := &Context{
					Raw:  c.raw,
					Kind: KindArg,
				}

				c.expectArg.Placeholders = attachDefaultPlaceholders(c.expectArg.Placeholders, ctx)
				c.expectFlag.Placeholders = attachDefaultPlaceholders(c.expectFlag.Placeholders, ctx)

				_, actual := RegularExpression.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectArg), actual)

				ctx.Kind = KindFlag

				_, actual = RegularExpression.Parse(nil, ctx)
				assert.Equal(t, errors.NewArgumentParsingErrorl(c.expectFlag), actual)
			})
		}
	})
}
