package restriction

import (
	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/permutil"
)

// =============================================================================
// EmbeddableError
// =====================================================================================

// EmbeddableError that formats differently, based on whether it is embedded in
// an any or all error, or whether it is used directly.
// Since the EmbeddableVersion is only used by any or all, it can safely expose
// its DefaultVersion using Unwrap, so that the errors.Handle will properly
// handle it.
type EmbeddableError struct {
	// EmbeddableVersion is the version used when embedded in an any or all
	// error.
	EmbeddableVersion *plugin.RestrictionError
	// DefaultVersion is the version returned if the error won't get embedded.
	DefaultVersion *plugin.RestrictionError
}

// Unwrap returns the error's DefaultVersion.
func (e *EmbeddableError) Unwrap() error {
	return e.DefaultVersion
}

func (e *EmbeddableError) Error() string { return e.DefaultVersion.Error() }

// =============================================================================
// ChannelTypesError
// =====================================================================================

// ChannelTypesError is the error returned if the command must be invoked in
// one of the Allowed plugin.ChannelTypes, but is not.
//
// It makes itself available as a plugin.RestrictionError via errors.As.
type ChannelTypesError struct {
	Allowed plugin.ChannelTypes

	underlying *plugin.RestrictionError
}

// NewChannelTypesError returns a new *ChannelTypesError created using the
// passed plugin.ChannelTypes.
func NewChannelTypesError(l *i18n.Localizer, allowed plugin.ChannelTypes) *ChannelTypesError {
	err := plugin.NewChannelTypeError(allowed)
	desc := err.Description(l)

	return &ChannelTypesError{
		Allowed:    allowed,
		underlying: plugin.NewRestrictionError(desc),
	}
}

// NewFatalChannelTypesError returns a fatal new *ChannelTypesError created
// using the passed plugin.ChannelTypes.
func NewFatalChannelTypesError(l *i18n.Localizer, allowed plugin.ChannelTypes) *ChannelTypesError {
	err := plugin.NewChannelTypeError(allowed)
	desc := err.Description(l)

	return &ChannelTypesError{
		Allowed:    allowed,
		underlying: plugin.NewFatalRestrictionError(desc),
	}
}

var _ error = new(ChannelTypesError)

func (e *ChannelTypesError) Error() string {
	return "restriction.ChannelTypesError"
}

func (e *ChannelTypesError) As(target interface{}) bool {
	switch err := target.(type) {
	case **plugin.RestrictionError:
		*err = e.underlying
		return true
	case *errors.Error:
		*err = e.underlying
		return true
	default:
		return false
	}
}

func (e *ChannelTypesError) AsRestrictionError() *plugin.RestrictionError {
	return e.underlying
}

// =============================================================================
// AllMissingRolesError
// =====================================================================================

// AllMissingRolesError is the error returned if the user needs all roles
// stored in Missing.
//
// It makes itself available as an *EmbeddableError and a
// *plugin.RestrictionError.
type AllMissingRolesError struct {
	Missing []discord.Role

	underlying *EmbeddableError
}

var _ error = new(AllMissingRolesError)

