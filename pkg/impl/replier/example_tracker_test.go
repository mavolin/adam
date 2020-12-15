package replier

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/plugin"
)

func ExampleTracker() {
	b, _ := bot.New(bot.Options{Token: "abc"})

	// A tracker is typically added to a Context through a middleware.
	// Make sure that the middleware replacing the default replier is executed
	// before any middlewares that could send replies.

	b.MustAddMiddleware(func(next bot.CommandFunc) bot.CommandFunc {
		return func(s *state.State, ctx *plugin.Context) error {
			t := NewTracker(s)
			ctx.Replier = t // replace the default replier

			err := next(s, ctx)
			if err != nil {
				return err
			}

			// do something with t.DMs() and t.GuildMessages()

			return nil
		}
	})
}
