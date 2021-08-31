package msgbuilder

import (
	"fmt"

	"github.com/mavolin/adam/internal/errorutil"
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

// NewActionRowError creates a new *ActionRowError.
func NewActionRowError(index int, typ string, cause error) *ActionRowError {
	return &ActionRowError{
		Index: index,
		Type:  typ,
		Cause: cause,
		s:     errorutil.GenerateStackTrace(1),
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

	s errorutil.StackTrace
}

// NewSelectError creates a new *ActionRowError.
func NewSelectError(index int, cause error) *SelectError {
	return &SelectError{
		Index: index,
		Cause: cause,
		s:     errorutil.GenerateStackTrace(1),
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
