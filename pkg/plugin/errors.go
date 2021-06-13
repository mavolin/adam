package plugin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/internal/shared"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
	"github.com/mavolin/adam/pkg/utils/permutil"
)

// =============================================================================
// ArgumentError
// =====================================================================================

// ArgumentError is the error used if an argument or flag a user supplied is
// invalid.
type ArgumentError struct {
	desc *i18nutil.Text
}

// NewArgumentError returns a new *ArgumentError with the passed
// description.
// The description mustn't be empty for this error to be handled properly.
func NewArgumentError(description string) *ArgumentError {
	return &ArgumentError{desc: i18nutil.NewText(description)}
}

// NewArgumentErrorl returns a new *ArgumentError using the passed *i18n.Config
// to generate a description.
func NewArgumentErrorl(description *i18n.Config) *ArgumentError {
	return &ArgumentError{desc: i18nutil.NewTextl(description)}
}

// NewArgumentErrorlt returns a new *ArgumentError using the passed term to
// generate a description.
func NewArgumentErrorlt(description i18n.Term) *ArgumentError {
	return NewArgumentErrorl(description.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *ArgumentError) Description(l *i18n.Localizer) (string, error) {
	return e.desc.Get(l)
}

func (e *ArgumentError) Error() string {
	return "argument error"
}

// Handle handles the ArgumentError.
// By default it sends an error Embed containing a description of which
// arg/flag was faulty in the channel the command was sent in.
func (e *ArgumentError) Handle(s *state.State, ctx *Context) error {
	return HandleArgumentError(e, s, ctx)
}

var HandleArgumentError = func(aerr *ArgumentError, _ *state.State, ctx *Context) error {
	desc, err := aerr.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := shared.ErrorEmbed.Clone().
		WithDescription(desc)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return err
}

// =============================================================================
// BotPermissionsError
// =====================================================================================

// BotPermissionsError is the error returned if the bot does not have
// sufficient permissions to execute a command.
type BotPermissionsError struct {
	// Missing are the missing permissions.
	Missing discord.Permissions
}

// DefaultBotPermissionsError is an *BotPermissionsError that displays a
// generic "missing permissions" error message, instead of listing the missing
// permissions.
var DefaultBotPermissionsError = new(BotPermissionsError)

// allPermissions except admin
var allPerms = discord.PermissionAll ^ discord.PermissionAdministrator

// NewBotPermissionsError creates a new *BotPermissionsError with the passed
// missing permissions.
//
// If the missing permissions contains discord.PermissionAdministrator, all
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
	if !errors.As(target, &typedTarget) {
		return false
	}

	return e.Missing == typedTarget.Missing
}

// Handle handles the BotPermissionsError.
// By default it sends an error Embed stating the missing permissions.
func (e *BotPermissionsError) Handle(s *state.State, ctx *Context) error {
	return HandleBotPermissionsError(e, s, ctx)
}

var HandleBotPermissionsError = func(perr *BotPermissionsError, _ *state.State, ctx *Context) error {
	// if this error arose because of a missing send messages permission,
	// do nothing, as we can't send an error message
	if perr.Missing.Has(discord.PermissionSendMessages) {
		return nil
	}

	embed := shared.ErrorEmbed.Clone().
		WithDescription(perr.Description(ctx.Localizer))

	if !perr.IsSinglePermission() {
		perms, err := ctx.Localize(insufficientPermissionsMissingPermissionsFieldName)
		if err != nil {
			return err
		}

		embed.WithField(perms, perr.PermissionList(ctx.Localizer))
	}

	_, err := ctx.ReplyEmbedBuilder(embed)
	return err
}

// =============================================================================
// ChannelTypeError
// =====================================================================================

// ChannelTypeError is the error returned if a command is invoked in a channel
// that is not supported by that command.
type ChannelTypeError struct {
	// Allowed are the plugin.ChannelTypes that the command supports.
	Allowed ChannelTypes
}

// NewChannelTypeError creates a new *ChannelTypeError with the passed allowed
// plugin.ChannelTypes.
func NewChannelTypeError(allowed ChannelTypes) *ChannelTypeError {
	return &ChannelTypeError{Allowed: allowed}
}

// Description returns the description containing the types of channels this
// command may be used in.
func (e *ChannelTypeError) Description(l *i18n.Localizer) (desc string) {
	switch {
	// ----- singles -----
	case e.Allowed == GuildTextChannels:
		desc, _ = l.Localize(channelTypeErrorGuildText)
	case e.Allowed == GuildNewsChannels:
		desc, _ = l.Localize(channelTypeErrorGuildNews)
	case e.Allowed == DirectMessages:
		desc, _ = l.Localize(channelTypeErrorDM)
	// ----- combinations -----
	case e.Allowed == GuildChannels:
		desc, _ = l.Localize(channelTypeErrorGuild)
	case e.Allowed == (DirectMessages | GuildTextChannels):
		desc, _ = l.Localize(channelTypeErrorDMAndGuildText)
	case e.Allowed == (DirectMessages | GuildNewsChannels):
		desc, _ = l.Localize(channelTypeErrorDMAndGuildNews)
	default:
		desc, _ = l.Localize(channelTypeErrorFallback)
	}

	return
}

func (e *ChannelTypeError) Error() string {
	return "channel type error"
}

func (e *ChannelTypeError) Is(target error) bool {
	var typedTarget *ChannelTypeError
	if !errors.As(target, &typedTarget) {
		return false
	}

	return e.Allowed == typedTarget.Allowed
}

