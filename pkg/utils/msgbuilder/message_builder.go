// Package msgbuilder provides builders to create messages, embeds and
// components.
// Its main purpose is to abstract away the handling of message components, as
// well as to ease working with localized messages.
package msgbuilder

import (
	"context"
	"io"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/internal/errorutil"
	"github.com/mavolin/adam/pkg/i18n"
)

// Builder is a message builder.
type Builder struct {
	s         *state.State
	localizer *i18n.Localizer

	content *i18n.Config
	embeds  *[]*EmbedBuilder

	components   *[]TopLevelComponentBuilder
	awaitIndexes map[int]struct{}

	tts bool

	allowedMentions *api.AllowedMentions
	reference       *discord.MessageReference

	flags *discord.MessageFlags

	attachments *[]discord.Attachment
	files       []sendpart.File

	channelID discord.ChannelID
	messageID discord.MessageID
}

// New creates a new *Builder.
func New(s *state.State, l *i18n.Localizer) *Builder {
	return &Builder{s: s, localizer: l, awaitIndexes: make(map[int]struct{})}
}

// WithContent sets the content of the message to the given content.
// It may be no longer than 2000 characters.
//
// Actions: create and edit
func (b *Builder) WithContent(content string) *Builder {
	return b.WithContentl(i18n.NewStaticConfig(content))
}

// WithContentlt sets the content of the message to the given content.
// It may be no longer than 2000 characters.
//
// Actions: create and edit
func (b *Builder) WithContentlt(content i18n.Term) *Builder {
	return b.WithContentl(content.AsConfig())
}

// WithContentl sets the content of the message to the given content.
// It may be no longer than 2000 characters.
//
// Actions: create and edit
func (b *Builder) WithContentl(content *i18n.Config) *Builder {
	b.content = content
	return b
}

// RemoveContent removes the content from the message.
//
// Actions: edit
func (b *Builder) RemoveContent() *Builder {
	b.content = i18n.EmptyConfig
	return b
}

// WithEmbed adds the passed embed to the message.
//
// Actions: create and edit
func (b *Builder) WithEmbed(embed *EmbedBuilder) *Builder {
	if b.embeds == nil {
		b.embeds = new([]*EmbedBuilder)
	}

	*b.embeds = append(*b.embeds, embed)
	return b
}

// RemoveEmbeds removes the embeds from the message.
//
// Actions: edit
func (b *Builder) RemoveEmbeds() *Builder {
	var nilEmbeds []*EmbedBuilder
	b.embeds = &nilEmbeds

	return b
}

// WithComponent adds the passed TopLevelComponentBuilder to the message.
//
// Actions: create and edit
func (b *Builder) WithComponent(component TopLevelComponentBuilder) *Builder {
	if b.components == nil {
		b.components = new([]TopLevelComponentBuilder)
	}

	*b.components = append(*b.components, component)
	return b
}

// WithAwaitedComponent adds the passed TopLevelComponentBuilder to the
// message, and waits for an interaction for that component when
// AwaitComponents or AwaitAllComponents is called.
//
// Actions: create and edit
func (b *Builder) WithAwaitedComponent(component TopLevelComponentBuilder) *Builder {
	if b.components == nil {
		b.components = new([]TopLevelComponentBuilder)
	}

	*b.components = append(*b.components, component)
	b.awaitIndexes[len(*b.components)-1] = struct{}{}

	return b
}

// RemoveComponents removes the components from the message.
//
// Actions: edit
func (b *Builder) RemoveComponents() *Builder {
	var nilComponents []TopLevelComponentBuilder
	b.components = &nilComponents

	return b
}

// AsTTS sends the message as a TTS message.
//
// Actions: create
func (b *Builder) AsTTS() *Builder {
	b.tts = true
	return b
}

// WithAllowedMentionTypes adds the passed api.AllowedMentionTypes to the
// allowed mentions of the message.
//
// Actions: create and edit
func (b *Builder) WithAllowedMentionTypes(allowed ...api.AllowedMentionType) *Builder {
	if b.allowedMentions == nil {
		b.allowedMentions = new(api.AllowedMentions)
	}

	b.allowedMentions.Parse = append(b.allowedMentions.Parse, allowed...)
	return b
}

// WithRoleMentions adds the passed role ids to the allowed mentions.
//
// Actions: create and edit
func (b *Builder) WithRoleMentions(roleIDs ...discord.RoleID) *Builder {
	if b.allowedMentions == nil {
		b.allowedMentions = new(api.AllowedMentions)
	}

	b.allowedMentions.Roles = append(b.allowedMentions.Roles, roleIDs...)
	return b
}

// WithUserMentions adds the passed user ids to the allowed mentions.
//
// Actions: create and edit
func (b *Builder) WithUserMentions(userIDs ...discord.UserID) *Builder {
	if b.allowedMentions == nil {
		b.allowedMentions = new(api.AllowedMentions)
	}

	b.allowedMentions.Users = append(b.allowedMentions.Users, userIDs...)
	return b
}

