package arg

import (
	"regexp"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/channelutil"
	"github.com/mavolin/adam/pkg/utils/discorderr"
	"github.com/mavolin/adam/pkg/utils/embedutil"
	emojiutil "github.com/mavolin/adam/pkg/utils/emoji"
	"github.com/mavolin/adam/pkg/utils/messageutil"
)

// TextChannelAllowIDs is a global flag that defines whether TextChannels may
// also be noted as plain Snowflakes.
var TextChannelAllowIDs = false

var (
	// ChooserCancelEmoji is the emoji used as cancel emoji in a chooser embed.
	ChooserCancelEmoji = emojiutil.CrossMarkButton
	// ChooserOptionEmojis are the emojis used as options in a chooser embed.
	// It must contain at leas 2 emojis.
	ChooserOptionEmojis = []api.Emoji{
		emojiutil.Keycap1, emojiutil.Keycap2, emojiutil.Keycap3, emojiutil.Keycap4, emojiutil.Keycap5,
		emojiutil.Keycap6, emojiutil.Keycap7, emojiutil.Keycap8, emojiutil.Keycap9, emojiutil.Keycap10,
	}

	// ChooserBuilder is the source embedutil.Builder used to create chooser
	// embeds.
	// When sending a chooser embed, title and description will be
	// set/overwritten and at most 2 fields will be added.
	ChooserBuilder = embedutil.NewBuilder()
)

var (
	// CategoryAllowSearch is a global flag that defines whether Categories may
	// be referenced by name.
	// If multiple matches are found, Category might ask the user to choose a
	// category through a reaction driven chooser embed.
	CategoryAllowSearch = true
	// CategorySearchTimeout is the amount of time the user has to choose the
	// desired category from the chooser embed.
	CategorySearchTimeout = 20 * time.Second
)

// =============================================================================
// TextChannel
// =====================================================================================

// TextChannel is the Type used for guild text channels and news channels.
// The channel must be on the same guild as the invoking one.
//
// TextChannel will always fail if used in a direct message.
//
// Go type: *discord.Channel
var TextChannel Type = new(textChannel)

type textChannel struct{}

func (t textChannel) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(textChannelName) // we have a fallback
	return name
}

func (t textChannel) Description(l *i18n.Localizer) string {
	if TextChannelAllowIDs {
		desc, err := l.Localize(textChannelDescriptionWithID)
		if err == nil {
			return desc
		}
	}

	desc, _ := l.Localize(textChannelDescriptionNoID) // we have a fallback
	return desc
}

var textChannelMentionRegexp = regexp.MustCompile(`^<#(?P<id>\d+)>$`)

func (t textChannel) Parse(s *state.State, ctx *Context) (interface{}, error) {
	err := restriction.ChannelTypes(plugin.GuildChannels)(s, ctx.Context)
	if err != nil {
		return nil, err
	}

	if matches := textChannelMentionRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		rawID := matches[1]

		id, err := discord.ParseSnowflake(rawID)
		if err != nil { // range err
			return nil, newArgParsingErr2(textChannelInvalidMentionErrorArg, textChannelInvalidMentionErrorFlag, ctx, nil)
		}

		c, err := s.Channel(discord.ChannelID(id))
		if err != nil {
			return nil, newArgParsingErr2(textChannelInvalidMentionErrorArg, textChannelInvalidMentionErrorFlag, ctx, nil)
		}

		if c.GuildID != ctx.GuildID {
			return nil, newArgParsingErr(textChannelGuildNotMatchingError, ctx, nil)
		} else if c.Type != discord.GuildText && c.Type != discord.GuildNews {
			return nil, newArgParsingErr(textChannelInvalidTypeError, ctx, nil)
		}

		return c, nil
	}

	if !TextChannelAllowIDs {
		return nil, newArgParsingErr(textChannelInvalidMentionWithRawError, ctx, nil)
	}

	id, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgParsingErr(textChannelInvalidError, ctx, nil)
	}

	c, err := s.Channel(discord.ChannelID(id))
	if err != nil {
		return nil, newArgParsingErr(channelIDInvalidError, ctx, nil)
	}

	if c.GuildID != ctx.GuildID {
		return nil, newArgParsingErr(textChannelIDGuildNotMatchingError, ctx, nil)
	} else if c.Type != discord.GuildText && c.Type != discord.GuildNews {
		return nil, newArgParsingErr(textChannelIDInvalidTypeError, ctx, nil)
	}

	return c, nil
}

