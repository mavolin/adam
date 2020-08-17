package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	expect := "abc"

	actual := New(expect)

	assert.Equal(t, expect, actual.Error())
}

func TestNewWithStack(t *testing.T) {
	expect := "abc"

	actual := NewWithStack(expect)

	assert.Equal(t, expect, actual.Error())
}

func TestNewWithStackf(t *testing.T) {
	expect := "abc def"

	actual := NewWithStackf("abc %s", "def")

	assert.Equal(t, expect, actual.Error())
}