// MentionRepliedUser allows the replied user to be mentioned.
//
// Actions: create and edit
func (b *Builder) MentionRepliedUser() *Builder {
	if b.allowedMentions == nil {
		b.allowedMentions = new(api.AllowedMentions)
	}

	b.allowedMentions.RepliedUser = option.True
	return b
}

// WithReference references the message with the passed id.
//
// Actions: create
func (b *Builder) WithReference(messageID discord.MessageID) *Builder {
	b.reference = &discord.MessageReference{MessageID: messageID}
	return b
}

// WithFlags sets the passed messages flags.
//
// Actions: edit
func (b *Builder) WithFlags(flags ...discord.MessageFlags) *Builder {
	if b.flags == nil {
		b.flags = new(discord.MessageFlags)
	}

	for _, f := range flags {
		*b.flags |= f
	}

	return b
}

// KeepAttachment adds the passed discord.Attachment to the list of
// attachments to keep.
//
// Actions: edit
func (b *Builder) KeepAttachment(attachment discord.Attachment) *Builder {
	if b.attachments == nil {
		b.attachments = new([]discord.Attachment)
	}

	*b.attachments = append(*b.attachments, attachment)
	return b
}

// RemoveAttachments removes all attachments from the message.
//
// Actions: edit
func (b *Builder) RemoveAttachments() *Builder {
	var nilAttachments []discord.Attachment
	b.attachments = &nilAttachments

	return b
}

// WithFile adds the passed file to the message.
//
// Actions: create and edit
func (b *Builder) WithFile(name string, reader io.Reader) *Builder {
	b.files = append(b.files, sendpart.File{Name: name, Reader: reader})
	return b
}

// disableAllComponents disables all components.
func (b *Builder) disableAllComponents() {
	if b.components != nil {
		for _, c := range *b.components {
			c.disable()
		}
	}
}

// Send creates a new message in the passed channel.
func (b *Builder) Send(channelID discord.ChannelID) (msg *discord.Message, err error) {
	data := api.SendMessageData{
		TTS:             b.tts,
		Files:           b.files,
		AllowedMentions: b.allowedMentions,
		Reference:       b.reference,
	}

	if b.content != nil {
		data.Content, err = b.localizer.Localize(b.content)
		if err != nil {
			return nil, err
		}
	}

	if b.embeds != nil && len(*b.embeds) > 0 {
		data.Embeds = make([]discord.Embed, len(*b.embeds))

		for i, embedBuilder := range *b.embeds {
			embed, err := embedBuilder.Build(b.localizer)
			if err != nil {
				return nil, err
			}

			data.Embeds[i] = embed
		}
	}

	if b.components != nil && len(*b.components) > 0 {
		data.Components = make([]discord.Component, len(*b.components))

		for i, componentBuilder := range *b.components {
			c, err := componentBuilder.Build(b.localizer)
			if err != nil {
				return nil, err
			}

			data.Components[i] = c
		}
	}

	msg, err = b.s.SendMessageComplex(channelID, data)
	if err == nil {
		b.channelID = msg.ChannelID
		b.messageID = msg.ID
	}

	return msg, errorutil.WithStack(err)
}

// Edit edits the message with the passed channel and message id.
func (b *Builder) Edit(channelID discord.ChannelID, messageID discord.MessageID) (*discord.Message, error) {
	b.channelID = channelID
	b.messageID = messageID

	data := api.EditMessageData{
		Content:         nil,
		Embeds:          nil,
		Components:      nil,
		AllowedMentions: b.allowedMentions,
		Attachments:     b.attachments,
		Flags:           b.flags,
		Files:           b.files,
	}

	if b.content != nil {
		content, err := b.localizer.Localize(b.content)
		if err != nil {
			return nil, err
		}

		data.Content = option.NewNullableString(content)
	}

	if b.embeds != nil && *b.embeds == nil {
		var nilEmbeds []discord.Embed
		data.Embeds = &nilEmbeds
	} else if b.embeds != nil && len(*b.embeds) > 0 {
		data.Embeds = new([]discord.Embed)
		*data.Embeds = make([]discord.Embed, len(*b.embeds))

		for i, embedBuilder := range *b.embeds {
			embed, err := embedBuilder.Build(b.localizer)
			if err != nil {
				return nil, err
			}

			(*data.Embeds)[i] = embed
		}
	}

	if b.components != nil && *b.components == nil {
		var nilComponents []discord.Component
		data.Components = &nilComponents
	} else if b.components != nil && len(*b.components) > 0 {
		data.Components = new([]discord.Component)
		*data.Components = make([]discord.Component, len(*b.components))

		for i, componentBuilder := range *b.components {
			c, err := componentBuilder.Build(b.localizer)
			if err != nil {
				return nil, err
			}

			(*data.Components)[i] = c
		}
	}

	msg, err := b.s.EditMessageComplex(channelID, messageID, data)
	return msg, errorutil.WithStack(err)
}

