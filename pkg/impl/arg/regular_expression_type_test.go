package arg

import (
	"regexp"
	resyntax "regexp/syntax"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestRegularExpression_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := regexp.MustCompile("abc")

		ctx := &plugin.ParseContext{Raw: expect.String()}

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
				ctx := &plugin.ParseContext{
					Raw:  c.raw,
					Kind: plugin.KindArg,
				}

				placeholders := map[string]interface{}{"expression": c.expression}

				expect := newArgumentError(c.expectArg, ctx, placeholders)

				_, actual := RegularExpression.Parse(nil, ctx)
				assert.Equal(t, expect, actual)

				ctx.Kind = plugin.KindFlag
				expect = newArgumentError(c.expectFlag, ctx, placeholders)

				_, actual = RegularExpression.Parse(nil, ctx)
				assert.Equal(t, expect, actual)
			})
		}
	})
}
