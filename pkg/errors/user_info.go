package errors

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/internal/constant"
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
func NewUserInfo(desc string) *UserInfo {
	return &UserInfo{
		descString: desc,
	}
}

// NewUserInfol creates a new UserInfo using the message generated from the
// passed localization.Config.
func NewUserInfol(desc localization.Config) *UserInfo {
	return &UserInfo{
		descConfig: desc,
	}
}

// NewUserInfolt creates a new UserInfo using the message generated from the
// passed term.
func NewUserInfolt(term string) *UserInfo {
	return NewUserInfol(localization.Config{
		Term: term,
	})
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
func (i *UserInfo) Handle(_ *state.State, ctx *plugin.Context) error {
	// we can ignore the error, because the fallback is set
	title, _ := ctx.Localizer.Localize(infoTitleConfig)

	desc, err := i.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	_, err = ctx.ReplyEmbed(discord.Embed{
		Title:       title,
		Description: desc,
		Color:       constant.InfoColor,
	})

	return err
}
