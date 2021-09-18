package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mockplugin "github.com/mavolin/adam/internal/mock/plugin"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/mock"
)

type asError struct {
	as Error
}

func (a *asError) As(target interface{}) bool {
	if err, ok := target.(*Error); ok {
		*err = a.as
		return true
	}

	return false
}

func (a *asError) Error() string {
	return "asError"
}

func TestWithStack(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		err := WithStack(nil)
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")

		ierr := WithStack(cause).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")

		ierr := WithStack(cause)

		require.Same(t, cause, ierr)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		cause := &asError{as: NewInformationalError("abc")}

		err := WithStack(cause)

		assert.Same(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		ierr := WithStack(cause).(*InternalError)

		assert.Equal(t, cause, ierr.cause)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})
}

func TestSilent(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := Silent(nil)
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")

		ierr := Silent(cause)

		require.Same(t, cause, ierr)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")

		ierr := Silent(cause)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		err := Silent(&asError{as: NewInformationalError("abc")})
		assert.Nil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		ierr := Silent(cause)

		assert.Equal(t, cause, ierr.cause)
		assert.Nil(t, ierr.desc)
	})
}

func TestMustInternal(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := MustInternal(nil)
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")

		ierr := MustInternal(cause)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")

		ierr := MustInternal(cause)

		require.Same(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, cause.desc, ierr.desc)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		ierr := MustInternal(cause)

		assert.Equal(t, cause, ierr.cause)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})
}

func TestMustSilent(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := MustSilent(nil)
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")

		ierr := MustSilent(cause)

		require.Same(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")

		ierr := MustSilent(cause)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		ierr := MustSilent(cause)

		assert.Equal(t, cause, ierr.cause)
		assert.Nil(t, ierr.desc)
	})
}

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := Wrap(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")
		expectMsg := "def"

		ierr := Wrap(cause, expectMsg).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")
		expectMsg := "def"

		ierr := Wrap(cause, expectMsg).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, cause.desc, ierr.desc)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		cause := &asError{as: NewInformationalError("abc")}

		ierr := Wrap(cause, "def")

		assert.Same(t, cause.as, ierr)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		expectMsg := "def"
		ierr := Wrap(cause, expectMsg).(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause}, ierr.cause)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})
}

func TestWrapSilent(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := Wrap(nil, "abc")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")

		expectMsg := "def"
		ierr := WrapSilent(cause, expectMsg).(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")

		expectMsg := "def"
		ierr := WrapSilent(cause, expectMsg).(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		err := WrapSilent(&asError{as: NewInformationalError("abc")}, "def")
		assert.Nil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		expectMsg := "def"
		ierr := WrapSilent(cause, expectMsg).(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause}, ierr.cause)
		assert.Nil(t, ierr.desc)
	})
}

func TestWrapf(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := Wrapf(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")

		expectMsg := "def ghi"

		ierr := Wrapf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")

		expectMsg := "def ghi"

		ierr := Wrapf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, cause.desc, ierr.desc)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		cause := &asError{as: NewInformationalError("abc")}

		err := Wrapf(cause, "def %s", "ghi")

		assert.Same(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		expectMsg := "def ghi"
		ierr := Wrapf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause}, ierr.cause)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})
}

func TestWrapSilentf(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := WrapSilentf(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")

		expectMsg := "def ghi"

		ierr := WrapSilentf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")

		expectMsg := "def ghi"

		ierr := WrapSilentf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		err := WrapSilentf(&asError{as: NewInformationalError("abc")}, "def %s", "ghi")
		assert.Nil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		expectMsg := "def ghi"
		ierr := WrapSilentf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause}, ierr.cause)
		assert.Nil(t, ierr.desc)
	})
}