// Handle handles the ChannelTypeError.
// By default it sends an error message stating the allowed channel types.
func (e *ChannelTypeError) Handle(s *state.State, ctx *Context) error {
	return HandleChannelTypeError(e, s, ctx)
}

var HandleChannelTypeError = func(cerr *ChannelTypeError, s *state.State, ctx *Context) error {
	embed := shared.ErrorEmbed.Clone().
		WithDescription(cerr.Description(ctx.Localizer))

	_, err := ctx.ReplyEmbedBuilder(embed)
	return err
}

// =============================================================================
// RestrictionError
// =====================================================================================

var (
	// DefaultRestrictionError is a restriction error with a default, generic
	// description.
	DefaultRestrictionError = NewRestrictionErrorl(defaultRestrictionDesc)
	// DefaultFatalRestrictionError is a restriction error with a default,
	// generic description and Fatal set to true.
	DefaultFatalRestrictionError = NewFatalRestrictionErrorl(defaultRestrictionDesc)
)

// RestrictionError is the error returned if a restriction fails.
// It contains a description stating the conditions that need to be fulfilled
// for a command to execute.
//
// Note that the description might contain mentions, which are intended not
// to ping anyone, e.g. "You need @role to use this command.".
// This means you should use allowed mentions if you are custom handling this
// error and not using an Embed, which suppresses mentions by default.
type RestrictionError struct {
	// description of the error
	desc *i18nutil.Text

	// Fatal defines if the RestrictionError is fatal.
	// A RestrictionError is fatal, if the user cannot prevent the error from
	// occurring again, without the action of another user, e.g. getting a
	// permission.
	Fatal bool
}

// NewRestrictionError creates a new *RestrictionError with the passed
// description.
func NewRestrictionError(description string) *RestrictionError {
	return &RestrictionError{desc: i18nutil.NewText(description)}
}

// NewRestrictionErrorl creates a new *RestrictionError using the message
// generated from the passed *i18n.Config as description.
func NewRestrictionErrorl(description *i18n.Config) *RestrictionError {
	return &RestrictionError{desc: i18nutil.NewTextl(description)}
}

// NewRestrictionErrorlt creates a new *RestrictionError using the message
// generated from the passed term as description.
func NewRestrictionErrorlt(description i18n.Term) *RestrictionError {
	return NewRestrictionErrorl(description.AsConfig())
}

// NewFatalRestrictionError creates a new fatal *RestrictionError with the
// passed description.
func NewFatalRestrictionError(description string) *RestrictionError {
	return &RestrictionError{
		desc:  i18nutil.NewText(description),
		Fatal: true,
	}
}

// NewFatalRestrictionErrorl creates a new fatal *RestrictionError using the
// message generated from the passed *i18n.Config as description.
func NewFatalRestrictionErrorl(description *i18n.Config) *RestrictionError {
	return &RestrictionError{
		desc:  i18nutil.NewTextl(description),
		Fatal: true,
	}
}

// NewFatalRestrictionErrorlt creates a new fatal *RestrictionError using the
// message generated from the passed term as description.
func NewFatalRestrictionErrorlt(description i18n.Term) *RestrictionError {
	return NewFatalRestrictionErrorl(description.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *RestrictionError) Description(l *i18n.Localizer) (string, error) {
	return e.desc.Get(l)
}

func (e *RestrictionError) Error() string { return "restriction error" }

// Handle handles the RestrictionError.
// By default it sends an error Embed with the description of the
// RestrictionError.
func (e *RestrictionError) Handle(s *state.State, ctx *Context) error {
	return HandleRestrictionError(e, s, ctx)
}

var HandleRestrictionError = func(rerr *RestrictionError, s *state.State, ctx *Context) error {
	desc, err := rerr.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := shared.ErrorEmbed.Clone().
		WithDescription(desc)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return err
}

// =============================================================================
// ThrottlingError
// =====================================================================================

// ThrottlingError is the error returned if a command gets throttled.
// It contains a description about when the command will become available
// again.
type ThrottlingError struct {
	// description of the error
	desc *i18nutil.Text
}

// NewThrottlingError creates a new *ThrottlingError with the passed
// description.
func NewThrottlingError(description string) *ThrottlingError {
	return &ThrottlingError{desc: i18nutil.NewText(description)}
}

// NewThrottlingErrorl creates a new *ThrottlingError using the message
// generated from the passed *i18n.Config as description.
func NewThrottlingErrorl(description *i18n.Config) *ThrottlingError {
	return &ThrottlingError{desc: i18nutil.NewTextl(description)}
}

// NewThrottlingErrorlt creates a new *ThrottlingError using the message
// generated from the passed term as description.
func NewThrottlingErrorlt(description i18n.Term) *ThrottlingError {
	return NewThrottlingErrorl(description.AsConfig())
}

// Description returns the description of the error and localizes it, if
// possible.
func (e *ThrottlingError) Description(l *i18n.Localizer) (string, error) {
	return e.desc.Get(l)
}

func (e *ThrottlingError) Error() string { return "throttling error" }

// Handle handles the ThrottlingError.
// By default it sends an info Embed with the description of the
// ThrottlingError.
func (e *ThrottlingError) Handle(s *state.State, ctx *Context) error {
	return HandleThrottlingError(e, s, ctx)
}

var HandleThrottlingError = func(terr *ThrottlingError, s *state.State, ctx *Context) error {
	desc, err := terr.Description(ctx.Localizer)
	if err != nil {
		return err
	}

	embed := shared.InfoEmbed.Clone().
		WithDescription(desc)

	_, err = ctx.ReplyEmbedBuilder(embed)
	return err
}
