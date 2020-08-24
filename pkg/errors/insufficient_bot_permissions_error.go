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
// InsufficientBotPermissionsError with the passed MissingPermissions
// discord.Permissions.
// If the missing permissions contain discord.PermissionAdministrator, all
// other permissions will be discarded, as they are included in Administrator.
func NewInsufficientBotPermissionsError(missing discord.Permissions) *InsufficientBotPermissionsError {
	// if we require Administrator, we will automatically receive all other
	// permissions once we get it
	if missing.Has(discord.PermissionAdministrator) {
		missing = discord.PermissionAdministrator
	}

	return &InsufficientBotPermissionsError{
		MissingPermissions: missing,
	}
}

// IsSinglePermission checks if only a single permission is missing.
func (e *InsufficientBotPermissionsError) IsSinglePermission() bool {
	return (e.MissingPermissions & (e.MissingPermissions - 1)) == 0
}

// Description returns the description of the error and localizes it, if
// possible.
// Note that if IsSinglePermission returns true, the Description will already contain
// the missing permissions, which would otherwise needed to be retrieved via
// PermissionList.
func (e *InsufficientBotPermissionsError) Description(l *localization.Localizer) (desc string) {
	if e.MissingPermissions == 0 {
		return
	}

	if e.IsSinglePermission() {
		missingNames := discordutil.PermissionNamesl(e.MissingPermissions, l)
		if len(missingNames) == 0 {
			return ""
		}

		// we can ignore this error, as there is a fallback
		desc, _ = l.Localize(insufficientBotPermissionsDescSingle.
			WithPlaceholders(insufficientBotPermissionsDescSinglePlaceholders{
				MissingPermission: "`" + discordutil.EscapeInlineCode(missingNames[0]) + "`",
			}))
	} else {
		// we can ignore this error, as there is a fallback
		desc, _ = l.Localize(insufficientBotPermissionsDescMulti)
	}

	return
}

// PermissionList returns a written bullet point list of the missing
// permissions, as used if multiple permissions are missing.
func (e *InsufficientBotPermissionsError) PermissionList(l *localization.Localizer) string {
	permNames := discordutil.PermissionNamesl(e.MissingPermissions, l)
	return "• " + strings.Join(permNames, "\n• ")
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
	embed := newErrorEmbedBuilder(ctx.Localizer).
		WithDescription(e.Description(ctx.Localizer))

	if !e.IsSinglePermission() {
		perms, _ := ctx.Localize(insufficientBotPermissionMissingMissingPermissionsFieldName)
		embed.WithField(perms, e.PermissionList(ctx.Localizer))
	}

	_, err = ctx.ReplyEmbedBuilder(embed)
	return
}
