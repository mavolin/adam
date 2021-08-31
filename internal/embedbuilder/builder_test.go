package embedbuilder

import (
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mocki18n "github.com/mavolin/adam/internal/mock/i18n"
	"github.com/mavolin/adam/pkg/i18n"
)

func TestEmbedBuilder_WithTitle(t *testing.T) {
	t.Parallel()

	title := "abc"

	expect := discord.Embed{Title: title}

	actual, err := New().
		WithTitle(title).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTitlelt(t *testing.T) {
	t.Parallel()

	title := "abc"

	expect := discord.Embed{Title: title}

	l := mocki18n.NewLocalizer(t).
		On("a", title).
		Build()

	actual, err := New().
		WithTitlelt("a").
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTitlel(t *testing.T) {
	t.Parallel()

	title := "abc"

	expect := discord.Embed{Title: title}

	l := mocki18n.NewLocalizer(t).
		On("a", title).
		Build()

	actual, err := New().
		WithTitlel(i18n.NewTermConfig("a")).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTitleURL(t *testing.T) {
	t.Parallel()

	url := "def"

	expect := discord.Embed{URL: url}

	actual, err := New().
		WithTitleURL(url).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithDescription(t *testing.T) {
	t.Parallel()

	description := "abc"

	expect := discord.Embed{Description: description}

	actual, err := New().
		WithDescription(description).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithDescriptionl(t *testing.T) {
	t.Parallel()

	description := "abc"

	expect := discord.Embed{Description: description}

	l := mocki18n.NewLocalizer(t).
		On("a", description).
		Build()

	actual, err := New().
		WithDescriptionl(i18n.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithDescriptionlt(t *testing.T) {
	t.Parallel()

	description := "abc"

	expect := discord.Embed{Description: description}

	l := mocki18n.NewLocalizer(t).
		On("a", description).
		Build()

	actual, err := New().
		WithDescriptionlt("a").
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTimestamp(t *testing.T) {
	t.Parallel()

	timestamp := discord.NowTimestamp()

	expect := discord.Embed{Timestamp: timestamp}

	actual, err := New().
		WithTimestamp(timestamp).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithColor(t *testing.T) {
	t.Parallel()

	var color discord.Color = 123

	expect := discord.Embed{Color: color}

	actual, err := New().
		WithColor(color).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooter(t *testing.T) {
	t.Parallel()

	text := "abc"

	expect := discord.Embed{Footer: &discord.EmbedFooter{Text: text}}

	actual, err := New().
		WithFooter(text).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooterlt(t *testing.T) {
	t.Parallel()

	text := "abc"

	expect := discord.Embed{Footer: &discord.EmbedFooter{Text: text}}

	l := mocki18n.NewLocalizer(t).
		On("a", text).
		Build()

	actual, err := New().
		WithFooterlt("a").
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooterl(t *testing.T) {
	t.Parallel()

	text := "abc"

	expect := discord.Embed{Footer: &discord.EmbedFooter{Text: text}}

	l := mocki18n.NewLocalizer(t).
		On("a", text).
		Build()

	actual, err := New().
		WithFooterl(i18n.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooterIcon(t *testing.T) {
	t.Parallel()

	icon := "def"

	expect := discord.Embed{Footer: &discord.EmbedFooter{Icon: icon}}

	actual, err := New().
		WithFooterIcon(icon).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithImage(t *testing.T) {
	t.Parallel()

	image := "abc"

	expect := discord.Embed{
		Image: &discord.EmbedImage{URL: image},
	}

	actual, err := New().
		WithImage(image).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithThumbnail(t *testing.T) {
	t.Parallel()

	thumbnail := "abc"

	expect := discord.Embed{
		Thumbnail: &discord.EmbedThumbnail{URL: thumbnail},
	}

	actual, err := New().
		WithThumbnail(thumbnail).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthor(t *testing.T) {
	t.Parallel()

	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{Name: name},
	}

	actual, err := New().
		WithAuthor(name).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithSimpleAuthorlt(t *testing.T) {
	t.Parallel()

	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{Name: name},
	}

	l := mocki18n.NewLocalizer(t).
		On("a", name).
		Build()

	actual, err := New().
		WithAuthorlt("a").
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthorl(t *testing.T) {
	t.Parallel()

	name := "abc"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{Name: name},
	}

	l := mocki18n.NewLocalizer(t).
		On("a", name).
		Build()

	actual, err := New().
		WithAuthorl(i18n.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthorURL(t *testing.T) {
	t.Parallel()

	url := "def"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{
			URL: url,
		},
	}

	actual, err := New().
		WithAuthorURL(url).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithAuthorIcon(t *testing.T) {
	t.Parallel()

	icon := "def"

	expect := discord.Embed{
		Author: &discord.EmbedAuthor{Icon: icon},
	}

	actual, err := New().
		WithAuthorIcon(icon).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithField(t *testing.T) {
	t.Parallel()

	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: false,
	}

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	actual, err := New().
		WithField(field.Name, field.Value).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFieldlt(t *testing.T) {
	t.Parallel()

	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: false,
	}

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	l := mocki18n.NewLocalizer(t).
		On("a", field.Name).
		On("b", field.Value).
		Build()

	actual, err := New().
		WithFieldlt("a", "b").
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFieldl(t *testing.T) {
	t.Parallel()

	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: false,
	}

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	l := mocki18n.NewLocalizer(t).
		On("a", field.Name).
		On("b", field.Value).
		Build()

	actual, err := New().
		WithFieldl(i18n.NewTermConfig("a"), i18n.NewTermConfig("b")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithInlinedField(t *testing.T) {
	t.Parallel()

	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: true,
	}

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	actual, err := New().
		WithInlinedField(field.Name, field.Value).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithInlinedFieldlt(t *testing.T) {
	t.Parallel()

	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: true,
	}

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	l := mocki18n.NewLocalizer(t).
		On("a", field.Name).
		On("b", field.Value).
		Build()

	actual, err := New().
		WithInlinedFieldlt("a", "b").
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithInlinedFieldl(t *testing.T) {
	t.Parallel()

	field := discord.EmbedField{
		Name:   "abc",
		Value:  "def",
		Inline: true,
	}

	expect := discord.Embed{Fields: []discord.EmbedField{field}}

	l := mocki18n.NewLocalizer(t).
		On("a", field.Name).
		On("b", field.Value).
		Build()

	actual, err := New().
		WithInlinedFieldl(i18n.NewTermConfig("a"), i18n.NewTermConfig("b")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestBuilder_Clone(t *testing.T) {
	t.Parallel()

	expectA := New().
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

	a := New().
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

	b.WithTitle("cba").
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
