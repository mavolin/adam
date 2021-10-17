package msgbuilder

import (
	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type (
	// EmbedBuilder is used to build embeds.
	EmbedBuilder struct {
		title       *i18n.Config
		description *i18n.Config

		url discord.URL

		timestamp discord.Timestamp
		color     discord.Color

		footer       *footer
		imageURL     discord.URL
		thumbnailURL discord.URL
		author       *author
		fields       []field
	}

	footer struct {
		text *i18n.Config
		icon discord.URL
	}

	author struct {
		name *i18n.Config
		icon discord.URL
		url  discord.URL
	}

	field struct {
		name    *i18n.Config
		value   *i18n.Config
		inlined bool
	}
)

// NewEmbed creates a new EmbedBuilder.
func NewEmbed() *EmbedBuilder {
	return new(EmbedBuilder)
}

// WithTitle sets the title (max. 256 characters) to the passed title.
func (b *EmbedBuilder) WithTitle(title string) *EmbedBuilder {
	return b.WithTitlel(i18n.NewStaticConfig(title))
}

// WithTitlel sets the title (max. 256 characters) to the passed title.
func (b *EmbedBuilder) WithTitlel(title *i18n.Config) *EmbedBuilder {
	b.title = title
	return b
}

// WithTitleURL assigns a discord.URL to the title.
func (b *EmbedBuilder) WithTitleURL(url discord.URL) *EmbedBuilder {
	b.url = url
	return b
}

// WithDescription sets the description (max. 2048 characters) to the passed
// description.
func (b *EmbedBuilder) WithDescription(description string) *EmbedBuilder {
	return b.WithDescriptionl(i18n.NewStaticConfig(description))
}

// WithDescriptionl sets the description (max. 2048 characters) to the passed
// description.
func (b *EmbedBuilder) WithDescriptionl(description *i18n.Config) *EmbedBuilder {
	b.description = description
	return b
}

// WithTimestamp sets the timestamp to the passed discord.Timestamp.
func (b *EmbedBuilder) WithTimestamp(timestamp discord.Timestamp) *EmbedBuilder {
	b.timestamp = timestamp
	return b
}

// WithTimestampNow sets the timestamp to a timestamp of the current time.
func (b *EmbedBuilder) WithTimestampNow() *EmbedBuilder {
	return b.WithTimestamp(discord.NowTimestamp())
}

// WithColor sets the color to the passed discord.Color.
func (b *EmbedBuilder) WithColor(color discord.Color) *EmbedBuilder {
	b.color = color
	return b
}

// WithFooter sets the text of the footer (max. 2048 characters) to the passed
// text.
func (b *EmbedBuilder) WithFooter(text string) *EmbedBuilder {
	return b.WithFooterl(i18n.NewStaticConfig(text))
}

// WithFooterl sets the text of the footer (max. 2048 characters) to the passed
// text.
func (b *EmbedBuilder) WithFooterl(text *i18n.Config) *EmbedBuilder {
	if b.footer == nil {
		b.footer = &footer{text: text}
	} else {
		b.footer.text = text
	}

	return b
}

// WithFooterIcon sets the icon of the footer to the passed icon url.
func (b *EmbedBuilder) WithFooterIcon(icon discord.URL) *EmbedBuilder {
	if b.footer == nil {
		b.footer = &footer{icon: icon}
	} else {
		b.footer.icon = icon
	}

	return b
}

// WithImage sets the image to the passed image url.
func (b *EmbedBuilder) WithImage(image discord.URL) *EmbedBuilder {
	b.imageURL = image
	return b
}

// WithThumbnail sets the thumbnail to the passed thumbnail url.
func (b *EmbedBuilder) WithThumbnail(thumbnail discord.URL) *EmbedBuilder {
	b.thumbnailURL = thumbnail

	return b
}

// WithAuthor sets the author's name (max. 256 characters) to the passed
// name.
func (b *EmbedBuilder) WithAuthor(name string) *EmbedBuilder {
	return b.WithAuthorl(i18n.NewStaticConfig(name))
}

// WithAuthorl sets the author's name (max. 256 characters) to the passed
// name.
func (b *EmbedBuilder) WithAuthorl(name *i18n.Config) *EmbedBuilder {
	if b.author == nil {
		b.author = &author{name: name}
	} else {
		b.author.name = name
	}

	return b
}

// WithAuthorURL assigns the author the passed discord.URL.
func (b *EmbedBuilder) WithAuthorURL(url discord.URL) *EmbedBuilder {
	if b.author == nil {
		b.author = &author{url: url}
	} else {
		b.author.url = url
	}

	return b
}

// WithAuthorIcon sets the icon of the author to the passed icon url.
func (b *EmbedBuilder) WithAuthorIcon(icon discord.URL) *EmbedBuilder {
	if b.author == nil {
		b.author = &author{icon: icon}
	} else {
		b.author.icon = icon
	}

	return b
}

// WithField adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *EmbedBuilder) WithField(name, value string) *EmbedBuilder {
	return b.WithFieldl(i18n.NewStaticConfig(name), i18n.NewStaticConfig(value))
}

// WithFieldl adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *EmbedBuilder) WithFieldl(name, value *i18n.Config) *EmbedBuilder {
	b.fields = append(b.fields, field{
		inlined: false,
		name:    name,
		value:   value,
	})

	return b
}

