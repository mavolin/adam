// Package discorderr provides utilities for interacting with Discord API errors.
package discorderr

import (
	"errors"

	"github.com/diamondburned/arikawa/v2/utils/httputil"
)

// As calls errors.As(*httputil.HTTPError) on the passed error.
// If errors.As returns true, As returns the httputil.HTTPError, otherwise As
// returns nil.
func As(err error) (herr *httputil.HTTPError) {
	if errors.As(err, &herr) {
		return
	}

	return nil
}

// Is is short for (err != nil && err.Code == code).
func Is(err *httputil.HTTPError, codes ...httputil.ErrorCode) bool {
	if err == nil {
		return false
	}

	for _, c := range codes {
		if c == err.Code {
			return true
		}
	}

	return false
}
