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

	"github.com/mavolin/adam/internal/embedbuilder"
	"github.com/mavolin/adam/internal/errorutil"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type EmbedBuilder = embedbuilder.Builder

func NewEmbed() *EmbedBuilder {
	return embedbuilder.New()
}

// Builder is a message builder.
// After creating or editing a message, it must not be used again.
type Builder struct {
	// ========= message =========

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

	// ========= internals =========

	state *state.State
	ctx   *plugin.Context

	userID discord.UserID

	channelID discord.ChannelID // NullChannelID if replier
	messageID discord.MessageID
	dm        bool // only used when using replier
}

// New creates a new *Builder that awaits a component interaction from the
// invoking user.
func New(s *state.State, ctx *plugin.Context) *Builder {
	return &Builder{
		awaitIndexes: make(map[int]struct{}),
		state:        s,
		ctx:          ctx,
		userID:       ctx.Author.ID,
	}
}

// =============================================================================
// Message Related Methods
// =====================================================================================

// WithContent sets the content of the message to the given content.
// It may be no longer than 2000 characters.
//
// Actions: send and edit
func (b *Builder) WithContent(content string) *Builder {
	return b.WithContentl(i18n.NewStaticConfig(content))
}

// WithContentlt sets the content of the message to the given content.
// It may be no longer than 2000 characters.
//
// Actions: send and edit
func (b *Builder) WithContentlt(content i18n.Term) *Builder {
	return b.WithContentl(content.AsConfig())
}

// WithContentl sets the content of the message to the given content.
// It may be no longer than 2000 characters.
//
// Actions: send and edit
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
// Actions: send and edit
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
// Actions: send and edit
func (b *Builder) WithComponent(component TopLevelComponentBuilder) *Builder {
	if b.components == nil {
		b.components = new([]TopLevelComponentBuilder)
	}

	*b.components = append(*b.components, component)
	return b
}

// WithAwaitedComponent adds the passed TopLevelComponentBuilder to the
// message, and waits for an interaction for that component when
// AwaitComponents is called.
//
// Actions: send and edit
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
// Actions: send
func (b *Builder) AsTTS() *Builder {
	b.tts = true
	return b
}

// WithAllowedMentionTypes adds the passed api.AllowedMentionTypes to the
// allowed mentions of the message.
//
// Actions: send and edit
func (b *Builder) WithAllowedMentionTypes(allowed ...api.AllowedMentionType) *Builder {
	if b.allowedMentions == nil {
		b.allowedMentions = new(api.AllowedMentions)
	}

	b.allowedMentions.Parse = append(b.allowedMentions.Parse, allowed...)
	return b
}

// WithRoleMentions adds the passed role ids to the allowed mentions.
//
// Actions: send and edit
func (b *Builder) WithRoleMentions(roleIDs ...discord.RoleID) *Builder {
	if b.allowedMentions == nil {
		b.allowedMentions = new(api.AllowedMentions)
	}

	b.allowedMentions.Roles = append(b.allowedMentions.Roles, roleIDs...)
	return b
}

// WithUserMentions adds the passed user ids to the allowed mentions.
//
// Actions: send and edit
func (b *Builder) WithUserMentions(userIDs ...discord.UserID) *Builder {
	if b.allowedMentions == nil {
		b.allowedMentions = new(api.AllowedMentions)
	}

	b.allowedMentions.Users = append(b.allowedMentions.Users, userIDs...)
	return b
}

// MentionRepliedUser allows the replied user to be mentioned.
//
// Actions: send and edit
func (b *Builder) MentionRepliedUser() *Builder {
	if b.allowedMentions == nil {
		b.allowedMentions = new(api.AllowedMentions)
	}

	b.allowedMentions.RepliedUser = option.True
	return b
}

// WithReference references the message with the passed id.
//
// Actions: send
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
// Actions: send and edit
func (b *Builder) WithFile(name string, reader io.Reader) *Builder {
	b.files = append(b.files, sendpart.File{Name: name, Reader: reader})
	return b
}

// =============================================================================
// Component Await Settings
// =====================================================================================

// AwaitingUser sets the id of the user that is awaited to the passed id.
func (b *Builder) AwaitingUser(userID discord.UserID) *Builder {
	b.userID = userID
	return b
}

// =============================================================================
// Message Sending
// =====================================================================================

func (b *Builder) sendMessageData() (data api.SendMessageData, err error) {
	data = api.SendMessageData{
		TTS:             b.tts,
		Files:           b.files,
		AllowedMentions: b.allowedMentions,
		Reference:       b.reference,
	}

	if b.content != nil {
		data.Content, err = b.ctx.Localize(b.content)
		if err != nil {
			return api.SendMessageData{}, nil
		}
	}

	if b.embeds != nil && len(*b.embeds) > 0 {
		data.Embeds = make([]discord.Embed, len(*b.embeds))

		for i, embedBuilder := range *b.embeds {
			embed, err := embedBuilder.Build(b.ctx.Localizer)
			if err != nil {
				return api.SendMessageData{}, nil
			}

			data.Embeds[i] = embed
		}
	}

	if b.components != nil && len(*b.components) > 0 {
		data.Components = make([]discord.Component, len(*b.components))

		for i, componentBuilder := range *b.components {
			c, err := componentBuilder.Build(b.ctx.Localizer)
			if err != nil {
				return api.SendMessageData{}, nil
			}

			data.Components[i] = c
		}
	}

	return data, err
}

// Reply sends the message using the plugin.Context's plugin.Replier.
func (b *Builder) Reply() (*discord.Message, error) {
	data, err := b.sendMessageData()
	if err != nil {
		return nil, err
	}

	msg, err := b.ctx.ReplyMessage(data)
	if err == nil {
		b.channelID = discord.NullChannelID
		b.messageID = msg.ID
	}

	return msg, err
}

