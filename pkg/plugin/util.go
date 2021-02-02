package plugin

import "github.com/diamondburned/arikawa/v2/discord"

func embedEmpty(e discord.Embed) bool {
	return !(len(e.Title) > 0 || len(e.Description) > 0 || len(e.URL) > 0 ||
		e.Timestamp.IsValid() || e.Color > 0 || e.Footer != nil ||
		e.Image != nil || e.Thumbnail != nil || e.Video != nil ||
		e.Provider != nil || e.Author != nil || len(e.Fields) > 0)
}
