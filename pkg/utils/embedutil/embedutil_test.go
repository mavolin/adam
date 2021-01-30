package embedutil

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/stretchr/testify/assert"
)

func TestCountChars(t *testing.T) {
	testCases := []struct {
		name string
		e    discord.Embed

		expect int
	}{
		{
			name:   "title",
			e:      discord.Embed{Title: "lorem"},
			expect: 5,
		},
		{
			name:   "description",
			e:      discord.Embed{Description: "ipsum"},
			expect: 5,
		},
		{
			name:   "footer",
			e:      discord.Embed{Footer: &discord.EmbedFooter{Text: "dolor"}},
			expect: 5,
		},
		{
			name:   "author",
			e:      discord.Embed{Author: &discord.EmbedAuthor{Name: "sit"}},
			expect: 3,
		},
		{
			name: "fields",
			e: discord.Embed{
				Fields: []discord.EmbedField{
					{Name: "amet", Value: "consectetur"},
					{Name: "adipisci", Value: "elit"},
				},
			},
			expect: 27,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			actual := CountChars(c.e)
			assert.Equal(t, c.expect, actual)
		})
	}
}
