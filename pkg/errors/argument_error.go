package errors //nolint:dupl

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
)

// ArgumentError is the error used if an argument or flag a user supplied is
// invalid.
type ArgumentError struct {
	desc *i18nutil.Text
}

var _ Error = new(ArgumentError)

// NewArgumentError returns a new ArgumentError with the passed
// description.
// The description mustn't be empty for this error to be handled properly.
func NewArgumentError(description string) *ArgumentError {
	return &ArgumentError{desc: i18nutil.NewText(description)}
}

// NewArgumentErrorl returns a new ArgumentError using the passed i18n.Config
// to generate a description.
func NewArgumentErrorl(description *i18n.Config) *ArgumentError {
	return &ArgumentError{desc: i18nutil.NewTextl(description)}
}

// NewArgumentErrorlt returns a new ArgumentError using the passed term to
// generate a description.
func NewArgumentErrorlt(description i18n.Term) *ArgumentError {
	return NewArgumentErrorl(description.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *ArgumentError) Description(l *i18n.Localizer) (string, error) {
	return e.desc.Get(l)
}

func (e *ArgumentError) Error() string {
	return "argument error"
}

// Handle handles the ArgumentError.
// By default it sends an error Embed containing a description of which
// arg/flag was faulty in the channel the command was sent in.
func (e *ArgumentError) Handle(s *state.State, ctx *plugin.Context) error {
	return HandleArgumentError(e, s, ctx)
}

var HandleArgumentError = func(aerr *ArgumentError, _ *state.State, ctx *plugin.Context) error {
	desc, err := aerr.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := ErrorEmbed.Clone().
		WithDescription(desc)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return err
}
