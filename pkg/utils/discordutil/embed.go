package discordutil

import (
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/localization"
)

// EmbedBuilder is a utility struct, used to create embeds.
type EmbedBuilder struct {
	e discord.Embed
}

// NewEmbedBuilder creates a new embed builder.
func NewEmbedBuilder() *EmbedBuilder {
	return new(EmbedBuilder)
}

// WithSimpleTitle adds a plain title (max. 256 characters) to the embed.
func (e *EmbedBuilder) WithSimpleTitle(title string) *EmbedBuilder {
	e.e.Title = title
	return e
}

// WithTitle adds a title (max. 256 characters) with a link to the embed.
func (e *EmbedBuilder) WithTitle(title string, url discord.URL) *EmbedBuilder {
	e.e.Title = title
	e.e.URL = url
	return e
}

// WithDescription adds a description (max. 2048 characters) to the embed.
func (e *EmbedBuilder) WithDescription(description string) *EmbedBuilder {
	e.e.Description = description
	return e
}

// WithTimestamp adds a discord.Timestamp to the embed.
func (e *EmbedBuilder) WithTimestamp(timestamp discord.Timestamp) *EmbedBuilder {
	e.e.Timestamp = timestamp
	return e
}

// WithTimestamp adds a timestamp of the current time to the embed.
func (e *EmbedBuilder) WithTimestampNow() *EmbedBuilder {
	return e.WithTimestamp(discord.NowTimestamp())
}

// WithColor sets the color of the embed to the passed discord.Color.
func (e *EmbedBuilder) WithColor(color discord.Color) *EmbedBuilder {
	e.e.Color = color
	return e
}

// WithSimpleFooter adds a plain footer (max. 2048 characters) to the embed.
func (e *EmbedBuilder) WithSimpleFooter(text string) *EmbedBuilder {
	e.e.Footer = &discord.EmbedFooter{
		Text: text,
	}

	return e
}

// WithFooter adds a footer (max. 2048 character) with an icon to the embed.
func (e *EmbedBuilder) WithFooter(text string, icon discord.URL) *EmbedBuilder {
	e.e.Footer = &discord.EmbedFooter{
		Text: text,
		Icon: icon,
	}

	return e
}

// WithImage adds an image to the embed.
func (e *EmbedBuilder) WithImage(image discord.URL) *EmbedBuilder {
	e.e.Image = &discord.EmbedImage{
		URL: image,
	}

	return e
}

// WithThumbnail adds a thumbnail to the embed.
func (e *EmbedBuilder) WithThumbnail(thumbnail discord.URL) *EmbedBuilder {
	e.e.Thumbnail = &discord.EmbedThumbnail{
		URL: thumbnail,
	}

	return e
}

// WithSimpleAuthor adds a plain author (max. 256 characters) to the embed.
func (e *EmbedBuilder) WithSimpleAuthor(name string) *EmbedBuilder {
	e.e.Author = &discord.EmbedAuthor{
		Name: name,
	}

	return e
}

// WithSimpleAuthorWithURL adds an author (max. 256 character) with a URL to
// the embed.
func (e *EmbedBuilder) WithSimpleAuthorWithURL(name string, url discord.URL) *EmbedBuilder {
	e.e.Author = &discord.EmbedAuthor{
		Name: name,
		URL:  url,
	}

	return e
}

// WithAuthor adds an author (max 256 characters) with an icon to the embed.
func (e *EmbedBuilder) WithAuthor(name string, icon discord.URL) *EmbedBuilder {
	e.e.Author = &discord.EmbedAuthor{
		Name: name,
		Icon: icon,
	}

	return e
}

// WithAuthorWithURL adds an author (max 256 characters) with an icon and a URL
// to the embed.
func (e *EmbedBuilder) WithAuthorWithURL(name string, icon, url discord.URL) *EmbedBuilder {
	e.e.Author = &discord.EmbedAuthor{
		Name: name,
		URL:  url,
		Icon: icon,
	}

	return e
}

// WithFieldt appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (e *EmbedBuilder) WithField(name, value string) *EmbedBuilder {
	e.e.Fields = append(e.e.Fields, discord.EmbedField{
		Name:   name,
		Value:  value,
		Inline: false,
	})

	return e
}

// WithInlinedField appends an inlined field (name: max. 256 characters, value: max 1024
// // characters) to the embed.
func (e *EmbedBuilder) WithInlinedField(name, value string) *EmbedBuilder {
	e.e.Fields = append(e.e.Fields, discord.EmbedField{
		Name:   name,
		Value:  value,
		Inline: true,
	})

	return e
}

