package errors

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/discordutil"
)

// InsufficientBotPermissionsError is the error returned if the bot does not
// have sufficient permissions to execute a command.
type InsufficientBotPermissionsError struct {
	MissingPermissions discord.Permissions
}

// NewInsufficientBotPermissionError creates a new
// InsufficientBotPermissionsError with the passed MissingPermissions discord.Permissions.
func NewInsufficientBotPermissionsError(missing discord.Permissions) *InsufficientBotPermissionsError {
	return &InsufficientBotPermissionsError{
		MissingPermissions: missing,
	}
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *InsufficientBotPermissionsError) Description(l *localization.Localizer) string {
	// we can ignore this error, as there is a fallback
	desc, _ := l.Localize(insufficientBotPermissionsDesc)
	return desc
}

func (e *InsufficientBotPermissionsError) Error() string {
	return fmt.Sprintf("missingPermissions bot permissions: %d", e.MissingPermissions)
}

func (e *InsufficientBotPermissionsError) Is(err error) bool {
	casted, ok := err.(*InsufficientBotPermissionsError)
	if !ok {
		return false
	}

	return e.MissingPermissions == casted.MissingPermissions
}

// Handle sends an error message stating the MissingPermissions permissions.
func (e *InsufficientBotPermissionsError) Handle(_ *state.State, ctx *plugin.Context) (err error) {
	permNames := discordutil.PermissionNamesl(e.MissingPermissions, ctx.Localizer)

	perms, _ := ctx.Localize(insufficientBotPermissionMissingPermissionFieldName)
	permsVal := "• " + strings.Join(permNames, "\n• ")

	embed := newErrorEmbedBuilder(ctx.Localizer).
		WithDescription(e.Description(ctx.Localizer)).
		WithField(perms, permsVal)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return
}
