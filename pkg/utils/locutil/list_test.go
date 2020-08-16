package locutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/mock"
)

func TestInterfacesToList(t *testing.T) {
	expect := "milk, eggs, orange juice and jam"

	list := []interface{}{"milk", "eggs", "orange juice", "jam"}

	l := mock.
		NewLocalizer().
		On("lang.lists.default_separator", ", ").
		On("lang.lists.last_separator", " and ").
		Build()

	actual := InterfacesToList(list, l)
	assert.Equal(t, expect, actual)
}

func TestConfigsToList(t *testing.T) {
	expect := "milk, eggs and orange juice"

	list := []localization.Config{
		localization.QuickConfig("milk"),
		localization.QuickConfig("eggs"),
		localization.QuickConfig("orange_juice"),
	}

	l := mock.
		NewLocalizer().
		On("lang.lists.default_separator", ", ").
		On("lang.lists.last_separator", " and ").
		On("milk", "milk").
		On("eggs", "eggs").
		On("orange_juice", "orange juice").
		Build()

	actual, err := ConfigsToList(list, l)
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestTermsToList(t *testing.T) {
	expect := "milk, eggs and orange juice"

	list := []string{"milk", "eggs", "orange_juice"}

	l := mock.
		NewLocalizer().
		On("lang.lists.default_separator", ", ").
		On("lang.lists.last_separator", " and ").
		On("milk", "milk").
		On("eggs", "eggs").
		On("orange_juice", "orange juice").
		Build()

	actual, err := TermsToList(list, l)
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}
