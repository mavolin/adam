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

// InsufficientPermissionsError is the error returned if the bot does not
// have sufficient permissions to execute a command.
type InsufficientPermissionsError struct {
	// MissingPermissions are the missing permissions.
	MissingPermissions discord.Permissions
}

var DefaultInsufficientPermissionsError = new(InsufficientPermissionsError)

// NewInsufficientPermissionError creates a new InsufficientPermissionsError
// with the passed missing permissions.
//
// If the missing permissions contain discord.PermissionAdministrator, all
// other permissions will be discarded, as they are included in Administrator.
//
// If missing is 0, a generic error message will be used.
func NewInsufficientPermissionsError(missing discord.Permissions) *InsufficientPermissionsError {
	// if we require Administrator, we will automatically receive all other
	// permissions once we get it
	if missing.Has(discord.PermissionAdministrator) {
		missing = discord.PermissionAdministrator
	}

	return &InsufficientPermissionsError{
		MissingPermissions: missing,
	}
}

// IsSinglePermission checks if only a single permission is missing.
func (e *InsufficientPermissionsError) IsSinglePermission() bool {
	return (e.MissingPermissions & (e.MissingPermissions - 1)) == 0
}

// Description returns the description of the error and localizes it, if
// possible.
// Note that if IsSinglePermission returns true, the description will already
// contain the missing permissions, which otherwise would need to be retrieved
// via PermissionList.
func (e *InsufficientPermissionsError) Description(l *i18n.Localizer) (desc string) {
	if e.MissingPermissions == 0 {
		// we can ignore this error, as there is a fallback
		desc, _ = l.Localize(insufficientPermissionsDefault)
		return desc
	}

	if e.IsSinglePermission() {
		missingNames := permutil.Namesl(e.MissingPermissions, l)
		if len(missingNames) == 0 {
			return ""
		}

		// we can ignore this error, as there is a fallback
		desc, _ = l.Localize(insufficientPermissionsDescSingle.
			WithPlaceholders(&insufficientBotPermissionsDescSinglePlaceholders{
				MissingPermission: missingNames[0],
			}))
	} else {
		// we can ignore this error, as there is a fallback
		desc, _ = l.Localize(insufficientPermissionsDescMulti)
	}

	return
}

// PermissionList returns a written bullet point list of the missing
// permissions, as used if multiple permissions are missing.
func (e *InsufficientPermissionsError) PermissionList(l *i18n.Localizer) string {
	permNames := permutil.Namesl(e.MissingPermissions, l)
	return "• " + strings.Join(permNames, "\n• ")
}

func (e *InsufficientPermissionsError) Error() string {
	return fmt.Sprintf("missingPermissions bot permissions: %d", e.MissingPermissions)
}

func (e *InsufficientPermissionsError) Is(target error) bool {
	var typedTarget *InsufficientPermissionsError
	if !As(target, &typedTarget) {
		return false
	}

	return e.MissingPermissions == typedTarget.MissingPermissions
}

// Handle handles the InsufficientPermissionsError.
// By default it sends an error embed stating the missing permissions.
func (e *InsufficientPermissionsError) Handle(s *state.State, ctx *plugin.Context) {
	HandleInsufficientPermissionsError(e, s, ctx)
}

var HandleInsufficientPermissionsError = func(
	ierr *InsufficientPermissionsError, _ *state.State, ctx *plugin.Context,
) {
	embed := ErrorEmbed.Clone().
		WithDescription(ierr.Description(ctx.Localizer))

	if !ierr.IsSinglePermission() {
		perms, err := ctx.Localize(insufficientPermissionsMissingPermissionsFieldName)
		if err != nil {
			return
		}

		embed.WithField(perms, ierr.PermissionList(ctx.Localizer))
	}

	_, _ = ctx.ReplyEmbedBuilder(embed)
}
