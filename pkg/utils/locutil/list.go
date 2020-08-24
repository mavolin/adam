package locutil

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mavolin/adam/pkg/localization"
)

// InterfaceToList creates a written list filled with fmt.Sprinted version of
// the passed data.
//
// Example:
//	1, 3 and 2
func InterfacesToList(list []interface{}, l *localization.Localizer) (s string) {
	var (
		// we can ignore the errors, as we have fallbacks
		defaultSep, _ = l.Localize(defaultSeparatorConfig)
		lastSep, _    = l.Localize(lastSeparatorConfig)
	)

	for i, elem := range list {
		s += fmt.Sprint(elem)

		if i < len(list)-2 {
			s += defaultSep
		} else if i == len(list)-2 {
			s += lastSep
		}
	}

	return
}

// InterfaceToSortedList creates a sorted written list filled with fmt.Sprinted
// version of the passed data.
//
// Example:
//	1, 2 and 3
func InterfacesToSortedList(list []interface{}, l *localization.Localizer) string {
	s := make([]string, len(list))

	for i, elem := range list {
		s[i] = fmt.Sprint(elem)
	}

	return StringsToSortedList(s, l)
}

// StringsToSortedList creates a written list filled the passed strings.
//
// Example:
//	1, 3 and 2
func StringsToList(list []string, l *localization.Localizer) string {
	var (
		// we can ignore the errors, as we have fallbacks
		defaultSep, _ = l.Localize(defaultSeparatorConfig)
		lastSep, _    = l.Localize(lastSeparatorConfig)
	)

	var b strings.Builder

	if len(list) > 2 {
		b.Grow((len(list) - 2) * len(defaultSep))
	}

	b.Grow(len(lastSep))

	for _, s := range list {
		b.Grow(len(s))
	}

	for i, s := range list {
		b.WriteString(s)

		if i < len(list)-2 {
			b.WriteString(defaultSep)
		} else if i == len(list)-2 {
			b.WriteString(lastSep)
		}
	}

	return b.String()
}

// StringsToSortedList creates a sorted written list filled with the passed
// strings.
//
// Example:
//	1, 2 and 3
func StringsToSortedList(list []string, l *localization.Localizer) string {
	sort.Strings(list)
	return StringsToList(list, l)
}

// ConfigsToList creates a written list filled with the passed configs.
//
// Example:
//	1, 3 and 2
func ConfigsToList(list []localization.Config, l *localization.Localizer) (s string, err error) {
	var (
		// we can ignore the errors, as we have fallbacks
		defaultSep, _ = l.Localize(defaultSeparatorConfig)
		lastSep, _    = l.Localize(lastSeparatorConfig)
	)

	var elem string

	for i, c := range list {
		elem, err = l.Localize(c)
		if err != nil {
			return
		}

		s += elem

		if i < len(list)-2 {
			s += defaultSep
		} else if i == len(list)-2 {
			s += lastSep
		}
	}

	return
}

// ConfigsToList creates a sorted written list filled with the passed configs.
//
// Example:
//	1, 2 and 3
func ConfigsToSortedList(list []localization.Config, l *localization.Localizer) (string, error) {
	s := make([]string, len(list))

	for i, c := range list {
		elem, err := l.Localize(c)
		if err != nil {
			return "", err
		}

		s[i] = elem
	}

	return StringsToSortedList(s, l), nil
}

// TermsToList creates a written list filled with the passed terms.
//
// Example:
//	1, 3 and 2
func TermsToList(list []localization.Term, l *localization.Localizer) (s string, err error) {
	var (
		// we can ignore the errors, as we have fallbacks
		defaultSep, _ = l.Localize(defaultSeparatorConfig)
		lastSep, _    = l.Localize(lastSeparatorConfig)
	)

	var elem string

	for i, t := range list {
		elem, err = l.LocalizeTerm(t)
		if err != nil {
			return
		}

		s += elem

		if i < len(list)-2 {
			s += defaultSep
		} else if i == len(list)-2 {
			s += lastSep
		}
	}

	return
}

// TermsToList creates a sorted written list filled with the passed terms.
//
// Example:
//	1, 2 and 3
func TermsToSortedList(list []localization.Term, l *localization.Localizer) (string, error) {
	s := make([]string, len(list))

	for i, t := range list {
		elem, err := l.LocalizeTerm(t)
		if err != nil {
			return "", err
		}

		s[i] = elem
	}

	return StringsToSortedList(s, l), nil
}
