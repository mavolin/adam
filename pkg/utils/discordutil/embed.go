package discordutil

import (
	"github.com/diamondburned/arikawa/discord"

	"github.com/mavolin/adam/pkg/localization"
)

type (
	// EmbedBuilder is a utility struct used to build embeds.
	EmbedBuilder struct {
		title       *text
		description *text

		url discord.URL

		timestamp discord.Timestamp
		color     discord.Color

		footer       *embedFooter
		imageURL     discord.URL
		thumbnailURL discord.URL
		author       *embedAuthor
		fields       []embedField
	}

	embedFooter struct {
		text text
		icon discord.URL
	}

	embedAuthor struct {
		name text
		icon discord.URL
		url  discord.URL
	}

	embedField struct {
		name    *text
		value   *text
		inlined bool
	}
)

// NewLocalizedEmbedBuilder creates a new EmbedBuilder.
func NewEmbedBuilder() *EmbedBuilder {
	return new(EmbedBuilder)
}

// WithSimpleTitle adds a plain title (max. 256 characters) to the embed.
func (b *EmbedBuilder) WithSimpleTitle(title string) *EmbedBuilder {
	b.title = stringText(title)
	return b
}

// WithSimpleTitlel adds a plain title (max. 256 characters) to the embed.
func (b *EmbedBuilder) WithSimpleTitlel(title localization.Config) *EmbedBuilder {
	b.title = configText(title)
	return b
}

// WithSimpleTitlelt adds a plain title (max. 256 characters) to the embed.
func (b *EmbedBuilder) WithSimpleTitlelt(title localization.Term) *EmbedBuilder {
	return b.WithSimpleTitlel(title.AsConfig())
}

// WithTitle adds a title (max. 256 characters) with a link to the embed.
func (b *EmbedBuilder) WithTitle(title string, url discord.URL) *EmbedBuilder {
	b.title = stringText(title)
	b.url = url

	return b
}

// WithTitlel adds a title (max. 256 characters) with a link to the embed.
func (b *EmbedBuilder) WithTitlel(title localization.Config, url discord.URL) *EmbedBuilder {
	b.title = configText(title)
	b.url = url

	return b
}

// WithTitlelt adds a title (max. 256 characters) with a link to the embed.
func (b *EmbedBuilder) WithTitlelt(title localization.Term, url discord.URL) *EmbedBuilder {
	return b.WithTitlel(title.AsConfig(), url)
}

// WithDescription adds a description (max. 2048 characters) to the embed.
func (b *EmbedBuilder) WithDescription(description string) *EmbedBuilder {
	b.description = stringText(description)
	return b
}

// WithDescriptionl adds a description (max. 2048 characters) to the embed.
func (b *EmbedBuilder) WithDescriptionl(description localization.Config) *EmbedBuilder {
	b.description = configText(description)
	return b
}

// WithDescriptionlt adds a description (max. 2048 characters) to the embed.
func (b *EmbedBuilder) WithDescriptionlt(description localization.Term) *EmbedBuilder {
	return b.WithDescriptionl(description.AsConfig())
}

// WithTimestamp adds a discord.Timestamp to the embed.
func (b *EmbedBuilder) WithTimestamp(timestamp discord.Timestamp) *EmbedBuilder {
	b.timestamp = timestamp
	return b
}

// WithTimestamp adds a timestamp of the current time to the embed.
func (b *EmbedBuilder) WithTimestampNow() *EmbedBuilder {
	return b.WithTimestamp(discord.NowTimestamp())
}

// WithColor sets the color of the embed to the passed discord.Color.
func (b *EmbedBuilder) WithColor(color discord.Color) *EmbedBuilder {
	b.color = color
	return b
}

// WithSimpleFooter adds a plain footer (max. 2048 characters) to the embed.
func (b *EmbedBuilder) WithSimpleFooter(text string) *EmbedBuilder {
	b.footer = &embedFooter{
		text: *stringText(text),
	}

	return b
}

// WithSimpleFooterl adds a plain footer (max. 2048 characters) to the embed.
func (b *EmbedBuilder) WithSimpleFooterl(text localization.Config) *EmbedBuilder {
	b.footer = &embedFooter{
		text: *configText(text),
	}

	return b
}

