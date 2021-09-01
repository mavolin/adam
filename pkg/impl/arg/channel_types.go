package arg

import (
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/channelutil"
	"github.com/mavolin/adam/pkg/utils/discorderr"
	"github.com/mavolin/adam/pkg/utils/msgbuilder"
)

// TextChannelAllowIDs is a global flag that defines whether TextChannels may
// also be noted as plain Snowflakes.
var TextChannelAllowIDs = false

var (
	// CategoryAllowSearch is a global flag that defines whether categories may
	// be referenced by name.
	// If multiple matches are found, Category might ask the user to choose a
	// category through a reaction driven chooser embed.
	CategoryAllowSearch = true
	// CategoryChooserTimeout is the amount of time the user has to choose the
	// desired category from the chooser embed.
	CategoryChooserTimeout = 20 * time.Second
)

var (
	// VoiceChannelAllowSearch is a global flag that defines whether a
	// voice channels may be referenced by name.
	// If multiple matches are found, VoiceChannel might ask the user to choose
	// a category through a reaction driven chooser embed.
	VoiceChannelAllowSearch = true
	// VoiceChannelChooserTimeout is the amount of time the user has to choose
	// the desired voice channel from the chooser embed.
	VoiceChannelChooserTimeout = 20 * time.Second
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
var TextChannel plugin.ArgType = new(textChannel)

type textChannel struct{}

func (t textChannel) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(textChannelName) // we have a fallback
	return name
}

