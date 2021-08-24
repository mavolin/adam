package restriction

import (
	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/permutil"
)

// newChannelTypesError returns a new *plugin.RestrictionError wrapping
// a plugin.ChannelTypeError.
func newChannelTypesError(allowed plugin.ChannelTypes, l *i18n.Localizer, fatal bool) error {
	err := plugin.NewChannelTypeError(allowed)
	desc := err.Description(l)

	if fatal {
		return plugin.NewFatalRestrictionError(desc)
	}

	return plugin.NewRestrictionError(desc)
}

// newAllMissingRolesError creates a new error containing an
// error message for missing roles.
// The name field of the roles must be set.
func newAllMissingRolesError(missing []discord.Role, l *i18n.Localizer) error {
	if len(missing) == 0 {
		return nil
	} else if len(missing) == 1 {
		return plugin.NewFatalRestrictionErrorl(
			missingRoleError.
				WithPlaceholders(&missingRoleErrorPlaceholders{
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
		EmbeddableVersion: plugin.NewFatalRestrictionError(embeddableDesc),
		DefaultVersion:    plugin.NewFatalRestrictionError(defaultDesc),
	}
}

// newAnyMissingRolesError creates a new error containing an
// error message for missing roles.
// The name field of the roles must be set.
func newAnyMissingRolesError(missing []discord.Role, l *i18n.Localizer) error {
	if len(missing) == 0 {
		return nil
	} else if len(missing) == 1 {
		return plugin.NewFatalRestrictionErrorl(
			missingRoleError.
				WithPlaceholders(&missingRoleErrorPlaceholders{
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
		EmbeddableVersion: plugin.NewFatalRestrictionError(embeddableDesc),
		DefaultVersion:    plugin.NewFatalRestrictionError(defaultDesc),
	}
}

// newChannelsError creates a new error containing an error
// message containing the allowed channels.
// The name field of the channels must be set.
func newChannelsError(allowed []discord.ChannelID, l *i18n.Localizer) error {
	if len(allowed) == 0 {
		return nil
	} else if len(allowed) == 1 {
		return plugin.NewRestrictionErrorl(
			blockedChannelErrorSingle.
				WithPlaceholders(&blockedChannelErrorSinglePlaceholders{
					Channel: allowed[0].Mention(),
				}))
	}

	desc, _ := l.Localize(blockedChannelErrorMulti)

	embeddableDesc := desc
	indent, _ := genIndent(1)

	for _, c := range allowed {
		embeddableDesc += "\n" + indent + entryPrefix + c.Mention()
	}

	defaultDesc := desc + "\n"

	for _, c := range allowed {
		defaultDesc += "\n" + entryPrefix + c.Mention()
	}

	return &EmbeddableError{
		EmbeddableVersion: plugin.NewRestrictionError(embeddableDesc),
		DefaultVersion:    plugin.NewRestrictionError(defaultDesc),
	}
}

// newUserPermissionsError returns a new error containing the
// missing permissions.
func newUserPermissionsError(missing discord.Permissions, l *i18n.Localizer) error {
	if missing == 0 {
		return nil
	}

	missingNames := permutil.Namesl(missing, l)

	if len(missingNames) == 0 {
		return nil
	} else if len(missingNames) == 1 {
		// we can ignore this error, as there is a fallback
		desc, _ := l.Localize(
			userPermissionsDescSingle.
				WithPlaceholders(&userPermissionsDescSinglePlaceholders{
					MissingPermission: missingNames[0],
				}))

		return plugin.NewFatalRestrictionError(desc)
	}

	desc, _ := l.Localize(userPermissionsDescMulti)

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
		EmbeddableVersion: plugin.NewFatalRestrictionError(embeddableDesc),
		DefaultVersion:    plugin.NewFatalRestrictionError(defaultDesc),
	}
}
