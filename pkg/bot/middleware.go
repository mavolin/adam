package bot

import (
	"sync"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

var ErrNotAMiddleware = errors.New("the passed func does not resemble a valid middleware")

type (
	// CommandFunc is the Invoke function of a plugin.Command.
	// It doesn't have the reply return value, as the reply will already be
	// sent, after the first next is called.
	CommandFunc func(s *state.State, ctx *plugin.Context) error
	// MiddlewareFunc is the function of a middleware.
	MiddlewareFunc func(next CommandFunc) CommandFunc
)

// Middlewarer is an abstraction of a plugin that provides middlewares.
// If a plugin does not implement the interface, it will be assumed, that the
// plugin does not provide any middlewares.
type Middlewarer interface {
	// Middlewares returns a copy of the MiddlewareFuncs of the plugin.
	Middlewares() []MiddlewareFunc
}

// MiddlewareManager is a struct that can be embedded in commands and modules
// to provide middleware capabilities.
// It implements Middlewarer.
//
// MiddlewareManagers zero value is guaranteed to function properly.
type MiddlewareManager struct {
	middlewares []MiddlewareFunc
	mutex       sync.RWMutex
}

// AddMiddleware adds the passed middleware to the MiddlewareManager.
// If the middleware's type is invalid, AddMiddleware will return
// ErrNotAMiddleware.
//
// Supported types are:
//		• func(*state.State, interface{})
//		• func(*state.State, interface{}) error
//		• func(*state.State, *state.Base)
//		• func(*state.State, *state.Base) error
//		• func(*state.State, *state.MessageCreateEvent)
//		• func(*state.State, *state.MessageCreateEvent) error
//		• func(next CommandFunc) CommandFunc
func (m *MiddlewareManager) AddMiddleware(f interface{}) error {
	var mf MiddlewareFunc

	switch f := f.(type) {
	case func(*state.State, interface{}):
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				f(s, ctx.MessageCreateEvent)
				return next(s, ctx)
			}
		}
	case func(*state.State, interface{}) error:
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				err := f(s, ctx.MessageCreateEvent)
				if err != nil {
					return err
				}

				return next(s, ctx)
			}
		}
	case func(*state.State, *state.Base):
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				f(s, ctx.MessageCreateEvent.Base)
				return next(s, ctx)
			}
		}
	case func(*state.State, *state.Base) error:
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				err := f(s, ctx.MessageCreateEvent.Base)
				if err != nil {
					return err
				}

				return next(s, ctx)
			}
		}
	case func(*state.State, *state.MessageCreateEvent):
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				f(s, ctx.MessageCreateEvent)
				return next(s, ctx)
			}
		}
	case func(*state.State, *state.MessageCreateEvent) error:
		mf = func(next CommandFunc) CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				err := f(s, ctx.MessageCreateEvent)
				if err != nil {
					return err
				}

				return next(s, ctx)
			}
		}
	case func(next CommandFunc) CommandFunc:
		mf = f
	case MiddlewareFunc:
		mf = f
	default:
		return errors.WithStack(ErrNotAMiddleware)
	}

	m.mutex.Lock()
	m.middlewares = append(m.middlewares, mf)
	m.mutex.Unlock()

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
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.middlewares) == 0 {
		return nil
	}

	cp = make([]MiddlewareFunc, len(m.middlewares))
	copy(cp, m.middlewares)

	return
}
