package msgbuilder

import (
	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/pkg/i18n"
)

type (
	// EmbedBuilder is a builder used to build embeds.
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

// WithTitlelt sets the title (max. 256 characters) to the passed title.
func (b *EmbedBuilder) WithTitlelt(title i18n.Term) *EmbedBuilder {
	return b.WithTitlel(title.AsConfig())
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

// WithDescriptionlt sets the description (max. 2048 characters) to the passed
// description.
func (b *EmbedBuilder) WithDescriptionlt(description i18n.Term) *EmbedBuilder {
	return b.WithDescriptionl(description.AsConfig())
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

// WithFooterlt sets the text of the footer (max. 2048 characters) to the
// passed text.
func (b *EmbedBuilder) WithFooterlt(text i18n.Term) *EmbedBuilder {
	return b.WithFooterl(text.AsConfig())
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

// WithAuthorlt sets the author's name (max. 256 characters) to the passed
// name.
func (b *EmbedBuilder) WithAuthorlt(name i18n.Term) *EmbedBuilder {
	return b.WithAuthorl(name.AsConfig())
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

// WithFieldlt adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *EmbedBuilder) WithFieldlt(name, value i18n.Term) *EmbedBuilder {
	return b.WithFieldl(name.AsConfig(), value.AsConfig())
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

// WithInlinedFieldlt adds an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
func (b *EmbedBuilder) WithInlinedFieldlt(name, value i18n.Term) *EmbedBuilder {
	return b.WithInlinedFieldl(name.AsConfig(), value.AsConfig())
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
//nolint:funlen
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
		e.Footer = &discord.EmbedFooter{
			Icon: b.footer.icon,
		}

		if b.footer.text != nil {
			e.Footer.Text, err = l.Localize(b.footer.text)
			if err != nil {
				return discord.Embed{}, err
			}
		}
	}

	if b.imageURL != "" {
		e.Image = &discord.EmbedImage{
			URL: b.imageURL,
		}
	}

	if b.thumbnailURL != "" {
		e.Thumbnail = &discord.EmbedThumbnail{
			URL: b.thumbnailURL,
		}
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

// MustBuild is the same as Build, but panics if Build returns an error.
func (b *EmbedBuilder) MustBuild(l *i18n.Localizer) discord.Embed {
	e, err := b.Build(l)
	if err != nil {
		panic(err)
	}

	return e
}
