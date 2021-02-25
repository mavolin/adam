package messageutil

import (
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/adam/pkg/errors"
)

// TimeoutError is an error that fulfills errors.As for *errors.UserError.
type TimeoutError struct {
	UserID discord.UserID
	// Cause contains the cause of the TimeoutError.
	// Currently, this will only be filled, if a context.Context expires, while
	// awaiting a message or reaction.
	// Should that be the case, Cause will hold ctx.Err().
	//
	// If Cause is nil, it can be assumed that the timeout was reached
	// regularly.
	Cause error
}

func (e *TimeoutError) Error() string {
	return "timeout"
}

// Unwrap returns the cause for the timeout.
// Refer to the documentation for Cause for more information.
func (e *TimeoutError) Unwrap() error {
	return e.Cause
}

func (e *TimeoutError) As(target interface{}) bool {
	switch err := target.(type) {
	case **errors.UserInfo:
		*err = errors.NewUserInfol(timeoutInfo.
			WithPlaceholders(timeoutInfoPlaceholders{
				Mention: e.UserID.Mention(),
			}))
		return true
	case *errors.Error:
		*err = errors.NewUserInfol(timeoutInfo.
			WithPlaceholders(timeoutInfoPlaceholders{
				Mention: e.UserID.Mention(),
			}))
		return true
	default:
		return false
	}
}
