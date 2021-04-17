package arg

import (
	"net/url"
	"regexp"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
)

// =============================================================================
// Text
// =====================================================================================

// Text is the Type for a string.
//
// Go type: string
type Text struct {
	// CustomName allows you to set a custom name for the id.
	// If not set, the default name will be used.
	CustomName *i18nutil.Text
	// CustomDescription allows you to set a custom description for the id.
	// If not set, the default description will be used.
	CustomDescription *i18nutil.Text

	// MinLength is the inclusive minimum length the text may have.
	MinLength uint
	// MaxLength is the inclusive maximum length the text may have.
	// If MaxLength is 0, the id won't have a maximum.
	MaxLength uint

	// Regexp is the regular expression the text must match.
	// If Regexp is set to nil/zero, any text within the bounds will pass.
	//
	// If matching fails, the corresponding RegexpErrorX will be returned.
	Regexp *regexp.Regexp
	// RegexpErrorArg is the error message used if an argument doesn't match
	// the regular expression defined.
	// If you want an unlocalized error, just fill Fallback.Other field of the
	// config.
	//
	// Available Placeholders are:
	//
	// 		• name - the name of the argument
	// 		• raw - the raw argument
	// 		• position - the position of the id (1-indexed)
	// 		• regexp - the regular expression that needs to be matched
	//
	// Defaults to: regexpNotMatchingErrorArg
	RegexpErrorArg *i18n.Config
	// RegexpErrorFlag is the error message used if a flag doesn't match the
	// regular expression defined.
	// If you want an unlocalized error, just fill Fallback.Other field of the
	// config.
	//
	// Available Placeholders are:
	//
	// 		• name - the full name of the flag
	// 		• used_name - the name of the flag the invoking user used
	// 		• raw - the raw flag without the flags name
	// 		• regexp - the regular expression that needs to be matched
	//
	// Defaults to: regexpNotMatchingErrorFlag
	RegexpErrorFlag *i18n.Config
}

var (
	// SimpleText is a Text with no length boundaries and no regular expression.
	SimpleText Type = new(Text)
	_          Type = Text{}
)

func (t Text) Name(l *i18n.Localizer) string {
	if t.CustomName != nil {
		name, err := t.CustomName.Get(l)
		if err == nil {
			return name
		}
	}

	name, _ := l.Localize(textName) // we have a fallback
	return name
}

func (t Text) Description(l *i18n.Localizer) string {
	if t.CustomDescription != nil {
		desc, err := t.CustomDescription.Get(l)
		if err == nil {
			return desc
		}
	}

	desc, _ := l.Localize(textDescription) // we have a fallback
	return desc
}

func (t Text) Parse(_ *state.State, ctx *Context) (interface{}, error) { //nolint:dupl
	if uint(len(ctx.Raw)) < t.MinLength {
		return nil, newArgumentError2(
			textBelowMinLengthErrorArg, textBelowMinLengthErrorFlag, ctx, map[string]interface{}{
				"min": t.MinLength,
			})
	} else if t.MaxLength > 0 && uint(len(ctx.Raw)) > t.MaxLength {
		return nil, newArgumentError2(
			textAboveMaxLengthErrorArg, textAboveMaxLengthErrorFlag, ctx, map[string]interface{}{
				"max": t.MaxLength,
			})
	}

	if t.Regexp != nil && !t.Regexp.MatchString(ctx.Raw) {
		if ctx.Kind == KindArg && t.RegexpErrorArg == nil {
			t.RegexpErrorArg = regexpNotMatchingErrorArg
		} else if ctx.Kind == KindFlag && t.RegexpErrorFlag == nil {
			t.RegexpErrorFlag = regexpNotMatchingErrorFlag
		}

		return nil, newArgumentError2(t.RegexpErrorArg, t.RegexpErrorFlag, ctx, map[string]interface{}{
			"regexp": t.Regexp.String(),
		})
	}

	return ctx.Raw, nil
}

func (t Text) Default() interface{} {
	return ""
}

// =============================================================================
// Link
// =====================================================================================

// Link is the Type used for URLs.
//
// Go type: string
type Link struct {
	// Validator checks if the passed *url.URL is valid.
	//
	// By default, Validator will check if the scheme is either 'http' or
	// 'https'.
	Validator func(u *url.URL) bool
	// ErrorArg is the error message used if an argument doesn't match
	// the regular expression defined, or does not pass url.ParseRequestURI.
	// If you want an unlocalized error, just fill Fallback.Other field of the
	// config.
	//
	// Available Placeholders are:
	//
	// 		• name - the name of the argument
	// 		• raw - the raw argument
	// 		• position - the position of the id (1-indexed)
	// 		• regexp - the regular expression that needs to be matched
	//
	// Defaults to: linkInvalidErrorArg
	ErrorArg *i18n.Config
	// ErrorFlag is the error message used if a flag doesn't match the
	// regular expression defined, or does not pass url.ParseRequestURI.
	// If you want an unlocalized error, just fill Fallback.Other field of the
	// config.
	//
	// Available Placeholders are:
	//
	// 		• name - the full name of the flag
	// 		• used_name - the name of the flag the invoking user used
	// 		• raw - the raw flag without the flags name
	// 		• regexp - the regular expression that needs to be matched
	//
	// Defaults to: linkInvalidErrorFlag
	ErrorFlag *i18n.Config
}

