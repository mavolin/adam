package duration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	s := "50min 3y 12 min6s3h"

	expect := 50*Minute + 3*Year + 12*Minute + 6*Second + 3*Hour

	actual, err := Parse(s)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}
