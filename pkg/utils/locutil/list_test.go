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
		On("common.lists.default_separator", ", ").
		On("common.lists.last_separator", " and ").
		Build()

	actual := InterfacesToList(list, l)
	assert.Equal(t, expect, actual)
}

func TestInterfacesToSortedList(t *testing.T) {
	expect := "eggs, jam, milk and orange juice"

	list := []interface{}{"milk", "eggs", "orange juice", "jam"}

	l := mock.
		NewLocalizer().
		On("common.lists.default_separator", ", ").
		On("common.lists.last_separator", " and ").
		Build()

	actual := InterfacesToSortedList(list, l)
	assert.Equal(t, expect, actual)
}

func TestStringsToList(t *testing.T) {
	expect := "def, jkl, ghi and abc"

	list := []string{"def", "jkl", "ghi", "abc"}

	l := mock.
		NewLocalizer().
		On("common.lists.default_separator", ", ").
		On("common.lists.last_separator", " and ").
		Build()

	actual := StringsToList(list, l)
	assert.Equal(t, expect, actual)
}

func TestStringsToSortedList(t *testing.T) {
	expect := "abc, def, ghi and jkl"

	list := []string{"def", "jkl", "ghi", "abc"}

	l := mock.
		NewLocalizer().
		On("common.lists.default_separator", ", ").
		On("common.lists.last_separator", " and ").
		Build()

	actual := StringsToSortedList(list, l)
	assert.Equal(t, expect, actual)
}

func TestConfigsToList(t *testing.T) {
	expect := "milk, eggs and orange juice"

	list := []localization.Config{
		localization.NewTermConfig("milk"),
		localization.NewTermConfig("eggs"),
		localization.NewTermConfig("orange_juice"),
	}

	l := mock.
		NewLocalizer().
		On("common.lists.default_separator", ", ").
		On("common.lists.last_separator", " and ").
		On("milk", "milk").
		On("eggs", "eggs").
		On("orange_juice", "orange juice").
		Build()

	actual, err := ConfigsToList(list, l)
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestConfigsToSortedList(t *testing.T) {
	expect := "eggs, milk and orange juice"

	list := []localization.Config{
		localization.NewTermConfig("milk"),
		localization.NewTermConfig("eggs"),
		localization.NewTermConfig("orange_juice"),
	}

	l := mock.
		NewLocalizer().
		On("common.lists.default_separator", ", ").
		On("common.lists.last_separator", " and ").
		On("milk", "milk").
		On("eggs", "eggs").
		On("orange_juice", "orange juice").
		Build()

	actual, err := ConfigsToSortedList(list, l)
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestTermsToList(t *testing.T) {
	expect := "milk, eggs and orange juice"

	list := []localization.Term{"milk", "eggs", "orange_juice"}

	l := mock.
		NewLocalizer().
		On("common.lists.default_separator", ", ").
		On("common.lists.last_separator", " and ").
		On("milk", "milk").
		On("eggs", "eggs").
		On("orange_juice", "orange juice").
		Build()

	actual, err := TermsToList(list, l)
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestTermsToSortedList(t *testing.T) {
	expect := "eggs, milk and orange juice"

	list := []localization.Term{"milk", "eggs", "orange_juice"}

	l := mock.
		NewLocalizer().
		On("common.lists.default_separator", ", ").
		On("common.lists.last_separator", " and ").
		On("milk", "milk").
		On("eggs", "eggs").
		On("orange_juice", "orange juice").
		Build()

	actual, err := TermsToSortedList(list, l)
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}
