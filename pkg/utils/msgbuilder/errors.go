package msgbuilder

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/internal/errorutil"
	"github.com/mavolin/adam/pkg/errors"
)

// =============================================================================
// ActionRowError
// =====================================================================================

// ActionRowError is the error returned if one of the ComponentBuilders inside
// an ActionRow failed to build.
type ActionRowError struct {
	Index int
	Type  string

	Cause error

	s errorutil.StackTrace
}

// newActionRowError creates a new *ActionRowError.
func newActionRowError(index int, typ string, cause error) *ActionRowError {
	return &ActionRowError{
		Index: index,
		Type:  typ,
		Cause: cause,
		s:     errors.GenerateStackTrace(1),
	}
}

func (e *ActionRowError) StackTrace() errorutil.StackTrace {
	return e.s
}

func (e *ActionRowError) Unwrap() error {
	return e.Cause
}

func (e *ActionRowError) Error() string {
	return fmt.Sprintf("msgbuilder: ActionRowBuilder: could not build the %s at index %d: %s",
		e.Type, e.Index, e.Cause.Error())
}

// =============================================================================
// SelectError
// =====================================================================================

// SelectError is the error returned if one of the SelectOptionBuilders inside
// an Select failed to build.
type SelectError struct {
	Index int

	Cause error

	s errors.StackTrace
}

// newSelectError creates a new *ActionRowError.
func newSelectError(index int, cause error) *SelectError {
	return &SelectError{
		Index: index,
		Cause: cause,
		s:     errors.GenerateStackTrace(1),
	}
}

func (e *SelectError) StackTrace() errorutil.StackTrace {
	return e.s
}

func (e *SelectError) Unwrap() error {
	return e.Cause
}

func (e *SelectError) Error() string {
	return fmt.Sprintf("msgbuilder: SelectBuilder: could not build the SelectOptionBuilder at index %d: %s",
		e.Index, e.Cause.Error())
}

// =============================================================================
// TimeoutError
// =====================================================================================

// TimeoutError is an error that fulfills errors.As for *errors.UserInfo.
// It is returned by Builder's await methods if the context is canceled before
// a matching component interaction happens.
type TimeoutError struct {
	UserID discord.UserID
	Cause  error
}

func (e *TimeoutError) Error() string {
	return "timeout"
}

// Unwrap returns the cause for the timeout.
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