// WithSimpleFooterlt adds a plain footer (max. 2048 characters) to the embed.
func (b *EmbedBuilder) WithSimpleFooterlt(text localization.Term) *EmbedBuilder {
	return b.WithSimpleFooterl(text.AsConfig())
}

// WithFooter adds a footer (max. 2048 character) with an icon to the embed.
func (b *EmbedBuilder) WithFooter(text string, icon discord.URL) *EmbedBuilder {
	b.footer = &embedFooter{
		text: *stringText(text),
		icon: icon,
	}

	return b
}

// WithFooterl adds a footer (max. 2048 character) with an icon to the embed.
func (b *EmbedBuilder) WithFooterl(text localization.Config, icon discord.URL) *EmbedBuilder {
	b.footer = &embedFooter{
		text: *configText(text),
		icon: icon,
	}

	return b
}

// WithFooterlt adds a footer (max. 2048 character) with an icon to the embed.
func (b *EmbedBuilder) WithFooterlt(text localization.Term, icon discord.URL) *EmbedBuilder {
	return b.WithFooterl(text.AsConfig(), icon)
}

// WithImage adds an image to the embed.
func (b *EmbedBuilder) WithImage(image discord.URL) *EmbedBuilder {
	b.imageURL = image
	return b
}

// WithThumbnail adds a thumbnail to the embed.
func (b *EmbedBuilder) WithThumbnail(thumbnail discord.URL) *EmbedBuilder {
	b.thumbnailURL = thumbnail

	return b
}

// WithSimpleAuthor adds a plain author (max. 256 characters) to the embed.
func (b *EmbedBuilder) WithSimpleAuthor(name string) *EmbedBuilder {
	b.author = &embedAuthor{
		name: *stringText(name),
	}

	return b
}

// WithSimpleAuthorl adds a plain author (max. 256 characters) to the embed.
func (b *EmbedBuilder) WithSimpleAuthorl(name localization.Config) *EmbedBuilder {
	b.author = &embedAuthor{
		name: *configText(name),
	}

	return b
}

// WithSimpleAuthorlt adds a plain author (max. 256 characters) to the embed.
func (b *EmbedBuilder) WithSimpleAuthorlt(name localization.Term) *EmbedBuilder {
	return b.WithSimpleAuthorl(name.AsConfig())
}

// WithSimpleAuthorWithURL adds an author (max. 256 character) with a URL to
// the embed.
func (b *EmbedBuilder) WithSimpleAuthorWithURL(name string, url discord.URL) *EmbedBuilder {
	b.author = &embedAuthor{
		name: *stringText(name),
		url:  url,
	}

	return b
}

// WithSimpleAuthorWithURLl adds an author (max. 256 character) with a URL to
// the embed.
func (b *EmbedBuilder) WithSimpleAuthorWithURLl(name localization.Config, url discord.URL) *EmbedBuilder {
	b.author = &embedAuthor{
		name: *configText(name),
		url:  url,
	}

	return b
}

// WithSimpleAuthorWithURLlt adds an author (max. 256 character) with a URL to
// the embed.
func (b *EmbedBuilder) WithSimpleAuthorWithURLlt(name localization.Term, url discord.URL) *EmbedBuilder {
	return b.WithSimpleAuthorWithURLl(name.AsConfig(), url)
}

// WithAuthor adds an author (max 256 characters) with an icon to the embed.
func (b *EmbedBuilder) WithAuthor(name string, icon discord.URL) *EmbedBuilder {
	b.author = &embedAuthor{
		name: *stringText(name),
		icon: icon,
	}

	return b
}

// WithAuthorl adds an author (max 256 characters) with an icon to the embed.
func (b *EmbedBuilder) WithAuthorl(name localization.Config, icon discord.URL) *EmbedBuilder {
	b.author = &embedAuthor{
		name: *configText(name),
		icon: icon,
	}

	return b
}

// WithAuthorlt adds an author (max 256 characters) with an icon to the embed.
func (b *EmbedBuilder) WithAuthorlt(name localization.Term, icon discord.URL) *EmbedBuilder {
	return b.WithAuthorl(name.AsConfig(), icon)
}

