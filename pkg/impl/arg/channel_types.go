package arg

import (
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/disstate/v3/pkg/state"

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
	// CategoryCancelEmoji is the emoji used as cancel emoji in a Category
	// chooser embed.
	CategoryCancelEmoji = discord.APIEmoji(emojiutil.CrossMarkButton)
	// CategoryOptionEmojis are the emojis used as options in a category
	// chooser embed.
	// It must contain at least 2 emojis.
	CategoryOptionEmojis = []discord.APIEmoji{
		discord.APIEmoji(emojiutil.Keycap1), discord.APIEmoji(emojiutil.Keycap2), discord.APIEmoji(emojiutil.Keycap3),
		discord.APIEmoji(emojiutil.Keycap4), discord.APIEmoji(emojiutil.Keycap5), discord.APIEmoji(emojiutil.Keycap6),
		discord.APIEmoji(emojiutil.Keycap7), discord.APIEmoji(emojiutil.Keycap8), discord.APIEmoji(emojiutil.Keycap9),
		discord.APIEmoji(emojiutil.Keycap10),
	}

	// CategoryChooserBuilder is the source *embedutil.Builder used to create
	// category chooser embeds.
	// When sending a chooser embed, title and description will be
	// set/overwritten and at most 2 fields will be added.
	CategoryChooserBuilder = embedutil.NewBuilder()

	// CategoryAllowSearch is a global flag that defines whether categories may
	// be referenced by name.
	// If multiple matches are found, Category might ask the user to choose a
	// category through a reaction driven chooser embed.
	CategoryAllowSearch = true
	// CategorySearchTimeout is the amount of time the user has to choose the
	// desired category from the chooser embed.
	CategorySearchTimeout = 20 * time.Second
)

var (
	// VoiceChannelCancelEmoji is the emoji used as cancel emoji in a
	// VoiceChannel chooser embed.
	VoiceChannelCancelEmoji = discord.APIEmoji(emojiutil.CrossMarkButton)
	// VoiceChannelOptionEmojis are the emojis used as options in a
	// VoiceChannel chooser embed.
	// It must contain at least 2 emojis.
	VoiceChannelOptionEmojis = []discord.APIEmoji{
		discord.APIEmoji(emojiutil.Keycap1), discord.APIEmoji(emojiutil.Keycap2), discord.APIEmoji(emojiutil.Keycap3),
		discord.APIEmoji(emojiutil.Keycap4), discord.APIEmoji(emojiutil.Keycap5), discord.APIEmoji(emojiutil.Keycap6),
		discord.APIEmoji(emojiutil.Keycap7), discord.APIEmoji(emojiutil.Keycap8), discord.APIEmoji(emojiutil.Keycap9),
		discord.APIEmoji(emojiutil.Keycap10),
	}

	// VoiceChannelChooserBuilder is the source *embedutil.Builder used to
	// create VoiceChannel chooser embeds.
	// When sending a chooser embed, title and description will be
	// set/overwritten and at most 2 fields will be added.
	VoiceChannelChooserBuilder = embedutil.NewBuilder()

	// VoiceChannelAllowSearch is a global flag that defines whether a
	// voice channels may be referenced by name.
	// If multiple matches are found, VoiceChannel might ask the user to choose
	// a category through a reaction driven chooser embed.
	VoiceChannelAllowSearch = true
	// VoiceChannelSearchTimeout is the amount of time the user has to choose
	// the desired voice channel from the chooser embed.
	VoiceChannelSearchTimeout = 20 * time.Second
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
			return nil, newArgumentError2(textChannelInvalidMentionErrorArg, textChannelInvalidMentionErrorFlag, ctx, nil)
		}

		c, err := s.Channel(discord.ChannelID(id))
		if err != nil {
			return nil, newArgumentError2(textChannelInvalidMentionErrorArg, textChannelInvalidMentionErrorFlag, ctx, nil)
		}

		if c.GuildID != ctx.GuildID {
			return nil, newArgumentError(textChannelGuildNotMatchingError, ctx, nil)
		} else if c.Type != discord.GuildText && c.Type != discord.GuildNews {
			return nil, newArgumentError(textChannelInvalidTypeError, ctx, nil)
		}

		return c, nil
	}

	if !TextChannelAllowIDs {
		return nil, newArgumentError(textChannelInvalidMentionWithRawError, ctx, nil)
	}

	id, err := discord.ParseSnowflake(ctx.Raw)
	if err != nil {
		return nil, newArgumentError(textChannelInvalidError, ctx, nil)
	}

	c, err := s.Channel(discord.ChannelID(id))
	if err != nil {
		return nil, newArgumentError(channelIDInvalidError, ctx, nil)
	}

	if c.GuildID != ctx.GuildID {
		return nil, newArgumentError(textChannelIDGuildNotMatchingError, ctx, nil)
	} else if c.Type != discord.GuildText && c.Type != discord.GuildNews {
		return nil, newArgumentError(textChannelIDInvalidTypeError, ctx, nil)
	}

	return c, nil
}

