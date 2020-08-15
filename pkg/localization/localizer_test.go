package localization

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			langFunc: func(term string, placeholders Placeholders, plural interface{}) (string, error) {
				if term != "abc" {
					panic(fmt.Sprint("unexpected term: ", term))
				}

				if !reflect.DeepEqual(placeholders, Placeholders{"def": "ghi"}) {
					panic(fmt.Sprint("unexpected placeholders: ", placeholders))
				}

				if plural != "jkl" {
					panic(fmt.Sprint("unexpected plural: ", plural))
				}

				return "abc", nil
			},
			config: Config{
				Term:         "abc",
				Placeholders: Placeholders{"def": "ghi"},
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
			name: "lang func error",
			langFunc: func(string, Placeholders, interface{}) (string, error) {
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
			f: func(t string, placeholders Placeholders, plural interface{}) (string, error) {
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
			f: func(t string, placeholders Placeholders, plural interface{}) (string, error) {
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
