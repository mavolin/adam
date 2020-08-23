package restriction

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/mock"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestALL(t *testing.T) {
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
			name:  "multiple restriction funcs - multiple errors",
			funcs: []plugin.RestrictionFunc{errorFunc1, errorFunc2},
			expect: &allError{
				restrictions: []*errors.RestrictionError{errorFuncReturn1, errorFuncReturn2},
			},
		},
		{
			name:   "multiple restriction funcs - default error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, defaultRestrictionErrorFunc},
			expect: errors.DefaultRestrictionError,
		},
		{
			name:   "multiple restriction funcs - unexpected error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, unexpectedErrorFunc},
			expect: unexpectedErrorFuncReturn,
		},
		{
			name:   "nested all - single error",
			funcs:  []plugin.RestrictionFunc{ALL(errorFunc1, nilFunc)},
			expect: errorFuncReturn1,
		},
		{
			name:  "nested all - multiple errors",
			funcs: []plugin.RestrictionFunc{ALL(errorFunc1, errorFunc2)},
			expect: &allError{
				restrictions: []*errors.RestrictionError{errorFuncReturn1, errorFuncReturn2},
			},
		},
		{
			name:  "restriction func and nested all",
			funcs: []plugin.RestrictionFunc{errorFunc1, ALL(errorFunc2, errorFunc3)},
			expect: &allError{
				restrictions: []*errors.RestrictionError{errorFuncReturn1, errorFuncReturn2, errorFuncReturn3},
			},
		},
		{
			name:  "nested any",
			funcs: []plugin.RestrictionFunc{ANY(errorFunc1, errorFunc2)},
			expect: &anyError{
				restrictions: []*errors.RestrictionError{errorFuncReturn1, errorFuncReturn2},
			},
		},
		{
			name:  "multiple nested anys",
			funcs: []plugin.RestrictionFunc{ANY(errorFunc1, errorFunc2), ANY(errorFunc3, errorFunc4)},
			expect: &allError{
				anys: []*anyError{
					{
						restrictions: []*errors.RestrictionError{errorFuncReturn1, errorFuncReturn2},
					},
					{
						restrictions: []*errors.RestrictionError{errorFuncReturn3, errorFuncReturn4},
					},
				},
			},
		},
		{
			name:  "restriction func and nested any",
			funcs: []plugin.RestrictionFunc{errorFunc1, ANY(errorFunc2, errorFunc3)},
			expect: &allError{
				restrictions: []*errors.RestrictionError{errorFuncReturn1},
				anys: []*anyError{
					{
						restrictions: []*errors.RestrictionError{errorFuncReturn2, errorFuncReturn3},
					},
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f := ALL(c.funcs...)

			actual := f(nil, &plugin.Context{Localizer: mock.NewNoOpLocalizer()})
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestALLf(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := errors.New("abc")

		f := ALLf(err, nilFunc, nilFunc)

		actual := f(nil, nil)
		assert.Nil(t, actual)
	})

	t.Run("failure", func(t *testing.T) {
		expect := errors.New("abc")

		f := ALLf(expect, errorFunc1, nilFunc)

		actual := f(nil, nil)
		assert.Equal(t, expect, actual)
	})
}

func TestANY(t *testing.T) {
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
			name:  "multiple restriction funcs - all errors",
			funcs: []plugin.RestrictionFunc{errorFunc1, errorFunc2},
			expect: &anyError{
				restrictions: []*errors.RestrictionError{errorFuncReturn1, errorFuncReturn2},
			},
		},
		{
			name:   "multiple restriction funcs - default error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, defaultRestrictionErrorFunc},
			expect: errors.DefaultRestrictionError,
		},
		{
			name:   "multiple restriction funcs - unexpected error",
			funcs:  []plugin.RestrictionFunc{errorFunc1, unexpectedErrorFunc},
			expect: unexpectedErrorFuncReturn,
		},
		{
			name:   "nested all - single error",
			funcs:  []plugin.RestrictionFunc{ALL(errorFunc1, nilFunc)},
			expect: errorFuncReturn1,
		},
		{
			name:  "nested all - multiple errors",
			funcs: []plugin.RestrictionFunc{ALL(errorFunc1, errorFunc2)},
			expect: &allError{
				restrictions: []*errors.RestrictionError{errorFuncReturn1, errorFuncReturn2},
			},
		},
		{
			name:  "restriction func and nested all",
			funcs: []plugin.RestrictionFunc{errorFunc1, ALL(errorFunc2, errorFunc3)},
			expect: &anyError{
				restrictions: []*errors.RestrictionError{errorFuncReturn1},
				alls: []*allError{
					{
						restrictions: []*errors.RestrictionError{errorFuncReturn2, errorFuncReturn3},
					},
				},
			},
		},
		{
			name:  "multiple nested anys",
			funcs: []plugin.RestrictionFunc{ANY(errorFunc1, errorFunc2), ANY(errorFunc3, errorFunc4)},
			expect: &anyError{
				restrictions: []*errors.RestrictionError{
					errorFuncReturn1, errorFuncReturn2, errorFuncReturn3, errorFuncReturn4,
				},
			},
		},
		{
			name:  "restriction func and nested any",
			funcs: []plugin.RestrictionFunc{errorFunc1, ALL(errorFunc2, errorFunc3, nilFunc)},
			expect: &anyError{
				restrictions: []*errors.RestrictionError{errorFuncReturn1},
				alls: []*allError{
					{
						restrictions: []*errors.RestrictionError{errorFuncReturn2, errorFuncReturn3},
					},
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			f := ANY(c.funcs...)

			actual := f(nil, &plugin.Context{Localizer: mock.NewNoOpLocalizer()})
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestANYf(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := errors.New("abc")

		f := ANYf(err, nilFunc, errorFunc1)

		actual := f(nil, nil)
		assert.Nil(t, actual)
	})

	t.Run("failure", func(t *testing.T) {
		expect := errors.New("abc")

		f := ANYf(expect, errorFunc1, errorFunc2)

		actual := f(nil, nil)
		assert.Equal(t, expect, actual)
	})
}

func Test_allError_format(t *testing.T) {
	testCases := []struct {
		name   string
		err    allError
		expect string
	}{
		{
			name: "only restrictions",
			err: allError{
				restrictions: []*errors.RestrictionError{
					errors.NewRestrictionError("abc"),
					errors.NewRestrictionError("def"),
				},
			},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def",
		},
		{
			name: "only anys",
			err: allError{
				anys: []*anyError{
					{
						restrictions: []*errors.RestrictionError{
							errors.NewRestrictionError("abc"),
							errors.NewRestrictionError("def"),
						},
					},
					{
						restrictions: []*errors.RestrictionError{
							errors.NewRestrictionError("ghi"),
							errors.NewRestrictionError("jkl"),
						},
					},
				},
			},
			expect: entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"    " + entryPrefix + "abc\n" +
				"    " + entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"    " + entryPrefix + "ghi\n" +
				"    " + entryPrefix + "jkl",
		},
		{
			name: "both",
			err: allError{
				restrictions: []*errors.RestrictionError{
					errors.NewRestrictionError("abc"),
					errors.NewRestrictionError("def"),
				},
				anys: []*anyError{
					{
						restrictions: []*errors.RestrictionError{
							errors.NewRestrictionError("ghi"),
							errors.NewRestrictionError("jkl"),
						},
					},
					{
						restrictions: []*errors.RestrictionError{
							errors.NewRestrictionError("mno"),
							errors.NewRestrictionError("pqr"),
						},
					},
				},
			},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"    " + entryPrefix + "ghi\n" +
				"    " + entryPrefix + "jkl\n" +
				entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"    " + entryPrefix + "mno\n" +
				"    " + entryPrefix + "pqr",
		},
		{
			name: "any with nested all",
			err: allError{
				anys: []*anyError{
					{
						restrictions: []*errors.RestrictionError{
							errors.NewRestrictionError("abc"),
						},
						alls: []*allError{
							{
								restrictions: []*errors.RestrictionError{
									errors.NewRestrictionError("def"),
									errors.NewRestrictionError("ghi"),
								},
								anys: []*anyError{
									{
										restrictions: []*errors.RestrictionError{
											errors.NewRestrictionError("jkl"),
											errors.NewRestrictionError("mno"),
										},
									},
								},
							},
							{
								restrictions: []*errors.RestrictionError{
									errors.NewRestrictionError("pqr"),
									errors.NewRestrictionError("stu"),
								},
							},
						},
					},
				},
			},
			expect: entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"    " + entryPrefix + "abc\n" +
				"    " + entryPrefix + "You need to fulfill all of these requirements:\n" +
				"        " + entryPrefix + "def\n" +
				"        " + entryPrefix + "ghi\n" +
				"        " + entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"            " + entryPrefix + "jkl\n" +
				"            " + entryPrefix + "mno\n" +
				"    " + entryPrefix + "You need to fulfill all of these requirements:\n" +
				"        " + entryPrefix + "pqr\n" +
				"        " + entryPrefix + "stu",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := c.err.format(0, mock.NewNoOpLocalizer())
			require.NoError(t, err)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func Test_anyError_format(t *testing.T) {
	testCases := []struct {
		name   string
		err    anyError
		expect string
	}{
		{
			name: "only restrictions",
			err: anyError{
				restrictions: []*errors.RestrictionError{
					errors.NewRestrictionError("abc"),
					errors.NewRestrictionError("def"),
				},
			},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def",
		},
		{
			name: "only alls",
			err: anyError{
				alls: []*allError{
					{
						restrictions: []*errors.RestrictionError{
							errors.NewRestrictionError("abc"),
							errors.NewRestrictionError("def"),
						},
					},
					{
						restrictions: []*errors.RestrictionError{
							errors.NewRestrictionError("ghi"),
							errors.NewRestrictionError("jkl"),
						},
					},
				},
			},
			expect: entryPrefix + "You need to fulfill all of these requirements:\n" +
				"    " + entryPrefix + "abc\n" +
				"    " + entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill all of these requirements:\n" +
				"    " + entryPrefix + "ghi\n" +
				"    " + entryPrefix + "jkl",
		},
		{
			name: "both",
			err: anyError{
				restrictions: []*errors.RestrictionError{
					errors.NewRestrictionError("abc"),
					errors.NewRestrictionError("def"),
				},
				alls: []*allError{
					{
						restrictions: []*errors.RestrictionError{
							errors.NewRestrictionError("ghi"),
							errors.NewRestrictionError("jkl"),
						},
					},
					{
						restrictions: []*errors.RestrictionError{
							errors.NewRestrictionError("mno"),
							errors.NewRestrictionError("pqr"),
						},
					},
				},
			},
			expect: entryPrefix + "abc\n" +
				entryPrefix + "def\n" +
				entryPrefix + "You need to fulfill all of these requirements:\n" +
				"    " + entryPrefix + "ghi\n" +
				"    " + entryPrefix + "jkl\n" +
				entryPrefix + "You need to fulfill all of these requirements:\n" +
				"    " + entryPrefix + "mno\n" +
				"    " + entryPrefix + "pqr",
		},
		{
			name: "any with nested all",
			err: anyError{
				alls: []*allError{
					{
						restrictions: []*errors.RestrictionError{
							errors.NewRestrictionError("abc"),
						},
						anys: []*anyError{
							{
								restrictions: []*errors.RestrictionError{
									errors.NewRestrictionError("def"),
									errors.NewRestrictionError("ghi"),
								},
								alls: []*allError{
									{
										restrictions: []*errors.RestrictionError{
											errors.NewRestrictionError("jkl"),
											errors.NewRestrictionError("mno"),
										},
									},
								},
							},
							{
								restrictions: []*errors.RestrictionError{
									errors.NewRestrictionError("pqr"),
									errors.NewRestrictionError("stu"),
								},
							},
						},
					},
				},
			},
			expect: entryPrefix + "You need to fulfill all of these requirements:\n" +
				"    " + entryPrefix + "abc\n" +
				"    " + entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"        " + entryPrefix + "def\n" +
				"        " + entryPrefix + "ghi\n" +
				"        " + entryPrefix + "You need to fulfill all of these requirements:\n" +
				"            " + entryPrefix + "jkl\n" +
				"            " + entryPrefix + "mno\n" +
				"    " + entryPrefix + "You need to fulfill at least one of these requirements:\n" +
				"        " + entryPrefix + "pqr\n" +
				"        " + entryPrefix + "stu",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := c.err.format(0, mock.NewNoOpLocalizer())
			require.NoError(t, err)
			assert.Equal(t, c.expect, actual)
		})
	}
}
