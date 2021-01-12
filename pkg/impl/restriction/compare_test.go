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
			expect: errorFuncReturn1,
		},
		{
			name:   "restriction error",
			funcs:  []plugin.RestrictionFunc{errorFunc1},
			expect: errorFuncReturn1,
		},
		{
			name:   "fatal restriction error",
			funcs:  []plugin.RestrictionFunc{fatalErrorFunc},
			expect: fatalErrorFuncReturn,
		},
		{
			name:   "single any func",
			funcs:  []plugin.RestrictionFunc{Any(errorFunc1, errorFunc1), nilFunc},
			expect: Any(errorFunc1, errorFunc1)(nil, nil),
		},
		{
			name:   "single embeddable error",
			funcs:  []plugin.RestrictionFunc{embeddableErrorFunc},
			expect: embeddableErrorFuncReturn,
		},
		{

			name:  "multiple fatal embeddable errors",
			funcs: []plugin.RestrictionFunc{fatalEmbeddableErrorFunc, embeddableErrorFunc},
			expect: &allError{
				restrictions: []*plugin.RestrictionError{
					fatalEmbeddableErrorFuncReturn.EmbeddableVersion,
					embeddableErrorFuncReturn.EmbeddableVersion,
				},
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
			expect: errorFuncReturn1,
		},
		{
			name:   "multiple restriction funcs - single embeddable error",
			funcs:  []plugin.RestrictionFunc{embeddableErrorFunc, nilFunc},
			expect: embeddableErrorFuncReturn,
		},
		{
			name:  "multiple restriction funcs - multiple errors",
			funcs: []plugin.RestrictionFunc{errorFunc1, errorFunc2},
			expect: &allError{
				restrictions: []*plugin.RestrictionError{errorFuncReturn1, errorFuncReturn2},
			},
		},
		{
			name:  "multiple restriction funcs - multiple fatal errors",
			funcs: []plugin.RestrictionFunc{errorFunc1, fatalErrorFunc},
			expect: &allError{
				restrictions: []*plugin.RestrictionError{errorFuncReturn1, fatalErrorFuncReturn},
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
			expect: unexpectedErrorFuncReturn,
		},
		{
			name:   "nested all - single error",
			funcs:  []plugin.RestrictionFunc{All(errorFunc1, nilFunc)},
			expect: errorFuncReturn1,
		},
		{
			name:  "nested all - multiple errors",
			funcs: []plugin.RestrictionFunc{All(errorFunc1, errorFunc2)},
			expect: &allError{
				restrictions: []*plugin.RestrictionError{errorFuncReturn1, errorFuncReturn2},
			},
		},
		{
			name:  "restriction func and nested all",
			funcs: []plugin.RestrictionFunc{errorFunc1, All(errorFunc2, errorFunc3)},
			expect: &allError{
				restrictions: []*plugin.RestrictionError{errorFuncReturn1, errorFuncReturn2, errorFuncReturn3},
			},
		},
		{
			name:  "nested any",
			funcs: []plugin.RestrictionFunc{Any(errorFunc1, errorFunc2)},
			expect: &anyError{
				restrictions: []*plugin.RestrictionError{errorFuncReturn1, errorFuncReturn2},
			},
		},
		{
			name:  "multiple nested anys",
			funcs: []plugin.RestrictionFunc{Any(errorFunc1, errorFunc2), Any(errorFunc3, errorFunc4)},
			expect: &allError{
				anys: []*anyError{
					{restrictions: []*plugin.RestrictionError{errorFuncReturn1, errorFuncReturn2}},
					{restrictions: []*plugin.RestrictionError{errorFuncReturn3, errorFuncReturn4}},
				},
			},
		},
		{
			name:  "restriction func and nested any",
			funcs: []plugin.RestrictionFunc{errorFunc1, Any(errorFunc2, errorFunc3)},
			expect: &allError{
				restrictions: []*plugin.RestrictionError{errorFuncReturn1},
				anys: []*anyError{
					{restrictions: []*plugin.RestrictionError{errorFuncReturn2, errorFuncReturn3}},
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f := All(c.funcs...)

			actual := f(nil, &plugin.Context{Localizer: i18n.FallbackLocalizer})
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
			expect: errorFuncReturn1,
		},
		{
			name:   "restriction error",
			funcs:  []plugin.RestrictionFunc{errorFunc1},
			expect: errorFuncReturn1,
		},
		{
			name:   "fatal restriction error",
			funcs:  []plugin.RestrictionFunc{fatalErrorFunc},
			expect: fatalErrorFuncReturn,
		},
		{
			name:   "single all func",
			funcs:  []plugin.RestrictionFunc{All(errorFunc1, errorFunc1)},
			expect: All(errorFunc1, errorFunc1)(nil, nil),
		},
		{
			name:   "single embeddable error",
			funcs:  []plugin.RestrictionFunc{embeddableErrorFunc},
			expect: embeddableErrorFuncReturn,
		},
		{

			name:  "multiple fatal embeddable errors",
			funcs: []plugin.RestrictionFunc{fatalEmbeddableErrorFunc, embeddableErrorFunc},
			expect: &anyError{
				restrictions: []*plugin.RestrictionError{
					fatalEmbeddableErrorFuncReturn.EmbeddableVersion,
					embeddableErrorFuncReturn.EmbeddableVersion,
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
				restrictions: []*plugin.RestrictionError{errorFuncReturn1, errorFuncReturn2},
			},
		},
		{
			name:  "multiple restriction funcs - multiple fatal errors",
			funcs: []plugin.RestrictionFunc{errorFunc1, fatalErrorFunc},
			expect: &anyError{
				restrictions: []*plugin.RestrictionError{errorFuncReturn1, fatalErrorFuncReturn},
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
			expect: unexpectedErrorFuncReturn,
		},
		{
			name:   "nested all - single error",
			funcs:  []plugin.RestrictionFunc{All(errorFunc1, nilFunc)},
			expect: errorFuncReturn1,
		},
		{
			name:  "nested all - multiple errors",
			funcs: []plugin.RestrictionFunc{All(errorFunc1, errorFunc2)},
			expect: &allError{
				restrictions: []*plugin.RestrictionError{errorFuncReturn1, errorFuncReturn2},
			},
		},
		{
			name:  "restriction func and nested all",
			funcs: []plugin.RestrictionFunc{errorFunc1, All(errorFunc2, errorFunc3)},
			expect: &anyError{
				restrictions: []*plugin.RestrictionError{errorFuncReturn1},
				alls: []*allError{
					{restrictions: []*plugin.RestrictionError{errorFuncReturn2, errorFuncReturn3}},
				},
			},
		},
		{
			name:  "multiple nested anys",
			funcs: []plugin.RestrictionFunc{Any(errorFunc1, errorFunc2), Any(errorFunc3, errorFunc4)},
			expect: &anyError{
				restrictions: []*plugin.RestrictionError{
					errorFuncReturn1, errorFuncReturn2, errorFuncReturn3, errorFuncReturn4,
				},
			},
		},
		{
			name:  "restriction func and nested any",
			funcs: []plugin.RestrictionFunc{errorFunc1, All(errorFunc2, errorFunc3, nilFunc)},
			expect: &anyError{
				restrictions: []*plugin.RestrictionError{errorFuncReturn1},
				alls: []*allError{
					{restrictions: []*plugin.RestrictionError{errorFuncReturn2, errorFuncReturn3}},
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f := Any(c.funcs...)

			actual := f(nil, &plugin.Context{Localizer: i18n.FallbackLocalizer})
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
		err    allError
		expect string
		fatal  bool
	}{
		{
			name: "only restrictions",
			err: allError{
				restrictions: []*plugin.RestrictionError{
					plugin.NewRestrictionError("abc"),
					plugin.NewRestrictionError("def"),
				},
			},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def",
			fatal: false,
		},
		{
			name: "only anys",
			err: allError{
				anys: []*anyError{
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("abc"),
							plugin.NewRestrictionError("def"),
						},
					},
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("ghi"),
							plugin.NewRestrictionError("jkl"),
						},
					},
				},
			},
			expect: entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"　　" + entryPrefix + "abc\n" +
				"　　" + entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"　　" + entryPrefix + "ghi\n" +
				"　　" + entryPrefix + "jkl",
			fatal: false,
		},
		{
			name: "restrictions and anys",
			err: allError{
				restrictions: []*plugin.RestrictionError{
					plugin.NewRestrictionError("abc"),
					plugin.NewRestrictionError("def"),
				},
				anys: []*anyError{
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("ghi"),
							plugin.NewRestrictionError("jkl"),
						},
					},
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("mno"),
							plugin.NewRestrictionError("pqr"),
						},
					},
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
			fatal: false,
		},
		{
			name: "any with fatals",
			err: allError{
				restrictions: []*plugin.RestrictionError{
					plugin.NewRestrictionError("abc"),
					plugin.NewRestrictionError("def"),
				},
				anys: []*anyError{
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewFatalRestrictionError("ghi"),
							plugin.NewFatalRestrictionError("jkl"),
						},
					},
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("mno"),
							plugin.NewRestrictionError("pqr"),
						},
					},
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
			fatal: true,
		},
		{
			name: "any with nested all",
			err: allError{
				anys: []*anyError{
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("abc"),
						},
						alls: []*allError{
							{
								restrictions: []*plugin.RestrictionError{
									plugin.NewRestrictionError("def"),
									plugin.NewRestrictionError("ghi"),
								},
								anys: []*anyError{
									{
										restrictions: []*plugin.RestrictionError{
											plugin.NewRestrictionError("jkl"),
											plugin.NewRestrictionError("mno"),
										},
									},
								},
							},
							{
								restrictions: []*plugin.RestrictionError{
									plugin.NewRestrictionError("pqr"),
									plugin.NewRestrictionError("stu"),
								},
							},
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
			fatal: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual, fatal, err := c.err.format(0, i18n.FallbackLocalizer)
			require.NoError(t, err)
			assert.Equal(t, c.fatal, fatal)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func Test_allError_Wrap(t *testing.T) {
	t.Run("fatal", func(t *testing.T) {
		expect := plugin.NewFatalRestrictionError("You need to fulfill all of these requirements to execute the " +
			"command:\n\n" +
			entryPrefix + "mno\n" +
			entryPrefix + "mno")

		ctx := &plugin.Context{
			Localizer: mock.NewLocalizer(t).
				On(allMessageHeader.Term, allMessageHeader.Fallback.Other).
				On(anyMessageInline.Term, anyMessageInline.Fallback.Other).
				Build(),
		}

		err := All(fatalErrorFunc, fatalErrorFunc)(nil, ctx)

		require.IsType(t, new(allError), err)

		allErr := err.(*allError)

		actual := allErr.Wrap(nil, ctx)
		assert.Equal(t, expect, actual)
	})

	t.Run("not fatal", func(t *testing.T) {
		expect := plugin.NewRestrictionError("You need to fulfill all of these requirements to execute the command:\n\n" +
			entryPrefix + "abc\n" +
			entryPrefix + "def")

		ctx := &plugin.Context{
			Localizer: mock.NewLocalizer(t).
				On(allMessageHeader.Term, allMessageHeader.Fallback.Other).
				On(anyMessageInline.Term, anyMessageInline.Fallback.Other).
				Build(),
		}

		err := All(errorFunc1, errorFunc2)(nil, ctx)

		require.IsType(t, new(allError), err)

		allErr := err.(*allError)

		actual := allErr.Wrap(nil, ctx)
		assert.Equal(t, expect, actual)
	})
}

func Test_anyError_format(t *testing.T) {
	testCases := []struct {
		name   string
		err    anyError
		expect string
		fatal  bool
	}{
		{
			name: "only restrictions",
			err: anyError{
				restrictions: []*plugin.RestrictionError{
					plugin.NewRestrictionError("abc"),
					plugin.NewRestrictionError("def"),
				},
			},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def",
			fatal: false,
		},
		{
			name: "fatal restrictions",
			err: anyError{
				restrictions: []*plugin.RestrictionError{
					plugin.NewFatalRestrictionError("abc"),
					plugin.NewFatalRestrictionError("def"),
				},
			},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def",
			fatal: true,
		},
		{
			name: "fatal alls",
			err: anyError{
				alls: []*allError{
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewFatalRestrictionError("abc"),
							plugin.NewFatalRestrictionError("def"),
						},
					},
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewFatalRestrictionError("ghi"),
							plugin.NewFatalRestrictionError("jkl"),
						},
					},
				},
			},
			expect: entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　" + entryPrefix + "abc\n" +
				"　　" + entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　" + entryPrefix + "ghi\n" +
				"　　" + entryPrefix + "jkl",
			fatal: true,
		},
		{
			name: "only alls",
			err: anyError{
				alls: []*allError{
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("abc"),
							plugin.NewRestrictionError("def"),
						},
					},
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("ghi"),
							plugin.NewRestrictionError("jkl"),
						},
					},
				},
			},
			expect: entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　" + entryPrefix + "abc\n" +
				"　　" + entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill all of these requirements:\n" +
				"　　" + entryPrefix + "ghi\n" +
				"　　" + entryPrefix + "jkl",
			fatal: false,
		},
		{
			name: "restrictions and alls",
			err: anyError{
				restrictions: []*plugin.RestrictionError{
					plugin.NewRestrictionError("abc"),
					plugin.NewRestrictionError("def"),
				},
				alls: []*allError{
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("ghi"),
							plugin.NewRestrictionError("jkl"),
						},
					},
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("mno"),
							plugin.NewRestrictionError("pqr"),
						},
					},
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
			fatal: false,
		},
		{
			name: "all with fatals",
			err: anyError{
				restrictions: []*plugin.RestrictionError{
					plugin.NewRestrictionError("abc"),
					plugin.NewRestrictionError("def"),
				},
				alls: []*allError{
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewFatalRestrictionError("ghi"),
							plugin.NewFatalRestrictionError("jkl"),
						},
					},
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("mno"),
							plugin.NewRestrictionError("pqr"),
						},
					},
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
			fatal: false,
		},
		{
			name: "any with nested all",
			err: anyError{
				alls: []*allError{
					{
						restrictions: []*plugin.RestrictionError{
							plugin.NewRestrictionError("abc"),
						},
						anys: []*anyError{
							{
								restrictions: []*plugin.RestrictionError{
									plugin.NewRestrictionError("def"),
									plugin.NewRestrictionError("ghi"),
								},
								alls: []*allError{
									{
										restrictions: []*plugin.RestrictionError{
											plugin.NewRestrictionError("jkl"),
											plugin.NewRestrictionError("mno"),
										},
									},
								},
							},
							{
								restrictions: []*plugin.RestrictionError{
									plugin.NewRestrictionError("pqr"),
									plugin.NewRestrictionError("stu"),
								},
							},
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
			fatal: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual, fatal, err := c.err.format(0, i18n.FallbackLocalizer)
			require.NoError(t, err)
			assert.Equal(t, c.fatal, fatal)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func Test_anyError_Wrap(t *testing.T) {
	t.Run("fatal", func(t *testing.T) {
		expect := plugin.NewFatalRestrictionError("You need to fulfill at least one of these requirements to " +
			"execute the command:\n\n" +
			entryPrefix + "mno\n" +
			entryPrefix + "mno")

		ctx := &plugin.Context{
			Localizer: mock.NewLocalizer(t).
				On(anyMessageHeader.Term, anyMessageHeader.Fallback.Other).
				On(allMessageInline.Term, allMessageInline.Fallback.Other).
				Build(),
		}

		err := Any(fatalErrorFunc, fatalErrorFunc)(nil, ctx)

		require.IsType(t, new(anyError), err)

		anyErr := err.(*anyError)

		actual := anyErr.Wrap(nil, ctx)
		assert.Equal(t, expect, actual)
	})

	t.Run("not fatal", func(t *testing.T) {
		expect := plugin.NewRestrictionError("You need to fulfill at least one of these requirements to execute " +
			"the command:\n\n" +
			entryPrefix + "abc\n" +
			entryPrefix + "def")

		ctx := &plugin.Context{
			Localizer: mock.NewLocalizer(t).
				On(anyMessageHeader.Term, anyMessageHeader.Fallback.Other).
				On(allMessageInline.Term, allMessageInline.Fallback.Other).
				Build(),
		}

		err := Any(errorFunc1, errorFunc2)(nil, ctx)

		require.IsType(t, new(anyError), err)

		anyErr := err.(*anyError)

		actual := anyErr.Wrap(nil, ctx)
		assert.Equal(t, expect, actual)
	})
}
