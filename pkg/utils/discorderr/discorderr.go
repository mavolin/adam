// Package discorderr provides utilities for interacting with Discord API errors.
package discorderr

import (
	"github.com/diamondburned/arikawa/utils/httputil"

	"github.com/mavolin/adam/pkg/errors"
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
func Is(err *httputil.HTTPError, code httputil.ErrorCode) bool {
	return err != nil && err.Code == code
}

// InRange checks if the passed httputil.HTTPError's code is in the passed
// CodeRange.
func InRange(err *httputil.HTTPError, r CodeRange) bool {
	if err == nil {
		return false
	}

	for _, r := range r {
		if err.Code >= r[0] && err.Code <= r[1] {
			return true
		}
	}

	return false
}
