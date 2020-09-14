package restriction

import (
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/locutil"
)

// newInvalidChannelTypeError returns a new errors.RestrictionError wrapping
// an errors.InvalidChannelTypeError.
func newInvalidChannelTypeError(allowed plugin.ChannelTypes, l *localization.Localizer, fatal bool) error {
	err := errors.NewInvalidChannelTypeError(allowed)
	desc := err.Description(l)

	if fatal {
		return errors.NewFatalRestrictionError(desc)
	}

	return errors.NewRestrictionError(desc)
}

// newAllMissingRolesError creates a new error containing an
// error message for missing roles.
// The name field of the roles must be set.
func newAllMissingRolesError(missing []discord.Role, l *localization.Localizer) error {
	if len(missing) == 0 {
		return nil
	} else if len(missing) == 1 {
		return errors.NewFatalRestrictionErrorl(
			missingRoleError.
				WithPlaceholders(missingRoleErrorPlaceholders{
					Role: missing[0].Mention(),
				}))
	}

	embeddableDesc, _ := l.Localize(missingRolesAllError)

	indent, _ := genIndent(1)

	for _, r := range missing {
		embeddableDesc += "\n" + indent + entryPrefix + r.Mention()
	}

	defaultDesc, _ := l.Localize(missingRolesAllError)
	defaultDesc += "\n"

	for _, r := range missing {
		defaultDesc += "\n" + entryPrefix + r.Mention()
	}

	return &EmbeddableError{
		EmbeddableVersion: errors.NewFatalRestrictionError(embeddableDesc),
		DefaultVersion:    errors.NewFatalRestrictionError(defaultDesc),
	}
}

// newAnyMissingRolesError creates a new error containing an
// error message for missing roles.
// The name field of the roles must be set.
func newAnyMissingRolesError(missing []discord.Role, l *localization.Localizer) error {
	if len(missing) == 0 {
		return nil
	} else if len(missing) == 1 {
		return errors.NewFatalRestrictionErrorl(
			missingRoleError.
				WithPlaceholders(missingRoleErrorPlaceholders{
					Role: missing[0].Mention(),
				}))
	}

	desc, _ := l.Localize(missingRolesAnyError)

	embeddableDesc := desc
	indent, _ := genIndent(1)

	for _, r := range missing {
		embeddableDesc += "\n" + indent + entryPrefix + r.Mention()
	}

	defaultDesc := desc + "\n"

	for _, r := range missing {
		defaultDesc += "\n" + entryPrefix + r.Mention()
	}

	return &EmbeddableError{
		EmbeddableVersion: errors.NewFatalRestrictionError(embeddableDesc),
		DefaultVersion:    errors.NewFatalRestrictionError(defaultDesc),
	}
}

// newChannelsError creates a new error containing an error
// message containing the allowed channels.
// The name field of the channels must be set.
func newChannelsError(allowed []discord.ChannelID, l *localization.Localizer) error {
	if len(allowed) == 0 {
		return nil
	} else if len(allowed) == 1 {
		return errors.NewRestrictionErrorl(
			blockedChannelErrorSingle.
				WithPlaceholders(blockedChannelErrorSinglePlaceholders{
					Channel: "<#" + allowed[0].String() + ">",
				}))
	}

	desc, _ := l.Localize(blockedChannelErrorMulti)

	embeddableDesc := desc
	indent, _ := genIndent(1)

	for _, c := range allowed {
		embeddableDesc += "\n" + indent + entryPrefix + "<#" + c.String() + ">"
	}

	defaultDesc := desc + "\n"

	for _, c := range allowed {
		defaultDesc += "\n" + entryPrefix + "<#" + c.String() + ">"
	}

	return &EmbeddableError{
		EmbeddableVersion: errors.NewRestrictionError(embeddableDesc),
		DefaultVersion:    errors.NewRestrictionError(defaultDesc),
	}

}

// newInsufficientBotPermissions creates a new error containing the missing
// bot permissions
func newInsufficientBotPermissionsError(missing discord.Permissions, l *localization.Localizer) error {
	if missing == 0 {
		return nil
	}

	err := errors.NewInsufficientBotPermissionsError(missing)

	desc := err.Description(l)
	if err.IsSinglePermission() {
		return errors.NewRestrictionError(desc)
	}

	missingNames := locutil.PermissionNamesl(err.MissingPermissions, l)

	embeddableDesc := desc
	indent, _ := genIndent(1)

	for _, p := range missingNames {
		embeddableDesc += indent + "\n" + entryPrefix + p
	}

	defaultDesc := desc + "\n\n" + err.PermissionList(l)

	return &EmbeddableError{
		EmbeddableVersion: errors.NewRestrictionError(embeddableDesc),
		DefaultVersion:    errors.NewRestrictionError(defaultDesc),
	}
}

// newInsufficientUserPermissionsError returns a new error containing the
// missing permissions.
func newInsufficientUserPermissionsError(missing discord.Permissions, l *localization.Localizer) error {
	if missing == 0 {
		return nil
	}

	missingNames := locutil.PermissionNamesl(missing, l)

	if len(missingNames) == 0 {
		return nil
	} else if len(missingNames) == 1 {
		// we can ignore this error, as there is a fallback
		desc, _ := l.Localize(
			insufficientUserPermissionsDescSingle.
				WithPlaceholders(insufficientUserPermissionsDescSinglePlaceholders{
					MissingPermission: missingNames[0],
				}))

		return errors.NewFatalRestrictionError(desc)
	}

	desc, _ := l.Localize(insufficientUserPermissionsDescMulti)

	embeddableDesc := desc
	indent, _ := genIndent(1)

	for _, p := range missingNames {
		embeddableDesc += indent + "\n" + entryPrefix + p
	}

	defaultDesc := desc + "\n"

	for _, p := range missingNames {
		defaultDesc += indent + "\n" + entryPrefix + p
	}

	return &EmbeddableError{
		EmbeddableVersion: errors.NewFatalRestrictionError(embeddableDesc),
		DefaultVersion:    errors.NewFatalRestrictionError(defaultDesc),
	}
}
