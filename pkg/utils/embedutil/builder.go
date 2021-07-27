package embedutil

import (
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/adam/pkg/i18n"
)

type (
	// Builder is a utility struct used to build embeds.
	Builder struct {
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

// NewBuilder creates a new Builder.
func NewBuilder() *Builder {
	return new(Builder)
}

// WithTitle sets the title (max. 256 characters) to the passed title.
func (b *Builder) WithTitle(title string) *Builder {
	return b.WithTitlel(i18n.NewStaticConfig(title))
}

// WithTitlelt sets the title (max. 256 characters) to the passed title.
func (b *Builder) WithTitlelt(title i18n.Term) *Builder {
	return b.WithTitlel(title.AsConfig())
}

// WithTitlel sets the title (max. 256 characters) to the passed title.
func (b *Builder) WithTitlel(title *i18n.Config) *Builder {
	b.title = title
	return b
}

// WithTitleURL assigns a discord.URL to the title.
func (b *Builder) WithTitleURL(url discord.URL) *Builder {
	b.url = url
	return b
}

// WithDescription sets the description (max. 2048 characters) to the passed
// description.
func (b *Builder) WithDescription(description string) *Builder {
	return b.WithDescriptionl(i18n.NewStaticConfig(description))
}

// WithDescriptionlt sets the description (max. 2048 characters) to the passed
// description.
func (b *Builder) WithDescriptionlt(description i18n.Term) *Builder {
	return b.WithDescriptionl(description.AsConfig())
}

// WithDescriptionl sets the description (max. 2048 characters) to the passed
// description.
func (b *Builder) WithDescriptionl(description *i18n.Config) *Builder {
	b.description = description
	return b
}

// WithTimestamp sets the timestamp to the passed discord.Timestamp.
func (b *Builder) WithTimestamp(timestamp discord.Timestamp) *Builder {
	b.timestamp = timestamp
	return b
}

// WithTimestampNow sets the timestamp to a timestamp of the current time.
func (b *Builder) WithTimestampNow() *Builder {
	return b.WithTimestamp(discord.NowTimestamp())
}

// WithColor sets the color to the passed discord.Color.
func (b *Builder) WithColor(color discord.Color) *Builder {
	b.color = color
	return b
}

// WithFooter sets the text of the footer (max. 2048 characters) to the passed
// text.
func (b *Builder) WithFooter(text string) *Builder {
	return b.WithFooterl(i18n.NewStaticConfig(text))
}

// WithFooterlt sets the text of the footer (max. 2048 characters) to the
// passed text.
func (b *Builder) WithFooterlt(text i18n.Term) *Builder {
	return b.WithFooterl(text.AsConfig())
}

// WithFooterl sets the text of the footer (max. 2048 characters) to the passed
// text.
func (b *Builder) WithFooterl(text *i18n.Config) *Builder {
	if b.footer == nil {
		b.footer = &footer{text: text}
	} else {
		b.footer.text = text
	}

	return b
}

// WithFooterIcon sets the icon of the footer to the passed icon url.
func (b *Builder) WithFooterIcon(icon discord.URL) *Builder {
	if b.footer == nil {
		b.footer = &footer{icon: icon}
	} else {
		b.footer.icon = icon
	}

	return b
}

// WithImage sets the image to the passed image url.
func (b *Builder) WithImage(image discord.URL) *Builder {
	b.imageURL = image
	return b
}

// WithThumbnail sets the thumbnail to the passed thumbnail url.
func (b *Builder) WithThumbnail(thumbnail discord.URL) *Builder {
	b.thumbnailURL = thumbnail

	return b
}

// WithAuthor sets the author's name (max. 256 characters) to the passed
// name.
func (b *Builder) WithAuthor(name string) *Builder {
	return b.WithAuthorl(i18n.NewStaticConfig(name))
}

// WithAuthorlt sets the author's name (max. 256 characters) to the passed
// name.
func (b *Builder) WithAuthorlt(name i18n.Term) *Builder {
	return b.WithAuthorl(name.AsConfig())
}

// WithAuthorl sets the author's name (max. 256 characters) to the passed
// name.
func (b *Builder) WithAuthorl(name *i18n.Config) *Builder {
	if b.author == nil {
		b.author = &author{name: name}
	} else {
		b.author.name = name
	}

	return b
}

// WithAuthorURL assigns the author the passed discord.URL.
func (b *Builder) WithAuthorURL(url discord.URL) *Builder {
	if b.author == nil {
		b.author = &author{url: url}
	} else {
		b.author.url = url
	}

	return b
}

// WithAuthorIcon sets the icon of the author to the passed icon url.
func (b *Builder) WithAuthorIcon(icon discord.URL) *Builder {
	if b.author == nil {
		b.author = &author{icon: icon}
	} else {
		b.author.icon = icon
	}

	return b
}

// WithField adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *Builder) WithField(name, value string) *Builder {
	return b.WithFieldl(i18n.NewStaticConfig(name), i18n.NewStaticConfig(value))
}

// WithFieldlt adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *Builder) WithFieldlt(name, value i18n.Term) *Builder {
	return b.WithFieldl(name.AsConfig(), value.AsConfig())
}

// WithFieldl adds a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *Builder) WithFieldl(name, value *i18n.Config) *Builder {
	b.fields = append(b.fields, field{
		inlined: false,
		name:    name,
		value:   value,
	})

	return b
}

// WithInlinedField adds an inlined field (name: max. 256 characters, value:
// max 1024 characters) to the embed.
func (b *Builder) WithInlinedField(name, value string) *Builder {
	return b.WithInlinedFieldl(i18n.NewStaticConfig(name), i18n.NewStaticConfig(value))
}

// WithInlinedFieldlt adds an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
func (b *Builder) WithInlinedFieldlt(name, value i18n.Term) *Builder {
	return b.WithInlinedFieldl(name.AsConfig(), value.AsConfig())
}

// WithInlinedFieldl adds an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
func (b *Builder) WithInlinedFieldl(name, value *i18n.Config) *Builder {
	b.fields = append(b.fields, field{
		inlined: true,
		name:    name,
		value:   value,
	})

	return b
}

// Clone creates a copy of the Builder.
func (b *Builder) Clone() *Builder {
	cp := new(Builder)

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
func (b *Builder) Build(l *i18n.Localizer) (e discord.Embed, err error) {
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
func (b *Builder) MustBuild(l *i18n.Localizer) discord.Embed {
	e, err := b.Build(l)
	if err != nil {
		panic(err)
	}

	return e
}
