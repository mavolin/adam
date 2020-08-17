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

	t.Run("not nil", func(t *testing.T) {
		cause := New("abc")

		err := Silent(cause)

		assert.IsType(t, new(SilentError), err)

		casted := err.(*SilentError)

		assert.Equal(t, cause, casted.cause)
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

		assert.Equal(t, fmt.Sprintf("%s: %s", message, cause.Error()), err.Error())
	})
}
