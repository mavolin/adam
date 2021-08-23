package errors

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"
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
	t.Run("nil", func(t *testing.T) {
		err := WithStack(nil)
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc").(*InternalError)

		ierr := WithStack(cause).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc")

		ierr := WithStack(cause)

		require.Same(t, cause, ierr)
	})

	t.Run("Error", func(t *testing.T) {
		cause := &asError{as: NewInformationalError("abc")}

		err := WithStack(cause)

		assert.Same(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		ierr := WithStack(cause).(*InternalError)

		assert.Equal(t, cause, ierr.cause)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})
}

func TestSilent(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Silent(nil)
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc")

		ierr := Silent(cause)

		require.Same(t, cause, ierr)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		ierr := Silent(cause).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("Error", func(t *testing.T) {
		err := Silent(&asError{as: NewInformationalError("abc")})
		assert.Nil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		ierr := Silent(cause).(*InternalError)

		assert.Equal(t, cause, ierr.cause)
		assert.Nil(t, ierr.desc)
	})
}

func TestMustInternal(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := MustInternal(nil)
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc").(*InternalError)

		ierr := MustInternal(cause).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		ierr := MustInternal(cause).(*InternalError)

		require.Same(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, cause.desc, ierr.desc)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		ierr := MustInternal(cause).(*InternalError)

		assert.Equal(t, cause, ierr.cause)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})
}

func TestMustSilent(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := MustSilent(nil)
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc").(*InternalError)

		ierr := MustSilent(cause).(*InternalError)

		require.Same(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		ierr := MustSilent(cause).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		ierr := MustSilent(cause).(*InternalError)

		assert.Equal(t, cause, ierr.cause)
		assert.Nil(t, ierr.desc)
	})
}

func TestWrap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrap(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc").(*InternalError)
		expectMsg := "def"

		ierr := Wrap(cause, expectMsg).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)
		expectMsg := "def"

		ierr := Wrap(cause, expectMsg).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, cause.desc, ierr.desc)
	})

	t.Run("Error", func(t *testing.T) {
		cause := &asError{as: NewInformationalError("abc")}

		ierr := Wrap(cause, "def")

		assert.Same(t, cause.as, ierr)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		expectMsg := "def"
		ierr := Wrap(cause, expectMsg).(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause}, ierr.cause)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})
}

func TestWrapSilent(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrap(nil, "abc")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc").(*InternalError)

		expectMsg := "def"
		ierr := WrapSilent(cause, expectMsg).(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		expectMsg := "def"
		ierr := WrapSilent(cause, expectMsg).(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("Error", func(t *testing.T) {
		err := WrapSilent(&asError{as: NewInformationalError("abc")}, "def")
		assert.Nil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		expectMsg := "def"
		ierr := WrapSilent(cause, expectMsg).(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause}, ierr.cause)
		assert.Nil(t, ierr.desc)
	})
}

func TestWrapf(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrapf(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc").(*InternalError)

		expectMsg := "def ghi"

		ierr := Wrapf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		expectMsg := "def ghi"

		ierr := Wrapf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Equal(t, cause.desc, ierr.desc)
	})

	t.Run("Error", func(t *testing.T) {
		cause := &asError{as: NewInformationalError("abc")}

		err := Wrapf(cause, "def %s", "ghi")

		assert.Same(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		expectMsg := "def ghi"
		ierr := Wrapf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause}, ierr.cause)
		assert.Equal(t, defaultInternalDesc, ierr.desc)
	})
}

func TestWrapSilentf(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WrapSilentf(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc").(*InternalError)

		expectMsg := "def ghi"

		ierr := WrapSilentf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		expectMsg := "def ghi"

		ierr := WrapSilentf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause.cause}, ierr.cause)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
		assert.Nil(t, ierr.desc)
	})

	t.Run("Error", func(t *testing.T) {
		err := WrapSilentf(&asError{as: NewInformationalError("abc")}, "def %s", "ghi")
		assert.Nil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		expectMsg := "def ghi"
		ierr := WrapSilentf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, &messageError{msg: expectMsg, cause: cause}, ierr.cause)
		assert.Nil(t, ierr.desc)
	})
}