// ReplyAndAwait sends the message using the plugin.Context's plugin.Replier
// and then waits for a component interaction.
// Upon returning, it disables all components.
func (b *Builder) ReplyAndAwait(timeout time.Duration) (*discord.Message, error) {
	msg, err := b.Reply()
	if err != nil {
		return nil, err
	}

	return msg, b.AwaitComponents(timeout, true)
}

// ReplyDM sends the message in a DM using the plugin.Context's plugin.Replier.
func (b *Builder) ReplyDM() (*discord.Message, error) {
	data, err := b.sendMessageData()
	if err != nil {
		return nil, err
	}

	msg, err := b.ctx.ReplyMessageDM(data)
	if err == nil {
		b.channelID = discord.NullChannelID
		b.messageID = msg.ID
		b.dm = true
	}

	return msg, err
}

// Send creates a new message in the passed channel.
func (b *Builder) Send(channelID discord.ChannelID) (msg *discord.Message, err error) {
	data, err := b.sendMessageData()
	if err != nil {
		return nil, err
	}

	msg, err = b.state.SendMessageComplex(channelID, data)
	if err == nil {
		b.channelID = msg.ChannelID
		b.messageID = msg.ID
	}

	return msg, errorutil.WithStack(err)
}

// =============================================================================
// Message Editing
// =====================================================================================

func (b *Builder) editMessageData() (api.EditMessageData, error) {
	data := api.EditMessageData{
		AllowedMentions: b.allowedMentions,
		Attachments:     b.attachments,
		Flags:           b.flags,
		Files:           b.files,
	}

	if b.content != nil {
		content, err := b.ctx.Localize(b.content)
		if err != nil {
			return api.EditMessageData{}, err
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
			embed, err := embedBuilder.Build(b.ctx.Localizer)
			if err != nil {
				return api.EditMessageData{}, err
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
			c, err := componentBuilder.Build(b.ctx.Localizer)
			if err != nil {
				return api.EditMessageData{}, err
			}

			(*data.Components)[i] = c
		}
	}
	return data, nil
}

// EditReply edits the message with the given id using the plugin.Context's
// plugin.Replier.
func (b *Builder) EditReply(messageID discord.MessageID) (*discord.Message, error) {
	b.messageID = messageID
	b.channelID = discord.NullChannelID

	data, err := b.editMessageData()
	if err != nil {
		return nil, err
	}

	return b.ctx.EditMessage(messageID, data)
}

// EditReplyDM edits the direct message with the given id using the
// plugin.Context's plugin.Replier.
func (b *Builder) EditReplyDM(messageID discord.MessageID) (*discord.Message, error) {
	b.messageID = messageID
	b.channelID = discord.NullChannelID
	b.dm = true

	data, err := b.editMessageData()
	if err != nil {
		return nil, err
	}

	return b.ctx.EditMessageDM(messageID, data)
}

// Edit edits the message with the passed channel and message id.
func (b *Builder) Edit(channelID discord.ChannelID, messageID discord.MessageID) (*discord.Message, error) {
	b.channelID = channelID
	b.messageID = messageID

	data, err := b.editMessageData()
	if err != nil {
		return nil, err
	}

	msg, err := b.state.EditMessageComplex(channelID, messageID, data)
	return msg, errorutil.WithStack(err)
}

// =============================================================================
// Await Components
// =====================================================================================

// AwaitComponents calls AwaitComponentsContext with a context with the given
// timeout.
func (b *Builder) AwaitComponents(timeout time.Duration, disable bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return b.AwaitComponentsContext(ctx, disable)
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
func (b *Builder) AwaitComponentsContext(ctx context.Context, disable bool) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eventChan := make(chan *event.InteractionCreate)
	rm := b.state.AddHandler(b.interactionCreateHandler(ctx, eventChan))
	defer rm()

	done := make(chan error, 1)

	go func() {
		var ok bool

		for {
			select {
			case <-ctx.Done():
				done <- &TimeoutError{UserID: b.userID, Cause: ctx.Err()}
				return
			case e := <-eventChan:
				for i, c := range *b.components {
					ok, err = c.handle(e.Data)
					if err != nil { // something went wrong, this takes precedence
						done <- err
						return
					}

					// check if the component matched the event
					if !ok {
						continue
					}

					if disable {
						if err := b.disableAllComponents(); err != nil {
							done <- err
							return
						}
					}

					err = b.state.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
						Type: api.DeferredMessageUpdate,
					})
					if err != nil {
						done <- err
						return
					}

					// the component matched the event and this is one of the
					// components we should wait for
					if _, await := b.awaitIndexes[i]; await {
						done <- nil
						return
					}
				}
			}
		}
	}()

	return <-done
}

// =============================================================================
// Utils
// =====================================================================================

// disableAllComponents disables all components and edits the message.
func (b *Builder) disableAllComponents() error {
	if b.components != nil {
		for _, c := range *b.components {
			c.disable()
		}
	}

	if b.channelID == discord.NullChannelID {
		if b.dm {
			_, err := b.EditReplyDM(b.messageID)
			return err
		}

		_, err := b.EditReply(b.messageID)
		return err
	}

	_, err := b.Edit(b.channelID, b.messageID)
	return err
}

func (b *Builder) interactionCreateHandler(
	ctx context.Context, eventChan chan *event.InteractionCreate,
) func(*state.State, *event.InteractionCreate) {
	return func(_ *state.State, e *event.InteractionCreate) {
		if e.Data == nil || e.Message.ID != b.messageID ||
			(e.User != nil && e.User.ID != b.userID) || (e.Member != nil && e.Member.User.ID != b.userID) {
			return
		}

		select {
		case <-ctx.Done():
		case eventChan <- e:
		}
	}
}