// Build builds the discord.Embed.
func (e *EmbedBuilder) Build() discord.Embed {
	return e.e
}

type (
	LocalizedEmbedBuilder struct {
		title       *localization.Config
		description *localization.Config

		url discord.URL

		timestamp discord.Timestamp
		color     discord.Color

		footer       *localizedFooter
		imageURL     discord.URL
		thumbnailURL discord.URL
		author       *localizedAuthor
		fields       []localizedField
	}

	localizedFooter struct {
		text localization.Config
		icon discord.URL
	}

	localizedAuthor struct {
		name localization.Config
		icon discord.URL
		url  discord.URL
	}

	localizedField struct {
		name    *localization.Config
		value   *localization.Config
		inlined bool
	}
)

// NewLocalizedEmbedBuilder creates a new LocalizedEmbedBuilder.
func NewLocalizedEmbedBuilder() *LocalizedEmbedBuilder {
	return new(LocalizedEmbedBuilder)
}

// WithSimpleTitle adds a plain title (max. 256 characters) to the embed.
func (b *LocalizedEmbedBuilder) WithSimpleTitle(title localization.Config) *LocalizedEmbedBuilder {
	b.title = &title
	return b
}

// WithSimpleTitlet adds a plain title (max. 256 characters) to the embed.
func (b *LocalizedEmbedBuilder) WithSimpleTitlet(titleTerm string) *LocalizedEmbedBuilder {
	return b.WithSimpleTitle(localization.Term(titleTerm))
}

// WithTitle adds a title (max. 256 characters) with a link to the embed.
func (b *LocalizedEmbedBuilder) WithTitle(title localization.Config, url discord.URL) *LocalizedEmbedBuilder {
	b.title = &title
	b.url = url
	return b
}

// WithTitlet adds a title (max. 256 characters) with a link to the embed.
func (b *LocalizedEmbedBuilder) WithTitlet(titleTerm string, url discord.URL) *LocalizedEmbedBuilder {
	return b.WithTitle(localization.Term(titleTerm), url)
}

// WithDescription adds a description (max. 2048 characters) to the embed.
func (b *LocalizedEmbedBuilder) WithDescription(description localization.Config) *LocalizedEmbedBuilder {
	b.description = &description
	return b
}

// WithDescriptiont adds a description (max. 2048 characters) to the embed.
func (b *LocalizedEmbedBuilder) WithDescriptiont(descriptionTerm string) *LocalizedEmbedBuilder {
	return b.WithDescription(localization.Term(descriptionTerm))
}

// WithTimestamp adds a discord.Timestamp to the embed.
func (b *LocalizedEmbedBuilder) WithTimestamp(timestamp discord.Timestamp) *LocalizedEmbedBuilder {
	b.timestamp = timestamp
	return b
}

// WithTimestamp adds a timestamp of the current time to the embed.
func (b *LocalizedEmbedBuilder) WithTimestampNow() *LocalizedEmbedBuilder {
	return b.WithTimestamp(discord.NowTimestamp())
}

// WithColor sets the color of the embed to the passed discord.Color.
func (b *LocalizedEmbedBuilder) WithColor(color discord.Color) *LocalizedEmbedBuilder {
	b.color = color
	return b
}

// WithSimpleFooter adds a plain footer (max. 2048 characters) to the embed.
func (b *LocalizedEmbedBuilder) WithSimpleFooter(text localization.Config) *LocalizedEmbedBuilder {
	b.footer = &localizedFooter{
		text: text,
	}

	return b
}

// WithSimpleFootert adds a plain footer (max. 2048 characters) to the embed.
func (b *LocalizedEmbedBuilder) WithSimpleFootert(textTerm string) *LocalizedEmbedBuilder {
	return b.WithSimpleFooter(localization.Term(textTerm))
}

// WithFooter adds a footer (max. 2048 character) with an icon to the embed.
func (b *LocalizedEmbedBuilder) WithFooter(text localization.Config, icon discord.URL) *LocalizedEmbedBuilder {
	b.footer = &localizedFooter{
		text: text,
		icon: icon,
	}

	return b
}

// WithFootert adds a footer (max. 2048 character) with an icon to the embed.
func (b *LocalizedEmbedBuilder) WithFootert(textTerm string, icon discord.URL) *LocalizedEmbedBuilder {
	return b.WithFooter(localization.Term(textTerm), icon)
}

// WithImage adds an image to the embed.
func (b *LocalizedEmbedBuilder) WithImage(image discord.URL) *LocalizedEmbedBuilder {
	b.imageURL = image
	return b
}