func TestWithDescription(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WithDescription(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc").(*InternalError)

		expectDesc := "def"
		ierr := WithDescription(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		expectDesc := "def"
		ierr := WithDescription(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("Error", func(t *testing.T) {
		cause := &asError{as: NewInformationalError("abc")}

		err := WithDescription(cause, "def")

		assert.Same(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		expectDesc := "def"
		ierr := WithDescription(cause, expectDesc).(*InternalError)

		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause, ierr.cause)
	})
}

func TestWithDescriptionf(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WithDescriptionf(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc").(*InternalError)

		expectDesc := "def"
		ierr := WithDescription(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		expectDesc := "def ghi"
		ierr := WithDescriptionf(cause, "def %s", "ghi").(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("Error", func(t *testing.T) {
		cause := &asError{as: NewInformationalError("abc")}

		err := WithDescriptionf(cause, "def %s", "ghi")

		assert.Same(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		expectDesc := "def ghi"

		ierr := WithDescriptionf(cause, "def %s", "ghi").(*InternalError)

		assert.Equal(t, i18n.NewStaticConfig(expectDesc), ierr.desc)
		assert.Equal(t, cause, ierr.cause)
	})
}

func TestWithDescriptionl(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WithDescriptionl(nil, nil)
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc").(*InternalError)

		expectDesc := i18n.NewTermConfig("def")
		ierr := WithDescriptionl(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, expectDesc, ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		expectDesc := i18n.NewTermConfig("def")

		ierr := WithDescriptionl(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, expectDesc, ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("Error", func(t *testing.T) {
		cause := &asError{as: NewInformationalError("abc")}

		err := WithDescriptionl(cause, i18n.NewTermConfig("def"))

		assert.Same(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		expectDesc := i18n.NewTermConfig("def")
		ierr := WithDescriptionl(cause, expectDesc).(*InternalError)

		assert.Equal(t, expectDesc, ierr.desc)
		assert.Equal(t, cause, ierr.cause)
	})
}

func TestWithDescriptionlt(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := WithDescriptionlt(nil, "")
		assert.Nil(t, err)
	})

	t.Run("silent internal error", func(t *testing.T) {
		cause := NewSilent("abc").(*InternalError)

		var expectDesc i18n.Term = "def"
		ierr := WithDescriptionlt(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, expectDesc.AsConfig(), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("non-silent internal error", func(t *testing.T) {
		cause := NewWithStack("abc").(*InternalError)

		var expectDesc i18n.Term = "def"
		ierr := WithDescriptionlt(cause, expectDesc).(*InternalError)

		require.NotSame(t, cause, ierr)
		assert.Equal(t, cause.cause, ierr.cause)
		assert.Equal(t, expectDesc.AsConfig(), ierr.desc)
		assert.Equal(t, cause.stackTrace, ierr.stackTrace)
	})

	t.Run("Error", func(t *testing.T) {
		cause := &asError{as: NewInformationalError("abc")}

		err := WithDescriptionlt(cause, "def")

		assert.Equal(t, cause.as, err)
	})

	t.Run("success", func(t *testing.T) {
		cause := New("abc")

		var expectDesc i18n.Term = "def"
		ierr := WithDescriptionlt(cause, expectDesc).(*InternalError)

		assert.Equal(t, expectDesc.AsConfig(), ierr.desc)
		assert.Equal(t, cause, ierr.cause)
	})
}

func TestInternalError_Description(t *testing.T) {
	t.Run("string description", func(t *testing.T) {
		expect := "abc"

		err := WithDescription(New("def"), expect)

		actual := err.(*InternalError).Description(i18n.NewFallbackLocalizer())
		assert.Equal(t, expect, actual)
	})

	t.Run("localized description", func(t *testing.T) {
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
		ierr := NewSilent("abc").(*InternalError)
		assert.Empty(t, ierr.Description(i18n.NewFallbackLocalizer()))
	})
}

func TestInternalError_Handle(t *testing.T) {
	t.Run("silent", func(t *testing.T) {
		m, s := state.NewMocker(t)
		defer m.Eval()

		ctx := &plugin.Context{
			Message: discord.Message{ChannelID: 123},
			Localizer: mock.NewLocalizer(t).
				On(internalErrorTitle.Term, "abc").
				Build(),
			InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{Name: "abc"}),
			Replier:        mockplugin.NewWrappedReplier(s, 123, 0),
		}

		e := NewSilent("abc")

		err := e.(*InternalError).Handle(s, ctx)
		require.NoError(t, err, "InternalError.Handle should never return an error")
	})

	t.Run("non-silent", func(t *testing.T) {
		expectDesc := "abc"

		m, s := state.NewMocker(t)
		defer m.Eval()

		ctx := &plugin.Context{
			Message: discord.Message{ChannelID: 123},
			Localizer: mock.NewLocalizer(t).
				On(internalErrorTitle.Term, "abc").
				Build(),
			InvokedCommand: mock.ResolveCommand(plugin.BuiltInSource, mock.Command{Name: "abc"}),
			Replier:        mockplugin.NewWrappedReplier(s, 123, 0),
		}

		embed := NewErrorEmbed().
			WithTitlelt(internalErrorTitle.Term).
			WithDescription(expectDesc).
			MustBuild(ctx.Localizer)

		m.SendEmbed(discord.Message{
			ChannelID: ctx.ChannelID,
			Embeds:    []discord.Embed{embed},
		})

		e := WithDescription(New(""), expectDesc)

		err := e.(*InternalError).Handle(s, ctx)
		require.NoError(t, err, "InternalError.Handle should never return an error")
	})
}
