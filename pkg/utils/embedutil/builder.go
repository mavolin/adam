package embedutil

import (
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/utils/locutil"
)

type (
	// Builder is a utility struct used to build embeds.
	Builder struct {
		title       locutil.Text
		description locutil.Text

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
		text locutil.Text
		icon discord.URL
	}

	author struct {
		name locutil.Text
		icon discord.URL
		url  discord.URL
	}

	field struct {
		name    locutil.Text
		value   locutil.Text
		inlined bool
	}
)

// NewLocalizedEmbedBuilder creates a new Builder.
func NewBuilder() *Builder {
	return new(Builder)
}

// WithSimpleTitle adds a plain title (max. 256 characters) to the embed.
func (b *Builder) WithSimpleTitle(title string) *Builder {
	b.title = locutil.NewStaticText(title)
	return b
}

// WithSimpleTitlel adds a plain title (max. 256 characters) to the embed.
func (b *Builder) WithSimpleTitlel(title localization.Config) *Builder {
	b.title = locutil.NewLocalizedText(title)
	return b
}

// WithSimpleTitlelt adds a plain title (max. 256 characters) to the embed.
func (b *Builder) WithSimpleTitlelt(title localization.Term) *Builder {
	return b.WithSimpleTitlel(title.AsConfig())
}

// WithTitle adds a title (max. 256 characters) with a link to the embed.
func (b *Builder) WithTitle(title string, url discord.URL) *Builder {
	b.title = locutil.NewStaticText(title)
	b.url = url

	return b
}

// WithTitlel adds a title (max. 256 characters) with a link to the embed.
func (b *Builder) WithTitlel(title localization.Config, url discord.URL) *Builder {
	b.title = locutil.NewLocalizedText(title)
	b.url = url

	return b
}

// WithTitlelt adds a title (max. 256 characters) with a link to the embed.
func (b *Builder) WithTitlelt(title localization.Term, url discord.URL) *Builder {
	return b.WithTitlel(title.AsConfig(), url)
}

// WithDescription adds a description (max. 2048 characters) to the embed.
func (b *Builder) WithDescription(description string) *Builder {
	b.description = locutil.NewStaticText(description)
	return b
}

// WithDescriptionl adds a description (max. 2048 characters) to the embed.
func (b *Builder) WithDescriptionl(description localization.Config) *Builder {
	b.description = locutil.NewLocalizedText(description)
	return b
}

// WithDescriptionlt adds a description (max. 2048 characters) to the embed.
func (b *Builder) WithDescriptionlt(description localization.Term) *Builder {
	return b.WithDescriptionl(description.AsConfig())
}

// WithTimestamp adds a discord.Timestamp to the embed.
func (b *Builder) WithTimestamp(timestamp discord.Timestamp) *Builder {
	b.timestamp = timestamp
	return b
}

// WithTimestamp adds a timestamp of the current time to the embed.
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
		text: locutil.NewStaticText(text),
	}

	return b
}

// WithSimpleFooterl adds a plain footer (max. 2048 characters) to the embed.
func (b *Builder) WithSimpleFooterl(text localization.Config) *Builder {
	b.footer = &footer{
		text: locutil.NewLocalizedText(text),
	}

	return b
}

// WithSimpleFooterlt adds a plain footer (max. 2048 characters) to the embed.
func (b *Builder) WithSimpleFooterlt(text localization.Term) *Builder {
	return b.WithSimpleFooterl(text.AsConfig())
}

// WithFooter adds a footer (max. 2048 character) with an icon to the embed.
func (b *Builder) WithFooter(text string, icon discord.URL) *Builder {
	b.footer = &footer{
		text: locutil.NewStaticText(text),
		icon: icon,
	}

	return b
}

// WithFooterl adds a footer (max. 2048 character) with an icon to the embed.
func (b *Builder) WithFooterl(text localization.Config, icon discord.URL) *Builder {
	b.footer = &footer{
		text: locutil.NewLocalizedText(text),
		icon: icon,
	}

	return b
}

// WithFooterlt adds a footer (max. 2048 character) with an icon to the embed.
func (b *Builder) WithFooterlt(text localization.Term, icon discord.URL) *Builder {
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
		name: locutil.NewStaticText(name),
	}

	return b
}

// WithSimpleAuthorl adds a plain author (max. 256 characters) to the embed.
func (b *Builder) WithSimpleAuthorl(name localization.Config) *Builder {
	b.author = &author{
		name: locutil.NewLocalizedText(name),
	}

	return b
}

// WithSimpleAuthorlt adds a plain author (max. 256 characters) to the embed.
func (b *Builder) WithSimpleAuthorlt(name localization.Term) *Builder {
	return b.WithSimpleAuthorl(name.AsConfig())
}

// WithSimpleAuthorWithURL adds an author (max. 256 character) with a URL to
// the embed.
func (b *Builder) WithSimpleAuthorWithURL(name string, url discord.URL) *Builder {
	b.author = &author{
		name: locutil.NewStaticText(name),
		url:  url,
	}

	return b
}

// WithSimpleAuthorWithURLl adds an author (max. 256 character) with a URL to
// the embed.
func (b *Builder) WithSimpleAuthorWithURLl(name localization.Config, url discord.URL) *Builder {
	b.author = &author{
		name: locutil.NewLocalizedText(name),
		url:  url,
	}

	return b
}

