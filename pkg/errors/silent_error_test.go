package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSilent(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Silent(nil)
		assert.Nil(t, err)
	})

	t.Run("silent error", func(t *testing.T) {
		expect := Silent(New("abc"))

		actual := Silent(expect)
		assert.Equal(t, expect, actual)
	})

	t.Run("internal error", func(t *testing.T) {
		expectCause := New("abc")

		cause := WithStack(expectCause)

		actual := Silent(cause)
		assert.Equal(t, expectCause, actual.(*SilentError).cause)
	})

	t.Run("normal error", func(t *testing.T) {
		cause := New("abc")

		err := Silent(cause)
		assert.Equal(t, cause, err.(*SilentError).cause)
	})
}

func TestWrapSilent(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrap(nil, "")
		assert.Nil(t, err)
	})

	t.Run("not nil", func(t *testing.T) {
		var (
			cause   = New("abc")
			message = "def"
		)

		err := WrapSilent(cause, message)

		//goland:noinspection GoNilness
		assert.Equal(t, fmt.Sprintf("%s: %s", message, cause.Error()), err.Error())
	})
}

func TestWrapSilentf(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WrapSilentf(nil, "")
		assert.Nil(t, err)
	})

	t.Run("not nil", func(t *testing.T) {
		var (
			cause   = New("abc")
			message = "def ghi"
		)

		err := WrapSilentf(cause, "def %s", "ghi")

		//goland:noinspection GoNilness
		assert.Equal(t, fmt.Sprintf("%s: %s", message, cause.Error()), err.Error())
	})
}
