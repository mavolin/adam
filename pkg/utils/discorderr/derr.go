// Package discorderr provides utilities for interacting with Discord API errors.
package discorderr

import (
	"github.com/diamondburned/arikawa/utils/httputil"

	"github.com/mavolin/adam/pkg/errors"
)

func As(err error) (herr *httputil.HTTPError) {
	if errors.As(err, &herr) {
		return
	}

	return nil
}

func Is(err error, code httputil.ErrorCode) bool {
	herr := As(err)
	if herr == nil {
		return false
	}

	return herr.Code == code
}

func InRange(err error, r CodeRange) bool {
	herr := As(err)
	if herr == nil {
		return false
	}

	for _, r := range r {
		if herr.Code >= r[0] && herr.Code <= r[1] {
			return true
		}
	}

	return false
}