// AwaitComponents calls AwaitComponentsContext with a context with the given
// timeout.
func (b *Builder) AwaitComponents(timeout time.Duration, userID discord.UserID, disable bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return b.AwaitComponentsContext(ctx, userID, disable)
}

// SendAndAwait sends the message in the given channel, and waits until the
// user with the given id interacts with the message, or the timeout elapses.
// Afterwards, it disables all components.
func (b *Builder) SendAndAwait(
	channelID discord.ChannelID, userID discord.UserID, timeout time.Duration,
) (*discord.Message, error) {
	msg, err := b.Send(channelID)
	if err != nil {
		return nil, err
	}

	return msg, b.AwaitComponents(timeout, userID, true)
}

// AwaitComponentsContext waits until the first awaited component is interacted
// with, or the passed context.Context is done, whichever is first.
//
// Subsequent calls to AwaitComponentsContext will await further interactions.
//
// If disable is set to true, all components will be disabled after the
// function returns, making subsequent calls impossible.
// When calling AwaitComponentsContext for the last time, disable should always
// be true.
//
// Interactions that happened before AwaitComponentsContext was called, or
// those that happened between calls will not be evaluated.
//nolint:gocognit
func (b *Builder) AwaitComponentsContext(ctx context.Context, userID discord.UserID, disable bool) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eventChan := make(chan *event.InteractionCreate)

	rm := b.s.AddHandler(func(_ *state.State, e *event.InteractionCreate) {
		if e.Data == nil || e.Message.ID != b.messageID ||
			(e.User != nil && e.User.ID != userID) || (e.Member != nil && e.Member.User.ID != userID) {
			return
		}

		select {
		case <-ctx.Done():
		case eventChan <- e:
		}
	})

	go func() {
		var ok bool

		for e := range eventChan {
			for i, c := range *b.components {
				ok, err = c.handle(e.Data)
				if err != nil { // something went wrong, this takes precedence
					cancel()
					return
				}

				// check if the component matched the event
				if !ok {
					continue
				}

				if disable {
					b.disableAllComponents()
					_, err = b.Edit(b.channelID, b.messageID)
					if err != nil {
						cancel()
						return
					}
				}

				err = b.s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
					Type: api.DeferredMessageUpdate,
				})
				if err != nil {
					cancel()
					return
				}

				// the component matched the event and this is one of the
				// components we should wait for
				if _, await := b.awaitIndexes[i]; await {
					cancel()
					return
				}
			}
		}
	}()

	<-ctx.Done()
	rm()
	return err
}

// AwaitAllComponents calls AwaitAllComponentsContext with a context with the
// given timeout.
func (b *Builder) AwaitAllComponents(s *state.State, userID discord.UserID, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return b.AwaitAllComponentsContext(ctx, userID, s)
}

// AwaitAllComponentsContext waits until all awaited components were interacted
// with once, or the context is done, whichever happens first.
// Before returning, it disables all components.
//
// Subsequent calls to AwaitAllComponentsContext will return immediately with
// a nil error.
// However, calling AwaitAllComponentsContext after AwaitComponentsContext was
// called, is allowed, and will behave as described above.
//
// Interactions that happened before AwaitAllComponentsContext was called, or
// those that happened between calls will not be evaluated.
//nolint:gocognit
func (b *Builder) AwaitAllComponentsContext(ctx context.Context, userID discord.UserID, s *state.State) (err error) {
	if len(b.awaitIndexes) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eventChan := make(chan *event.InteractionCreate)

	rm := s.AddHandler(func(_ *state.State, e *event.InteractionCreate) {
		if e.Data == nil || e.Message.ID != b.messageID ||
			(e.User != nil && e.User.ID != userID) || (e.Member != nil && e.Member.User.ID != userID) {
			return
		}

		select {
		case <-ctx.Done():
		case eventChan <- e:
		}
	})

	go func() {
		var ok bool

		for e := range eventChan {
			for i, c := range *b.components {
				ok, err = c.handle(e.Data)
				if err != nil { // something went wrong, this takes precedence
					cancel()
					return
				}

				// check that the component matched the event
				if !ok {
					continue
				}

				err = s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
					Type: api.DeferredMessageUpdate,
				})
				if err != nil {
					cancel()
					return
				}

				// the component matched the event and this is one of the
				// components we should wait for
				if _, await := b.awaitIndexes[i]; await {
					// remove the component from the awaited components
					delete(b.awaitIndexes, i)

					// if all components we waited for were interacted with,
					// return
					if len(b.awaitIndexes) == 0 {
						b.disableAllComponents()
						_, err = b.Edit(b.channelID, b.messageID)
						if err != nil {
							cancel()
							return
						}
					}

					err = s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
						Type: api.DeferredMessageUpdate,
					})
					if err != nil {
						cancel()
						return
					}

					if len(b.awaitIndexes) == 0 {
						cancel()
						return
					}

					// otherwise, break and wait for the next event
					break
				}
			}
		}
	}()

	<-ctx.Done()
	rm()
	return err
}
