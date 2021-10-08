package msgbuilder

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

	actual, err := NewEmbed().
		WithTitle(title).
		Build(mocki18n.NewLocalizer(t).Build())

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

	actual, err := NewEmbed().
		WithTitlel(i18n.NewTermConfig("a")).
		Build(l)

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTitleURL(t *testing.T) {
	t.Parallel()

	url := "def"

	expect := discord.Embed{URL: url}

	actual, err := NewEmbed().
		WithTitleURL(url).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithDescription(t *testing.T) {
	t.Parallel()

	description := "abc"

	expect := discord.Embed{Description: description}

	actual, err := NewEmbed().
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

	actual, err := NewEmbed().
		WithDescriptionl(i18n.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithTimestamp(t *testing.T) {
	t.Parallel()

	timestamp := discord.NowTimestamp()

	expect := discord.Embed{Timestamp: timestamp}

	actual, err := NewEmbed().
		WithTimestamp(timestamp).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithColor(t *testing.T) {
	t.Parallel()

	var color discord.Color = 123

	expect := discord.Embed{Color: color}

	actual, err := NewEmbed().
		WithColor(color).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooter(t *testing.T) {
	t.Parallel()

	text := "abc"

	expect := discord.Embed{Footer: &discord.EmbedFooter{Text: text}}

	actual, err := NewEmbed().
		WithFooter(text).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooterl(t *testing.T) {
	t.Parallel()

	text := "abc"

	expect := discord.Embed{Footer: &discord.EmbedFooter{Text: text}}

	l := mocki18n.NewLocalizer(t).
		On("a", text).
		Build()

	actual, err := NewEmbed().
		WithFooterl(i18n.NewTermConfig("a")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestEmbedBuilder_WithFooterIcon(t *testing.T) {
	t.Parallel()

	icon := "def"

	expect := discord.Embed{Footer: &discord.EmbedFooter{Icon: icon}}

	actual, err := NewEmbed().
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

	actual, err := NewEmbed().
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

	actual, err := NewEmbed().
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

	actual, err := NewEmbed().
		WithAuthor(name).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
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

	actual, err := NewEmbed().
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

	actual, err := NewEmbed().
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

	actual, err := NewEmbed().
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

	actual, err := NewEmbed().
		WithField(field.Name, field.Value).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
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

	actual, err := NewEmbed().
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

	actual, err := NewEmbed().
		WithInlinedField(field.Name, field.Value).
		Build(mocki18n.NewLocalizer(t).Build())

	assert.NoError(t, err)
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

	actual, err := NewEmbed().
		WithInlinedFieldl(i18n.NewTermConfig("a"), i18n.NewTermConfig("b")).
		Build(l)

	require.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func TestBuilder_Clone(t *testing.T) {
	t.Parallel()

	expectA := NewEmbed().
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

	a := NewEmbed().
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