// WithSimpleAuthorWithURLlt adds an author (max. 256 character) with a URL to
// the embed.
func (b *Builder) WithSimpleAuthorWithURLlt(name localization.Term, url discord.URL) *Builder {
	return b.WithSimpleAuthorWithURLl(name.AsConfig(), url)
}

// WithAuthor adds an author (max 256 characters) with an icon to the embed.
func (b *Builder) WithAuthor(name string, icon discord.URL) *Builder {
	b.author = &author{
		name: locutil.NewStaticText(name),
		icon: icon,
	}

	return b
}

// WithAuthorl adds an author (max 256 characters) with an icon to the embed.
func (b *Builder) WithAuthorl(name localization.Config, icon discord.URL) *Builder {
	b.author = &author{
		name: locutil.NewLocalizedText(name),
		icon: icon,
	}

	return b
}

// WithAuthorlt adds an author (max 256 characters) with an icon to the embed.
func (b *Builder) WithAuthorlt(name localization.Term, icon discord.URL) *Builder {
	return b.WithAuthorl(name.AsConfig(), icon)
}

// WithAuthorWithURL adds an author (max 256 characters) with an icon and a URL
// to the embed.
func (b *Builder) WithAuthorWithURL(name string, icon, url discord.URL) *Builder {
	b.author = &author{
		name: locutil.NewStaticText(name),
		icon: icon,
		url:  url,
	}

	return b
}

// WithAuthorWithURLl adds an author (max 256 characters) with an icon and a
// URL to the embed.
func (b *Builder) WithAuthorWithURLl(name localization.Config, icon, url discord.URL) *Builder {
	b.author = &author{
		name: locutil.NewLocalizedText(name),
		icon: icon,
		url:  url,
	}

	return b
}

// WithAuthorWithURLlt adds an author (max 256 characters) with an icon and a
// URL to the embed.
func (b *Builder) WithAuthorWithURLlt(name localization.Term, icon, url discord.URL) *Builder {
	return b.WithAuthorWithURLl(name.AsConfig(), icon, url)
}

// WithFieldt appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *Builder) WithField(name, value string) *Builder {
	b.withField(name, value, false)
	return b
}

// WithFieldl appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *Builder) WithFieldl(name, value localization.Config) *Builder {
	b.withFieldl(name, value, false)
	return b
}

// WithFieldlt appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *Builder) WithFieldlt(name, value localization.Term) *Builder {
	return b.WithFieldl(name.AsConfig(), value.AsConfig())
}

// WithInlinedField appends an inlined field (name: max. 256 characters, value:
// max 1024 characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *Builder) WithInlinedField(name, value string) *Builder {
	b.withField(name, value, true)
	return b
}

// WithInlinedFieldl appends an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *Builder) WithInlinedFieldl(name, value localization.Config) *Builder {
	b.withFieldl(name, value, true)
	return b
}

// WithInlinedFieldlt appends an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *Builder) WithInlinedFieldlt(name, value localization.Term) *Builder {
	return b.WithInlinedFieldl(name.AsConfig(), value.AsConfig())
}

func (b *Builder) withField(name, value string, inlined bool) {
	f := field{
		inlined: inlined,
	}

	f.name = locutil.NewStaticText(name)
	f.value = locutil.NewStaticText(value)

	b.fields = append(b.fields, f)
}

func (b *Builder) withFieldl(name, value localization.Config, inlined bool) {
	f := field{
		inlined: inlined,
	}

	f.name = locutil.NewLocalizedText(name)
	f.value = locutil.NewLocalizedText(value)

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
func (b *Builder) Build(l *localization.Localizer) (e discord.Embed, err error) {
	if !b.title.IsEmpty() {
		e.Title, err = b.title.Get(l)
		if err != nil {
			return
		}
	}

	if !b.description.IsEmpty() {
		e.Description, err = b.description.Get(l)
		if err != nil {
			return
		}
	}

	e.URL = b.url
	e.Timestamp = b.timestamp
	e.Color = b.color

	if b.footer != nil {
		e.Footer = &discord.EmbedFooter{
			Icon: b.footer.icon,
		}

		e.Footer.Text, err = b.footer.text.Get(l)
		if err != nil {
			return
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

		e.Author.Name, err = b.author.name.Get(l)
		if err != nil {
			return
		}
	}

	if len(b.fields) > 0 {
		e.Fields = make([]discord.EmbedField, len(b.fields))
	}

	for i, f := range b.fields {
		var name string

		if !f.name.IsEmpty() {
			name, err = f.name.Get(l)
			if err != nil {
				return
			}
		}

		var value string

		if !f.value.IsEmpty() {
			value, err = f.value.Get(l)
			if err != nil {
				return
			}
		}

		e.Fields[i] = discord.EmbedField{
			Name:   name,
			Value:  value,
			Inline: f.inlined,
		}
	}

	return
}

// MustBuild is the same as Build, but panics if Build returns an error.
func (b *Builder) MustBuild(l *localization.Localizer) discord.Embed {
	e, err := b.Build(l)
	if err != nil {
		panic(err)
	}

	return e
}
