package restriction

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

func TestAll(t *testing.T) {
	testCases := []struct {
		name   string
		funcs  []plugin.RestrictionFunc
		expect error
	}{
		{
			name:   "no funcs",
			expect: nil,
		},
		{
			name:   "single func",
			funcs:  []plugin.RestrictionFunc{errorFunc1},
			expect: errorFunc1(nil, nil),
		},
		{
			name:   "restriction error",
			funcs:  []plugin.RestrictionFunc{errorFunc1},
			expect: errorFunc1(nil, nil),
		},
		{
			name:   "fatal restriction error",
			funcs:  []plugin.RestrictionFunc{fatalErrorFunc},
			expect: fatalErrorFunc(nil, nil),
		},
		{
			name:  "single any func",
			funcs: []plugin.RestrictionFunc{Any(errorFunc1, errorFunc2), nilFunc},
			expect: &anyError{
				restrictions: []string{errorFunc1Description, errorFunc2Description},
			},
		},
		{
			name:  "single embeddable error - not fatal",
			funcs: []plugin.RestrictionFunc{embeddableErrorFunc},
			expect: &EmbeddableError{
				EmbeddableVersion: plugin.NewRestrictionError(embeddableErrorFuncEmbeddableDescription),
				DefaultVersion:    plugin.NewRestrictionError(embeddableErrorFuncDefaultDescription),
			},
		},
		{

			name:  "multiple embeddable errors - fatal",
			funcs: []plugin.RestrictionFunc{fatalEmbeddableErrorFunc, embeddableErrorFunc},
			expect: &allError{
				restrictions: []string{
					fatalEmbeddableErrorFuncEmbeddableDescription,
					embeddableErrorFuncEmbeddableDescription,
				},
				fatal: true,
			},
		},
		{
			name:   "pass",
			funcs:  []plugin.RestrictionFunc{nilFunc, nilFunc},
			expect: nil,
		},
		{
			name:   "multiple restriction funcs - single error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, nilFunc},
			expect: errorFunc1(nil, nil),
		},
		{
			name:  "multiple restriction funcs - single embeddable error",
			funcs: []plugin.RestrictionFunc{embeddableErrorFunc, nilFunc},
			expect: &EmbeddableError{
				EmbeddableVersion: plugin.NewRestrictionError(embeddableErrorFuncEmbeddableDescription),
				DefaultVersion:    plugin.NewRestrictionError(embeddableErrorFuncDefaultDescription),
			},
		},
		{
			name:  "multiple restriction funcs - multiple errors - not fatal",
			funcs: []plugin.RestrictionFunc{errorFunc1, errorFunc2},
			expect: &allError{
				restrictions: []string{errorFunc1Description, errorFunc2Description},
			},
		},
		{
			name:  "multiple restriction funcs - multiple errors - fatal",
			funcs: []plugin.RestrictionFunc{errorFunc1, fatalErrorFunc},
			expect: &allError{
				restrictions: []string{errorFunc1Description, fatalErrorFuncDescription},
				fatal:        true,
			},
		},
		{
			name:   "multiple restriction funcs - default error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, defaultRestrictionErrorFunc},
			expect: plugin.DefaultRestrictionError,
		},
		{
			name:   "multiple restriction funcs - default fatal error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, defaultFatalRestrictionErrorFunc},
			expect: plugin.DefaultFatalRestrictionError,
		},
		{
			name:   "multiple restriction funcs - unexpected error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, unexpectedErrorFunc},
			expect: errUnexpectedErrorFuncReturn,
		},
		{
			name:   "nested all - single error",
			funcs:  []plugin.RestrictionFunc{All(errorFunc1, nilFunc)},
			expect: errorFunc1(nil, nil),
		},
		{
			name:  "nested all - multiple errors",
			funcs: []plugin.RestrictionFunc{All(errorFunc1, errorFunc2)},
			expect: &allError{
				restrictions: []string{errorFunc1Description, errorFunc2Description},
			},
		},
		{
			name:  "restriction func and nested all",
			funcs: []plugin.RestrictionFunc{errorFunc1, All(errorFunc2, errorFunc3)},
			expect: &allError{
				restrictions: []string{errorFunc1Description, errorFunc2Description, errorFunc3Description},
			},
		},
		{
			name:  "nested any",
			funcs: []plugin.RestrictionFunc{Any(errorFunc1, errorFunc2)},
			expect: &anyError{
				restrictions: []string{errorFunc1Description, errorFunc2Description},
			},
		},
		{
			name:  "multiple nested anys",
			funcs: []plugin.RestrictionFunc{Any(errorFunc1, errorFunc2), Any(errorFunc3, errorFunc4)},
			expect: &allError{
				anys: []*anyError{
					{restrictions: []string{errorFunc1Description, errorFunc2Description}},
					{restrictions: []string{errorFunc3Description, errorFunc4Description}},
				},
			},
		},
		{
			name:  "restriction func and nested any",
			funcs: []plugin.RestrictionFunc{errorFunc1, Any(errorFunc2, errorFunc3)},
			expect: &allError{
				restrictions: []string{errorFunc1Description},
				anys: []*anyError{
					{restrictions: []string{errorFunc2Description, errorFunc3Description}},
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ctx := &plugin.Context{Localizer: i18n.NewFallbackLocalizer()}

			fillHeaderAndInline(c.expect, ctx.Localizer)

			f := All(c.funcs...)

			actual := f(nil, ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestAllf(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := errors.New("abc")

		f := Allf(err, nilFunc, nilFunc)

		actual := f(nil, nil)
		assert.Nil(t, actual)
	})

	t.Run("failure", func(t *testing.T) {
		expect := errors.New("abc")

		f := Allf(expect, errorFunc1, nilFunc)

		actual := f(nil, nil)
		assert.Equal(t, expect, actual)
	})
}

func TestAny(t *testing.T) {
	testCases := []struct {
		name   string
		funcs  []plugin.RestrictionFunc
		expect error
	}{
		{
			name:   "no funcs",
			expect: nil,
		},
		{
			name:   "single func",
			funcs:  []plugin.RestrictionFunc{errorFunc1},
			expect: errorFunc1(nil, nil),
		},
		{
			name:   "restriction error",
			funcs:  []plugin.RestrictionFunc{errorFunc1},
			expect: errorFunc1(nil, nil),
		},
		{
			name:   "fatal restriction error",
			funcs:  []plugin.RestrictionFunc{fatalErrorFunc},
			expect: fatalErrorFunc(nil, nil),
		},
		{
			name:  "single all func",
			funcs: []plugin.RestrictionFunc{All(errorFunc1, errorFunc2)},
			expect: &allError{
				restrictions: []string{errorFunc1Description, errorFunc2Description},
			},
		},
		{
			name:  "single embeddable error",
			funcs: []plugin.RestrictionFunc{embeddableErrorFunc},
			expect: &EmbeddableError{
				EmbeddableVersion: plugin.NewRestrictionError(embeddableErrorFuncEmbeddableDescription),
				DefaultVersion:    plugin.NewRestrictionError(embeddableErrorFuncDefaultDescription),
			},
		},
		{

			name:  "multiple fatal embeddable errors",
			funcs: []plugin.RestrictionFunc{fatalEmbeddableErrorFunc, embeddableErrorFunc},
			expect: &anyError{
				restrictions: []string{
					fatalEmbeddableErrorFuncEmbeddableDescription,
					embeddableErrorFuncEmbeddableDescription,
				},
			},
		},
		{
			name:   "nil errors",
			funcs:  []plugin.RestrictionFunc{nilFunc, nilFunc},
			expect: nil,
		},
		{
			name:   "multiple restriction funcs - single error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, nilFunc},
			expect: nil,
		},
		{
			name:  "multiple restriction funcs - multiple errors",
			funcs: []plugin.RestrictionFunc{errorFunc1, errorFunc2},
			expect: &anyError{
				restrictions: []string{errorFunc1Description, errorFunc2Description},
			},
		},
		{
			name:  "multiple restriction funcs - multiple fatal errors",
			funcs: []plugin.RestrictionFunc{errorFunc1, fatalErrorFunc},
			expect: &anyError{
				restrictions: []string{errorFunc1Description, fatalErrorFuncDescription},
			},
		},
		{
			name:   "multiple restriction funcs - default error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, defaultRestrictionErrorFunc},
			expect: plugin.DefaultRestrictionError,
		},
		{
			name:   "multiple restriction funcs - default fatal error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, defaultFatalRestrictionErrorFunc},
			expect: plugin.DefaultFatalRestrictionError,
		},
		{
			name:   "multiple restriction funcs - unexpected error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, unexpectedErrorFunc},
			expect: errUnexpectedErrorFuncReturn,
		},
		{
			name:   "nested all - single error",
			funcs:  []plugin.RestrictionFunc{All(errorFunc1, nilFunc)},
			expect: errorFunc1(nil, nil),
		},
		{
			name:  "nested all - multiple errors",
			funcs: []plugin.RestrictionFunc{All(errorFunc1, errorFunc2)},
			expect: &allError{
				restrictions: []string{errorFunc1Description, errorFunc2Description},
			},
		},
		{
			name:  "restriction func and nested all",
			funcs: []plugin.RestrictionFunc{errorFunc1, All(errorFunc2, errorFunc3)},
			expect: &anyError{
				restrictions: []string{errorFunc1Description},
				alls: []*allError{
					{restrictions: []string{errorFunc2Description, errorFunc3Description}},
				},
			},
		},
		{
			name:  "multiple nested anys",
			funcs: []plugin.RestrictionFunc{Any(errorFunc1, errorFunc2), Any(errorFunc3, errorFunc4)},
			expect: &anyError{
				restrictions: []string{
					errorFunc1Description, errorFunc2Description, errorFunc3Description, errorFunc4Description,
				},
			},
		},
		{
			name:  "restriction func and nested all",
			funcs: []plugin.RestrictionFunc{errorFunc1, All(errorFunc2, errorFunc3, nilFunc)},
			expect: &anyError{
				restrictions: []string{errorFunc1Description},
				alls: []*allError{
					{restrictions: []string{errorFunc2Description, errorFunc3Description}},
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ctx := &plugin.Context{Localizer: i18n.NewFallbackLocalizer()}

			fillHeaderAndInline(c.expect, ctx.Localizer)

			f := Any(c.funcs...)

			actual := f(nil, ctx)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestAnyf(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := errors.New("abc")

		f := Anyf(err, nilFunc, errorFunc1)

		actual := f(nil, nil)
		assert.Nil(t, actual)
	})

	t.Run("failure", func(t *testing.T) {
		expect := errors.New("abc")

		f := Anyf(expect, errorFunc1, errorFunc2)

		actual := f(nil, nil)
		assert.Equal(t, expect, actual)
	})
}

func Test_allError_format(t *testing.T) {
	testCases := []struct {
		name   string
		err    *allError
		expect string
	}{
		{
			name: "only restrictions",
			err:  &allError{restrictions: []string{"abc", "def"}},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def",
		},
		{
			name: "only anys",
			err: &allError{
				anys: []*anyError{
					{restrictions: []string{"abc", "def"}},
					{restrictions: []string{"ghi", "jkl"}},
				},
			},
			expect: entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"　　" + entryPrefix + "abc\n" +
				"　　" + entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"　　" + entryPrefix + "ghi\n" +
				"　　" + entryPrefix + "jkl",
		},
		{
			name: "restrictions and anys",
			err: &allError{
				restrictions: []string{"abc", "def"},
				anys: []*anyError{
					{restrictions: []string{"ghi", "jkl"}},
					{restrictions: []string{"mno", "pqr"}},
				},
			},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"　　" + entryPrefix + "ghi\n" +
				"　　" + entryPrefix + "jkl\n" +
				entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"　　" + entryPrefix + "mno\n" +
				"　　" + entryPrefix + "pqr",
		},
		{
			name: "any with nested all",
			err: &allError{
				anys: []*anyError{
					{
						restrictions: []string{"abc"},
						alls: []*allError{
							{
								restrictions: []string{"def", "ghi"},
								anys:         []*anyError{{restrictions: []string{"jkl", "mno"}}},
							},
							{restrictions: []string{"pqr", "stu"}},
						},
					},
				},
			},
			expect: entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"　　" + entryPrefix + "abc\n" +
				"　　" + entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　　　" + entryPrefix + "def\n" +
				"　　　　" + entryPrefix + "ghi\n" +
				"　　　　" + entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"　　　　　　" + entryPrefix + "jkl\n" +
				"　　　　　　" + entryPrefix + "mno\n" +
				"　　" + entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　　　" + entryPrefix + "pqr\n" +
				"　　　　" + entryPrefix + "stu",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			fillHeaderAndInline(c.err, i18n.NewFallbackLocalizer())
			c.expect = c.err.header + "\n\n" + c.expect

			actual := c.err.format(0)

			assert.Equal(t, c.expect, actual)
		})
	}
}

func Test_allError_As(t *testing.T) {
	testCases := []struct {
		name string
		desc string
		fun  plugin.RestrictionFunc
	}{
		{
			name: "fatal",
			desc: "You need to fulfill all of these requirements to execute the command:\n\n" +
				entryPrefix + fatalErrorFuncDescription + "\n" +
				entryPrefix + errorFunc1Description,
			fun: All(fatalErrorFunc, errorFunc1),
		},
		{
			name: "not fatal",
			desc: "You need to fulfill all of these requirements to execute the command:\n\n" +
				entryPrefix + errorFunc1Description + "\n" +
				entryPrefix + errorFunc2Description,
			fun: All(errorFunc1, errorFunc2),
		},
	}

	ctx := &plugin.Context{
		Localizer: mock.NewLocalizer(t).
			On(allMessageHeader.Term, allMessageHeader.Fallback.Other).
			On(anyMessageInline.Term, anyMessageInline.Fallback.Other).
			On(allMessageHeader.Term, allMessageHeader.Fallback.Other).
			Build(),
	}

	t.Run("errors.Error", func(t *testing.T) {
		for _, c := range testCases {
			t.Run(c.name, func(t *testing.T) {
				err := c.fun(nil, ctx)
				require.IsType(t, new(allError), err)

				var expect *plugin.RestrictionError

				if allErr := err.(*allError); allErr.fatal {
					expect = plugin.NewFatalRestrictionError(c.desc)
				} else {
					expect = plugin.NewRestrictionError(c.desc)
				}

				var actual errors.Error
				assert.ErrorAs(t, err, &actual)
				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("*plugin.RestrictionError", func(t *testing.T) {
		for _, c := range testCases {
			t.Run(c.name, func(t *testing.T) {
				err := c.fun(nil, ctx)
				require.IsType(t, new(allError), err)

				var expect *plugin.RestrictionError

				if allErr := err.(*allError); allErr.fatal {
					expect = plugin.NewFatalRestrictionError(c.desc)
				} else {
					expect = plugin.NewRestrictionError(c.desc)
				}

				actual := new(plugin.RestrictionError)
				assert.ErrorAs(t, err, &actual)
				assert.Equal(t, expect, actual)
			})
		}
	})
}

func Test_anyError_format(t *testing.T) {
	testCases := []struct {
		name   string
		err    *anyError
		expect string
	}{
		{
			name: "only restrictions",
			err:  &anyError{restrictions: []string{"abc", "def"}},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def",
		},
		{
			name: "fatal restrictions",
			err:  &anyError{restrictions: []string{"abc", "def"}},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def",
		},
		{
			name: "fatal alls",
			err: &anyError{
				alls: []*allError{
					{restrictions: []string{"abc", "def"}},
					{restrictions: []string{"ghi", "jkl"}},
				},
			},
			expect: entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　" + entryPrefix + "abc\n" +
				"　　" + entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　" + entryPrefix + "ghi\n" +
				"　　" + entryPrefix + "jkl",
		},
		{
			name: "only alls",
			err: &anyError{
				alls: []*allError{
					{restrictions: []string{"abc", "def"}},
					{restrictions: []string{"ghi", "jkl"}},
				},
			},
			expect: entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　" + entryPrefix + "abc\n" +
				"　　" + entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　" + entryPrefix + "ghi\n" +
				"　　" + entryPrefix + "jkl",
		},
		{
			name: "restrictions and alls",
			err: &anyError{
				restrictions: []string{"abc", "def"},
				alls: []*allError{
					{restrictions: []string{"ghi", "jkl"}},
					{restrictions: []string{"mno", "pqr"}},
				},
			},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　" + entryPrefix + "ghi\n" +
				"　　" + entryPrefix + "jkl\n" +
				entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　" + entryPrefix + "mno\n" +
				"　　" + entryPrefix + "pqr",
		},
		{
			name: "any with nested all",
			err: &anyError{
				alls: []*allError{
					{
						restrictions: []string{"abc"},
						anys: []*anyError{
							{
								restrictions: []string{"def", "ghi"},
								alls:         []*allError{{restrictions: []string{"jkl", "mno"}}},
							},
							{restrictions: []string{"pqr", "stu"}},
						},
					},
				},
			},
			expect: entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　" + entryPrefix + "abc\n" +
				"　　" + entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"　　　　" + entryPrefix + "def\n" +
				"　　　　" + entryPrefix + "ghi\n" +
				"　　　　" + entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　　　　　" + entryPrefix + "jkl\n" +
				"　　　　　　" + entryPrefix + "mno\n" +
				"　　" + entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"　　　　" + entryPrefix + "pqr\n" +
				"　　　　" + entryPrefix + "stu",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			fillHeaderAndInline(c.err, i18n.NewFallbackLocalizer())
			c.expect = c.err.header + "\n\n" + c.expect

			actual := c.err.format(0)

			assert.Equal(t, c.expect, actual)
		})
	}
}

func Test_anyError_As(t *testing.T) {
	testCases := []struct {
		name string
		desc string
		fun  plugin.RestrictionFunc
	}{
		{
			name: "fatal",
			desc: "You need to fulfill at least one of these requirements to execute the command:\n\n" +
				entryPrefix + fatalErrorFuncDescription + "\n" +
				entryPrefix + fatalErrorFuncDescription,
			fun: Any(fatalErrorFunc, fatalErrorFunc),
		},
		{
			name: "not fatal",
			desc: "You need to fulfill at least one of these requirements to execute the command:\n\n" +
				entryPrefix + errorFunc1Description + "\n" +
				entryPrefix + errorFunc2Description,
			fun: Any(errorFunc1, errorFunc2),
		},
	}

	ctx := &plugin.Context{
		Localizer: mock.NewLocalizer(t).
			On(allMessageHeader.Term, allMessageHeader.Fallback.Other).
			On(anyMessageInline.Term, anyMessageInline.Fallback.Other).
			On(anyMessageHeader.Term, anyMessageHeader.Fallback.Other).
			Build(),
	}

	t.Run("errors.Error", func(t *testing.T) {
		for _, c := range testCases {
			t.Run(c.name, func(t *testing.T) {
				err := c.fun(nil, ctx)
				require.IsType(t, new(anyError), err)

				var expect *plugin.RestrictionError

				if allErr := err.(*anyError); allErr.fatal {
					expect = plugin.NewFatalRestrictionError(c.desc)
				} else {
					expect = plugin.NewRestrictionError(c.desc)
				}

				var actual errors.Error
				assert.ErrorAs(t, err, &actual)
				assert.Equal(t, expect, actual)
			})
		}
	})

	t.Run("*plugin.RestrictionError", func(t *testing.T) {
		for _, c := range testCases {
			t.Run(c.name, func(t *testing.T) {
				err := c.fun(nil, ctx)
				require.IsType(t, new(anyError), err)

				var expect *plugin.RestrictionError

				if allErr := err.(*anyError); allErr.fatal {
					expect = plugin.NewFatalRestrictionError(c.desc)
				} else {
					expect = plugin.NewRestrictionError(c.desc)
				}

				actual := new(plugin.RestrictionError)
				assert.ErrorAs(t, err, &actual)
				assert.Equal(t, expect, actual)
			})
		}
	})
}

func fillHeaderAndInline(err error, l *i18n.Localizer) {
	if allErr, ok := err.(*allError); ok {
		allErr.header = l.MustLocalize(allMessageHeader)

		if len(allErr.anys) > 0 {
			allErr.anyMessage = l.MustLocalize(anyMessageInline)

			for _, any := range allErr.anys {
				fillHeaderAndInline(any, l)
			}
		}
	} else if anyErr, ok := err.(*anyError); ok {
		anyErr.header = l.MustLocalize(anyMessageHeader)

		if len(anyErr.alls) > 0 {
			anyErr.allMessage = l.MustLocalize(allMessageInline)

			for _, all := range anyErr.alls {
				fillHeaderAndInline(all, l)
			}
		}
	}
}
