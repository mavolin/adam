package bot

import "fmt"

// ReplyTypeError is the error used if a reply returned by
// plugin.Command.Invoke is not of a supported types
type ReplyTypeError struct {
	Reply interface{}
}

func (r *ReplyTypeError) Error() string {
	return fmt.Sprintf("bot: cannot use %T for reply", r.Reply)
}
