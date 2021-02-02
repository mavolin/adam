package i18n

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalizer_WithDefaultPlaceholder(t *testing.T) {
	t.Run("new map", func(t *testing.T) {
		k, v := "abc", "def"

		expect := map[string]interface{}{
			k: v,
		}

		l := new(Localizer)
		l.WithPlaceholder(k, v)

		assert.Equal(t, expect, l.defaultPlaceholders)
	})

	t.Run("append map", func(t *testing.T) {
		k, v := "ghi", "jkl"

		expect := map[string]interface{}{
			"abc": "def",
			k:     v,
		}

		l := &Localizer{
			defaultPlaceholders: map[string]interface{}{
				"abc": "def",
			},
		}

		l.WithPlaceholder(k, v)

		assert.Equal(t, expect, l.defaultPlaceholders)
	})

	t.Run("overwrite map", func(t *testing.T) {
		k, v := "abc", "ghi"

		expect := map[string]interface{}{
			k: v,
		}

		l := &Localizer{
			defaultPlaceholders: map[string]interface{}{
				k: "def",
			},
		}

		l.WithPlaceholder(k, v)

		assert.Equal(t, expect, l.defaultPlaceholders)
	})
}

func TestLocalizer_WithDefaultPlaceholders(t *testing.T) {
	t.Run("new map", func(t *testing.T) {
		m := map[string]interface{}{
			"abc": 123,
			"def": "ghi",
		}

		l := new(Localizer)
		l.WithPlaceholders(m)

		assert.Equal(t, m, l.defaultPlaceholders)
	})

	t.Run("append map", func(t *testing.T) {
		m := map[string]interface{}{
			"ghi": 123,
			"jkl": "mno",
		}

		expect := map[string]interface{}{
			"abc": "def",
			"ghi": 123,
			"jkl": "mno",
		}

		l := &Localizer{
			defaultPlaceholders: map[string]interface{}{
				"abc": "def",
			},
		}

		l.WithPlaceholders(m)

		assert.Equal(t, expect, l.defaultPlaceholders)
	})

	t.Run("overwrite map", func(t *testing.T) {
		m := map[string]interface{}{
			"abc": 123,
			"def": "ghi",
		}

		expect := map[string]interface{}{
			"abc": 123,
			"def": "ghi",
		}

		l := &Localizer{
			defaultPlaceholders: map[string]interface{}{
				"abc": "def",
			},
		}

		l.WithPlaceholders(m)

		assert.Equal(t, expect, l.defaultPlaceholders)
	})
}

func TestLocalizer_Localize(t *testing.T) {
	successCases := []struct {
		name                string
		defaultPlaceholders map[string]interface{}
		langFunc            func(*testing.T) Func
		config              *Config
		expect              string
	}{
		{
			name: "lang func",
			langFunc: func(t *testing.T) Func {
				return func(term Term, placeholders map[string]interface{}, plural interface{}) (string, error) {
					var (
						expectTerm         Term = "abc"
						expectPlaceholders      = map[string]interface{}{"def": "ghi"}
						expectPlural            = "jkl"
					)

					assert.Equal(t, expectTerm, term)
					assert.Equal(t, expectPlaceholders, placeholders)
					assert.Equal(t, expectPlural, plural)

					return "abc", nil
				}
			},
			config: &Config{
				Term:         "abc",
				Placeholders: map[string]interface{}{"def": "ghi"},
				Plural:       "jkl",
			},
			expect: "abc",
		},
		{
			name: "fallback",
			config: &Config{
				Fallback: Fallback{
					Other: "abc",
				},
			},
			expect: "abc",
		},
		{
			name: "default placeholders",
			defaultPlaceholders: map[string]interface{}{
				"def": "ghi",
			},
			config: &Config{
				Fallback: Fallback{
					Other: "abc {{.def}}",
				},
			},
			expect: "abc ghi",
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				if c.langFunc == nil {
					c.langFunc = func(t *testing.T) Func { return nil }
				}

				l := &Localizer{
					f:                   c.langFunc(t),
					defaultPlaceholders: c.defaultPlaceholders,
				}

				actual, err := l.Localize(c.config)
				require.NoError(t, err)
				assert.Equal(t, c.expect, actual)
			})
		}
	})

	failureCases := []struct {
		name     string
		langFunc Func
		config   *Config
	}{
		{
			name:     "placeholders error",
			langFunc: nil,
			config: &Config{
				Placeholders: []string{},
			},
		},
		{
			name: "lang func error",
			langFunc: func(Term, map[string]interface{}, interface{}) (string, error) {
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
			config: &Config{
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

				if c.config == nil {
					c.config = NewTermConfig("term")
				}

				actual, err := l.Localize(c.config)
				assert.Equal(t, c.config.Term, Term(actual))
				assert.Error(t, err)
			})
		}
	})
}

func TestLocalizer_LocalizeTerm(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var expectTerm Term = "abc"

		expect := "def"

		l := &Localizer{
			f: func(actualTerm Term, placeholders map[string]interface{}, plural interface{}) (string, error) {
				assert.Equal(t, expectTerm, actualTerm)
				assert.Nil(t, placeholders)
				assert.Nil(t, plural)

				return expect, nil
			},
		}

		actual, err := l.LocalizeTerm(expectTerm)
		require.NoError(t, err)
		assert.Equal(t, expect, actual)
	})

	t.Run("failure", func(t *testing.T) {
		var term Term = "abc"

		l := &Localizer{
			f: nil,
		}

		actual, err := l.LocalizeTerm(term)
		assert.Equal(t, term, Term(actual))
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
			actual = l.MustLocalize(&Config{
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
			l.MustLocalize(&Config{})
		})
	})
}

func TestLocalizer_MustLocalizeTerm(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var expectTerm Term = "abc"

		expect := "def"

		l := &Localizer{
			f: func(actualTerm Term, placeholders map[string]interface{}, plural interface{}) (string, error) {
				assert.Equal(t, expectTerm, actualTerm)
				assert.Nil(t, placeholders)
				assert.Nil(t, plural)

				return expect, nil
			},
		}

		var actual string

		require.NotPanics(t, func() {
			actual = l.MustLocalizeTerm(expectTerm)
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