func (t textChannel) Default() interface{} {
	return (*discord.Channel)(nil)
}

// =============================================================================
// Category
// =====================================================================================

// Category is the type used for channels of type category.
// A category can either be referenced by id or through name matching, if
// CategoryAllowSearch is true.
//
// If multiple categories match the given name, a reaction based chooser embed
// will be sent, that contains up to len(ChooserOptionEmojis) categories that
// match the search.
// The user will then be asked to choose by clicking the corresponding
// reaction.
//
// Go type: *discord.Channel
var Category Type = new(category)

type category struct{}

func (c category) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(categoryName) // we have a fallback
	return name
}

func (c category) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(categoryDescription) // we have a fallback
	return desc
}

type match struct {
	channel discord.Channel
	pos     int
}

func (c category) Parse(s *state.State, ctx *Context) (interface{}, error) {
	id, err := discord.ParseSnowflake(ctx.Raw)
	if err == nil {
		channel, err := c.handleID(s, ctx, discord.ChannelID(id))
		if channel != nil || err != nil {
			return channel, err
		}
	}

	if !CategoryAllowSearch {
		return nil, newArgParsingErr2(categoryIDInvalidErrorArg, categoryIDInvalidErrorFlag, ctx, nil)
	}

	return c.handleName(s, ctx)
}

// handleID attempts to fetch Channel with the passed id.
// It returns nil, nil if no such channel exists.
func (c category) handleID(s *state.State, ctx *Context, id discord.ChannelID) (*discord.Channel, error) {
	channel, err := s.Channel(id)
	if err == nil {
		if channel.Type != discord.GuildCategory {
			return nil, newArgParsingErr(categoryIDInvalidTypeError, ctx, nil)
		}

		return channel, err
	}

	// the channel name might be a num, and the arg we received was
	// therefore not an id
	if discorderr.Is(discorderr.As(err), discorderr.UnknownChannel) {
		return nil, nil
	}

	// something else went wrong, capture this
	return nil, errors.WithStack(err)
}

// handleName attempts to find a category that matches ctx.Raw partially or
// fully.
// It ignores case.
func (c category) handleName(s *state.State, ctx *Context) (*discord.Channel, error) {
	channels, err := s.Channels(ctx.GuildID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resolved := channelutil.ResolveCategories(channels)

	var (
		fullMatches = make([]match, 0, len(ChooserOptionEmojis))

		partialMatches  = make([]match, 0, len(ChooserOptionEmojis))
		partialOverflow = false
	)

	lowerRaw := strings.ToLower(ctx.Raw)

	for i, channels := range resolved[1:] {
		lowerName := strings.ToLower(channels[0].Name)

		if lowerName == lowerRaw {
			if len(fullMatches) >= len(ChooserOptionEmojis) {
				return nil, newArgParsingErr(categoryTooManyMatchesError, ctx, nil)
			}

			fullMatches = append(fullMatches, match{
				channel: channels[0],
				pos:     i,
			})
		} else if strings.Contains(lowerName, lowerRaw) {
			if len(partialMatches) >= len(ChooserOptionEmojis) {
				partialOverflow = true
				continue
			}

			partialMatches = append(partialMatches, match{
				channel: channels[0],
				pos:     i,
			})
		}
	}

	switch {
	case len(fullMatches) == 0 && len(partialMatches) == 0:
		return nil, newArgParsingErr(categoryNotFoundError, ctx, nil)
	case len(fullMatches) == 0 && partialOverflow:
		return nil, newArgParsingErr(categoryTooManyPartialMatchesError, ctx, nil)
	case len(fullMatches) == 1 && len(partialMatches) == 0:
		return &fullMatches[0].channel, nil
	case len(fullMatches) == 0 && len(partialMatches) == 1:
		return &partialMatches[0].channel, nil
	default:
		return c.sendChooser(s, ctx, fullMatches, partialMatches)
	}
}

func (c category) sendChooser(
	s *state.State, ctx *Context, fullMatches, partialMatches []match,
) (*discord.Channel, error) {
	chooser, numMatches, err := c.genChooserEmbed(ctx, fullMatches, partialMatches)
	if err != nil {
		return nil, err
	}

	msg, err := ctx.ReplyEmbedBuilder(chooser)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := s.DeleteMessage(msg.ChannelID, msg.ID)
		if discorderr.InRange(discorderr.As(err), discorderr.UnknownResource) {
			ctx.HandleErrorSilent(err)
		}
	}()

	choice, err := messageutil.NewReactionWaiter(s, ctx.Context, msg.ID).
		WithReactions(ChooserOptionEmojis[:numMatches]...).
		WithCancelReactions(ChooserCancelEmoji).
		NoAutoDelete(). // we are gonna delete the whole message anyway
		Await(CategorySearchTimeout)
	if err != nil {
		var terr *messageutil.TimeoutError
		if errors.As(err, &terr) {
			err = errors.Abort
		}

		return nil, errors.WithStack(err)
	}

	var i int

	for i = 0; i < numMatches; i++ {
		if ChooserOptionEmojis[i] == choice {
			break
		}
	}

	if i < len(fullMatches) {
		return &fullMatches[i].channel, nil
	}

	return &partialMatches[i-len(fullMatches)].channel, nil
}

