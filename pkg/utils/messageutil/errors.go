package messageutil

import (
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/errors"
)

// TimeoutError is an error that fulfills errors.As for *errors.UserError.
type TimeoutError struct {
	UserID discord.UserID
}

func (e *TimeoutError) Error() string {
	return "timeout"
}

func (e *TimeoutError) As(target interface{}) bool {
	switch err := target.(type) {
	case **errors.UserInfo:
		*err = errors.NewUserInfol(timeoutInfo.
			WithPlaceholders(timeoutInfoPlaceholders{
				ResponseUserMention: e.UserID.Mention(),
			}))
		return true
	case *errors.Interface:
		*err = errors.NewUserInfol(timeoutInfo.
			WithPlaceholders(timeoutInfoPlaceholders{
				ResponseUserMention: e.UserID.Mention(),
			}))
		return true
	default:
		return false
	}
}