func TestWithDescription(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := WithDescription(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")

		expectDesc := "def"
		ierr := WithDescription(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")

		expectDesc := "def"
		ierr := WithDescription(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		cause := &asError{as: NewInformationalError("abc")}

		err := WithDescription(cause, "def")

		assert.Same(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		expectDesc := "def"
		ierr := WithDescription(cause, expectDesc).(*InternalError)

		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause, ierr.cause)
	})
}

func TestWithDescriptionf(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := WithDescriptionf(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")

		expectDesc := "def"
		ierr := WithDescription(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")

		expectDesc := "def ghi"
		ierr := WithDescriptionf(cause, "def %s", "ghi").(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		cause := &asError{as: NewInformationalError("abc")}

		err := WithDescriptionf(cause, "def %s", "ghi")

		assert.Same(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		expectDesc := "def ghi"

		ierr := WithDescriptionf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause, ierr.cause)
	})
}

func TestWithDescriptionl(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := WithDescriptionl(nil, nil)
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")

		expectDesc := i18n.NewTermConfig("def")
		ierr := WithDescriptionl(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, expectDesc, ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")

		expectDesc := i18n.NewTermConfig("def")

		ierr := WithDescriptionl(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, expectDesc, ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		cause := &asError{as: NewInformationalError("abc")}

		err := WithDescriptionl(cause, i18n.NewTermConfig("def"))

		assert.Same(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		expectDesc := i18n.NewTermConfig("def")
		ierr := WithDescriptionl(cause, expectDesc).(*InternalError)

		assert.Equal(t, expectDesc, ierr.desc)
		assert.Equal(t, cause, ierr.cause)
	})
}

func TestWithDescriptionlt(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		err := WithDescriptionlt(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewSilent("abc")

		var expectDesc i18n.Term = "def"
		ierr := WithDescriptionlt(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, expectDesc.AsConfig(), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		t.Parallel()

		cause := NewWithStack("abc")

		var expectDesc i18n.Term = "def"
		ierr := WithDescriptionlt(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, expectDesc.AsConfig(), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		cause := &asError{as: NewInformationalError("abc")}

		err := WithDescriptionlt(cause, "def")

		assert.Equal(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cause := New("abc")

		var expectDesc i18n.Term = "def"
		ierr := WithDescriptionlt(cause, expectDesc).(*InternalError)

		assert.Equal(t, expectDesc.AsConfig(), ierr.desc)
		assert.Equal(t, cause, ierr.cause)
	})
}

func TestInternalError_Description(t *testing.T) {
	t.Parallel()

	t.Run("string description", func(t *testing.T) {
		t.Parallel()

		expect := "abc"

		err := WithDescription(New("def"), expect)

		actual := err.(*InternalError).Description(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})

	t.Run("localized description", func(t *testing.T) {
		t.Parallel()

		var term i18n.Term = "abc"

		expect := "def"

		l := mock.
			NewLocalizer(t).
			On(term, expect).
			Build()

		err := WithDescriptionlt(New("ghi"), term)

		actual := err.(*InternalError).Description(l)
		assert.Equal(t, expect, actual)
	})

	t.Run("no description", func(t *testing.T) {
		t.Parallel()

		ierr := NewSilent("abc")
		assert.Empty(t, ierr.Description(i18n.NewFallbackLocalizer()))
	})
}

func TestInternalError_Handle(t *testing.T) {
	t.Parallel()

	t.Run("silent", func(t *testing.T) {
		t.Parallel()

		_, s := state.NewMocker(t)

		ctx := &plugin.Context{
			Message: discord.Message{ChannelID: 123},
			Localizer: mock.NewLocalizer(t).
				On(internalErrorTitle.Term, "abc").
				Build(),
			InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{Name: "abc"}),
			Replier:        mockplugin.NewWrappedReplier(s, 123, 0),
		}

		e := NewSilent("abc")

		err := e.Handle(s, ctx)
		require.NoError(t, err, "InternalError.Handle should never return an error")
	})

	t.Run("non-silent", func(t *testing.T) {
		t.Parallel()

		expectDesc := "abc"

		m, s := state.NewMocker(t)

		ctx := &plugin.Context{
			Message: discord.Message{ChannelID: 123},
			Localizer: mock.NewLocalizer(t).
				On(internalErrorTitle.Term, "abc").
				Build(),
			InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{Name: "abc"}),
			Replier:        mockplugin.NewWrappedReplier(s, 123, 0),
		}

		embed, err := NewErrorEmbed().
			WithTitlelt(internalErrorTitle.Term).
			WithDescription(expectDesc).
			Build(ctx.Localizer)
		require.NoError(t, err)

		m.SendEmbeds(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := WithDescription(New(""), expectDesc)

		err = e.(*InternalError).Handle(s, ctx)
		require.NoError(t, err, "InternalError.Handle should never return an error")
	})
}