func (c category) genChooserEmbed(
	ctx *Context, fullMatches, partialMatches []match,
) (chooser *embedutil.Builder, numMatches int, err error) {
	chooser = ChooserBuilder.Clone().
		WithSimpleTitlel(categoryChooserTitle).
		WithDescriptionl(categoryChooserDescription.
			WithPlaceholders(categoryChooserDescriptionPlaceholders{
				CancelEmoji: ChooserCancelEmoji,
			}))

	var b strings.Builder
	b.Grow(1024) // max field size

	if len(fullMatches) > 0 {
		fullMatchesName, err := ctx.Localize(categoryChooserFullMatchesName)
		if err != nil {
			return nil, 0, err
		}

		for i, m := range fullMatches {
			if i > 0 {
				b.WriteRune('\n')
			}

			match, err := ctx.Localize(categoryChooserMatch.
				WithPlaceholders(categoryChooserMatchPlaceholders{
					Emoji:           ChooserOptionEmojis[numMatches],
					ChannelName:     m.channel.Name,
					ChannelPosition: m.pos + 1,
				}))
			if err != nil {
				return nil, 0, err
			}

			b.WriteString(match)

			numMatches++
		}

		chooser.WithField(fullMatchesName, b.String())
		b.Reset()
	}

	if len(partialMatches) > len(ChooserOptionEmojis)-numMatches && len(partialMatches) != 0 {
		chooser.WithFieldl(
			categoryChooserPartialMatchesName,
			categoryChooserTooManyPartialMatches.
				WithPlaceholders(categoryChooserTooManyPartialMatchesPlaceholders{
					NumPartialMatches: len(partialMatches),
				}))
	} else if len(partialMatches) <= len(ChooserOptionEmojis)-numMatches && len(partialMatches) != 0 {
		// add partialMatches if there are any, and there is still room
		partialMatchesName, err := ctx.Localize(categoryChooserPartialMatchesName)
		if err != nil {
			return nil, 0, err
		}

		for i, m := range partialMatches {
			if i > 0 {
				b.WriteRune('\n')
			}

			match, err := ctx.Localize(categoryChooserMatch.
				WithPlaceholders(categoryChooserMatchPlaceholders{
					Emoji:           ChooserOptionEmojis[numMatches],
					ChannelName:     m.channel.Name,
					ChannelPosition: m.pos + 1,
				}))
			if err != nil {
				return nil, 0, err
			}

			b.WriteString(match)

			numMatches++
		}

		chooser.WithField(partialMatchesName, b.String())
	}

	return chooser, numMatches, nil
}

func (c category) Default() interface{} {
	return (*discord.Channel)(nil)
}