var (
	// SimpleLink is a link that uses no custom regular expression.
	SimpleLink Type = new(Link)
	_          Type = Link{}
)

func (l Link) Name(loc *i18n.Localizer) string {
	name, _ := loc.Localize(linkName) // we have a fallback
	return name
}

func (l Link) Description(loc *i18n.Localizer) string {
	desc, _ := loc.Localize(linkDescription)
	return desc
}

func (l Link) Parse(_ *state.State, ctx *Context) (interface{}, error) {
	if l.Validator == nil {
		l.Validator = defaultLinkValidator
	}

	u, err := url.ParseRequestURI(ctx.Raw)
	if err != nil || !l.Validator(u) {
		if (ctx.Kind == KindArg && l.ErrorArg == nil) || (ctx.Kind == KindFlag && l.ErrorFlag == nil) {
			return nil, newArgumentError2(linkInvalidErrorArg, linkInvalidErrorFlag, ctx, nil)
		}

		return nil, newArgumentError2(l.ErrorArg, l.ErrorFlag, ctx, nil)
	}

	return ctx.Raw, nil
}

func (l Link) Default() interface{} {
	return ""
}

func defaultLinkValidator(u *url.URL) bool {
	return u.Scheme == "http" || u.Scheme == "https"
}

// =============================================================================
// AlphanumericID
// =====================================================================================

// AlphanumericID is a Type for alphanumeric ids.
// By default AlphanumericIDs share the same name and description as a
// NumericID, simply their definition differs.
//
// In contrast to a NumericID, a AlphanumericID returns strings and can handle
// numbers that exceed NumericIDs 64 bit limit.
//
// Go type: string
type AlphanumericID struct {
	// CustomName allows you to set a custom name for the id.
	// If not set, the default name will be used.
	CustomName *i18nutil.Text
	// CustomDescription allows you to set a custom description for the id.
	// If not set, the default description will be used.
	CustomDescription *i18nutil.Text

	// MinLength is the inclusive minimum length the ID may have.
	MinLength uint
	// MaxLength is the inclusive maximum length the id may have.
	// If MaxLength is 0, the id won't have a maximum.
	MaxLength uint

	// Regexp is the regular expression the id needs to match to pass.
	// If Regexp is set to nil/zero, any id within the bounds will pass.
	//
	// If matching fails, the corresponding RegexpErrorX will be returned.
	Regexp *regexp.Regexp
	// RegexpErrorArg is the error message used if an argument doesn't match
	// the regular expression defined.
	// If you want an unlocalized error, just fill Fallback.Other field of the
	// config.
	//
	// Available Placeholders are:
	//
	// 		• name - the name of the argument
	// 		• raw - the raw argument
	// 		• position - the position of the id (1-indexed)
	// 		• regexp - the regular expression that needs to be matched
	//
	// Defaults to: idRegexpNotMatchingErrorArg
	RegexpErrorArg *i18n.Config
	// RegexpErrorFlag is the error message used if a flag doesn't match the
	// regular expression defined.
	// If you want an unlocalized error, just fill Fallback.Other field of the
	// config.
	//
	// Available Placeholders are:
	//
	// 		• name - the full name of the flag
	// 		• used_name - the name of the flag the invoking user used
	// 		• raw - the raw flag without the flags name
	// 		• regexp - the regular expression that needs to be matched
	//
	// Defaults to: idRegexpNotMatchingErrorFlag
	RegexpErrorFlag *i18n.Config
}

var (
	SimpleAlphanumericID Type = new(AlphanumericID)
	_                    Type = AlphanumericID{}
)

func (id AlphanumericID) Name(l *i18n.Localizer) string {
	if id.CustomName != nil {
		name, err := id.CustomName.Get(l)
		if err == nil {
			return name
		}
	}

	name, _ := l.Localize(idName) // we have a fallback
	return name
}

func (id AlphanumericID) Description(l *i18n.Localizer) string {
	if id.CustomDescription != nil {
		desc, err := id.CustomDescription.Get(l)
		if err == nil {
			return desc
		}
	}

	desc, _ := l.Localize(idDescription) // we have a fallback
	return desc
}

func (id AlphanumericID) Parse(_ *state.State, ctx *Context) (interface{}, error) { //nolint:dupl
	if uint(len(ctx.Raw)) < id.MinLength {
		return nil, newArgumentError2(
			idBelowMinLengthErrorArg, idBelowMinLengthErrorFlag, ctx, map[string]interface{}{
				"min": id.MinLength,
			})
	} else if id.MaxLength > 0 && uint(len(ctx.Raw)) > id.MaxLength {
		return nil, newArgumentError2(
			idAboveMaxLengthErrorArg, idAboveMaxLengthErrorFlag, ctx, map[string]interface{}{
				"max": id.MaxLength,
			})
	}

	if id.Regexp != nil && !id.Regexp.MatchString(ctx.Raw) {
		if ctx.Kind == KindArg && id.RegexpErrorArg == nil {
			id.RegexpErrorArg = regexpNotMatchingErrorArg
		} else if ctx.Kind == KindFlag && id.RegexpErrorFlag == nil {
			id.RegexpErrorFlag = regexpNotMatchingErrorFlag
		}

		return nil, newArgumentError2(id.RegexpErrorArg, id.RegexpErrorFlag, ctx, map[string]interface{}{
			"regexp": id.Regexp.String(),
		})
	}

	return ctx.Raw, nil
}

func (id AlphanumericID) Default() interface{} {
	return ""
}
