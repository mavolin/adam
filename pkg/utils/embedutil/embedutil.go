// Package embedutil provides utilities to generate and interact with embeds.
package embedutil

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

// MaxChars is the maximum amount of characters Discord allows an embed to
// hold.
const MaxChars = 6000

// CountChars returns the number of characters in the embed.
func CountChars(e discord.Embed) int {
	sum := len([]rune(e.Title)) + len([]rune(e.Description))

	if e.Footer != nil {
		sum += len([]rune(e.Footer.Text))
	}

	if e.Author != nil {
		sum += len([]rune(e.Author.Name))
	}

	for _, f := range e.Fields {
		sum += len([]rune(f.Name)) + len([]rune(f.Value))
	}

	return sum
}
