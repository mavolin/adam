package errors

import (
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/locutil"
)

// ThrottlingError is the error returned if a command gets throttled.
// It contains a description about when the command will become available
// again.
type ThrottlingError struct {
	// description of the error
	desc locutil.Text
}

// NewThrottlingError creates a new ThrottlingError with the passed
// description.
func NewThrottlingError(description string) *ThrottlingError {
	return &ThrottlingError{
		desc: locutil.NewStaticText(description),
	}
}

// NewThrottlingErrorl creates a new ThrottlingError using the message
// generated from the passed localization.Config as description.
func NewThrottlingErrorl(description localization.Config) *ThrottlingError {
	return &ThrottlingError{
		desc: locutil.NewLocalizedText(description),
	}
}

// NewThrottlingErrorlt creates a new ThrottlingError using the message
// generated from the passed term as description.
func NewThrottlingErrorlt(description localization.Term) *ThrottlingError {
	return NewThrottlingErrorl(description.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *ThrottlingError) Description(l *localization.Localizer) (string, error) {
	return e.desc.Get(l)
}

func (e *ThrottlingError) Error() string { return "throttling error" }

// Handle sends an info embed with the description of the ThrottlingError.
func (e *ThrottlingError) Handle(_ *state.State, ctx *plugin.Context) error {
	desc, err := e.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := InfoEmbed.Clone().
		WithDescription(desc)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return err
}
