package discordutil

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/mock"
)

func TestEmbedBuilder_WithSimpleTitle(t *testing.T) {
	title := "abc"

	expect := discord.Embed{
		Title: title,
	}

	actual, err := NewEmbedBuilder().
		WithSimpleTitle(title).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleTitlel(t *testing.T) {
	title := "abc"

	expect := discord.Embed{
		Title: title,
	}

	l := mock.
		NewLocalizer().
		On("a", title).
		Build()

	actual, err := NewEmbedBuilder().
		WithSimpleTitlel(localization.NewTermConfig("a")).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleTitlelt(t *testing.T) {
	title := "abc"

	expect := discord.Embed{
		Title: title,
	}

	l := mock.
		NewLocalizer().
		On("a", title).
		Build()

	actual, err := NewEmbedBuilder().
		WithSimpleTitlelt("a").
		Build(l)

	assert.NoError(t, err)
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

	actual, err := NewEmbedBuilder().
		WithTitle(title, url).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTitlel(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithTitlel(localization.NewTermConfig("a"), url).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTitlelt(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithTitlelt("a", url).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithDescription(t *testing.T) {
	description := "abc"

	expect := discord.Embed{
		Description: description,
	}

	actual, err := NewEmbedBuilder().
		WithDescription(description).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithDescriptionl(t *testing.T) {
	description := "abc"

	expect := discord.Embed{
		Description: description,
	}

	l := mock.
		NewLocalizer().
		On("a", description).
		Build()

	actual, err := NewEmbedBuilder().
		WithDescriptionl(localization.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithDescriptionlt(t *testing.T) {
	description := "abc"

	expect := discord.Embed{
		Description: description,
	}

	l := mock.
		NewLocalizer().
		On("a", description).
		Build()

	actual, err := NewEmbedBuilder().
		WithDescriptionlt("a").
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTimestamp(t *testing.T) {
	timestamp := discord.NowTimestamp()

	expect := discord.Embed{
		Timestamp: timestamp,
	}

	actual, err := NewEmbedBuilder().
		WithTimestamp(timestamp).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithColor(t *testing.T) {
	var color discord.Color = 123

	expect := discord.Embed{
		Color: color,
	}

	actual, err := NewEmbedBuilder().
		WithColor(color).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleFooter(t *testing.T) {
	text := "abc"

	expect := discord.Embed{
		Footer: &discord.EmbedFooter{
			Text: text,
		},
	}

	actual, err := NewEmbedBuilder().
		WithSimpleFooter(text).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleFooterl(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithSimpleFooterl(localization.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleFooterlt(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithSimpleFooterlt("a").
		Build(l)

	require.NoError(t, err)
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

	actual, err := NewEmbedBuilder().
		WithFooter(text, icon).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooterl(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithFooterl(localization.NewTermConfig("a"), icon).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooterlt(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithFooterlt("a", icon).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithImage(t *testing.T) {
	image := "abc"

	expect := discord.Embed{
		Image: &discord.EmbedImage{
			URL: image,
		},
	}

	actual, err := NewEmbedBuilder().
		WithImage(image).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithThumbnail(t *testing.T) {
	thumbnail := "abc"

	expect := discord.Embed{
		Thumbnail: &discord.EmbedThumbnail{
			URL: thumbnail,
		},
	}

	actual, err := NewEmbedBuilder().
		WithThumbnail(thumbnail).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleAuthor(t *testing.T) {
	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
		},
	}

	actual, err := NewEmbedBuilder().
		WithSimpleAuthor(name).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleAuthorl(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithSimpleAuthorl(localization.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleAuthorlt(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithSimpleAuthorlt("a").
		Build(l)

	require.NoError(t, err)
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

	actual, err := NewEmbedBuilder().
		WithSimpleAuthorWithURL(name, url).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleAuthorWithURLl(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithSimpleAuthorWithURLl(localization.NewTermConfig("a"), url).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleAuthorWithURLlt(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithSimpleAuthorWithURLlt("a", url).
		Build(l)

	require.NoError(t, err)
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

	actual, err := NewEmbedBuilder().
		WithAuthor(name, icon).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthorl(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithAuthorl(localization.NewTermConfig("a"), icon).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthorlt(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithAuthorlt("a", icon).
		Build(l)

	require.NoError(t, err)
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

	actual, err := NewEmbedBuilder().
		WithAuthorWithURL(name, icon, url).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthorWithURLl(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithAuthorWithURLl(localization.NewTermConfig("a"), icon, url).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthorWithURLlt(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithAuthorWithURLlt("a", icon, url).
		Build(l)

	require.NoError(t, err)
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

	actual, err := NewEmbedBuilder().
		WithField(field.Name, field.Value).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFieldl(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithFieldl(localization.NewTermConfig("a"), localization.NewTermConfig("b")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFieldlt(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithFieldlt("a", "b").
		Build(l)

	require.NoError(t, err)
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

	actual, err := NewEmbedBuilder().
		WithInlinedField(field.Name, field.Value).
		Build(nil)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithInlinedFieldl(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithInlinedFieldl(localization.NewTermConfig("a"), localization.NewTermConfig("b")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithInlinedFieldlt(t *testing.T) {
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

	actual, err := NewEmbedBuilder().
		WithInlinedFieldlt("a", "b").
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_withField(t *testing.T) {
	t.Run("all filled", func(t *testing.T) {
		var (
			name    = "abc"
			value   = "def"
			inlined = true
		)

		expect := &EmbedBuilder{
			fields: []embedField{
				{
					name: &text{
						string: name,
					},
					value: &text{
						string: value,
					},
					inlined: inlined,
				},
			},
		}

		actual := NewEmbedBuilder()
		actual.withField(name, value, inlined)

		assert.Equal(t, expect, actual)
	})

	t.Run("name filled", func(t *testing.T) {
		var (
			name    = "abc"
			inlined = false
		)

		expect := &EmbedBuilder{
			fields: []embedField{
				{
					name: &text{
						string: name,
					},
					value:   nil,
					inlined: inlined,
				},
			},
		}

		actual := NewEmbedBuilder()
		actual.withField(name, "", inlined)

		assert.Equal(t, expect, actual)
	})

	t.Run("value filled", func(t *testing.T) {
		var (
			value   = "def"
			inlined = true
		)

		expect := &EmbedBuilder{
			fields: []embedField{
				{
					name: nil,
					value: &text{
						string: value,
					},
					inlined: inlined,
				},
			},
		}

		actual := NewEmbedBuilder()
		actual.withField("", value, inlined)

		assert.Equal(t, expect, actual)
	})
}

func TestEmbedBuilder_withFieldl(t *testing.T) {
	t.Run("all filled", func(t *testing.T) {
		var (
			name    = localization.NewTermConfig("abc")
			value   = localization.NewTermConfig("def")
			inlined = true
		)

		expect := &EmbedBuilder{
			fields: []embedField{
				{
					name: &text{
						config: name,
					},
					value: &text{
						config: value,
					},
					inlined: inlined,
				},
			},
		}

		actual := NewEmbedBuilder()
		actual.withFieldl(name, value, inlined)

		assert.Equal(t, expect, actual)
	})

	t.Run("name filled", func(t *testing.T) {
		var (
			name    = localization.NewTermConfig("abc")
			inlined = false
		)

		expect := &EmbedBuilder{
			fields: []embedField{
				{
					name: &text{
						config: name,
					},
					value:   nil,
					inlined: inlined,
				},
			},
		}

		actual := NewEmbedBuilder()
		actual.withFieldl(name, localization.Config{}, inlined)

		assert.Equal(t, expect, actual)
	})

	t.Run("value filled", func(t *testing.T) {
		var (
			value   = localization.NewTermConfig("def")
			inlined = true
		)

		expect := &EmbedBuilder{
			fields: []embedField{
				{
					name: nil,
					value: &text{
						config: value,
					},
					inlined: inlined,
				},
			},
		}

		actual := NewEmbedBuilder()
		actual.withFieldl(localization.Config{}, value, inlined)

		assert.Equal(t, expect, actual)
	})
}
