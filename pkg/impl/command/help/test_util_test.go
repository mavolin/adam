package help

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mockHideFunc(t *testing.T) {
	t.Parallel()

	testCases := []HiddenLevel{Show, HideList, Hide}

	for _, c := range testCases {
		actual := mockHideFunc(c)(nil, nil, nil)
		assert.Equal(t, c, actual)
	}
}
