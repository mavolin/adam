// Package bot provides the Bot handling all commands.
package bot

type Bot struct {
	Options
	*MiddlewareManager
}

type Options struct {
	Token string
}

func New(o Options) *Bot {
	return &Bot{
		Options:           o,
		MiddlewareManager: new(MiddlewareManager),
	}
}
