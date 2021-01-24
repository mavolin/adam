package help

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkHideFuncs(t *testing.T) {
	testCases := []struct {
		name  string
		funcs []HideFunc

		expect HiddenLevel
	}{
		{
			name:   "success",
			funcs:  []HideFunc{mockHideFunc(Show), mockHideFunc(HideList), mockHideFunc(Show)},
			expect: HideList,
		},
		{
			name:   "hide is max",
			funcs:  []HideFunc{mockHideFunc(Show), mockHideFunc(HideList), mockHideFunc(Hide + 1)},
			expect: Hide,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := checkHideFuncs(nil, nil, nil, c.funcs...)
			assert.Equal(t, c.expect, actual)
		})
	}
}
