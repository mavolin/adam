// Package emoji provides utilities for interacting with unicode emojis.
//
// All emojis found in this package are also available as default emojis on
// discord.
package emoji

import "strings"

// discord uses Unicode Emoji version 13.0
//go:generate go run ../../../tools/codegen/emoji/fetch.go ../../../tools/codegen/emoji/stringutil.go ../../../tools/codegen/emoji/main.go 13.1

type (
	// Emoji is an emoji.
	Emoji = string

	// SkinToneEmoji is an emoji that has different skin tones.
	SkinToneEmoji struct {
		// NeutralSkin is the emoji with neutral skin color.
		NeutralSkin Emoji
		// LightSkin is the emoji with light skin color.
		LightSkin Emoji
		// MediumLightSkin is the emoji with medium light skin color.
		MediumLightSkin Emoji
		// MediumSkin is the emoji with Medium skin color.
		MediumSkin Emoji
		// MediumDarkSkin is the emoji with medium dark skin color.
		MediumDarkSkin Emoji
		// DarkSkin is the emoji with dark skin color.
		DarkSkin Emoji
	}
)

// IsValid checks if the passed emoji is a valid emoji as used by discord.
func IsValid(emoji string) bool {
	_, ok := emojis[emoji]
	return ok
}

// CountryFlag returns the Emoji for the country with the passed ISO 3166-1
// Alpha 2 country code.
// If the code is invalid or there is no flag for the passed code, CountryFlag
// returns an empty string.
func CountryFlag(code string) Emoji {
	if len(code) != 2 {
		return ""
	}

	code = strings.ToLower(code)
	flag := countryCodeLetter(code[0]) + countryCodeLetter(code[1])
	if !IsValid(flag) {
		return ""
	}

	return flag
}

// countryCodeLetter shifts given letter byte as flagBaseIndex.
func countryCodeLetter(l byte) string {
	const flagBaseIndex = '\U0001F1E6' - 'a'

	return string(rune(l) + flagBaseIndex)
}
