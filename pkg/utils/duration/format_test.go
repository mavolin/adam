package duration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	t.Parallel()

	d := 10*Year + 70*Minute

	expect := "10y 1h 10min"

	actual := Format(d)
	assert.Equal(t, expect, actual)
}
