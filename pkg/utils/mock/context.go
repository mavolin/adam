package mock

import (
	"reflect"
	"sync"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/adam/pkg/plugin"
)

// DiscordDataProvider is the mock implementation of
// plugin.DiscordDataProvider.
type DiscordDataProvider struct {
	ChannelReturn *discord.Channel
	ChannelError  error

	GuildReturn *discord.Guild
	GuildError  error

	SelfReturn *discord.Member
	SelfError  error
}

var _ plugin.DiscordDataProvider = DiscordDataProvider{}

func (d DiscordDataProvider) ChannelAsync() func() (*discord.Channel, error) {
	return func() (*discord.Channel, error) {
		return d.ChannelReturn, d.ChannelError
	}
}

func (d DiscordDataProvider) GuildAsync() func() (*discord.Guild, error) {
	return func() (*discord.Guild, error) {
		return d.GuildReturn, d.GuildError
	}
}

func (d DiscordDataProvider) SelfAsync() func() (*discord.Member, error) {
	return func() (*discord.Member, error) {
		return d.SelfReturn, d.SelfError
	}
}

type ErrorHandler struct {
	t *testing.T

	mut          sync.Mutex
	expectErr    []error
	expectSilent []error
}

var _ plugin.ErrorHandler = new(ErrorHandler)

func NewErrorHandler(t *testing.T) *ErrorHandler {
	h := &ErrorHandler{t: t}
	t.Cleanup(h.eval)

	return h
}

func (h *ErrorHandler) ExpectError(err error) *ErrorHandler {
	h.expectErr = append(h.expectErr, err)
	return h
}

func (h *ErrorHandler) ExpectSilentError(err error) *ErrorHandler {
	h.expectSilent = append(h.expectSilent, err)
	return h
}

func (h *ErrorHandler) HandleError(err error) {
	h.mut.Lock()
	defer h.mut.Unlock()

	for i, expect := range h.expectErr {
		if reflect.DeepEqual(err, expect) {
			h.expectErr = append(h.expectErr[:i], h.expectErr[i+1:]...)
			return
		}

		err2 := err

		type unwrapper interface{ Unwrap() error }

		for uerr, ok := err2.(unwrapper); ok; uerr, ok = err2.(unwrapper) { //nolint:errorlint
			err2 = uerr.Unwrap()

			if reflect.DeepEqual(err2, expect) {
				h.expectErr = append(h.expectErr[:i], h.expectErr[i+1:]...)
				return
			}
		}
	}

	h.t.Errorf("unexpected call to plugin.ErrorHandler.HandleError: %+v", err)
}

func (h *ErrorHandler) HandleErrorSilently(err error) {
	h.mut.Lock()
	defer h.mut.Unlock()

	for i, expect := range h.expectSilent {
		if reflect.DeepEqual(err, expect) {
			h.expectSilent = append(h.expectSilent[:i], h.expectSilent[i+1:]...)
			return
		}

		err2 := err

		type unwrapper interface{ Unwrap() error }

		for uerr, ok := err2.(unwrapper); ok; uerr, ok = err2.(unwrapper) { //nolint:errorlint
			err2 = uerr.Unwrap()

			if reflect.DeepEqual(err2, expect) {
				h.expectSilent = append(h.expectSilent[:i], h.expectSilent[i+1:]...)
				return
			}
		}
	}

	h.t.Errorf("unexpected call to plugin.ErrorHandler.HandleErrorSilently: %+v", err)
}

func (h *ErrorHandler) eval() {
	if len(h.expectErr) > 0 {
		h.t.Errorf("there are unhandled errors: %+v", h.expectErr)
	}

	if len(h.expectSilent) > 0 {
		h.t.Errorf("there are unhandled silent errors: %+v", h.expectSilent)
	}
}
