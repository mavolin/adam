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

// WithSimpleTitle adds a plain title (max. 256 characters) to the embed.
func (b *Builder) WithSimpleTitle(title string) *Builder {
	b.title = i18n.NewStaticConfig(title)
	return b
}

// WithSimpleTitlel adds a plain title (max. 256 characters) to the embed.
func (b *Builder) WithSimpleTitlel(title *i18n.Config) *Builder {
	b.title = (*i18n.Config)(title)
	return b
}

// WithSimpleTitlelt adds a plain title (max. 256 characters) to the embed.
func (b *Builder) WithSimpleTitlelt(title i18n.Term) *Builder {
	return b.WithSimpleTitlel(title.AsConfig())
}

// WithTitle adds a title (max. 256 characters) with a link to the embed.
func (b *Builder) WithTitle(title string, url discord.URL) *Builder {
	b.title = i18n.NewStaticConfig(title)
	b.url = url

	return b
}

// WithTitlel adds a title (max. 256 characters) with a link to the embed.
func (b *Builder) WithTitlel(title *i18n.Config, url discord.URL) *Builder {
	b.title = (*i18n.Config)(title)
	b.url = url

	return b
}

// WithTitlelt adds a title (max. 256 characters) with a link to the embed.
func (b *Builder) WithTitlelt(title i18n.Term, url discord.URL) *Builder {
	return b.WithTitlel(title.AsConfig(), url)
}

// WithDescription adds a description (max. 2048 characters) to the embed.
func (b *Builder) WithDescription(description string) *Builder {
	b.description = i18n.NewStaticConfig(description)
	return b
}

// WithDescriptionl adds a description (max. 2048 characters) to the embed.
func (b *Builder) WithDescriptionl(description *i18n.Config) *Builder {
	b.description = (*i18n.Config)(description)
	return b
}

// WithDescriptionlt adds a description (max. 2048 characters) to the embed.
func (b *Builder) WithDescriptionlt(description i18n.Term) *Builder {
	return b.WithDescriptionl(description.AsConfig())
}

// WithTimestamp adds a discord.Timestamp to the embed.
func (b *Builder) WithTimestamp(timestamp discord.Timestamp) *Builder {
	b.timestamp = timestamp
	return b
}

// WithTimestampNow adds a timestamp of the current time to the embed.
func (b *Builder) WithTimestampNow() *Builder {
	return b.WithTimestamp(discord.NowTimestamp())
}

// WithColor sets the color of the embed to the passed discord.Color.
func (b *Builder) WithColor(color discord.Color) *Builder {
	b.color = color
	return b
}

// WithSimpleFooter adds a plain footer (max. 2048 characters) to the embed.
func (b *Builder) WithSimpleFooter(text string) *Builder {
	b.footer = &footer{
		text: i18n.NewStaticConfig(text),
	}

	return b
}

// WithSimpleFooterl adds a plain footer (max. 2048 characters) to the embed.
func (b *Builder) WithSimpleFooterl(text *i18n.Config) *Builder {
	b.footer = &footer{
		text: (*i18n.Config)(text),
	}

	return b
}

// WithSimpleFooterlt adds a plain footer (max. 2048 characters) to the embed.
func (b *Builder) WithSimpleFooterlt(text i18n.Term) *Builder {
	return b.WithSimpleFooterl(text.AsConfig())
}

// WithFooter adds a footer (max. 2048 character) with an icon to the embed.
func (b *Builder) WithFooter(text string, icon discord.URL) *Builder {
	b.footer = &footer{
		text: i18n.NewStaticConfig(text),
		icon: icon,
	}

	return b
}

// WithFooterl adds a footer (max. 2048 character) with an icon to the embed.
func (b *Builder) WithFooterl(text *i18n.Config, icon discord.URL) *Builder {
	b.footer = &footer{
		text: (*i18n.Config)(text),
		icon: icon,
	}

	return b
}

// WithFooterlt adds a footer (max. 2048 character) with an icon to the embed.
func (b *Builder) WithFooterlt(text i18n.Term, icon discord.URL) *Builder {
	return b.WithFooterl(text.AsConfig(), icon)
}

// WithImage adds an image to the embed.
func (b *Builder) WithImage(image discord.URL) *Builder {
	b.imageURL = image
	return b
}

// WithThumbnail adds a thumbnail to the embed.
func (b *Builder) WithThumbnail(thumbnail discord.URL) *Builder {
	b.thumbnailURL = thumbnail

	return b
}

// WithSimpleAuthor adds a plain author (max. 256 characters) to the embed.
func (b *Builder) WithSimpleAuthor(name string) *Builder {
	b.author = &author{
		name: i18n.NewStaticConfig(name),
	}

	return b
}

// WithSimpleAuthorl adds a plain author (max. 256 characters) to the embed.
func (b *Builder) WithSimpleAuthorl(name *i18n.Config) *Builder {
	b.author = &author{
		name: (*i18n.Config)(name),
	}

	return b
}

