package restriction

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/mock"
	"github.com/mavolin/adam/pkg/plugin"
)

func TestWrap(t *testing.T) {
	testCases := []struct {
		name   string
		in     plugin.RestrictionFunc
		expect error
	}{
		{
			name: "all error",
			in:   ALL(errorFunc1, errorFunc2),
			expect: errors.NewRestrictionError("You need to fulfill all of these requirements:\n\n" +
				entryPrefix + "abc\n" +
				entryPrefix + "def"),
		},
		{
			name: "any error",
			in:   ANY(errorFunc1, errorFunc2),
			expect: errors.NewRestrictionError("You need to fulfill at least one of these requirements:\n\n" +
				entryPrefix + "abc\n" +
				entryPrefix + "def"),
		},
		{
			name:   "other",
			in:     unexpectedErrorFunc,
			expect: unexpectedErrorFuncReturn,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			wrapped := Wrap(c.in)

			actual := wrapped(nil, &plugin.Context{Localizer: mock.NewNoOpLocalizer()})
			assert.Equal(t, c.expect, actual)
		})
	}
}
