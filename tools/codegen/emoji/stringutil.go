package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	// nonLatinReplacer is the strings.Replacer used to replace abbreviations
	// and characters that are not replaced during normalization.
	nonLatinReplacer = strings.NewReplacer(
		"*", "asterisk",
		"#", "hash",
		"&", "and",
		"U.S.", "US",
		"1st", "first",
		"2nd", "second",
		"3rd", "third",
		// the below characters are not replaced during normalization and need
		// to be replaced manually
		"Đ", "Dj",
		"đ", "dj",
		"Æ", "A",
		"Ç", "C",
		"Ø", "O",
		"Þ", "B",
		"ß", "ss",
		"æ", "a",
		"ø", "o",
	)

	// nonLatinRegexp is the regular expression used in the last step of name
	// sanitization.
	// It matches all non-letter and -digits.
	nonLatinRegex = regexp.MustCompile(`\W+`)
)

func sanitizeName(name string) string {
	name = nonLatinReplacer.Replace(name)

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	name, _, err := transform.String(t, name)
	if err != nil {
		panic(err)
	}

	return nonLatinRegex.ReplaceAllString(name, " ")
}

func genName(desc string) string {
	desc = sanitizeName(desc)
	return strcase.ToCamel(desc)
}

const (
	lightSkin       = '\U0001F3FB'
	mediumLightSkin = '\U0001F3FC'
	mediumSkin      = '\U0001F3FD'
	mediumDarkSkin  = '\U0001F3FE'
	darkSkin        = '\U0001F3FF'
)

func withSkinTone(base string, tone rune) string {
	chars := []rune(base)

	chars = append(chars, 0)
	copy(chars[2:], chars[1:])

	chars[1] = tone

	return string(chars)
}

func escapedUnicodeSequence(s string) string {
	var b strings.Builder

	b.Grow(len(s) * 10) // max space used per char

	for _, r := range s {
		if r <= 0xFFFF {
			b.WriteString(`\u`)
			b.WriteString(fmt.Sprintf("%.4x", r))
		} else {
			b.WriteString(`\U`)
			b.WriteString(fmt.Sprintf("%.8x", r))
		}
	}

	return b.String()
}