// WithSimpleAuthorlt adds a plain author (max. 256 characters) to the embed.
func (b *Builder) WithSimpleAuthorlt(name i18n.Term) *Builder {
	return b.WithSimpleAuthorl(name.AsConfig())
}

// WithSimpleAuthorWithURL adds an author (max. 256 character) with a URL to
// the embed.
func (b *Builder) WithSimpleAuthorWithURL(name string, url discord.URL) *Builder {
	b.author = &author{
		name: i18n.NewStaticConfig(name),
		url:  url,
	}

	return b
}

// WithSimpleAuthorWithURLl adds an author (max. 256 character) with a URL to
// the embed.
func (b *Builder) WithSimpleAuthorWithURLl(name *i18n.Config, url discord.URL) *Builder {
	b.author = &author{
		name: (*i18n.Config)(name),
		url:  url,
	}

	return b
}

// WithSimpleAuthorWithURLlt adds an author (max. 256 character) with a URL to
// the embed.
func (b *Builder) WithSimpleAuthorWithURLlt(name i18n.Term, url discord.URL) *Builder {
	return b.WithSimpleAuthorWithURLl(name.AsConfig(), url)
}

// WithAuthor adds an author (max 256 characters) with an icon to the embed.
func (b *Builder) WithAuthor(name string, icon discord.URL) *Builder {
	b.author = &author{
		name: i18n.NewStaticConfig(name),
		icon: icon,
	}

	return b
}

// WithAuthorl adds an author (max 256 characters) with an icon to the embed.
func (b *Builder) WithAuthorl(name *i18n.Config, icon discord.URL) *Builder {
	b.author = &author{
		name: (*i18n.Config)(name),
		icon: icon,
	}

	return b
}

// WithAuthorlt adds an author (max 256 characters) with an icon to the embed.
func (b *Builder) WithAuthorlt(name i18n.Term, icon discord.URL) *Builder {
	return b.WithAuthorl(name.AsConfig(), icon)
}

// WithAuthorWithURL adds an author (max 256 characters) with an icon and a URL
// to the embed.
func (b *Builder) WithAuthorWithURL(name string, icon, url discord.URL) *Builder {
	b.author = &author{
		name: i18n.NewStaticConfig(name),
		icon: icon,
		url:  url,
	}

	return b
}

// WithAuthorWithURLl adds an author (max 256 characters) with an icon and a
// URL to the embed.
func (b *Builder) WithAuthorWithURLl(name *i18n.Config, icon, url discord.URL) *Builder {
	b.author = &author{
		name: (*i18n.Config)(name),
		icon: icon,
		url:  url,
	}

	return b
}

// WithAuthorWithURLlt adds an author (max 256 characters) with an icon and a
// URL to the embed.
func (b *Builder) WithAuthorWithURLlt(name i18n.Term, icon, url discord.URL) *Builder {
	return b.WithAuthorWithURLl(name.AsConfig(), icon, url)
}

// WithField appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *Builder) WithField(name, value string) *Builder {
	b.withField(name, value, false)
	return b
}

// WithFieldl appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *Builder) WithFieldl(name, value *i18n.Config) *Builder {
	b.withFieldl(name, value, false)
	return b
}

// WithFieldlt appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *Builder) WithFieldlt(name, value i18n.Term) *Builder {
	return b.WithFieldl(name.AsConfig(), value.AsConfig())
}

// WithInlinedField appends an inlined field (name: max. 256 characters, value:
// max 1024 characters) to the embed.
func (b *Builder) WithInlinedField(name, value string) *Builder {
	b.withField(name, value, true)
	return b
}

// WithInlinedFieldl appends an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
func (b *Builder) WithInlinedFieldl(name, value *i18n.Config) *Builder {
	b.withFieldl(name, value, true)
	return b
}

// WithInlinedFieldlt appends an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
func (b *Builder) WithInlinedFieldlt(name, value i18n.Term) *Builder {
	return b.WithInlinedFieldl(name.AsConfig(), value.AsConfig())
}

func (b *Builder) withField(name, value string, inlined bool) {
	f := field{
		inlined: inlined,
	}

	if len(name) > 0 {
		f.name = i18n.NewStaticConfig(name)
	}

	if len(value) > 0 {
		f.value = i18n.NewStaticConfig(value)
	}

	b.fields = append(b.fields, f)
}

func (b *Builder) withFieldl(name, value *i18n.Config, inlined bool) {
	f := field{
		inlined: inlined,
	}

	f.name = (*i18n.Config)(name)
	f.value = (*i18n.Config)(value)

	b.fields = append(b.fields, f)
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
func (b *Builder) Build(l *i18n.Localizer) (e discord.Embed, err error) { //nolint:funlen
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

		e.Footer.Text, err = l.Localize(b.footer.text)
		if err != nil {
			return discord.Embed{}, err
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

		e.Author.Name, err = l.Localize(b.author.name)
		if err != nil {
			return discord.Embed{}, err
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
