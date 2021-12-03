package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTermConfig(t *testing.T) {
	t.Parallel()

	var term Term = "abc"

	expect := term.AsConfig()

	actual := NewTermConfig(term)
	assert.Equal(t, expect, actual)
}

func TestNewStaticConfig(t *testing.T) {
	t.Parallel()

	expect := "{{abc}}"

	c := NewStaticConfig(expect)

	actual, err := NewFallbackLocalizer().Localize(c)
	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestNewFallbackConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		term     Term
		fallback string

		expect *Config
	}{
		{
			name:     "fallback",
			term:     "abc",
			fallback: "def",
			expect: &Config{
				Term:     "abc",
				Fallback: Fallback{Other: "def"},
			},
		},
		{
			name:     "empty term and empty fallback",
			term:     "",
			fallback: "",
			expect:   &Config{},
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := NewFallbackConfig(c.term, c.fallback)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestConfig_WithPlaceholders(t *testing.T) {
	t.Parallel()

	a := Config{Term: "abc"}

	b := a.WithPlaceholders(map[string]interface{}{"def": "ghi"})

	assert.NotEqual(t, a, b)
	assert.Equal(t, a.Term, b.Term)
}

func TestConfig_placeholdersToMap(t *testing.T) {
	t.Parallel()

	successCases := []struct {
		name         string
		placeholders interface{}
		expect       map[string]interface{}
	}{
		{
			name:         "nil",
			placeholders: nil,
			expect:       nil,
		},
		{
			name: "map",
			placeholders: map[string]interface{}{
				"abc": true,
				"def": 123,
			},
			expect: map[string]interface{}{
				"abc": true,
				"def": 123,
			},
		},
		{
			name: "struct",
			placeholders: &struct {
				ThisIsFieldNumber1 string
				JSONData           int
				EvenMore           string
			}{
				ThisIsFieldNumber1: "abc",
				JSONData:           123,
				EvenMore:           "def",
			},
			expect: map[string]interface{}{
				"this_is_field_number_1": "abc",
				"json_data":              123,
				"even_more":              "def",
			},
		},
		{
			name: "handle unexported fields",
			placeholders: struct {
				Exported   int
				unexported string
			}{
				Exported:   123,
				unexported: "def",
			},
			expect: map[string]interface{}{
				"exported": 123,
			},
		},
		{
			name: "pointer to struct",
			placeholders: &struct {
				Field1 string
				Field2 int
				Field3 string
			}{
				Field1: "abc",
				Field2: 123,
				Field3: "def",
			},
			expect: map[string]interface{}{
				"field_1": "abc",
				"field_2": 123,
				"field_3": "def",
			},
		},
		{
			name: "struct tags",
			placeholders: struct {
				Field1 string `i18n:"wow_a_custom_name"`
				Field2 bool   `i18n:"so_many_possibilities"`
				Field3 int
			}{
				Field1: "abc",
				Field2: false,
				Field3: 123,
			},
			expect: map[string]interface{}{
				"wow_a_custom_name":     "abc",
				"so_many_possibilities": false,
				"field_3":               123,
			},
		},
	}

	for _, c := range successCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			cfg := Config{Placeholders: c.placeholders}

			actual, err := cfg.placeholdersToMap()
			require.NoError(t, err)
			assert.Equal(t, c.expect, actual)
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		t.Parallel()

		cfg := Config{Placeholders: []string{}}

		_, err := cfg.placeholdersToMap()
		assert.Error(t, err)
	})
}

func TestFallback_genTranslation(t *testing.T) {
	t.Parallel()

	const (
		expectOne   = "abc"
		expectOther = "def"
	)

	successCases := []struct {
		name string

		plural interface{}
		expect string
	}{
		{
			name:   "plural is 1",
			plural: 1,
			expect: expectOne,
		},
		{
			name:   "plural not 1",
			plural: 0,
			expect: expectOther,
		},
		{
			name:   "no plural",
			plural: nil,
			expect: expectOther,
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		for _, c := range successCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				f := Fallback{One: expectOne, Other: expectOther}

				actual, err := f.genTranslation(nil, c.plural)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	failureCases := []struct {
		name       string
		one, other string
		plural     interface{}
	}{
		{
			name:   "invalid plural type",
			plural: []int{1, 2, 3},
		},
		{
			name:   "invalid one template",
			one:    "{{{.Error}}",
			plural: 1,
		},
		{
			name:   "invalid other template",
			one:    "",
			other:  "{{{.Error}}",
			plural: nil,
		},
	}

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		for _, c := range failureCases {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				f := Fallback{One: c.one, Other: c.other}

				_, err := f.genTranslation(nil, c.plural)
				assert.Error(t, err)
			})
		}
	})
}
