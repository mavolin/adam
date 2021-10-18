// Package msgbuilder provides builders to create messages, embeds and
// components.
// Its main purpose is to abstract away the handling of message components, as
// well as to ease working with localized messages.
package msgbuilder

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/mavolin/disstate/v4/pkg/event"
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// typingInterval is the interval in which the client of the user sends the
// typing event, if the user is continuously typing.
//
// It has been observed that the first follow-up event is received after about
// 9.5 seconds, all successive events are received in an interval of
// approximately 8.25 seconds. Additionally, we add a 1.5 second margin for
// network delays.
const typingInterval = 11 * time.Second

type responseMiddlewaresKey struct{}

// AddResponseMiddlewares adds the passed middlewares to the list of the
// middlewares that are automatically added when awaiting a response from the
// user.
func AddResponseMiddlewares(ctx *plugin.Context, middlewares ...interface{}) {
	if currentMiddlewares, ok := ctx.Get(responseMiddlewaresKey{}).([]interface{}); ok {
		ctx.Set(responseMiddlewaresKey{}, append(currentMiddlewares, middlewares...))
	} else {
		ctx.Set(responseMiddlewaresKey{}, middlewares)
	}
}

// Builder is a message builder.
// After creating or editing a message, it must not be used again.
type Builder struct {
	// ========= Message =========

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

	// ========= Response Await =========

	responseMsgPtr         *discord.Message
	initialResponseTimeout time.Duration
	typingResponseTimeout  time.Duration

	responseMiddlewares []interface{}

	// ========= Internals =========

	state     *state.State
	pluginCtx *plugin.Context

	userID discord.UserID

	channelID discord.ChannelID // NullChannelID if replier
	messageID discord.MessageID
	dm        bool // only used when using replier

}

