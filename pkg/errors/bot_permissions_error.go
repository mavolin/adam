package errors

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/permutil"
)

// BotPermissionsError is the error returned if the bot does not have
// sufficient permissions to execute a command.
type BotPermissionsError struct {
	// Missing are the missing permissions.
	Missing discord.Permissions
}

var DefaultBotPermissionsError Error = new(BotPermissionsError)

// allPermissions except admin
var allPerms = discord.PermissionAll ^ discord.PermissionAdministrator

// NewBotPermissionsError creates a new BotPermissionsError with the passed
// missing permissions.
//
// If the missing permissions contain discord.PermissionAdministrator, all
// other permissions will be discarded, as they are included in Administrator.
//
// If missing is 0 or invalid, a generic error message will be used.
func NewBotPermissionsError(missing discord.Permissions) *BotPermissionsError {
	return &BotPermissionsError{Missing: missing}
}

// IsSinglePermission checks if only a single permission is missing.
func (e *BotPermissionsError) IsSinglePermission() bool {
	return e.Missing.Has(discord.PermissionAdministrator) || e.Missing.Has(allPerms) || (e.Missing&(e.Missing-1)) == 0
}

// Description returns the description of the error and localizes it, if
// possible.
// Note that if IsSinglePermission returns true, the description will already
// contain the missing permissions, which otherwise would need to be retrieved
// via PermissionList.
func (e *BotPermissionsError) Description(l *i18n.Localizer) (desc string) {
	if e.Missing == 0 {
		// we can ignore this error, as there is a fallback
		desc, _ = l.Localize(insufficientPermissionsDefault)
		return desc
	}

	missing := e.Missing

	// if we require Administrator, we will automatically receive all other
	// permissions once we get it
	if missing.Has(discord.PermissionAdministrator) || missing.Has(allPerms) {
		missing = discord.PermissionAdministrator
	}

	if e.IsSinglePermission() {
		missingNames := permutil.Namesl(missing, l)
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

	return desc
}

// PermissionList returns a written bullet point list of the missing
// permissions, as used if multiple permissions are missing.
func (e *BotPermissionsError) PermissionList(l *i18n.Localizer) string {
	permNames := permutil.Namesl(e.Missing, l)
	return "• " + strings.Join(permNames, "\n• ")
}

func (e *BotPermissionsError) Error() string {
	return fmt.Sprintf("missing bot permissions: %d", e.Missing)
}

func (e *BotPermissionsError) Is(target error) bool {
	var typedTarget *BotPermissionsError
	if !As(target, &typedTarget) {
		return false
	}

	return e.Missing == typedTarget.Missing
}

// Handle handles the BotPermissionsError.
// By default it sends an error Embed stating the missing permissions.
func (e *BotPermissionsError) Handle(s *state.State, ctx *plugin.Context) error {
	return HandleBotPermissionsError(e, s, ctx)
}

var HandleBotPermissionsError = func(
	ierr *BotPermissionsError, _ *state.State, ctx *plugin.Context,
) error {
	// if this error arose because of a missing send messages permission,
	// do nothing, as we can't send an error message
	if ierr.Missing.Has(discord.PermissionSendMessages) {
		return nil
	}

	embed := ErrorEmbed.Clone().
		WithDescription(ierr.Description(ctx.Localizer))

	if !ierr.IsSinglePermission() {
		perms, err := ctx.Localize(insufficientPermissionsMissingPermissionsFieldName)
		if err != nil {
			return err
		}

		embed.WithField(perms, ierr.PermissionList(ctx.Localizer))
	}

	_, err := ctx.ReplyEmbedBuilder(embed)
	return err
}
