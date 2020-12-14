package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSilent(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Silent(nil)
		assert.Nil(t, err)
	})

	t.Run("silent error", func(t *testing.T) {
		cause := NewSilent("abc").(*SilentError)

		serr := Silent(cause).(*SilentError)

		assert.Equal(t, cause.cause, serr.cause)
		assert.Equal(t, cause.stack, serr.stack)
	})

	t.Run("internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		serr := Silent(cause).(*SilentError)

		assert.Equal(t, cause.cause, serr.cause)
		assert.Equal(t, cause.stack, serr.stack)
	})

	t.Run("Error", func(t *testing.T) {
		err := Silent(&asError{as: NewInformationalError("abc")})
		assert.Nil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		serr := Silent(cause).(*SilentError)

		assert.Equal(t, cause, serr.cause)
	})
}

func TestMustSilent(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := MustSilent(nil)
		assert.Nil(t, err)
	})

	t.Run("silent error", func(t *testing.T) {
		cause := NewSilent("abc").(*SilentError)

		serr := MustSilent(cause).(*SilentError)

		assert.Equal(t, cause.cause, serr.cause)
		assert.Equal(t, cause.stack, serr.stack)
	})

	t.Run("internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		serr := MustSilent(cause).(*SilentError)

		assert.Equal(t, cause.cause, serr.cause)
		assert.Equal(t, cause.stack, serr.stack)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		serr := MustSilent(cause).(*SilentError)

		assert.Equal(t, cause, serr.cause)
	})
}

func TestWrapSilent(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrap(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent error", func(t *testing.T) {
		cause := NewSilent("abc").(*SilentError)

		expectMsg := "def"

		serr := WrapSilent(cause, expectMsg).(*SilentError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, serr.cause)
		assert.Equal(t, cause.stack, serr.stack)
	})

	t.Run("internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		expectMsg := "def"

		serr := WrapSilent(cause, expectMsg).(*SilentError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, serr.cause)
		assert.Equal(t, cause.stack, serr.stack)
	})

	t.Run("Error", func(t *testing.T) {
		err := WrapSilent(&asError{as: NewInformationalError("abc")}, "def")

		assert.Nil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		expectMsg := "def"

		serr := WrapSilent(cause, expectMsg).(*SilentError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause}, serr.cause)
	})
}

func TestWrapSilentf(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WrapSilentf(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent error", func(t *testing.T) {
		cause := NewSilent("abc").(*SilentError)

		expectMsg := "def ghi"

		serr := WrapSilentf(cause, "def %s", "ghi").(*SilentError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, serr.cause)
		assert.Equal(t, cause.stack, serr.stack)
	})

	t.Run("internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		expectMsg := "def ghi"

		serr := WrapSilentf(cause, "def %s", "ghi").(*SilentError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, serr.cause)
		assert.Equal(t, cause.stack, serr.stack)
	})

	t.Run("Error", func(t *testing.T) {
		err := WrapSilentf(&asError{as: NewInformationalError("abc")}, "def %s", "ghi")

		assert.Nil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		expectMsg := "def ghi"

		serr := WrapSilentf(cause, "def %s", "ghi").(*SilentError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause}, serr.cause)
	})
}