// New creates a new *Builder that awaits a component interaction from the
// invoking user.
//
// It adds the middlewares stored in the context under the
// ResponseMiddlewareKey to the builder's response middlewares, to be used when
// awaiting a response.
func New(s *state.State, ctx *plugin.Context) *Builder {
	b := &Builder{
		awaitIndexes: make(map[int]struct{}),
		state:        s,
		pluginCtx:    ctx,
		userID:       ctx.Author.ID,
	}

	middlewares := ctx.Get(responseMiddlewaresKey{})
	if middlewares, ok := middlewares.([]interface{}); ok && middlewares != nil {
		b.WithResponseMiddlewares(middlewares...)
	}

	return b
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
// Await is called.
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

// MentionRepliedUser mentions the user that is being replied.
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

// WithAwaitedResponse will cause the message to await a response from the
// user when Await or AwaitContext is called.
//
// The builder will wait for a response until the initial timeout expires and
// the user is not typing, or until the user stops typing and the typing
// timeout is reached.
//
// If typingTimeout is set to 0, typing will not be monitored.
//
// Note that you need the typing intent to monitor typing.
// If typingTimeout is > 0 and the current shard has no typing intent for
// guilds or direct messages, wherever the response is being awaited, Await and
// AwaitContext will return with an error.
//
// If one of the timeouts is reached, a *TimeoutError will be returned.
//
// The typing timeout will start after the user first starts typing.
// Because Discord sends the typing event in an interval of about 10 seconds,
// Await and AwaitContext might, in the worst case, only notice that the user
// ceased typing 10 seconds later.
// In the best case, it will be noticed almost immediately.
//
// The timeout given to Await, or the cancellation of the context given to
// AwaitContext serves as a maximum timeout.
// If it is reached, the await functions will return no matter if the user is
// still typing.
// It will also serve as the only timeout, if the typing timeout was disabled.
//
// Besides that, a reply can also be canceled through a middleware.
func (b *Builder) WithAwaitedResponse(
	responseVar *discord.Message, initialTimeout, typingTimeout time.Duration,
) *Builder {
	b.responseMsgPtr = responseVar
	b.initialResponseTimeout = initialTimeout
	b.typingResponseTimeout = typingTimeout

	return b
}

// WithResponseMiddlewares adds the passed middlewares to the builder to be
// executed before every message create event processed.
//
// Any errors returned by the middlewares will be returned by the await
// function that awaited a response.
//
// All middlewares of invalid type will be discarded.
//
// The following types are permitted:
// 	• func(*state.State, interface{})
//	• func(*state.State, interface{}) error
//	• func(*state.State, *event.Base)
//	• func(*state.State, *event.Base) error
//	• func(*state.State, *state.MessageCreateEvent)
//	• func(*state.State, *state.MessageCreateEvent) error
func (b *Builder) WithResponseMiddlewares(middlewares ...interface{}) *Builder {
	if len(b.responseMiddlewares) == 0 {
		b.responseMiddlewares = make([]interface{}, 0, len(middlewares))
	}

	for _, m := range middlewares {
		switch m.(type) { // check that the middleware is of a valid type
		case func(*state.State, interface{}):
		case func(*state.State, interface{}) error:
		case func(*state.State, *event.Base):
		case func(*state.State, *event.Base) error:
		case func(*state.State, *event.MessageCreate):
		case func(*state.State, *event.MessageCreate) error:
		default:
			continue
		}

		b.responseMiddlewares = append(b.responseMiddlewares, m)
	}

	return b
}

// AwaitingUser sets the id of the user, that is expected to interact with the
// components and to respond to the message, to the passed id.
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
		data.Content, err = b.pluginCtx.Localize(b.content)
		if err != nil {
			return api.SendMessageData{}, nil
		}
	}

	if b.embeds != nil && len(*b.embeds) > 0 {
		data.Embeds = make([]discord.Embed, len(*b.embeds))

		for i, embedBuilder := range *b.embeds {
			embed, err := embedBuilder.Build(b.pluginCtx.Localizer)
			if err != nil {
				return api.SendMessageData{}, nil
			}

			data.Embeds[i] = embed
		}
	}

	if b.components != nil && len(*b.components) > 0 {
		data.Components = make([]discord.Component, len(*b.components))

		for i, componentBuilder := range *b.components {
			c, err := componentBuilder.Build(b.pluginCtx.Localizer)
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

	msg, err := b.pluginCtx.ReplyMessage(data)
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

	return msg, b.Await(timeout, true)
}

// ReplyDM sends the message in a DM using the plugin.Context's plugin.Replier.
func (b *Builder) ReplyDM() (*discord.Message, error) {
	data, err := b.sendMessageData()
	if err != nil {
		return nil, err
	}

	msg, err := b.pluginCtx.ReplyMessageDM(data)
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

	return msg, errors.WithStack(err)
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
		content, err := b.pluginCtx.Localize(b.content)
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
			embed, err := embedBuilder.Build(b.pluginCtx.Localizer)
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
			c, err := componentBuilder.Build(b.pluginCtx.Localizer)
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

	return b.pluginCtx.EditMessage(messageID, data)
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

	return b.pluginCtx.EditMessageDM(messageID, data)
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
	return msg, errors.WithStack(err)
}

// =============================================================================
// Await
// =====================================================================================

// Await calls AwaitContext with a context with the given timeout.
func (b *Builder) Await(timeout time.Duration, disable bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return b.AwaitContext(ctx, disable)
}

// AwaitContext waits until the first awaited component is interacted
// with, or, if awaiting a response, a response is sent.
//
// Subsequent calls to AwaitContext will await further interactions/responses.
//
// To use await, the message that this builder builds must have been sent.
//
// If disable is set to true, all components will be disabled after the
// function returns, making subsequent calls impossible.
// You should do this when calling AwaitContext for the last time, to make
// clear that no further interactions with the components in the message are
// possible.
//
// Errors that occur when disabling components will be handled silently and
// will not be returned, as disabling the components is just cosmetic and has
// no influence over the successful execution of AwaitContext.
//
// If disable is set to true, and AwaitContext returns with an error, all
// components will still be disabled.
//
// Interactions that happened before AwaitContext was called, or
// those that happened between calls will not be evaluated.
func (b *Builder) AwaitContext(ctx context.Context, disable bool) (err error) {
	if disable {
		defer func() {
			if err := b.DisableComponents(); err != nil {
				b.pluginCtx.HandleErrorSilently(err)
			}
		}()
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	doneChan := make(chan error)

	if b.components != nil && len(*b.components) > 0 {
		interactionCleanup := b.handleInteractions(ctx, doneChan)
		defer interactionCleanup()
	}

	if b.responseMsgPtr != nil {
		perms, err := b.pluginCtx.SelfPermissions()
		if err != nil {
			return err
		}

		// make sure we have permission to send messages and create reactions
		// if time extensions are enabled
		if !perms.Has(discord.PermissionSendMessages) {
			return plugin.NewBotPermissionsError(discord.PermissionSendMessages)
		}

		msgCleanup := b.handleMessages(ctx, doneChan)
		defer msgCleanup()

		if b.typingResponseTimeout != 0 {
			typingIntent := gateway.IntentGuildMessageTyping
			if b.pluginCtx.GuildID == 0 {
				typingIntent = gateway.IntentDirectMessageTyping
			}

			if !b.state.GatewayFromGuildID(b.pluginCtx.GuildID).HasIntents(typingIntent) {
				return errors.NewWithStackf("msgbuilder: need typing intent to use typing timeout when awaiting response")
			}

			timeoutCleanup := b.watchTyping(ctx, b.initialResponseTimeout, b.typingResponseTimeout, doneChan)
			defer timeoutCleanup()
		}
	}

	select {
	case <-ctx.Done():
		return &TimeoutError{UserID: b.userID, Cause: ctx.Err()}
	case err := <-doneChan:
		return err
	}
}

func (b *Builder) handleInteractions(ctx context.Context, doneChan chan<- error) func() {
	ctx, cancel := context.WithCancel(ctx)

	var mut sync.Mutex

	return b.state.AddHandler(func(_ *state.State, e *event.InteractionCreate) {
		if e.Data == nil || e.Message.ID != b.messageID ||
			(e.User != nil && e.User.ID != b.userID) || (e.Member != nil && e.Member.User.ID != b.userID) {
			return
		}

		mut.Lock()
		defer mut.Unlock()

		select {
		case <-ctx.Done():
			return
		default:
		}

		for i, c := range *b.components {
			ok, err := c.handle(e.Data)
			if err != nil { // something went wrong, this takes precedence
				sendDone(ctx, doneChan, err)

				// Make sure only a single event is handled by preventing that
				// others pass the above select.
				// We need this, because the mutex lock might queue up events
				// that might get processed in between sending into done and
				// the parent context being canceled.
				// Cancelling by ourselves closes this gap.
				cancel()

				return
			}

			// check if the component matched the event
			if !ok {
				continue
			}

			err = b.state.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
				Type: api.DeferredMessageUpdate,
			})
			if err != nil {
				sendDone(ctx, doneChan, err)
				cancel()
				return
			}

			// the component matched the event and this is one of the
			// components we should wait for
			if _, await := b.awaitIndexes[i]; await {
				sendDone(ctx, doneChan, err)
				cancel()
				return
			}
		}
	})
}