// WithInlinedField adds an inlined field (name: max. 256 characters, value:
// max 1024 characters) to the embed.
func (b *EmbedBuilder) WithInlinedField(name, value string) *EmbedBuilder {
	return b.WithInlinedFieldl(i18n.NewStaticConfig(name), i18n.NewStaticConfig(value))
}

// WithInlinedFieldl adds an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
func (b *EmbedBuilder) WithInlinedFieldl(name, value *i18n.Config) *EmbedBuilder {
	b.fields = append(b.fields, field{
		inlined: true,
		name:    name,
		value:   value,
	})

	return b
}

// Clone creates a copy of the EmbedBuilder.
func (b *EmbedBuilder) Clone() *EmbedBuilder {
	cp := new(EmbedBuilder)

	*cp = *b

	cp.title = b.title
	cp.description = b.description

	if b.footer != nil {
		cp.footer = &footer{
			text: b.footer.text,
			icon: b.footer.icon,
		}
	}

	if b.author != nil {
		cp.author = &author{
			name: b.author.name,
			icon: b.author.icon,
			url:  b.author.url,
		}
	}

	if b.fields != nil {
		cp.fields = make([]field, len(b.fields))
		copy(cp.fields, b.fields)
	}

	return cp
}

// Build builds the discord.Embed.
func (b *EmbedBuilder) Build(l *i18n.Localizer) (e discord.Embed, err error) {
	if b.title != nil {
		e.Title, err = l.Localize(b.title)
		if err != nil {
			return discord.Embed{}, err
		}
	}

	if b.description != nil {
		e.Description, err = l.Localize(b.description)
		if err != nil {
			return discord.Embed{}, err
		}
	}

	e.URL = b.url
	e.Timestamp = b.timestamp
	e.Color = b.color

	if b.footer != nil {
		e.Footer = &discord.EmbedFooter{Icon: b.footer.icon}

		if b.footer.text != nil {
			e.Footer.Text, err = l.Localize(b.footer.text)
			if err != nil {
				return discord.Embed{}, err
			}
		}
	}

	if b.imageURL != "" {
		e.Image = &discord.EmbedImage{URL: b.imageURL}
	}

	if b.thumbnailURL != "" {
		e.Thumbnail = &discord.EmbedThumbnail{URL: b.thumbnailURL}
	}

	if b.author != nil {
		e.Author = &discord.EmbedAuthor{
			URL:  b.author.url,
			Icon: b.author.icon,
		}

		if b.author.name != nil {
			e.Author.Name, err = l.Localize(b.author.name)
			if err != nil {
				return discord.Embed{}, err
			}
		}
	}

	if len(b.fields) > 0 {
		e.Fields = make([]discord.EmbedField, len(b.fields))
	}

	for i, f := range b.fields {
		var name string

		if f.name != nil {
			name, err = l.Localize(f.name)
			if err != nil {
				return discord.Embed{}, err
			}
		}

		var value string

		if f.value != nil {
			value, err = l.Localize(f.value)
			if err != nil {
				return discord.Embed{}, err
			}
		}

		e.Fields[i] = discord.EmbedField{
			Name:   name,
			Value:  value,
			Inline: f.inlined,
		}
	}

	return e, err
}

// =============================================================================
// Utils
// =====================================================================================

// BuildEmbeds builds the passed *EmbedBuilders.
func BuildEmbeds(l *i18n.Localizer, builders ...*EmbedBuilder) ([]discord.Embed, error) {
	embeds := make([]discord.Embed, len(builders))

	for i, builder := range builders {
		embed, err := builder.Build(l)
		if err != nil {
			return nil, err
		}

		embeds[i] = embed
	}

	return embeds, nil
}

// ReplyEmbedBuilders builds the discord.Embeds from the passed
// *EmbedBuilders and sends them in the channel the command was sent in.
func ReplyEmbedBuilders(ctx *plugin.Context, builders ...*EmbedBuilder) (*discord.Message, error) {
	embeds, err := BuildEmbeds(ctx.Localizer, builders...)
	if err != nil {
		return nil, err
	}

	return ctx.ReplyEmbeds(embeds...)
}

// ReplyEmbedBuildersDM builds the discord.Embeds from the passed
// *EmbedBuilders and sends them in a direct message to the invoking user.
func ReplyEmbedBuildersDM(ctx *plugin.Context, builders ...*EmbedBuilder) (*discord.Message, error) {
	embeds, err := BuildEmbeds(ctx.Localizer, builders...)
	if err != nil {
		return nil, err
	}

	return ctx.ReplyEmbeds(embeds...)
}

// EditEmbedBuilders builds the discord.Embeds from the passed
// *EmbedBuilders, and replaces the embeds of the message with the passed id in
// the invoking channel with them.
func EditEmbedBuilders(
	ctx *plugin.Context, messageID discord.MessageID,
	builders ...*EmbedBuilder,
) (*discord.Message, error) {
	embeds, err := BuildEmbeds(ctx.Localizer, builders...)
	if err != nil {
		return nil, err
	}

	return ctx.EditEmbeds(messageID, embeds...)
}

// EditEmbedBuildersDM replaces the embeds of the message with the passed id in
// the direct message channel with the invoking user with those built from the
// passed *EmbedBuilders.
func EditEmbedBuildersDM(
	ctx *plugin.Context, messageID discord.MessageID,
	builders ...*EmbedBuilder,
) (*discord.Message, error) {
	embeds, err := BuildEmbeds(ctx.Localizer, builders...)
	if err != nil {
		return nil, err
	}

	return ctx.EditEmbeds(messageID, embeds...)
}
