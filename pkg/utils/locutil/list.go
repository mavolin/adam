package locutil

import (
	"fmt"

	"github.com/mavolin/adam/pkg/localization"
)

// InterfaceToList creates a written list filled with fmt.Sprinted version of
// the passed data.
func InterfacesToList(list []interface{}, l *localization.Localizer) (s string) {
	var (
		// we can ignore the errors, as we have fallbacks
		defaultSep, _ = l.Localize(defaultSeparatorConfig)
		lastSep, _    = l.Localize(lastSepartatorConfig)
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

// ConfigsToList creates a written list filled with the passed configs.
//
// Example:
//	1, 2 and 3
func ConfigsToList(list []localization.Config, l *localization.Localizer) (s string, err error) {
	var (
		// we can ignore the errors, as we have fallbacks
		defaultSep, _ = l.Localize(defaultSeparatorConfig)
		lastSep, _    = l.Localize(lastSepartatorConfig)
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

// TermsToList creates a written list filled with the passed terms.
//
// Example:
//	1, 2 and 3
func TermsToList(list []string, l *localization.Localizer) (s string, err error) {
	var (
		// we can ignore the errors, as we have fallbacks
		defaultSep, _ = l.Localize(defaultSeparatorConfig)
		lastSep, _    = l.Localize(lastSepartatorConfig)
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