// WithAuthorWithURL adds an author (max 256 characters) with an icon and a URL
// to the embed.
func (b *EmbedBuilder) WithAuthorWithURL(name string, icon, url discord.URL) *EmbedBuilder {
	b.author = &embedAuthor{
		name: *stringText(name),
		icon: icon,
		url:  url,
	}

	return b
}

// WithAuthorWithURLl adds an author (max 256 characters) with an icon and a
// URL to the embed.
func (b *EmbedBuilder) WithAuthorWithURLl(name localization.Config, icon, url discord.URL) *EmbedBuilder {
	b.author = &embedAuthor{
		name: *configText(name),
		icon: icon,
		url:  url,
	}

	return b
}

// WithAuthorWithURLlt adds an author (max 256 characters) with an icon and a
// URL to the embed.
func (b *EmbedBuilder) WithAuthorWithURLlt(name localization.Term, icon, url discord.URL) *EmbedBuilder {
	return b.WithAuthorWithURLl(name.AsConfig(), icon, url)
}

// WithFieldt appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *EmbedBuilder) WithField(name, value string) *EmbedBuilder {
	b.withField(name, value, false)
	return b
}

// WithFieldl appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *EmbedBuilder) WithFieldl(name, value localization.Config) *EmbedBuilder {
	b.withFieldl(name, value, false)
	return b
}

// WithFieldlt appends a field (name: max. 256 characters, value: max 1024
// characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *EmbedBuilder) WithFieldlt(name, value localization.Term) *EmbedBuilder {
	return b.WithFieldl(name.AsConfig(), value.AsConfig())
}

// WithInlinedField appends an inlined field (name: max. 256 characters, value:
// max 1024 characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *EmbedBuilder) WithInlinedField(name, value string) *EmbedBuilder {
	b.withField(name, value, true)
	return b
}

// WithInlinedFieldl appends an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *EmbedBuilder) WithInlinedFieldl(name, value localization.Config) *EmbedBuilder {
	b.withFieldl(name, value, true)
	return b
}

// WithInlinedFieldlt appends an inlined field (name: max. 256 characters,
// value: max 1024 characters) to the embed.
// Name or value may be empty, in which case the field won't have a name or
// value.
func (b *EmbedBuilder) WithInlinedFieldlt(name, value localization.Term) *EmbedBuilder {
	return b.WithInlinedFieldl(name.AsConfig(), value.AsConfig())
}

func (b *EmbedBuilder) withField(name, value string, inlined bool) {
	f := embedField{
		inlined: inlined,
	}

	if name != "" {
		f.name = stringText(name)
	}

	if value != "" {
		f.value = stringText(value)
	}

	b.fields = append(b.fields, f)
}

func (b *EmbedBuilder) withFieldl(name, value localization.Config, inlined bool) {
	f := embedField{
		inlined: inlined,
	}

	if name.IsValid() {
		f.name = configText(name)
	}

	if value.IsValid() {
		f.value = configText(value)
	}

	b.fields = append(b.fields, f)
}

// Build builds the discord.Embed.
func (b *EmbedBuilder) Build(l *localization.Localizer) (e discord.Embed, err error) {
	if b.title != nil {
		e.Title, err = b.title.localize(l)
		if err != nil {
			return
		}
	}

	if b.description != nil {
		e.Description, err = b.description.localize(l)
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

		e.Footer.Text, err = b.footer.text.localize(l)
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

		e.Author.Name, err = b.author.name.localize(l)
		if err != nil {
			return
		}
	}

	if len(b.fields) > 0 {
		e.Fields = make([]discord.EmbedField, len(b.fields))
	}

	for i, f := range b.fields {
		var name string

		if f.name != nil {
			name, err = f.name.localize(l)
			if err != nil {
				return
			}
		}

		var value string

		if f.value != nil {
			value, err = f.value.localize(l)
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
func (b *EmbedBuilder) MustBuild(l *localization.Localizer) discord.Embed {
	e, err := b.Build(l)
	if err != nil {
		panic(err)
	}

	return e
}

type text struct {
	string string
	config localization.Config
}

func stringText(src string) *text {
	return &text{
		string: src,
	}
}

func configText(src localization.Config) *text {
	return &text{
		config: src,
	}
}

func (t text) localize(l *localization.Localizer) (s string, err error) {
	s = t.string

	if s == "" && t.config.IsValid() {
		return l.Localize(t.config)
	}

	return
}
