package localization

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"text/template"
)

// ErrNaN gets returned, if a plural value is neither a number type nor a
// string containing a number
var ErrNaN = errors.New("the used plural value is not a number")

// isOne checks if the passed plural is a number or a string of such that is
// 1 or -1.
func isOne(plural interface{}) (bool, error) {
	switch plural := plural.(type) {
	case uint:
		return plural == 1, nil
	case uint8:
		return plural == 1, nil
	case uint16:
		return plural == 1, nil
	case uint32:
		return plural == 1, nil
	case uint64:
		return plural == 1, nil

	case int:
		return plural == 1 || plural == -1, nil
	case int8:
		return plural == 1 || plural == -1, nil
	case int16:
		return plural == 1 || plural == -1, nil
	case int32:
		return plural == 1 || plural == -1, nil
	case int64:
		return plural == 1 || plural == -1, nil

	case float32:
		return plural == 1 || plural == -1, nil
	case float64:
		return plural == 1 || plural == -1, nil

	case string:
		num, err := strconv.ParseFloat(plural, 64)
		if err != nil {
			return false, ErrNaN
		}

		return num == 1 || num == -1, nil
	default:
		return false, ErrNaN
	}
}

// fillTemplate is a helper to execute templates.
func fillTemplate(tmpl string, placeholders Placeholders) (string, error) {
	if !strings.Contains(tmpl, "{{") { // no need for parsing
		return tmpl, nil
	}

	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer

	err = t.Execute(&b, placeholders)
	if err != nil {
		return "", err
	}

	return b.String(), err
}
