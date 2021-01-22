package bot

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

// ErrMiddleware is the error returned if a middleware given to
// MiddlewareManager.AddMiddleware is not a valid middleware type.
var ErrMiddleware = errors.New("the passed function does not resemble a valid middleware")

type (
	// CommandFunc is the signature of the Invoke function of a plugin.Command,
	// without the reply (interface{}) return.
	CommandFunc func(s *state.State, ctx *plugin.Context) error
	// MiddlewareFunc is the function of a middleware.
	MiddlewareFunc func(next CommandFunc) CommandFunc
)

// Middlewarer is an abstraction of a plugin that provides middlewares.
// If a plugin does not implement the interface, it will be assumed that the
// plugin does not provide any middlewares.
type Middlewarer interface {
	// Middlewares returns a copy of the MiddlewareFuncs of the plugin.
	Middlewares() []MiddlewareFunc
}

// MiddlewareManager is a struct that can be embedded in commands and modules
// to provide middleware capabilities.
// It implements Middlewarer.
//
// MiddlewareManagers zero value is an empty MiddlewareManager.
type MiddlewareManager struct {
	middlewares []MiddlewareFunc
}

// AddMiddleware adds the passed middleware to the MiddlewareManager.
// If the middleware's type is invalid, AddMiddleware will return
// ErrMiddleware.
//
// Valid middleware types are:
//		• func(*state.State, interface{})
//		• func(*state.State, interface{}) error
//		• func(*state.State, *state.Base)
//		• func(*state.State, *state.Base) error
//		• func(*state.State, *state.MessageCreateEvent)
//		• func(*state.State, *state.MessageCreateEvent) error
//		• func(*state.State, *state.MessageUpdateEvent)
//		• func(*state.State, *state.MessageUpdateEvent) error
//		• func(next CommandFunc) CommandFunc
func (m *MiddlewareManager) AddMiddleware(f interface{}) error { //nolint:funlen,gocognit
	var mf MiddlewareFunc

	switch f := f.(type) {
	case func(*state.State, interface{}):
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				if !ctx.Message.EditedTimestamp.IsValid() {
					f(s, newMessageCreateEvent(ctx))
				} else {
					f(s, newMessageUpdateEvent(ctx))
				}

				return next(s, ctx)
			}
		}
	case func(*state.State, interface{}) error:
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				if !ctx.Message.EditedTimestamp.IsValid() {
					if err := f(s, newMessageCreateEvent(ctx)); err != nil {
						return err
					}
				} else {
					if err := f(s, newMessageUpdateEvent(ctx)); err != nil {
						return err
					}
				}

				return next(s, ctx)
			}
		}
	case func(*state.State, *state.Base):
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				f(s, ctx.Base)
				return next(s, ctx)
			}
		}
	case func(*state.State, *state.Base) error:
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				err := f(s, ctx.Base)
				if err != nil {
					return err
				}

				return next(s, ctx)
			}
		}
	case func(*state.State, *state.MessageCreateEvent):
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				if !ctx.Message.EditedTimestamp.IsValid() {
					f(s, newMessageCreateEvent(ctx))
				}

				return next(s, ctx)
			}
		}
	case func(*state.State, *state.MessageCreateEvent) error:
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				if !ctx.Message.EditedTimestamp.IsValid() {
					if err := f(s, newMessageCreateEvent(ctx)); err != nil {
						return err
					}
				}

				return next(s, ctx)
			}
		}
	case func(*state.State, *state.MessageUpdateEvent):
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				if ctx.Message.EditedTimestamp.IsValid() {
					f(s, newMessageUpdateEvent(ctx))
				}

				return next(s, ctx)
			}
		}
	case func(*state.State, *state.MessageUpdateEvent) error:
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				if ctx.Message.EditedTimestamp.IsValid() {
					if err := f(s, newMessageUpdateEvent(ctx)); err != nil {
						return err
					}
				}

				return next(s, ctx)
			}
		}
	case func(next CommandFunc) CommandFunc:
		mf = f
	case MiddlewareFunc:
		mf = f
	default:
		return errors.WithStack(ErrMiddleware)
	}

	m.middlewares = append(m.middlewares, mf)
	return nil
}

// MustAddMiddleware is the same as AddMiddleware, but panics if AddMiddleware
// returns an error.
func (m *MiddlewareManager) MustAddMiddleware(f interface{}) {
	err := m.AddMiddleware(f)
	if err != nil {
		panic(err)
	}
}

// Middlewares returns a copy of the middlewares of the manager.
func (m *MiddlewareManager) Middlewares() (cp []MiddlewareFunc) {
	return m.middlewares
}
