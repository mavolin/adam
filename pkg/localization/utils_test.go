package localization

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_isOne(t *testing.T) {
	successCases := []struct {
		// the number 1 in that type;
		// expected to produce true
		one interface{}
		// the number -1 in that type, nil if out of type range;
		// expected to produce true
		minusOne interface{}
		// a number other than 1;
		// expected to produce false
		other interface{}
	}{
		// uint types
		{uint(1), nil, uint64(0)},
		{uint8(1), nil, uint8(0)},
		{uint16(1), nil, uint16(0)},
		{uint32(1), nil, uint32(0)},
		{uint64(1), nil, uint(0)},

		// int types
		{int(1), int(-1), int(0)},
		{int8(1), int8(-1), int8(0)},
		{int16(1), int16(-1), int16(0)},
		{int32(1), int32(-1), int32(0)},
		{int64(1), int64(-1), int64(0)},

		// float types
		{float32(1), float32(-1), float32(1.0001)},
		{float64(1), float64(-1), float64(1.0001)},

		// string
		{"1", "-1", "0"},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			name := reflect.TypeOf(c.one).Name()

			t.Run(name, func(t *testing.T) {
				plural, err := isOne(c.one)
				if assert.NoError(t, err) {
					assert.True(t, plural)
				}

				if c.minusOne != nil {
					plural, err = isOne(c.minusOne)
					if assert.NoError(t, err) {
						assert.True(t, plural)
					}
				}

				plural, err = isOne(c.other)
				if assert.NoError(t, err) {
					assert.False(t, plural)
				}
			})
		}
	})

	failureCases := []struct {
		name   string
		plural interface{}
	}{
		{
			name:   "invalid string",
			plural: "abc",
		},
		{
			name:   "invalid type",
			plural: []int{1, 2, 3},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				_, err := isOne(c.plural)
				assert.Error(t, err)
			})
		}
	})
}

func Test_fillName(t *testing.T) {
	successCases := []struct {
		name         string
		tmpl         string
		placeholders Placeholders
		expect       string
	}{
		{
			name:   "no template",
			tmpl:   "abc",
			expect: "abc",
		},
		{
			name: "template",
			tmpl: "this is a {{.Test.Type}} test",
			placeholders: Placeholders{
				"Test": Placeholders{
					"Type": "unit",
				},
			},
			expect: "this is a unit test",
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				actual, err := fillTemplate(c.tmpl, c.placeholders)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("syntax error", func(t *testing.T) {
			tmpl := "{{{.Error}}"

			actual, err := fillTemplate(tmpl, nil)
			assert.Errorf(t, err, "succeeded with return '%s'", actual)
		})
	})
}
