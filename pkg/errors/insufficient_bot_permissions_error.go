package errors

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/permutil"
)

// InsufficientBotPermissionsError is the error returned if the bot does not
// have sufficient permissions to execute a command.
type InsufficientBotPermissionsError struct {
	// MissingPermissions are the missing permissions.
	MissingPermissions discord.Permissions
}

var _ Interface = new(InsufficientBotPermissionsError)

// NewInsufficientBotPermissionError creates a new
// InsufficientBotPermissionsError with the passed missing permissions.
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
// Note that if IsSinglePermission returns true, the description will already
// contain the missing permissions, which otherwise would need to be retrieved
// via PermissionList.
func (e *InsufficientBotPermissionsError) Description(l *i18n.Localizer) (desc string) {
	if e.MissingPermissions == 0 {
		return
	}

	if e.IsSinglePermission() {
		missingNames := permutil.Namesl(e.MissingPermissions, l)
		if len(missingNames) == 0 {
			return ""
		}

		// we can ignore this error, as there is a fallback
		desc, _ = l.Localize(insufficientBotPermissionsDescSingle.
			WithPlaceholders(&insufficientBotPermissionsDescSinglePlaceholders{
				MissingPermission: missingNames[0],
			}))
	} else {
		// we can ignore this error, as there is a fallback
		desc, _ = l.Localize(insufficientBotPermissionsDescMulti)
	}

	return
}

// PermissionList returns a written bullet point list of the missing
// permissions, as used if multiple permissions are missing.
func (e *InsufficientBotPermissionsError) PermissionList(l *i18n.Localizer) string {
	permNames := permutil.Namesl(e.MissingPermissions, l)
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

// Handle sends an error message stating the missing permissions.
func (e *InsufficientBotPermissionsError) Handle(_ *state.State, ctx *plugin.Context) (err error) {
	embed := ErrorEmbed.Clone().
		WithDescription(e.Description(ctx.Localizer))

	if !e.IsSinglePermission() {
		perms, _ := ctx.Localize(insufficientBotPermissionMissingMissingPermissionsFieldName)
		embed.WithField(perms, e.PermissionList(ctx.Localizer))
	}

	_, err = ctx.ReplyEmbedBuilder(embed)
	return
}
