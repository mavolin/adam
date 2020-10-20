package arg

import (
	"regexp"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
)

// =============================================================================
// User
// =====================================================================================

// User is the Type used to specify users globally.
// The User doesn't have to be on the same guild as the invoking one.
// In contrast to member, this can also be used in direct messages.
// A User can either be a mention, or an id.
//
// Gp type: *discord.User
var User = new(user)

type user struct{}

func (u user) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(userName) // we have a fallback
	return name
}

func (u user) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(userDescription) // we have a fallback
	return desc
}

var userMentionRegexp = regexp.MustCompile(`^<@!?(?P<id>\d+)>$`)

func (u user) Parse(s *state.State, ctx *Context) (interface{}, error) {
	if matches := userMentionRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		rawID := matches[1]

		id, err := discord.ParseSnowflake(rawID)
		if err != nil { // range err
			return nil, newArgParsingErr(userInvalidMentionArg, userInvalidMentionFlag, ctx, nil)
		}

		for _, m := range ctx.Mentions {
			if m.ID == discord.UserID(id) {
				return &m.User, nil
			}
		}

		user, err := s.User(discord.UserID(id))
		if err != nil {
			return nil, newArgParsingErr(userInvalidMentionArg, userInvalidMentionFlag, ctx, nil)
		}

		return user, nil
	}

	id, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgParsingErr(userInvalidIDWithRaw, userInvalidIDWithRaw, ctx, nil)
	}

	user, err := s.User(discord.UserID(id))
	if err != nil {
		return nil, newArgParsingErr(userInvalidIDArg, userInvalidIDFlag, ctx, nil)
	}

	return user, nil
}

func (u user) Default() interface{} {
	return (*discord.User)(nil)
}

// =============================================================================
// UserID
// =====================================================================================

// UserID is the same as a User, but it only accepts ids.
//
// Go type: *discord.User
var UserID = new(userID)

type userID struct{}

func (u userID) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(userIDName) // we have fallback
	return name
}

func (u userID) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(userIDDescription) // we have a fallback
	return desc
}

func (u userID) Parse(s *state.State, ctx *Context) (interface{}, error) {
	uid, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgParsingErr(userInvalidIDWithRaw, userInvalidIDWithRaw, ctx, nil)
	}

	user, err := s.User(discord.UserID(uid))
	if err != nil {
		return nil, newArgParsingErr(userInvalidIDArg, userInvalidIDFlag, ctx, nil)
	}

	return user, nil
}

func (u userID) Default() interface{} {
	return (*discord.User)(nil)
}