func (b *Builder) handleMessages(ctx context.Context, doneChan chan<- error) func() {
	var mut sync.Mutex

	ctx, cancel := context.WithCancel(ctx)

	channelID := b.channelID
	if channelID == discord.NullChannelID {
		channelID = b.pluginCtx.ChannelID
	}

	return b.state.AddHandler(func(s *state.State, e *event.MessageCreate) {
		// not the message we are waiting for
		if e.ChannelID != channelID || e.Author.ID != b.userID {
			return
		}

		mut.Lock()
		defer mut.Unlock()

		select {
		case <-ctx.Done():
			return
		default:
		}

		if err := b.invokeResponseMiddlewares(s, e); err != nil {
			sendDone(ctx, doneChan, err)
			cancel()
			return
		}

		*b.responseMsgPtr = e.Message
		sendDone(ctx, doneChan, nil)
		cancel()
	})
}

func (b *Builder) watchTyping(
	ctx context.Context, initialTimeout, typingTimeout time.Duration, doneChan chan<- error,
) func() {
	if typingTimeout <= 0 {
		return func() {}
	}

	t := time.NewTimer(initialTimeout)

	rm := b.state.AddHandler(func(s *state.State, e *event.TypingStart) {
		if e.ChannelID != b.channelID || e.UserID != b.userID {
			return
		}

		// this should always return true, except if timer expired after
		// the typing event was received and this handler was called
		if t.Stop() {
			t.Reset(typingTimeout + typingInterval)
		}
	})

	go func() {
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
				sendDone(ctx, doneChan, &TimeoutError{UserID: b.userID})
				return
			}
		}
	}()

	return rm
}

// =============================================================================
// Utils
// =====================================================================================

// DisableComponents disables all components, by editing the sent message.
func (b *Builder) DisableComponents() error {
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

func (b *Builder) invokeResponseMiddlewares(s *state.State, e *event.MessageCreate) error {
	for _, m := range b.responseMiddlewares {
		switch m := m.(type) {
		case func(*state.State, interface{}):
			m(s, e)
		case func(*state.State, interface{}) error:
			if err := m(s, e); err != nil {
				return err
			}
		case func(*state.State, *event.Base):
			m(s, e.Base)
		case func(*state.State, *event.Base) error:
			if err := m(s, e.Base); err != nil {
				return err
			}
		case func(*state.State, *event.MessageCreate):
			m(s, e)
		case func(*state.State, *event.MessageCreate) error:
			if err := m(s, e); err != nil {
				return err
			}
		}
	}

	return nil
}

func sendDone(ctx context.Context, doneChan chan<- error, err error) {
	select {
	case <-ctx.Done():
		return
	case doneChan <- err:
	}
}
