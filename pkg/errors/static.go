package errors

// Abort is similar to a break in a for-loop.
// It stops the execution of a command silently, without producing neither a
// logged error nor a message to the calling user.
//
// It is intended to be used, if the user signals to cancel a command early
// and is therefore just a signaling error, rather than an actual exception.
var Abort = New("abort")
