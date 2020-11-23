package errors //nolint: dupl

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
)

// ThrottlingError is the error returned if a command gets throttled.
// It contains a description about when the command will become available
// again.
type ThrottlingError struct {
	// description of the error
	desc *i18nutil.Text
}

var _ Interface = new(ThrottlingError)

// NewThrottlingError creates a new ThrottlingError with the passed
// description.
func NewThrottlingError(description string) *ThrottlingError {
	return &ThrottlingError{
		desc: i18nutil.NewText(description),
	}
}

// NewThrottlingErrorl creates a new ThrottlingError using the message
// generated from the passed i18n.Config as description.
func NewThrottlingErrorl(description *i18n.Config) *ThrottlingError {
	return &ThrottlingError{
		desc: i18nutil.NewTextl(description),
	}
}

// NewThrottlingErrorlt creates a new ThrottlingError using the message
// generated from the passed term as description.
func NewThrottlingErrorlt(description i18n.Term) *ThrottlingError {
	return NewThrottlingErrorl(description.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *ThrottlingError) Description(l *i18n.Localizer) (string, error) {
	return e.desc.Get(l)
}

func (e *ThrottlingError) Error() string { return "throttling error" }

// Handle handles the ThrottlingError.
// By default it sends an info embed with the description of the
// ThrottlingError.
func (e *ThrottlingError) Handle(s *state.State, ctx *plugin.Context) {
	HandleThrottlingError(e, s, ctx)
}

var HandleThrottlingError = func(terr *ThrottlingError, s *state.State, ctx *plugin.Context) {
	desc, err := terr.Description(ctx.Localizer)
	if err != nil {
		return
	}

	embed := InfoEmbed.Clone().
		WithDescription(desc)

	_, _ = ctx.ReplyEmbedBuilder(embed)
}
