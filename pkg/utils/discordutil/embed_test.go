package discordutil

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/mock"
)

func TestEmbedBuilder_WithSimpleTitle(t *testing.T) {
	title := "abc"

	expect := discord.Embed{
		Title: title,
	}

	actual := NewEmbedBuilder().
		WithSimpleTitle(title).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTitle(t *testing.T) {
	var (
		title = "abc"
		url   = "def"
	)

	expect := discord.Embed{
		Title: title,
		URL:   url,
	}

	actual := NewEmbedBuilder().
		WithTitle(title, url).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithDescription(t *testing.T) {
	description := "abc"

	expect := discord.Embed{
		Description: description,
	}

	actual := NewEmbedBuilder().
		WithDescription(description).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTimestamp(t *testing.T) {
	timestamp := discord.NowTimestamp()

	expect := discord.Embed{
		Timestamp: timestamp,
	}

	actual := NewEmbedBuilder().
		WithTimestamp(timestamp).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithColor(t *testing.T) {
	var color discord.Color = 123

	expect := discord.Embed{
		Color: color,
	}

	actual := NewEmbedBuilder().
		WithColor(color).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleFooter(t *testing.T) {
	text := "abc"

	expect := discord.Embed{
		Footer: &discord.EmbedFooter{
			Text: text,
		},
	}

	actual := NewEmbedBuilder().
		WithSimpleFooter(text).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooter(t *testing.T) {
	var (
		text = "abc"
		icon = "def"
	)

	expect := discord.Embed{
		Footer: &discord.EmbedFooter{
			Text: text,
			Icon: icon,
		},
	}

	actual := NewEmbedBuilder().
		WithFooter(text, icon).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithImage(t *testing.T) {
	image := "abc"

	expect := discord.Embed{
		Image: &discord.EmbedImage{
			URL: image,
		},
	}

	actual := NewEmbedBuilder().
		WithImage(image).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithThumbnail(t *testing.T) {
	thumbnail := "abc"

	expect := discord.Embed{
		Thumbnail: &discord.EmbedThumbnail{
			URL: thumbnail,
		},
	}

	actual := NewEmbedBuilder().
		WithThumbnail(thumbnail).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleAuthor(t *testing.T) {
	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
		},
	}

	actual := NewEmbedBuilder().
		WithSimpleAuthor(name).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleAuthorWithURL(t *testing.T) {
	var (
		name = "abc"
		url  = "def"
	)

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
			URL:  url,
		},
	}

	actual := NewEmbedBuilder().
		WithSimpleAuthorWithURL(name, url).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthor(t *testing.T) {
	var (
		name = "abc"
		icon = "def"
	)

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
			Icon: icon,
		},
	}

	actual := NewEmbedBuilder().
		WithAuthor(name, icon).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthorWithURL(t *testing.T) {
	var (
		name = "abc"
		icon = "def"
		url  = "ghi"
	)

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
			Icon: icon,
			URL:  url,
		},
	}

	actual := NewEmbedBuilder().
		WithAuthorWithURL(name, icon, url).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithField(t *testing.T) {
	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: false,
	}

	expect := discord.Embed{
		Fields: []discord.EmbedField{field},
	}

	actual := NewEmbedBuilder().
		WithField(field.Name, field.Value).
		Build()

	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithInlinedField(t *testing.T) {
	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: true,
	}

	expect := discord.Embed{
		Fields: []discord.EmbedField{field},
	}

	actual := NewEmbedBuilder().
		WithInlinedField(field.Name, field.Value).
		Build()

	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithSimpleTitle(t *testing.T) {
	title := "abc"

	expect := discord.Embed{
		Title: title,
	}

	l := mock.
		NewLocalizer().
		On("a", title).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithSimpleTitle(localization.QuickConfig("a")).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithSimpleTitlet(t *testing.T) {
	title := "abc"

	expect := discord.Embed{
		Title: title,
	}

	l := mock.
		NewLocalizer().
		On("a", title).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithSimpleTitlet("a").
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithTitle(t *testing.T) {
	var (
		title = "abc"
		url   = "def"
	)

	expect := discord.Embed{
		Title: title,
		URL:   url,
	}

	l := mock.
		NewLocalizer().
		On("a", title).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithTitle(localization.QuickConfig("a"), url).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithTitlet(t *testing.T) {
	var (
		title = "abc"
		url   = "def"
	)

	expect := discord.Embed{
		Title: title,
		URL:   url,
	}

	l := mock.
		NewLocalizer().
		On("a", title).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithTitlet("a", url).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithDescription(t *testing.T) {
	description := "abc"

	expect := discord.Embed{
		Description: description,
	}

	l := mock.
		NewLocalizer().
		On("a", description).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithDescription(localization.QuickConfig("a")).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithDescriptiont(t *testing.T) {
	description := "abc"

	expect := discord.Embed{
		Description: description,
	}

	l := mock.
		NewLocalizer().
		On("a", description).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithDescriptiont("a").
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithTimestamp(t *testing.T) {
	timestamp := discord.NowTimestamp()

	expect := discord.Embed{
		Timestamp: timestamp,
	}

	l := mock.NewLocalizer().Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithTimestamp(timestamp).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithColor(t *testing.T) {
	var color discord.Color = 123

	expect := discord.Embed{
		Color: color,
	}

	l := mock.NewLocalizer().Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithColor(color).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithSimpleFooter(t *testing.T) {
	text := "abc"

	expect := discord.Embed{
		Footer: &discord.EmbedFooter{
			Text: text,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", text).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithSimpleFooter(localization.QuickConfig("a")).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithSimpleFootert(t *testing.T) {
	text := "abc"

	expect := discord.Embed{
		Footer: &discord.EmbedFooter{
			Text: text,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", text).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithSimpleFootert("a").
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithFooter(t *testing.T) {
	var (
		text = "abc"
		icon = "def"
	)

	expect := discord.Embed{
		Footer: &discord.EmbedFooter{
			Text: text,
			Icon: icon,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", text).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithFooter(localization.QuickConfig("a"), icon).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithFootert(t *testing.T) {
	var (
		text = "abc"
		icon = "def"
	)

	expect := discord.Embed{
		Footer: &discord.EmbedFooter{
			Text: text,
			Icon: icon,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", text).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithFootert("a", icon).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithImage(t *testing.T) {
	image := "abc"

	expect := discord.Embed{
		Image: &discord.EmbedImage{
			URL: image,
		},
	}

	l := mock.NewLocalizer().Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithImage(image).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithThumbnail(t *testing.T) {
	thumbnail := "abc"

	expect := discord.Embed{
		Thumbnail: &discord.EmbedThumbnail{
			URL: thumbnail,
		},
	}

	l := mock.NewLocalizer().Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithThumbnail(thumbnail).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithSimpleAuthor(t *testing.T) {
	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", name).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithSimpleAuthor(localization.QuickConfig("a")).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithSimpleAuthort(t *testing.T) {
	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", name).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithSimpleAuthort("a").
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithSimpleAuthorWithURL(t *testing.T) {
	var (
		name = "abc"
		url  = "def"
	)

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
			URL:  url,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", name).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithSimpleAuthorWithURL(localization.QuickConfig("a"), url).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithSimpleAuthorWithURLt(t *testing.T) {
	var (
		name = "abc"
		url  = "def"
	)

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
			URL:  url,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", name).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithSimpleAuthorWithURLt("a", url).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithAuthor(t *testing.T) {
	var (
		name = "abc"
		icon = "def"
	)

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
			Icon: icon,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", name).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithAuthor(localization.QuickConfig("a"), icon).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithAuthort(t *testing.T) {
	var (
		name = "abc"
		icon = "def"
	)

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
			Icon: icon,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", name).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithAuthort("a", icon).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithAuthorWithURL(t *testing.T) {
	var (
		name = "abc"
		icon = "def"
		url  = "ghi"
	)

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
			Icon: icon,
			URL:  url,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", name).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithAuthorWithURL(localization.QuickConfig("a"), icon, url).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithAuthorWithURLt(t *testing.T) {
	var (
		name = "abc"
		icon = "def"
		url  = "ghi"
	)

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
			Icon: icon,
			URL:  url,
		},
	}

	l := mock.
		NewLocalizer().
		On("a", name).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithAuthorWithURLt("a", icon, url).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithField(t *testing.T) {
	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: false,
	}

	expect := discord.Embed{
		Fields: []discord.EmbedField{field},
	}

	l := mock.
		NewLocalizer().
		On("a", field.Name).
		On("b", field.Value).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithField(localization.QuickConfig("a"), localization.QuickConfig("b")).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithFieldt(t *testing.T) {
	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: false,
	}

	expect := discord.Embed{
		Fields: []discord.EmbedField{field},
	}

	l := mock.
		NewLocalizer().
		On("a", field.Name).
		On("b", field.Value).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithFieldt("a", "b").
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithInlinedField(t *testing.T) {
	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: true,
	}

	expect := discord.Embed{
		Fields: []discord.EmbedField{field},
	}

	l := mock.
		NewLocalizer().
		On("a", field.Name).
		On("b", field.Value).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithInlinedField(localization.QuickConfig("a"), localization.QuickConfig("b")).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestLocalizedEmbedBuilder_WithInlinedFieldt(t *testing.T) {
	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: true,
	}

	expect := discord.Embed{
		Fields: []discord.EmbedField{field},
	}

	l := mock.
		NewLocalizer().
		On("a", field.Name).
		On("b", field.Value).
		Build()

	actual, err := NewLocalizedEmbedBuilder().
		WithInlinedFieldt("a", "b").
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}
