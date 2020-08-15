package localization

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewManager(t *testing.T) {
	t.Run("nil Func", func(t *testing.T) {
		m := NewManager(nil)

		assert.NotNil(t, m.f)
	})
}

func TestLocalizer(t *testing.T) {
	var mocker mock.Mock

	expect := &Localizer{
		f:    nil,
		Lang: "de_DE",
	}

	f := Func(func(lang string) LangFunc {
		return mocker.Called(lang).Get(0).(LangFunc)
	})

	mocker.
		On("func1", expect.Lang).
		Return(expect.f)

	m := NewManager(f)
	actual := m.Localizer(expect.Lang)

	assert.Equal(t, expect, actual)

	mocker.AssertExpectations(t)
}
