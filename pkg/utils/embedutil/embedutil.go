// Package embedutil provides utilities to generate and interact with embeds.
package embedutil

import (
	"github.com/diamondburned/arikawa/v2/discord"
)

// MaxChars is the maximum amount of characters Discord allows an embed to
// hold.
const MaxChars = 6000

// CountChars returns the number of characters in the embed.
func CountChars(e discord.Embed) int {
	sum := len(e.Title) + len(e.Description)

	if e.Footer != nil {
		sum += len(e.Footer.Text)
	}

	if e.Author != nil {
		sum += len(e.Author.Name)
	}

	for _, f := range e.Fields {
		sum += len(f.Name) + len(f.Value)
	}

	return sum
}
