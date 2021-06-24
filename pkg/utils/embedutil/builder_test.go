package embedutil

import (
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
)

func TestEmbedBuilder_WithSimpleTitle(t *testing.T) {
	title := "abc"

	expect := discord.Embed{Title: title}

	actual, err := NewBuilder().
		WithSimpleTitle(title).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleTitlel(t *testing.T) {
	title := "abc"

	expect := discord.Embed{Title: title}

	l := newMockedLocalizer(t).
		on("a", title).
		build()

	actual, err := NewBuilder().
		WithSimpleTitlel(i18n.NewTermConfig("a")).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleTitlelt(t *testing.T) {
	title := "abc"

	expect := discord.Embed{Title: title}

	l := newMockedLocalizer(t).
		on("a", title).
		build()

	actual, err := NewBuilder().
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

	actual, err := NewBuilder().
		WithTitle(title, url).
		Build(newMockedLocalizer(t).build())

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

	l := newMockedLocalizer(t).
		on("a", title).
		build()

	actual, err := NewBuilder().
		WithTitlel(i18n.NewTermConfig("a"), url).
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

	l := newMockedLocalizer(t).
		on("a", title).
		build()

	actual, err := NewBuilder().
		WithTitlelt("a", url).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithDescription(t *testing.T) {
	description := "abc"

	expect := discord.Embed{Description: description}

	actual, err := NewBuilder().
		WithDescription(description).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithDescriptionl(t *testing.T) {
	description := "abc"

	expect := discord.Embed{Description: description}

	l := newMockedLocalizer(t).
		on("a", description).
		build()

	actual, err := NewBuilder().
		WithDescriptionl(i18n.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithDescriptionlt(t *testing.T) {
	description := "abc"

	expect := discord.Embed{Description: description}

	l := newMockedLocalizer(t).
		on("a", description).
		build()

	actual, err := NewBuilder().
		WithDescriptionlt("a").
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTimestamp(t *testing.T) {
	timestamp := discord.NowTimestamp()

	expect := discord.Embed{Timestamp: timestamp}

	actual, err := NewBuilder().
		WithTimestamp(timestamp).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithColor(t *testing.T) {
	var color discord.Color = 123

	expect := discord.Embed{Color: color}

	actual, err := NewBuilder().
		WithColor(color).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleFooter(t *testing.T) {
	text := "abc"

	expect := discord.Embed{
		Footer: &discord.EmbedFooter{Text: text},
	}

	actual, err := NewBuilder().
		WithSimpleFooter(text).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleFooterl(t *testing.T) {
	text := "abc"

	expect := discord.Embed{
		Footer: &discord.EmbedFooter{Text: text},
	}

	l := newMockedLocalizer(t).
		on("a", text).
		build()

	actual, err := NewBuilder().
		WithSimpleFooterl(i18n.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleFooterlt(t *testing.T) {
	text := "abc"

	expect := discord.Embed{
		Footer: &discord.EmbedFooter{Text: text},
	}

	l := newMockedLocalizer(t).
		on("a", text).
		build()

	actual, err := NewBuilder().
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

	actual, err := NewBuilder().
		WithFooter(text, icon).
		Build(newMockedLocalizer(t).build())

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

	l := newMockedLocalizer(t).
		on("a", text).
		build()

	actual, err := NewBuilder().
		WithFooterl(i18n.NewTermConfig("a"), icon).
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

	l := newMockedLocalizer(t).
		on("a", text).
		build()

	actual, err := NewBuilder().
		WithFooterlt("a", icon).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithImage(t *testing.T) {
	image := "abc"

	expect := discord.Embed{
		Image: &discord.EmbedImage{URL: image},
	}

	actual, err := NewBuilder().
		WithImage(image).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithThumbnail(t *testing.T) {
	thumbnail := "abc"

	expect := discord.Embed{
		Thumbnail: &discord.EmbedThumbnail{URL: thumbnail},
	}

	actual, err := NewBuilder().
		WithThumbnail(thumbnail).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleAuthor(t *testing.T) {
	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{Name: name},
	}

	actual, err := NewBuilder().
		WithSimpleAuthor(name).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleAuthorl(t *testing.T) {
	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{Name: name},
	}

	l := newMockedLocalizer(t).
		on("a", name).
		build()

	actual, err := NewBuilder().
		WithSimpleAuthorl(i18n.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleAuthorlt(t *testing.T) {
	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{Name: name},
	}

	l := newMockedLocalizer(t).
		on("a", name).
		build()

	actual, err := NewBuilder().
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

	actual, err := NewBuilder().
		WithSimpleAuthorWithURL(name, url).
		Build(newMockedLocalizer(t).build())

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

	l := newMockedLocalizer(t).
		on("a", name).
		build()

	actual, err := NewBuilder().
		WithSimpleAuthorWithURLl(i18n.NewTermConfig("a"), url).
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

	l := newMockedLocalizer(t).
		on("a", name).
		build()

	actual, err := NewBuilder().
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

	actual, err := NewBuilder().
		WithAuthor(name, icon).
		Build(newMockedLocalizer(t).build())

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

	l := newMockedLocalizer(t).
		on("a", name).
		build()

	actual, err := NewBuilder().
		WithAuthorl(i18n.NewTermConfig("a"), icon).
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

	l := newMockedLocalizer(t).
		on("a", name).
		build()

	actual, err := NewBuilder().
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

	actual, err := NewBuilder().
		WithAuthorWithURL(name, icon, url).
		Build(newMockedLocalizer(t).build())

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

	l := newMockedLocalizer(t).
		on("a", name).
		build()

	actual, err := NewBuilder().
		WithAuthorWithURLl(i18n.NewTermConfig("a"), icon, url).
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

	l := newMockedLocalizer(t).
		on("a", name).
		build()

	actual, err := NewBuilder().
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

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	actual, err := NewBuilder().
		WithField(field.Name, field.Value).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFieldl(t *testing.T) {
	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: false,
	}

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	l := newMockedLocalizer(t).
		on("a", field.Name).
		on("b", field.Value).
		build()

	actual, err := NewBuilder().
		WithFieldl(i18n.NewTermConfig("a"), i18n.NewTermConfig("b")).
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

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	l := newMockedLocalizer(t).
		on("a", field.Name).
		on("b", field.Value).
		build()

	actual, err := NewBuilder().
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

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	actual, err := NewBuilder().
		WithInlinedField(field.Name, field.Value).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithInlinedFieldl(t *testing.T) {
	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: true,
	}

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	l := newMockedLocalizer(t).
		on("a", field.Name).
		on("b", field.Value).
		build()

	actual, err := NewBuilder().
		WithInlinedFieldl(i18n.NewTermConfig("a"), i18n.NewTermConfig("b")).
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

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	l := newMockedLocalizer(t).
		on("a", field.Name).
		on("b", field.Value).
		build()

	actual, err := NewBuilder().
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

		expect := &Builder{
			fields: []field{
				{
					name:    i18n.NewStaticConfig(name),
					value:   i18n.NewStaticConfig(value),
					inlined: inlined,
				},
			},
		}

		actual := NewBuilder()
		actual.withField(name, value, inlined)

		assert.Equal(t, expect, actual)
	})

	t.Run("name filled", func(t *testing.T) {
		var (
			name    = "abc"
			inlined = false
		)

		expect := &Builder{
			fields: []field{
				{
					name:    i18n.NewStaticConfig(name),
					inlined: inlined,
				},
			},
		}

		actual := NewBuilder()
		actual.withField(name, "", inlined)

		assert.Equal(t, expect, actual)
	})

	t.Run("value filled", func(t *testing.T) {
		var (
			value   = "def"
			inlined = true
		)

		expect := &Builder{
			fields: []field{
				{
					value:   i18n.NewStaticConfig(value),
					inlined: inlined,
				},
			},
		}

		actual := NewBuilder()
		actual.withField("", value, inlined)

		assert.Equal(t, expect, actual)
	})
}

func TestEmbedBuilder_withFieldl(t *testing.T) {
	t.Run("all filled", func(t *testing.T) {
		var (
			name    = i18n.NewTermConfig("abc")
			value   = i18n.NewTermConfig("def")
			inlined = true
		)

		expect := &Builder{
			fields: []field{
				{
					name:    (*i18n.Config)(name),
					value:   (*i18n.Config)(value),
					inlined: inlined,
				},
			},
		}

		actual := NewBuilder()
		actual.withFieldl(name, value, inlined)

		assert.Equal(t, expect, actual)
	})

	t.Run("name filled", func(t *testing.T) {
		var (
			name    = i18n.NewTermConfig("abc")
			inlined = false
		)

		expect := &Builder{
			fields: []field{
				{
					name:    (*i18n.Config)(name),
					inlined: inlined,
				},
			},
		}

		actual := NewBuilder()
		actual.withFieldl(name, nil, inlined)

		assert.Equal(t, expect, actual)
	})

	t.Run("value filled", func(t *testing.T) {
		var (
			value   = i18n.NewTermConfig("def")
			inlined = true
		)

		expect := &Builder{
			fields: []field{
				{
					value:   (*i18n.Config)(value),
					inlined: inlined,
				},
			},
		}

		actual := NewBuilder()
		actual.withFieldl(nil, value, inlined)

		assert.Equal(t, expect, actual)
	})
}

func TestBuilder_Clone(t *testing.T) {
	expectA := NewBuilder().
		WithTitle("abc", "def").
		WithDescription("ghi").
		WithTimestamp(discord.NewTimestamp(time.Unix(0, 0))).
		WithColor(123).
		WithFooter("jkl", "mno").
		WithImage("pqr").
		WithThumbnail("stu").
		WithAuthorWithURL("vwx", "yza", "bcd").
		WithField("efg", "hij")

	a := NewBuilder().
		WithTitle("abc", "def").
		WithDescription("ghi").
		WithTimestamp(discord.NewTimestamp(time.Unix(0, 0))).
		WithColor(123).
		WithFooter("jkl", "mno").
		WithImage("pqr").
		WithThumbnail("stu").
		WithAuthorWithURL("vwx", "yza", "bcd").
		WithField("efg", "hij")

	b := a.Clone()

	assert.Equal(t, a, b)

	b.
		WithTitle("cba", "fed").
		WithDescription("ihg").
		WithTimestamp(discord.NowTimestamp()).
		WithColor(123).
		WithFooter("lkj", "onm").
		WithImage("rqp").
		WithThumbnail("uts").
		WithAuthorWithURL("xwv", "azy", "dcb").
		WithField("gfe", "jih")

	assert.NotEqual(t, a, b)
	assert.Equal(t, expectA, a)
}