func (t textChannel) Default() interface{} {
	return (*discord.Channel)(nil)
}

// =============================================================================
// CategoryName
// =====================================================================================

// Category is the type used for channels of type category.
// A category can either be referenced by id or through name matching, if
// CategoryAllowSearch is true.
//
// If multiple categories match the given name, a reaction based chooser embed
// will be sent, that contains up to len(CategoryOptionEmojis) categories that
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

func (c category) Parse(s *state.State, ctx *Context) (interface{}, error) {
	id, err := discord.ParseSnowflake(ctx.Raw)
	if err == nil {
		channel, err := c.handleID(s, ctx, discord.ChannelID(id))
		if channel != nil || err != nil {
			return channel, err
		}
	}

	if !CategoryAllowSearch {
		return nil, newArgumentError2(categoryIDInvalidErrorArg, categoryIDInvalidErrorFlag, ctx, nil)
	}

	return c.handleName(s, ctx)
}

// handleID attempts to fetch the category with the passed id.
// It returns nil, nil if no such channel exists.
func (c category) handleID(s *state.State, ctx *Context, id discord.ChannelID) (*discord.Channel, error) {
	channel, err := s.Channel(id)
	if err == nil {
		if channel.Type != discord.GuildCategory {
			return nil, newArgumentError(categoryIDInvalidTypeError, ctx, nil)
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

type categoryMatch struct {
	channel *discord.Channel
	pos     int
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
		fullMatches = make([]categoryMatch, 0, len(CategoryOptionEmojis))

		partialMatches  = make([]categoryMatch, 0, len(CategoryOptionEmojis))
		partialOverflow = false
	)

	lowerRaw := strings.ToLower(ctx.Raw)

	for i, categories := range resolved[1:] {
		lowerName := strings.ToLower(categories[0].Name)

		if lowerName == lowerRaw {
			if len(fullMatches) >= len(CategoryOptionEmojis) {
				return nil, newArgumentError(categoryTooManyMatchesError, ctx, nil)
			}

			fullMatches = append(fullMatches, categoryMatch{
				channel: &categories[0],
				pos:     i,
			})
		} else if strings.Contains(lowerName, lowerRaw) {
			if len(partialMatches) >= len(CategoryOptionEmojis) {
				partialOverflow = true
				continue
			}

			partialMatches = append(partialMatches, categoryMatch{
				channel: &categories[0],
				pos:     i,
			})
		}
	}

	switch {
	case len(fullMatches) == 0 && len(partialMatches) == 0:
		return nil, newArgumentError(categoryNotFoundError, ctx, nil)
	case len(fullMatches) == 0 && partialOverflow:
		return nil, newArgumentError(categoryTooManyPartialMatchesError, ctx, nil)
	case len(fullMatches) == 1 && len(partialMatches) == 0:
		return fullMatches[0].channel, nil
	case len(fullMatches) == 0 && len(partialMatches) == 1:
		return partialMatches[0].channel, nil
	default:
		return c.sendChooser(s, ctx, fullMatches, partialMatches)
	}
}

func (c category) sendChooser( //nolint:dupl
	s *state.State, ctx *Context, fullMatches, partialMatches []categoryMatch,
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
		if err != nil && !discorderr.InRange(discorderr.As(err), discorderr.UnknownResource) {
			ctx.HandleErrorSilent(err)
		}
	}()

	choice, err := messageutil.NewReactionWaiter(s, ctx.Context, msg.ID).
		WithReactions(CategoryOptionEmojis[:numMatches]...).
		WithCancelReactions(CategoryCancelEmoji).
		NoAutoDelete(). // we will delete the whole message anyway
		Await(CategorySearchTimeout)
	if err != nil {
		return nil, err
	}

	var i int

	for i = 0; i < numMatches; i++ {
		if CategoryOptionEmojis[i] == choice {
			break
		}
	}

	if i < len(fullMatches) {
		return fullMatches[i].channel, nil
	}

	return partialMatches[i-len(fullMatches)].channel, nil
}

func (c category) genChooserEmbed(
	ctx *Context, fullMatches, partialMatches []categoryMatch,
) (chooser *embedutil.Builder, numMatches int, err error) {
	chooser = CategoryChooserBuilder.Clone().
		WithSimpleTitlel(categoryChooserTitle).
		WithDescriptionl(categoryChooserDescription.
			WithPlaceholders(categoryChooserDescriptionPlaceholders{
				CancelEmoji: CategoryCancelEmoji,
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
					Emoji:        CategoryOptionEmojis[numMatches],
					CategoryName: m.channel.Name,
					Position:     m.pos + 1,
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

	if len(partialMatches) > len(CategoryOptionEmojis)-numMatches && len(partialMatches) != 0 {
		chooser.WithFieldl(
			categoryChooserPartialMatchesName,
			categoryChooserTooManyPartialMatches.
				WithPlaceholders(categoryChooserTooManyPartialMatchesPlaceholders{
					NumPartialMatches: len(partialMatches),
				}))
	} else if len(partialMatches) <= len(CategoryOptionEmojis)-numMatches && len(partialMatches) != 0 {
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
					Emoji:        CategoryOptionEmojis[numMatches],
					CategoryName: m.channel.Name,
					Position:     m.pos + 1,
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

// =============================================================================
// VoiceChannel
// =====================================================================================

// VoiceChannel is the type used for channels of type voice.
// A VoiceChannel can either be referenced by id or through name matching, if
// VoiceChannelAllowSearch is true.
//
// If multiple voice channels match the given name, a reaction based chooser
// embed will be sent, that contains up to len(VoiceChannelOptionEmojis)
// categories that match the search.
// The user will then be asked to choose by clicking the corresponding
// reaction.
//
// Go type: *discord.Channel
var VoiceChannel Type = new(voiceChannel)

type voiceChannel struct{}

func (c voiceChannel) Name(l *i18n.Localizer) string {
	name, _ := l.Localize(voiceChannelName) // we have a fallback
	return name
}

func (c voiceChannel) Description(l *i18n.Localizer) string {
	desc, _ := l.Localize(voiceChannelDescription) // we have a fallback
	return desc
}

func (c voiceChannel) Parse(s *state.State, ctx *Context) (interface{}, error) {
	id, err := discord.ParseSnowflake(ctx.Raw)
	if err == nil {
		channel, err := c.handleID(s, ctx, discord.ChannelID(id))
		if channel != nil || err != nil {
			return channel, err
		}
	}

	if !VoiceChannelAllowSearch {
		return nil, newArgumentError2(voiceChannelIDInvalidErrorArg, voiceChannelIDInvalidErrorFlag, ctx, nil)
	}

	return c.handleName(s, ctx)
}

// handleID attempts to fetch the voice channel with the passed id.
// It returns nil, nil if no such channel exists.
func (c voiceChannel) handleID(s *state.State, ctx *Context, id discord.ChannelID) (*discord.Channel, error) {
	channel, err := s.Channel(id)
	if err == nil {
		if channel.Type != discord.GuildVoice {
			return nil, newArgumentError(voiceChannelIDInvalidTypeError, ctx, nil)
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

type voiceMatch struct {
	categoryName string // no parent, if len(categoryName) == 0
	channel      *discord.Channel
	pos          int
}

// handleName attempts to find a voice channel that matches ctx.Raw partially or
// fully.
// It ignores case.
func (c voiceChannel) handleName(s *state.State, ctx *Context) (*discord.Channel, error) {
	channels, err := s.Channels(ctx.GuildID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resolved := channelutil.ResolveCategories(channels)

	var (
		fullMatches = make([]voiceMatch, 0, len(VoiceChannelOptionEmojis))

		partialMatches  = make([]voiceMatch, 0, len(VoiceChannelOptionEmojis))
		partialOverflow = false
	)

	lowerRaw := strings.ToLower(ctx.Raw)

	for i, categories := range resolved {
		vcStart := findVoiceStart(categories)

		for j, c := range categories[vcStart:] {
			categoryName := ""
			if i > 0 {
				categoryName = categories[0].Name
			}

			lowerName := strings.ToLower(c.Name)

			if lowerName == lowerRaw {
				if len(fullMatches) >= len(VoiceChannelOptionEmojis) {
					return nil, newArgumentError(voiceChannelTooManyMatchesError, ctx, nil)
				}

				fullMatches = append(fullMatches, voiceMatch{
					categoryName: categoryName,
					channel:      &categories[vcStart+j],
					pos:          j,
				})
			} else if strings.Contains(lowerName, lowerRaw) {
				if len(partialMatches) >= len(CategoryOptionEmojis) {
					partialOverflow = true
					continue
				}

				partialMatches = append(partialMatches, voiceMatch{
					categoryName: categoryName,
					channel:      &categories[vcStart+j],
					pos:          j,
				})
			}
		}
	}

	switch {
	case len(fullMatches) == 0 && len(partialMatches) == 0:
		return nil, newArgumentError(voiceChannelNotFoundError, ctx, nil)
	case len(fullMatches) == 0 && partialOverflow:
		return nil, newArgumentError(voiceChannelTooManyPartialMatchesError, ctx, nil)
	case len(fullMatches) == 1 && len(partialMatches) == 0:
		return fullMatches[0].channel, nil
	case len(fullMatches) == 0 && len(partialMatches) == 1:
		return partialMatches[0].channel, nil
	default:
		return c.sendChooser(s, ctx, fullMatches, partialMatches)
	}
}

func (c voiceChannel) sendChooser( //nolint:dupl
	s *state.State, ctx *Context, fullMatches, partialMatches []voiceMatch,
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
		if err != nil && !discorderr.InRange(discorderr.As(err), discorderr.UnknownResource) {
			ctx.HandleErrorSilent(err)
		}
	}()

	choice, err := messageutil.NewReactionWaiter(s, ctx.Context, msg.ID).
		WithReactions(VoiceChannelOptionEmojis[:numMatches]...).
		WithCancelReactions(VoiceChannelCancelEmoji).
		NoAutoDelete(). // we will delete the whole message anyway
		Await(VoiceChannelSearchTimeout)
	if err != nil {
		return nil, err
	}

	var i int

	for i = 0; i < numMatches; i++ {
		if VoiceChannelOptionEmojis[i] == choice {
			break
		}
	}

	if i < len(fullMatches) {
		return fullMatches[i].channel, nil
	}

	return partialMatches[i-len(fullMatches)].channel, nil
}

func (c voiceChannel) genChooserEmbed( //nolint:dupl,funlen
	ctx *Context, fullMatches, partialMatches []voiceMatch,
) (chooser *embedutil.Builder, numMatches int, err error) {
	chooser = VoiceChannelChooserBuilder.Clone().
		WithSimpleTitlel(voiceChannelChooserTitle).
		WithDescriptionl(voiceChannelChooserDescription.
			WithPlaceholders(voiceChannelChooserDescriptionPlaceholders{
				CancelEmoji: VoiceChannelCancelEmoji,
			}))

	var b strings.Builder
	b.Grow(1024) // max field size

	if len(fullMatches) > 0 {
		fullMatchesName, err := ctx.Localize(voiceChannelChooserFullMatchesName)
		if err != nil {
			return nil, 0, err
		}

		for i, m := range fullMatches {
			if i > 0 {
				b.WriteRune('\n')
			}

			match, err := ctx.Localize(voiceChannelChooserMatch.
				WithPlaceholders(voiceChannelChooserMatchPlaceholders{
					Emoji:        VoiceChannelOptionEmojis[numMatches],
					CategoryName: m.categoryName,
					ChannelName:  m.channel.Name,
					Position:     m.pos + 1,
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

	if len(partialMatches) > len(VoiceChannelOptionEmojis)-numMatches && len(partialMatches) != 0 {
		chooser.WithFieldl(
			voiceChannelChooserPartialMatchesName,
			voiceChannelChooserTooManyPartialMatches.
				WithPlaceholders(voiceChannelChooserTooManyPartialMatchesPlaceholders{
					NumPartialMatches: len(partialMatches),
				}))
	} else if len(partialMatches) <= len(VoiceChannelOptionEmojis)-numMatches && len(partialMatches) != 0 {
		// add partialMatches if there are any, and there is still room
		partialMatchesName, err := ctx.Localize(voiceChannelChooserPartialMatchesName)
		if err != nil {
			return nil, 0, err
		}

		for i, m := range partialMatches {
			if i > 0 {
				b.WriteRune('\n')
			}

			match, err := ctx.Localize(voiceChannelChooserMatch.
				WithPlaceholders(voiceChannelChooserMatchPlaceholders{
					Emoji:        VoiceChannelOptionEmojis[numMatches],
					CategoryName: m.categoryName,
					ChannelName:  m.channel.Name,
					Position:     m.pos + 1,
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

func findVoiceStart(c []discord.Channel) int {
	return sort.Search(len(c), func(i int) bool {
		return c[i].Type == discord.GuildVoice
	})
}

func (c voiceChannel) Default() interface{} {
	return (*discord.Channel)(nil)
}
