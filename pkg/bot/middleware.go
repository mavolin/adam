package bot

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

// ErrMiddleware is the error returned if a middleware given to
// MiddlewareManager.TryAddMiddleware is not a valid middleware type.
var ErrMiddleware = errors.New("the passed function does not resemble a valid middleware")

type (
	// CommandFunc is the signature of the Invoke function of a plugin.Command,
	// without the reply (interface{}) return.
	CommandFunc func(s *state.State, ctx *plugin.Context) error
	// Middleware is a middleware function.
	Middleware func(next CommandFunc) CommandFunc
)

// Middlewarer is an abstraction of a plugin that provides middlewares.
// If a plugin does not implement the interface, it will be assumed that the
// plugin does not provide any middlewares.
type Middlewarer interface {
	// Middlewares returns the MiddlewareFuncs of the plugin.
	Middlewares() []Middleware
}

// MiddlewareManager is a struct that can be embedded in commands and modules
// to provide middleware capabilities.
// It implements Middlewarer.
//
// MiddlewareManagers zero value is an empty MiddlewareManager.
type MiddlewareManager struct {
	middlewares []Middleware
}

// TryAddMiddleware adds the passed middleware to the MiddlewareManager.
// If the middleware's type is invalid, TryAddMiddleware will return
// ErrMiddleware.
//
// Valid middleware types are:
//	• func(*state.State, interface{})
//	• func(*state.State, interface{}) error
//	• func(*state.State, *event.Base)
//	• func(*state.State, *event.Base) error
//	• func(*state.State, *state.MessageCreateEvent)
//	• func(*state.State, *state.MessageCreateEvent) error
//	• func(*state.State, *state.MessageUpdateEvent)
//	• func(*state.State, *state.MessageUpdateEvent) error
//	• func(next CommandFunc) CommandFunc
//nolint:funlen,gocognit
func (m *MiddlewareManager) TryAddMiddleware(f interface{}) error {
	var mf Middleware

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
	case func(*state.State, *event.Base):
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				f(s, ctx.Base)
				return next(s, ctx)
			}
		}
	case func(*state.State, *event.Base) error:
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				err := f(s, ctx.Base)
				if err != nil {
					return err
				}

				return next(s, ctx)
			}
		}
	case func(*state.State, *event.MessageCreate):
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				if !ctx.Message.EditedTimestamp.IsValid() {
					f(s, newMessageCreateEvent(ctx))
				}

				return next(s, ctx)
			}
		}
	case func(*state.State, *event.MessageCreate) error:
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
	case func(*state.State, *event.MessageUpdate):
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				if ctx.Message.EditedTimestamp.IsValid() {
					f(s, newMessageUpdateEvent(ctx))
				}

				return next(s, ctx)
			}
		}
	case func(*state.State, *event.MessageUpdate) error:
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
	case Middleware:
		mf = f
	default:
		return errors.WithStack(ErrMiddleware)
	}

	m.middlewares = append(m.middlewares, mf)
	return nil
}

// newMessageCreateEvent creates a new state.MessageCreateEvent from the passed
// *plugin.Context.
func newMessageCreateEvent(ctx *plugin.Context) *event.MessageCreate {
	return &event.MessageCreate{
		MessageCreateEvent: &gateway.MessageCreateEvent{
			Message: ctx.Message,
			Member:  ctx.Member,
		},
		Base: ctx.Base,
	}
}

// newMessageUpdateEvent creates a new state.MessageUpdateEvent from the passed
// *plugin.Context.
func newMessageUpdateEvent(ctx *plugin.Context) *event.MessageUpdate {
	return &event.MessageUpdate{
		MessageUpdateEvent: &gateway.MessageUpdateEvent{
			Message: ctx.Message,
			Member:  ctx.Member,
		},
		Base: ctx.Base,
	}
}

// AddMiddleware is the same as TryAddMiddleware, but panics if TryAddMiddleware
// returns an error.
func (m *MiddlewareManager) AddMiddleware(f interface{}) {
	if err := m.TryAddMiddleware(f); err != nil {
		panic(err)
	}
}

// Middlewares returns the middlewares of the MiddlewareManager.
func (m *MiddlewareManager) Middlewares() []Middleware {
	if m == nil {
		return nil
	}

	return m.middlewares
}
