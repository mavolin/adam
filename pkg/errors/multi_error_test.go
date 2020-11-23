package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppend(t *testing.T) {
	err1 := NewWithStack("abc")
	err2 := NewWithStack("def")
	err3 := NewWithStack("ghi")
	err4 := NewWithStack("jkl")

	testCases := []struct {
		name       string
		err1, err2 error
		expect     error
	}{
		{
			name:   "no multiErrors",
			err1:   err1,
			err2:   err2,
			expect: multiError{err1, err2},
		},
		{
			name:   "err1 is multiError",
			err1:   multiError{err1, err2},
			err2:   err3,
			expect: multiError{err1, err2, err3},
		},
		{
			name:   "err2 is multiError",
			err1:   err1,
			err2:   multiError{err2, err3},
			expect: multiError{err1, err2, err3},
		},
		{
			name:   "both are multiErrors",
			err1:   multiError{err1, err2},
			err2:   multiError{err3, err4},
			expect: multiError{err1, err2, err3, err4},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Append(c.err1, c.err2)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestAppendSilent(t *testing.T) {
	t.Run("no multiErrors", func(t *testing.T) {
		err1 := New("abc")
		err2 := New("def")

		actual := AppendSilent(err1, err2)
		require.IsType(t, multiError{}, actual)

		actualTyped := actual.(multiError)
		require.Len(t, actualTyped, 2)

		require.IsType(t, new(SilentError), actualTyped[0])
		actual0Typed := actualTyped[0].(*SilentError)
		assert.Equal(t, err1, actual0Typed.Unwrap())

		require.IsType(t, new(SilentError), actualTyped[1])
		actual1Typed := actualTyped[1].(*SilentError)
		assert.Equal(t, err2, actual1Typed.Unwrap())
	})

	t.Run("err1 is multiError", func(t *testing.T) {
		err1 := multiError{New("abc"), New("def")}
		err2 := New("ghi")

		actual := AppendSilent(err1, err2)
		require.IsType(t, multiError{}, actual)

		actualTyped := actual.(multiError)
		require.Len(t, actualTyped, 3)

		assert.Equal(t, err1, actualTyped[:2])

		require.IsType(t, new(SilentError), actualTyped[2])
		actual3Typed := actualTyped[2].(*SilentError)
		assert.Equal(t, err2, actual3Typed.Unwrap())
	})

	t.Run("err2 is multiError", func(t *testing.T) {
		err1 := New("abc")
		err2 := multiError{New("def"), New("ghi")}

		actual := AppendSilent(err1, err2)
		require.IsType(t, multiError{}, actual)

		actualTyped := actual.(multiError)
		require.Len(t, actualTyped, 3)

		require.IsType(t, new(SilentError), actualTyped[0])
		actual0Typed := actualTyped[0].(*SilentError)
		assert.Equal(t, err1, actual0Typed.Unwrap())

		assert.Equal(t, err2, actualTyped[1:3])
	})

	t.Run("both are multiErrors", func(t *testing.T) {
		err1 := multiError{New("abc"), New("def")}
		err2 := multiError{New("ghi"), New("jkl")}

		actual := AppendSilent(err1, err2)
		require.IsType(t, multiError{}, actual)

		actualTyped := actual.(multiError)
		require.Len(t, actualTyped, 4)

		assert.Equal(t, err1, actualTyped[:2])
		assert.Equal(t, err2, actualTyped[2:4])
	})
}

func TestCombine(t *testing.T) {
	err1 := NewWithStack("abc")
	err2 := NewWithStack("def")
	err3 := NewWithStack("ghi")

	testCases := []struct {
		name   string
		errs   []error
		expect error
	}{
		{
			name:   "no errors",
			errs:   []error{},
			expect: nil,
		},
		{
			name:   "single error",
			errs:   []error{err1},
			expect: err1,
		},
		{
			name:   "no multiErrors",
			errs:   []error{err1, err2},
			expect: multiError{err1, err2},
		},
		{
			name:   "single error and multiError",
			errs:   []error{err1, multiError{err2, err3}},
			expect: multiError{err1, err2, err3},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := Combine(c.errs...)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestCombineSilent(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		var expect error = nil

		actual := CombineSilent()
		assert.Equal(t, expect, actual)
	})

	t.Run("single error", func(t *testing.T) {
		err := New("abc")

		actual := CombineSilent(err)
		require.IsType(t, new(SilentError), actual)
		actualTyped := actual.(*SilentError)

		assert.Equal(t, err, actualTyped.Unwrap())
	})

	t.Run("no multiErrors", func(t *testing.T) {
		err1 := New("abc")
		err2 := New("def")

		actual := CombineSilent(err1, err2)

		require.IsType(t, multiError{}, actual)
		actualTyped := actual.(multiError)

		require.Len(t, actualTyped, 2)

		require.IsType(t, new(SilentError), actualTyped[0])
		actual0Typed := actualTyped[0].(*SilentError)
		assert.Equal(t, err1, actual0Typed.Unwrap())

		require.IsType(t, new(SilentError), actualTyped[1])
		actual1Typed := actualTyped[1].(*SilentError)
		assert.Equal(t, err2, actual1Typed.Unwrap())
	})

	t.Run("single error and multiError", func(t *testing.T) {
		err1 := New("abc")
		err2 := multiError{New("def"), New("ghi")}

		actual := CombineSilent(err1, err2)

		require.IsType(t, multiError{}, actual)
		actualTyped := actual.(multiError)

		require.Len(t, actualTyped, 3)

		require.IsType(t, new(SilentError), actualTyped[0])
		actual0Typed := actualTyped[0].(*SilentError)
		assert.Equal(t, err1, actual0Typed.Unwrap())

		assert.Equal(t, err2, actualTyped[1:])
	})
}

func TestRetrieveErrors(t *testing.T) {
	testCases := []struct {
		name   string
		err    error
		expect []error
	}{
		{
			name:   "no multiError",
			err:    New("abc"),
			expect: []error{New("abc")},
		},
		{
			name:   "multiError",
			err:    multiError{New("abc"), New("def")},
			expect: []error{New("abc"), New("def")},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := RetrieveMultiError(c.err)
			assert.Equal(t, c.expect, actual)
		})
	}
}
