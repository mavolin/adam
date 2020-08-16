package localization

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuickConfig(t *testing.T) {
	term := "abc"

	expect := Config{
		Term: term,
	}

	actual := QuickConfig(term)
	assert.Equal(t, expect, actual)
}

func TestQuickFallbackConfig(t *testing.T) {
	var (
		term     = "abc"
		fallback = "def"
	)

	expect := Config{
		Term: term,
		Fallback: Fallback{
			Other: fallback,
		},
	}

	actual := QuickFallbackConfig(term, fallback)
	assert.Equal(t, expect, actual)
}

func TestConfig_placeholdersToMap(t *testing.T) {
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
			placeholders: struct {
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
				Field1 string `localization:"wow_a_custom_name"`
				Field2 bool   `localization:"so_many_possibilities"`
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
		t.Run(c.name, func(t *testing.T) {
			cfg := Config{
				Placeholders: c.placeholders,
			}

			actual, err := cfg.placeholdersToMap()
			assert.NoError(t, err)
			assert.Equal(t, c.expect, actual)
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		cfg := Config{
			Placeholders: []string{},
		}

		_, err := cfg.placeholdersToMap()
		assert.Error(t, err)
	})
}

func TestFallback_genTranslation(t *testing.T) {
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
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				f := Fallback{
					One:   expectOne,
					Other: expectOther,
				}

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
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				f := Fallback{
					One:   c.one,
					Other: c.other,
				}

				_, err := f.genTranslation(nil, c.plural)
				assert.Error(t, err)
			})
		}
	})
}

func TestLocalizer_Localize(t *testing.T) {
	successCases := []struct {
		name     string
		langFunc LangFunc
		config   Config
		expect   string
	}{
		{
			name: "lang func",
			langFunc: func(term string, placeholders map[string]interface{}, plural interface{}) (string, error) {
				if term != "abc" {
					panic(fmt.Sprint("unexpected term: ", term))
				}

				if !reflect.DeepEqual(placeholders, map[string]interface{}{"def": "ghi"}) {
					panic(fmt.Sprint("unexpected placeholders: ", placeholders))
				}

				if plural != "jkl" {
					panic(fmt.Sprint("unexpected plural: ", plural))
				}

				return "abc", nil
			},
			config: Config{
				Term:         "abc",
				Placeholders: map[string]interface{}{"def": "ghi"},
				Plural:       "jkl",
			},
			expect: "abc",
		},
		{
			name:     "fallback",
			langFunc: nil,
			config: Config{
				Fallback: Fallback{
					Other: "abc",
				},
			},
			expect: "abc",
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				l := &Localizer{
					f: c.langFunc,
				}

				actual, err := l.Localize(c.config)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	failureCases := []struct {
		name     string
		langFunc LangFunc
		config   Config
	}{
		{
			name:     "placeholders error",
			langFunc: nil,
			config: Config{
				Placeholders: []string{},
			},
		},
		{
			name: "lang func error",
			langFunc: func(string, map[string]interface{}, interface{}) (string, error) {
				return "", errors.New("something went wrong")
			},
		},
		{
			name:     "no lang func and fallback",
			langFunc: nil,
		},
		{
			name:     "fallback error",
			langFunc: nil,
			config: Config{
				Fallback: Fallback{
					Other: "{{{.Error}}",
				},
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				l := &Localizer{
					f: c.langFunc,
				}

				if c.config.Term == "" {
					c.config.Term = "term"
				}

				actual, err := l.Localize(c.config)
				assert.Equal(t, c.config.Term, actual)
				assert.Error(t, err)
			})
		}
	})
}

func TestLocalizer_LocalizeTerm(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		term := "abc"

		expect := "def"

		l := &Localizer{
			f: func(t string, placeholders map[string]interface{}, plural interface{}) (string, error) {
				if t != term {
					panic(fmt.Sprint("unexpected term: ", term))
				}

				if placeholders != nil {
					panic(fmt.Sprint("unexpected placeholders: ", placeholders))
				}

				if plural != nil {
					panic(fmt.Sprint("unexpected plural: ", plural))
				}

				return expect, nil
			},
		}

		actual, err := l.LocalizeTerm(term)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		term := "abc"

		l := &Localizer{
			f: nil,
		}

		actual, err := l.LocalizeTerm(term)
		assert.Equal(t, term, actual)
		assert.True(t, errors.Is(err, &NoTranslationGeneratedError{
			Term: term,
		}), "unexpected error")
	})
}

func TestLocalizer_MustLocalize(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expect := "abc"

		l := &Localizer{
			f: nil,
		}

		var actual string

		require.NotPanics(t, func() {
			actual = l.MustLocalize(Config{
				Fallback: Fallback{
					Other: expect,
				},
			})
		})
		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		l := &Localizer{
			f: nil,
		}

		assert.Panics(t, func() {
			l.MustLocalize(Config{})
		})
	})
}

func TestLocalizer_MustLocalizeTerm(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		term := "abc"

		expect := "def"

		l := &Localizer{
			f: func(t string, placeholders map[string]interface{}, plural interface{}) (string, error) {
				if t != term {
					panic(fmt.Sprint("unexpected term: ", term))
				}

				if placeholders != nil {
					panic(fmt.Sprint("unexpected placeholders: ", placeholders))
				}

				if plural != nil {
					panic(fmt.Sprint("unexpected plural: ", plural))
				}

				return expect, nil
			},
		}

		var actual string

		assert.NotPanics(t, func() {
			actual = l.MustLocalizeTerm(term)
		})
		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		l := &Localizer{
			f: nil,
		}

		assert.Panics(t, func() {
			l.MustLocalizeTerm("abc")
		})
	})
}
