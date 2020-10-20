package main

import (
	"os"
	"strconv"
)

func main() {
	emojis, err := fetchEmojis()
	if err != nil {
		panic(err)
	}

	version, err := strconv.ParseFloat(os.Args[1], 32)
	if err != nil {
		panic(err)
	}

	emojis = filterVersion(emojis, float32(version))

	err = generateConstants(emojis)
	if err != nil {
		panic(err)
	}

	err = generateSet(emojis)
	if err != nil {
		panic(err)
	}
}
