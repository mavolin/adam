package embedutil

import (
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mavolin/adam/pkg/i18n"
)

func TestEmbedBuilder_WithTitle(t *testing.T) {
	title := "abc"

	expect := discord.Embed{Title: title}

	actual, err := NewBuilder().
		WithTitle(title).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTitlelt(t *testing.T) {
	title := "abc"

	expect := discord.Embed{Title: title}

	l := newMockedLocalizer(t).
		on("a", title).
		build()

	actual, err := NewBuilder().
		WithTitlelt("a").
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTitlel(t *testing.T) {
	title := "abc"

	expect := discord.Embed{Title: title}

	l := newMockedLocalizer(t).
		on("a", title).
		build()

	actual, err := NewBuilder().
		WithTitlel(i18n.NewTermConfig("a")).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTitleURL(t *testing.T) {
	url := "def"

	expect := discord.Embed{URL: url}

	actual, err := NewBuilder().
		WithTitleURL(url).
		Build(newMockedLocalizer(t).build())

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

func TestEmbedBuilder_WithFooter(t *testing.T) {
	text := "abc"

	expect := discord.Embed{Footer: &discord.EmbedFooter{Text: text}}

	actual, err := NewBuilder().
		WithFooter(text).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooterlt(t *testing.T) {
	text := "abc"

	expect := discord.Embed{Footer: &discord.EmbedFooter{Text: text}}

	l := newMockedLocalizer(t).
		on("a", text).
		build()

	actual, err := NewBuilder().
		WithFooterlt("a").
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooterl(t *testing.T) {
	text := "abc"

	expect := discord.Embed{Footer: &discord.EmbedFooter{Text: text}}

	l := newMockedLocalizer(t).
		on("a", text).
		build()

	actual, err := NewBuilder().
		WithFooterl(i18n.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooterIcon(t *testing.T) {
	icon := "def"

	expect := discord.Embed{Footer: &discord.EmbedFooter{Icon: icon}}

	actual, err := NewBuilder().
		WithFooterIcon(icon).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
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

func TestEmbedBuilder_WithAuthor(t *testing.T) {
	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{Name: name},
	}

	actual, err := NewBuilder().
		WithAuthor(name).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
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
		WithAuthorlt("a").
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthorl(t *testing.T) {
	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{Name: name},
	}

	l := newMockedLocalizer(t).
		on("a", name).
		build()

	actual, err := NewBuilder().
		WithAuthorl(i18n.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthorURL(t *testing.T) {
	url := "def"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			URL: url,
		},
	}

	actual, err := NewBuilder().
		WithAuthorURL(url).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthorIcon(t *testing.T) {
	icon := "def"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			Icon: icon,
		},
	}

	actual, err := NewBuilder().
		WithAuthorIcon(icon).
		Build(newMockedLocalizer(t).build())

	assert.NoError(t, err)
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

func TestBuilder_Clone(t *testing.T) {
	expectA := NewBuilder().
		WithTitle("abc").
		WithTitleURL("def").
		WithDescription("ghi").
		WithTimestamp(discord.NewTimestamp(time.Unix(0, 0))).
		WithColor(123).
		WithFooter("jkl").
		WithFooterIcon("mno").
		WithImage("pqr").
		WithThumbnail("stu").
		WithAuthor("vwx").
		WithAuthorIcon("yza").
		WithAuthorURL("bcd").
		WithField("efg", "hij")

	a := NewBuilder().
		WithTitle("abc").
		WithTitleURL("def").
		WithDescription("ghi").
		WithTimestamp(discord.NewTimestamp(time.Unix(0, 0))).
		WithColor(123).
		WithFooter("jkl").
		WithFooterIcon("mno").
		WithImage("pqr").
		WithThumbnail("stu").
		WithAuthor("vwx").
		WithAuthorIcon("yza").
		WithAuthorURL("bcd").
		WithField("efg", "hij")

	b := a.Clone()

	assert.Equal(t, a, b)

	b.
		WithTitle("cba").
		WithTitleURL("fed").
		WithDescription("ihg").
		WithTimestamp(discord.NowTimestamp()).
		WithColor(123).
		WithFooter("lkj").
		WithFooterIcon("onm").
		WithImage("rqp").
		WithThumbnail("uts").
		WithAuthor("xwv").
		WithAuthorIcon("azy").
		WithAuthorURL("dcb").
		WithField("gfe", "jih")

	assert.NotEqual(t, a, b)
	assert.Equal(t, expectA, a)
}
