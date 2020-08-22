package errors

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
)

// UserInfo is less sever error on the user-side.
// The error will reported to the user via a message containing a detailed
// description of the problem.
// It won't be logged or captured by sentry.
type UserInfo struct {
	// description of the info, either is set
	descString string
	descConfig localization.Config
}

// NewUserInfo creates a new UserInfo using the passed description.
// The description mustn't be empty for this error to be handled properly.
func NewUserInfo(desc string) *UserInfo {
	return &UserInfo{
		descString: desc,
	}
}

// NewUserInfol creates a new UserInfo using the message generated from the
// passed localization.Config.
func NewUserInfol(description localization.Config) *UserInfo {
	return &UserInfo{
		descConfig: description,
	}
}

// NewUserInfolt creates a new UserInfo using the message generated from the
// passed term.
func NewUserInfolt(description localization.Term) *UserInfo {
	return NewUserInfol(description.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (i *UserInfo) Description(l *localization.Localizer) (string, error) {
	if i.descString != "" {
		return i.descString, nil
	}

	return l.Localize(i.descConfig)
}

func (i *UserInfo) Error() string { return "user info" }

// Handle sends an info embed with the description of the UserInfo.
func (i *UserInfo) Handle(_ *state.State, ctx *plugin.Context) (err error) {
	desc, err := i.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := newInfoEmbedBuilder(ctx.Localizer).
		WithDescription(desc)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return
}