func (t textChannel) GetDescription(l *i18n.Localizer) string {
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

func (t textChannel) Parse(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	err := restriction.ChannelTypes(plugin.GuildChannels)(s, ctx.Context)
	if err != nil {
		return nil, err
	}

	if matches := textChannelMentionRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		rawID := matches[1]

		id, err := discord.ParseSnowflake(rawID)
		if err != nil { // range err
			return nil,
				newArgumentError2(textChannelInvalidMentionErrorArg, textChannelInvalidMentionErrorFlag, ctx, nil)
		}

		c, err := s.Channel(discord.ChannelID(id))
		if err != nil {
			return nil,
				newArgumentError2(textChannelInvalidMentionErrorArg, textChannelInvalidMentionErrorFlag, ctx, nil)
		}

		if c.GuildID != ctx.GuildID || (c.Type != discord.GuildText && c.Type != discord.GuildNews) {
			return nil,
				newArgumentError2(textChannelInvalidMentionErrorArg, textChannelInvalidMentionErrorFlag, ctx, nil)
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

func (t textChannel) GetDefault() interface{} {
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
var Category plugin.ArgType = new(category)

type category struct{}

func (c category) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(categoryName) // we have a fallback
	return name
}

func (c category) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(categoryDescription) // we have a fallback
	return desc
}

func (c category) Parse(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	id, err := discord.ParseSnowflake(ctx.Raw)
	if err == nil {
		channel, err := c.handleID(s, ctx, discord.ChannelID(id))
		if channel != nil || err != nil {
			return channel, err
		}
	}

	//goland:noinspection GoBoolExpressions
	if !CategoryAllowSearch {
		return nil, newArgumentError2(categoryIDInvalidErrorArg, categoryIDInvalidErrorFlag, ctx, nil)
	}

	return c.handleName(s, ctx)
}

// handleID attempts to fetch the category with the passed id.
// It returns nil, nil if no such channel exists.
func (c category) handleID(s *state.State, ctx *plugin.ParseContext, id discord.ChannelID) (*discord.Channel, error) {
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

const maxCategoryMatches = 24

// handleName attempts to find a category that matches ctx.Raw partially or
// fully.
// It ignores case.
func (c category) handleName(s *state.State, ctx *plugin.ParseContext) (*discord.Channel, error) {
	channels, err := s.Channels(ctx.GuildID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resolved := channelutil.ResolveCategories(channels)

	var (
		fullMatches = make([]categoryMatch, 0, maxCategoryMatches)

		partialMatches  = make([]categoryMatch, 0, maxCategoryMatches)
		partialOverflow = false
	)

	lowerRaw := strings.ToLower(ctx.Raw)

	for i, categories := range resolved[1:] {
		lowerName := strings.ToLower(categories[0].Name)

		if lowerName == lowerRaw {
			if len(fullMatches) >= maxCategoryMatches {
				return nil, newArgumentError(categoryTooManyMatchesError, ctx, nil)
			}

			fullMatches = append(fullMatches, categoryMatch{
				channel: &categories[0],
				pos:     i,
			})
		} else if strings.Contains(lowerName, lowerRaw) {
			if len(partialMatches) >= maxCategoryMatches {
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

//nolint:dupl
func (c category) sendChooser(
	s *state.State, ctx *plugin.ParseContext, fullMatches, partialMatches []categoryMatch,
) (*discord.Channel, error) {
	content, err := ctx.Localize(categoryChooserContent)
	if err != nil {
		return nil, err
	}

	if len(fullMatches)+len(partialMatches) > maxCategoryMatches {
		contentAmend, err := ctx.Localize(categoryChooserPartialMatchesAddition.
			WithPlaceholders(categoryChooserPartialMatchesAdditionPlaceholders{
				NumPartialMatches: len(partialMatches),
			}))
		if err != nil {
			return nil, err
		}

		content += "\n" + contentAmend
	}

	var result *discord.Channel

	selectBuilder := msgbuilder.NewSelect(&result).
		WithDefault(msgbuilder.NewSelectOptionl(categoryChooserCancel, (*discord.Channel)(nil)))

	for _, match := range fullMatches {
		label := categoryChooserMatch.
			WithPlaceholders(categoryChooserMatchPlaceholders{
				CategoryName: match.channel.Name,
				Position:     match.pos,
			})

		selectBuilder.With(msgbuilder.NewSelectOptionl(label, match.channel))
	}

	if maxCategoryMatches-len(fullMatches) >= len(partialMatches) {
		for _, match := range partialMatches {
			label := categoryChooserMatch.
				WithPlaceholders(categoryChooserMatchPlaceholders{
					CategoryName: match.channel.Name,
					Position:     match.pos,
				})

			selectBuilder.With(msgbuilder.NewSelectOptionl(label, match.channel))
		}
	}

	_, err = msgbuilder.New(s, ctx.Context).
		WithContent(content).
		WithAwaitedComponent(selectBuilder).
		ReplyAndAwait(CategoryChooserTimeout)
	return result, err
}

func (c category) GetDefault() interface{} {
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
var VoiceChannel plugin.ArgType = new(voiceChannel)

type voiceChannel struct{}

func (c voiceChannel) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(voiceChannelName) // we have a fallback
	return name
}

func (c voiceChannel) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(voiceChannelDescription) // we have a fallback
	return desc
}

func (c voiceChannel) Parse(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	id, err := discord.ParseSnowflake(ctx.Raw)
	if err == nil {
		channel, err := c.handleID(s, ctx, discord.ChannelID(id))
		if channel != nil || err != nil {
			return channel, err
		}
	}

	//goland:noinspection GoBoolExpressions
	if !VoiceChannelAllowSearch {
		return nil, newArgumentError2(voiceChannelIDInvalidErrorArg, voiceChannelIDInvalidErrorFlag, ctx, nil)
	}

	return c.handleName(s, ctx)
}

// handleID attempts to fetch the voice channel with the passed id.
// It returns nil, nil if no such channel exists.
func (c voiceChannel) handleID(
	s *state.State, ctx *plugin.ParseContext, id discord.ChannelID,
) (*discord.Channel, error) {
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

const maxVoiceMatches = 24

// handleName attempts to find a voice channel that matches ctx.Raw partially or
// fully.
// It ignores case.
func (c voiceChannel) handleName(s *state.State, ctx *plugin.ParseContext) (*discord.Channel, error) {
	channels, err := s.Channels(ctx.GuildID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resolved := channelutil.ResolveCategories(channels)

	var (
		fullMatches = make([]voiceMatch, 0, maxVoiceMatches)

		partialMatches  = make([]voiceMatch, 0, maxVoiceMatches)
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
				if len(fullMatches) >= maxVoiceMatches {
					return nil, newArgumentError(voiceChannelTooManyMatchesError, ctx, nil)
				}

				fullMatches = append(fullMatches, voiceMatch{
					categoryName: categoryName,
					channel:      &categories[vcStart+j],
					pos:          j,
				})
			} else if strings.Contains(lowerName, lowerRaw) {
				if len(partialMatches) >= maxVoiceMatches {
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

//nolint:dupl
func (c voiceChannel) sendChooser(
	s *state.State, ctx *plugin.ParseContext, fullMatches, partialMatches []voiceMatch,
) (*discord.Channel, error) {
	content, err := ctx.Localize(voiceChannelChooserContent)
	if err != nil {
		return nil, err
	}

	if len(fullMatches)+len(partialMatches) > maxCategoryMatches {
		contentAmend, err := ctx.Localize(voiceChannelChooserPartialMatchesAddition.
			WithPlaceholders(voiceChannelChooserPartialMatchesAdditionPlaceholders{
				NumPartialMatches: len(partialMatches),
			}))
		if err != nil {
			return nil, err
		}

		content += "\n" + contentAmend
	}

	var result *discord.Channel

	selectBuilder := msgbuilder.NewSelect(&result).
		WithDefault(msgbuilder.NewSelectOptionl(voiceChannelChooserCancel, (*discord.Channel)(nil)))

	for _, match := range fullMatches {
		if match.categoryName == "" {
			label := voiceChannelChooserRootMatch.
				WithPlaceholders(voiceChannelChooserMatchPlaceholders{
					ChannelName: match.channel.Name,
					Position:    match.pos,
				})
			selectBuilder.With(msgbuilder.NewSelectOptionl(label, match.channel))
		} else {
			label := voiceChannelChooserNestedMatch.
				WithPlaceholders(voiceChannelChooserMatchPlaceholders{
					CategoryName: match.categoryName,
					ChannelName:  match.channel.Name,
					Position:     match.pos,
				})
			selectBuilder.With(msgbuilder.NewSelectOptionl(label, match.channel))
		}
	}

	if maxCategoryMatches-len(fullMatches) >= len(partialMatches) {
		for _, match := range partialMatches {
			if match.categoryName == "" {
				label := voiceChannelChooserRootMatch.
					WithPlaceholders(voiceChannelChooserMatchPlaceholders{
						ChannelName: match.channel.Name,
						Position:    match.pos,
					})
				selectBuilder.With(msgbuilder.NewSelectOptionl(label, match.channel))
			} else {
				label := voiceChannelChooserNestedMatch.
					WithPlaceholders(voiceChannelChooserMatchPlaceholders{
						CategoryName: match.categoryName,
						ChannelName:  match.channel.Name,
						Position:     match.pos,
					})
				selectBuilder.With(msgbuilder.NewSelectOptionl(label, match.channel))
			}
		}
	}

	_, err = msgbuilder.New(s, ctx.Context).
		WithContent(content).
		WithAwaitedComponent(selectBuilder).
		ReplyAndAwait(VoiceChannelChooserTimeout)
	return result, err
}

func findVoiceStart(c []discord.Channel) int {
	return sort.Search(len(c), func(i int) bool {
		return c[i].Type == discord.GuildVoice
	})
}

func (c voiceChannel) GetDefault() interface{} {
	return (*discord.Channel)(nil)
}
