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

func TestInterfacesToSortedList(t *testing.T) {
	expect := "eggs, jam, milk and orange juice"

	list := []interface{}{"milk", "eggs", "orange juice", "jam"}

	l := mock.
		NewLocalizer().
		On("lang.lists.default_separator", ", ").
		On("lang.lists.last_separator", " and ").
		Build()

	actual := InterfacesToSortedList(list, l)
	assert.Equal(t, expect, actual)
}

func TestConfigsToList(t *testing.T) {
	expect := "milk, eggs and orange juice"

	list := []localization.Config{
		localization.Term("milk"),
		localization.Term("eggs"),
		localization.Term("orange_juice"),
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

func TestConfigsToSortedList(t *testing.T) {
	expect := "eggs, milk and orange juice"

	list := []localization.Config{
		localization.Term("milk"),
		localization.Term("eggs"),
		localization.Term("orange_juice"),
	}

	l := mock.
		NewLocalizer().
		On("lang.lists.default_separator", ", ").
		On("lang.lists.last_separator", " and ").
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

func TestTermsToSortedList(t *testing.T) {
	expect := "eggs, milk and orange juice"

	list := []string{"milk", "eggs", "orange_juice"}

	l := mock.
		NewLocalizer().
		On("lang.lists.default_separator", ", ").
		On("lang.lists.last_separator", " and ").
		On("milk", "milk").
		On("eggs", "eggs").
		On("orange_juice", "orange juice").
		Build()

	actual, err := TermsToSortedList(list, l)
	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func Test_stringsToSortedList(t *testing.T) {
	expect := "abc, def, ghi and jkl"

	list := []string{"def", "jkl", "ghi", "abc"}

	actual := stringsToSortedList(list, ", ", " and ")
	assert.Equal(t, expect, actual)
}
