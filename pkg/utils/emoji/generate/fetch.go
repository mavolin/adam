package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type gemoji struct {
	Emoji          string  `json:"emoji"`
	Description    string  `json:"description"`
	Category       string  `json:"category"`
	UnicodeVersion float32 `json:"unicode_version"`
	SkinTones      bool    `json:"skin_tones,omitempty"`
}

func (g *gemoji) UnmarshalJSON(bytes []byte) error {
	var emoji struct {
		Emoji          string `json:"emoji"`
		Description    string `json:"description"`
		Category       string `json:"category"`
		UnicodeVersion string `json:"unicode_version"`
		SkinTones      bool   `json:"skin_tones,omitempty"`
	}

	err := json.Unmarshal(bytes, &emoji)
	if err != nil {
		return err
	}

	g.Emoji = emoji.Emoji
	g.Description = emoji.Description
	g.Category = emoji.Category
	g.SkinTones = emoji.SkinTones

	if len(emoji.UnicodeVersion) > 0 {
		ver, err := strconv.ParseFloat(emoji.UnicodeVersion, 32)
		if err != nil {
			return err
		}

		g.UnicodeVersion = float32(ver)
	}
	// some older emojis don't have a version, leaving them at ver 0 is just
	// fine

	return nil
}

func fetchEmojis() ([]gemoji, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json")
	if err != nil {
		return nil, err
	}

	var emojis []gemoji

	err = json.NewDecoder(resp.Body).Decode(&emojis)
	if err != nil {
		return nil, err
	}

	return emojis, resp.Body.Close()
}

func filterVersion(g []gemoji, maxVersion float32) []gemoji {
	for i := 0; i < len(g); i++ {
		e := g[i]

		if e.UnicodeVersion > maxVersion {
			g = append(g[:i], g[i+1:]...)
			i--
		}
	}

	return g
}
