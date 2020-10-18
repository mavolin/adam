package arg

import (
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
// A User can either be a mention, or an ID.
//
// Gp type: *discord.User
var User = &user{}

type user struct{}

func (u user) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(userName) // we have a fallback
	return name
}

func (u user) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(userDescription) // we have a fallback
	return desc
}

func (u user) Parse(s *state.State, ctx *Context) (interface{}, error) {
	if matches := userMentionRegexp.FindStringSubmatch(ctx.Raw); len(matches) > 1 {
		id := matches[1]

		uid, err := discord.ParseSnowflake(id)
		if err != nil { // range err
			return nil, newArgParsingErr(userInvalidMentionArg, userInvalidMentionFlag, ctx, nil)
		}

		for _, m := range ctx.Mentions {
			if m.ID == discord.UserID(uid) {
				return &m.User, nil
			}
		}

		user, err := s.User(discord.UserID(uid))
		if err != nil {
			return nil, newArgParsingErr(userInvalidMentionArg, userInvalidMentionFlag, ctx, nil)
		}

		return user, nil
	}

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

func (u user) Default() interface{} {
	return (*discord.User)(nil)
}
