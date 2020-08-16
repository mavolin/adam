package errorutil

import "runtime"

// stackDepth is the maximum depth of a Stack.
const stackDepth = 32

// Stack is a collection of program counters that show the stack trace of an
// error.
type Stack []uintptr

// GenerateStackTrace generates a Stack.
func GenerateStackTrace() Stack {
	pcs := make([]uintptr, stackDepth)

	n := runtime.Callers(3, pcs)
	return pcs[0:n]
}
