package main

// emojiVersion is the Unicode Emoji version Discord uses.
const emojiVersion = 12.0

func main() {
	emojis, err := fetchEmojis()
	if err != nil {
		panic(err)
	}

	emojis = filterVersion(emojis, emojiVersion)

	err = generateConstants(emojis)
	if err != nil {
		panic(err)
	}

	err = generateSet(emojis)
	if err != nil {
		panic(err)
	}
}
