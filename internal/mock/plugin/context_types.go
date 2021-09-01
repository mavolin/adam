package plugin

import (
	"reflect"
	"sync"
	"testing"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

// =============================================================================
// DiscordDataProvider
// =====================================================================================

// DiscordDataProvider is the mock implementation of
// plugin.DiscordDataProvider.
type DiscordDataProvider struct {
	ChannelReturn *discord.Channel
	ChannelError  error

	ParentChannelReturn *discord.Channel
	ParentChannelError  error

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

func (d DiscordDataProvider) ParentChannelAsync() func() (*discord.Channel, error) {
	return func() (*discord.Channel, error) {
		return d.ParentChannelReturn, d.ParentChannelError
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

// =============================================================================
// ErrorHandler
// =====================================================================================

type ErrorHandler struct {
	t *testing.T

	mut          sync.Mutex
	expectErr    []error
	expectSilent []error
}

var _ plugin.ErrorHandler = new(ErrorHandler)

func NewErrorHandler(t *testing.T) *ErrorHandler { //nolint:thelper
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
	h.t.Helper()

	h.mut.Lock()
	defer h.mut.Unlock()

	for i, expect := range h.expectErr {
		if reflect.DeepEqual(err, expect) {
			h.expectErr = append(h.expectErr[:i], h.expectErr[i+1:]...)
			return
		}

		err2 := err

		type unwrapper interface{ Unwrap() error }

		//nolint:errorlint
		for uerr, ok := err2.(unwrapper); ok; uerr, ok = err2.(unwrapper) {
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
	h.t.Helper()

	h.mut.Lock()
	defer h.mut.Unlock()

	for i, expect := range h.expectSilent {
		if reflect.DeepEqual(err, expect) {
			h.expectSilent = append(h.expectSilent[:i], h.expectSilent[i+1:]...)
			return
		}

		err2 := err

		type unwrapper interface{ Unwrap() error }

		//nolint:errorlint
		for uerr, ok := err2.(unwrapper); ok; uerr, ok = err2.(unwrapper) {
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

// =============================================================================
// WrappedReplier
// =====================================================================================

// WrappedReplier is a copy of replier.WrappedReplier.
// This type is not meant to be expose in mock, as outside users will always
// have access to replier.WrappedReplier.
type WrappedReplier struct {
	s         *state.State
	channelID discord.ChannelID

	userID discord.UserID
	dmID   discord.ChannelID
}

func NewWrappedReplier(s *state.State, channelID discord.ChannelID, userID discord.UserID) *WrappedReplier {
	return &WrappedReplier{s: s, channelID: channelID, userID: userID}
}

func (r *WrappedReplier) Reply(_ *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	return r.s.SendMessageComplex(r.channelID, data)
}

func (r *WrappedReplier) ReplyDM(_ *plugin.Context, data api.SendMessageData) (*discord.Message, error) {
	if !r.dmID.IsValid() {
		c, err := r.s.CreatePrivateChannel(r.userID)
		if err != nil {
			return nil, err
		}

		r.dmID = c.ID
	}

	return r.s.SendMessageComplex(r.dmID, data)
}

func (r *WrappedReplier) Edit(
	_ *plugin.Context, messageID discord.MessageID, data api.EditMessageData,
) (*discord.Message, error) {
	return r.s.EditMessageComplex(r.channelID, messageID, data)
}

func (r *WrappedReplier) EditDM(
	_ *plugin.Context, messageID discord.MessageID, data api.EditMessageData,
) (*discord.Message, error) {
	if !r.dmID.IsValid() {
		c, err := r.s.CreatePrivateChannel(r.userID)
		if err != nil {
			return nil, err
		}

		r.dmID = c.ID
	}

	return r.s.EditMessageComplex(r.dmID, messageID, data)
}