// NewAllMissingRolesError creates a new *AllMissingRolesError.
func NewAllMissingRolesError(l *i18n.Localizer, missing ...discord.Role) *AllMissingRolesError {
	if len(missing) == 0 {
		return nil
	} else if len(missing) == 1 {
		underlying := plugin.NewFatalRestrictionErrorl(
			missingRoleError.
				WithPlaceholders(&missingRoleErrorPlaceholders{
					Role: missing[0].Mention(),
				}))

		return &AllMissingRolesError{
			Missing: missing,
			underlying: &EmbeddableError{
				EmbeddableVersion: underlying,
				DefaultVersion:    underlying,
			},
		}
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

	return &AllMissingRolesError{
		Missing: missing,
		underlying: &EmbeddableError{
			EmbeddableVersion: plugin.NewFatalRestrictionError(embeddableDesc),
			DefaultVersion:    plugin.NewFatalRestrictionError(defaultDesc),
		},
	}
}

func (e *AllMissingRolesError) Error() string {
	return "restriction.AllMissingRolesError"
}

func (e *AllMissingRolesError) As(target interface{}) bool {
	switch err := target.(type) {
	case **EmbeddableError:
		*err = e.underlying
		return true
	case **plugin.RestrictionError:
		*err = e.underlying.DefaultVersion
		return true
	case *errors.Error:
		*err = e.underlying.DefaultVersion
		return true
	default:
		return false
	}
}

func (e *AllMissingRolesError) AsEmbeddableError() *EmbeddableError {
	return e.underlying
}

// =============================================================================
// AnyMissingRolesError
// =====================================================================================

// AnyMissingRolesError is the error returned if the user needs any of the
// roles stored in Missing.
//
// It makes itself available as an *EmbeddableError and a
// *plugin.RestrictionError.
type AnyMissingRolesError struct {
	Missing []discord.Role

	underlying *EmbeddableError
}

// NewAnyMissingRolesError creates a new error containing an
// error message for missing roles.
func NewAnyMissingRolesError(l *i18n.Localizer, missing ...discord.Role) *AnyMissingRolesError {
	if len(missing) == 0 {
		return nil
	} else if len(missing) == 1 {
		underlying := plugin.NewFatalRestrictionErrorl(
			missingRoleError.
				WithPlaceholders(&missingRoleErrorPlaceholders{
					Role: missing[0].Mention(),
				}))

		return &AnyMissingRolesError{
			Missing: missing,
			underlying: &EmbeddableError{
				EmbeddableVersion: underlying,
				DefaultVersion:    underlying,
			},
		}
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

	return &AnyMissingRolesError{
		Missing: missing,
		underlying: &EmbeddableError{
			EmbeddableVersion: plugin.NewFatalRestrictionError(embeddableDesc),
			DefaultVersion:    plugin.NewFatalRestrictionError(defaultDesc),
		},
	}
}

func (e *AnyMissingRolesError) Error() string {
	return "restriction.AnyMissingRolesError"
}

func (e *AnyMissingRolesError) As(target interface{}) bool {
	switch err := target.(type) {
	case **EmbeddableError:
		*err = e.underlying
		return true
	case **plugin.RestrictionError:
		*err = e.underlying.DefaultVersion
		return true
	case *errors.Error:
		*err = e.underlying.DefaultVersion
		return true
	default:
		return false
	}
}

func (e *AnyMissingRolesError) AsEmbeddableError() *EmbeddableError {
	return e.underlying
}

// =============================================================================
// ChannelsError
// =====================================================================================

// ChannelsError is the error returned if the user needs to invoke the command
// in any of channels whose ids are stored in Allowed.
//
// It makes itself available as an *EmbeddableError and a
// *plugin.RestrictionError.
type ChannelsError struct {
	Allowed []discord.ChannelID

	underlying *EmbeddableError
}

// NewChannelsError creates a new error containing an error
// message containing the allowed channels.
func NewChannelsError(l *i18n.Localizer, allowed ...discord.ChannelID) *ChannelsError {
	if len(allowed) == 0 {
		return nil
	} else if len(allowed) == 1 {
		underlying := plugin.NewRestrictionErrorl(
			blockedChannelErrorSingle.
				WithPlaceholders(&blockedChannelErrorSinglePlaceholders{
					Channel: allowed[0].Mention(),
				}))

		return &ChannelsError{
			Allowed: allowed,
			underlying: &EmbeddableError{
				EmbeddableVersion: underlying,
				DefaultVersion:    underlying,
			},
		}
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

	return &ChannelsError{
		Allowed: allowed,
		underlying: &EmbeddableError{
			EmbeddableVersion: plugin.NewRestrictionError(embeddableDesc),
			DefaultVersion:    plugin.NewRestrictionError(defaultDesc),
		},
	}
}

func (e *ChannelsError) Error() string {
	return "restriction.AnyMissingRolesError"
}

func (e *ChannelsError) As(target interface{}) bool {
	switch err := target.(type) {
	case **EmbeddableError:
		*err = e.underlying
		return true
	case **plugin.RestrictionError:
		*err = e.underlying.DefaultVersion
		return true
	case *errors.Error:
		*err = e.underlying.DefaultVersion
		return true
	default:
		return false
	}
}

func (e *ChannelsError) AsEmbeddableError() *EmbeddableError {
	return e.underlying
}

// =============================================================================
// UserPermissionsError
// =====================================================================================

// UserPermissionsError is the error returned if a user needs all permissions
// stored in Missing.
//
// It makes itself available as an *EmbeddableError and a
// *plugin.RestrictionError.
type UserPermissionsError struct {
	Missing discord.Permissions

	underlying *EmbeddableError
}

// NewUserPermissionsError returns a new error containing the
// missing permissions.
func NewUserPermissionsError(l *i18n.Localizer, missing discord.Permissions) *UserPermissionsError {
	if missing == 0 {
		return nil
	}

	missingNames := permutil.Names(l, missing)

	if len(missingNames) == 0 {
		return nil
	} else if len(missingNames) == 1 {
		underlying := plugin.NewRestrictionErrorl(userPermissionsDescSingle.
			WithPlaceholders(&userPermissionsDescSinglePlaceholders{
				MissingPermission: missingNames[0],
			}))

		return &UserPermissionsError{
			Missing: missing,
			underlying: &EmbeddableError{
				EmbeddableVersion: underlying,
				DefaultVersion:    underlying,
			},
		}
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

	return &UserPermissionsError{
		Missing: missing,
		underlying: &EmbeddableError{
			EmbeddableVersion: plugin.NewFatalRestrictionError(embeddableDesc),
			DefaultVersion:    plugin.NewFatalRestrictionError(defaultDesc),
		},
	}
}

func (e *UserPermissionsError) Error() string {
	return "restriction.AnyMissingRolesError"
}

func (e *UserPermissionsError) As(target interface{}) bool {
	switch err := target.(type) {
	case **EmbeddableError:
		*err = e.underlying
		return true
	case **plugin.RestrictionError:
		*err = e.underlying.DefaultVersion
		return true
	case *errors.Error:
		*err = e.underlying.DefaultVersion
		return true
	default:
		return false
	}
}

func (e *UserPermissionsError) AsEmbeddableError() *EmbeddableError {
	return e.underlying
}