// WithThumbnail adds a thumbnail to the embed.
func (b *LocalizedEmbedBuilder) WithThumbnail(thumbnail discord.URL) *LocalizedEmbedBuilder {
	b.thumbnailURL = thumbnail

	return b
}

// WithSimpleAuthor adds a plain author (max. 256 characters) to the embed.
func (b *LocalizedEmbedBuilder) WithSimpleAuthor(name localization.Config) *LocalizedEmbedBuilder {
	b.author = &localizedAuthor{
		name: name,
	}

	return b
}

// WithSimpleAuthort adds a plain author (max. 256 characters) to the embed.
func (b *LocalizedEmbedBuilder) WithSimpleAuthort(nameTerm string) *LocalizedEmbedBuilder {
	return b.WithSimpleAuthor(localization.Term(nameTerm))
}

// WithSimpleAuthorWithURL adds an author (max. 256 character) with a URL to
// the embed.
func (b *LocalizedEmbedBuilder) WithSimpleAuthorWithURL(name localization.Config, url discord.URL) *LocalizedEmbedBuilder {
	b.author = &localizedAuthor{
		name: name,
		url:  url,
	}

	return b
}

// WithSimpleAuthorWithURLt adds an author (max. 256 character) with a URL to
// the embed.
func (b *LocalizedEmbedBuilder) WithSimpleAuthorWithURLt(nameTerm string, url discord.URL) *LocalizedEmbedBuilder {
	return b.WithSimpleAuthorWithURL(localization.Term(nameTerm), url)
}

// WithAuthor adds an author (max 256 characters) with an icon to the embed.
func (b *LocalizedEmbedBuilder) WithAuthor(name localization.Config, icon discord.URL) *LocalizedEmbedBuilder {
	b.author = &localizedAuthor{
		name: name,
		icon: icon,
	}

	return b
}

// WithAuthort adds an author (max 256 characters) with an icon to the embed.
func (b *LocalizedEmbedBuilder) WithAuthort(nameTerm string, icon discord.URL) *LocalizedEmbedBuilder {
	return b.WithAuthor(localization.Term(nameTerm), icon)
}

// WithAuthorWithURLt adds an author (max 256 characters) with an icon and a URL
// to the embed.
func (b *LocalizedEmbedBuilder) WithAuthorWithURL(name localization.Config, icon, url discord.URL) *LocalizedEmbedBuilder {
	b.author = &localizedAuthor{
		name: name,
		icon: icon,
		url:  url,
	}

	return b
}

// WithAuthorWithURLt adds an author (max 256 characters) with an icon and a URL
// to the embed.
func (b *LocalizedEmbedBuilder) WithAuthorWithURLt(nameTerm string, icon, url discord.URL) *LocalizedEmbedBuilder {
	return b.WithAuthorWithURL(localization.Term(nameTerm), icon, url)
}

// WithField appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *LocalizedEmbedBuilder) WithField(name, value localization.Config) *LocalizedEmbedBuilder {
	b.fields = append(b.fields, localizedField{
		name:    &name,
		value:   &value,
		inlined: false,
	})

	return b
}

// WithFieldt appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *LocalizedEmbedBuilder) WithFieldt(nameTerm, valueTerm string) *LocalizedEmbedBuilder {
	return b.WithField(localization.Term(nameTerm), localization.Term(valueTerm))
}

// WithInlinedField appends an inlined field (name: max. 256 characters, value: max 1024
// // characters) to the embed.
func (b *LocalizedEmbedBuilder) WithInlinedField(name, value localization.Config) *LocalizedEmbedBuilder {
	b.fields = append(b.fields, localizedField{
		name:    &name,
		value:   &value,
		inlined: true,
	})

	return b
}

// WithInlinedFieldt appends an inlined field (name: max. 256 characters, value: max 1024
// characters) to the embed.
func (b *LocalizedEmbedBuilder) WithInlinedFieldt(nameTerm, valueTerm string) *LocalizedEmbedBuilder {
	return b.WithInlinedField(localization.Term(nameTerm), localization.Term(valueTerm))
}

// Build builds the discord.Embed.
func (b *LocalizedEmbedBuilder) Build(l *localization.Localizer) (e discord.Embed, err error) {
	if b.title != nil {
		e.Title, err = l.Localize(*b.title)
		if err != nil {
			return discord.Embed{}, err
		}
	}

	if b.description != nil {
		e.Description, err = l.Localize(*b.description)
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
			name, err = l.Localize(*f.name)
			if err != nil {
				return discord.Embed{}, err
			}
		}

		var value string

		if f.value != nil {
			value, err = l.Localize(*f.value)
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

	return
}
