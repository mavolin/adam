package main

import (
	"encoding/json"
	"log"
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
	type rawGemoji gemoji

	emoji := struct {
		rawGemoji
		UnicodeVersion string `json:"unicode_version"`
	}{}

	err := json.Unmarshal(bytes, &emoji)
	if err != nil {
		return err
	}

	*g = gemoji(emoji.rawGemoji)

	// some older emojis don't have a version, leaving them at ver 0 is just
	// fine
	if len(emoji.UnicodeVersion) > 0 {
		ver, err := strconv.ParseFloat(emoji.UnicodeVersion, 32)
		if err != nil {
			return err
		}

		g.UnicodeVersion = float32(ver)
	}

	return nil
}

func fetchEmojis() ([]gemoji, error) {
	log.Println("fetching emojis")

	resp, err := http.Get("https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json")
	if err != nil {
		return nil, err
	}

	var emojis []gemoji

	err = json.NewDecoder(resp.Body).Decode(&emojis)
	if err != nil {
		return nil, err
	}

	log.Printf("fetched %d emojis\n", len(emojis))

	return emojis, resp.Body.Close()
}

func filterVersion(g []gemoji, maxVersion float32) []gemoji {
	log.Printf("filtering to only include emojis with version <= %.1f\n", maxVersion)

	for i := 0; i < len(g); i++ {
		e := g[i]

		if e.UnicodeVersion > maxVersion {
			g = append(g[:i], g[i+1:]...)
			i--
		}
	}

	log.Printf("%d emojis remain\n", len(g))

	return g
}
