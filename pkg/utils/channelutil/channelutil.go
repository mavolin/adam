// Package channelutil provides utilities for interacting with channels.
package channelutil

import (
	"sort"

	"github.com/diamondburned/arikawa/discord"
)

// ResolvePositions resolves the position of the channels, as displayed in the
// client.
func ResolvePositions(c []discord.Channel) []discord.Channel {
	if len(c) == 0 {
		return nil
	}

	resolved, categories := sortChannels(c)

	// resolved has the capacity anyway, so use it
	for _, c := range categories {
		resolved = append(resolved, c.parent)
		resolved = append(resolved, c.children...)
	}

	return resolved
}

// ResolveCategories extracts the individual categories and returns them.
// ResolveCategories(c)[0] will always contain all those channels that are not
// in a category.
// If there are no such channels, ResolveCategories(c)[0] will be nil.
// For all remaining values the following applies:
//
// ResolveCategories(c)[n][0] will be the n-th category category.
// ResolveCategories(c)[n][1:] will contain the channels in n-th category.
func ResolveCategories(c []discord.Channel) [][]discord.Channel {
	if len(c) == 0 {
		return nil
	}

	topLevel, categories := sortChannels(c)

	resolved := make([][]discord.Channel, len(categories)+1)
	if len(topLevel) == 0 {
		resolved[0] = nil
	} else {
		resolved[0] = topLevel
	}

	for i, c := range categories {
		resolved[i+1] = append([]discord.Channel{c.parent}, c.children...)
	}

	return resolved
}

type category struct {
	parent   discord.Channel
	children []discord.Channel
}

// sortChannels sorts the passed channel, the same as the client does.
// It returns two slices, one containing all top-level text and voice channels
// and a second one containing all categories with ordered children.
//
// The first slice will have sufficient capacity to store all elements in c.
func sortChannels(c []discord.Channel) (topLevel []discord.Channel, categories []*category) {
	topLevel = make([]discord.Channel, 0, len(c))
	categoriesMap := make(map[discord.ChannelID]*category)

	for _, c := range c {
		switch {
		case c.Type == discord.GuildCategory:
			if v, ok := categoriesMap[c.ID]; ok {
				v.parent = c
			} else {
				categoriesMap[c.ID] = &category{parent: c}
			}
		case c.CategoryID.IsValid():
			if v, ok := categoriesMap[c.CategoryID]; ok {
				i := sort.Search(len(v.children), channelSearchFunc(c, v.children))
				v.children = insertChannel(c, i, v.children)
			} else {
				categoriesMap[c.CategoryID] = &category{children: []discord.Channel{c}}
			}
		default:
			i := sort.Search(len(topLevel), channelSearchFunc(c, topLevel))
			topLevel = insertChannel(c, i, topLevel)
		}
	}

	return topLevel, sortCategoriesMap(categoriesMap)
}

// insertChannel is a helper function that inserts the passed channel at the
// passed position in the passed slice.
// It returns the modified slice.
func insertChannel(c discord.Channel, i int, channels []discord.Channel) []discord.Channel {
	if i >= len(channels) {
		channels = append(channels, c)
	} else {
		channels = append(channels, discord.Channel{})
		copy(channels[i+1:], channels[i:])

		channels[i] = c
	}

	return channels
}

// sortCategoriesMap sorts the passed categories into a category slice.
func sortCategoriesMap(categories map[discord.ChannelID]*category) []*category {
	sorted := make([]*category, 0, len(categories))

	for _, c := range categories {
		i := sort.Search(len(sorted), func(i int) bool {
			target := sorted[i]

			// sort by category id on equal position
			return target.parent.Position > c.parent.Position ||
				(target.parent.Position == c.parent.Position && target.parent.ID > c.parent.ID)
		})

		if i >= len(sorted) {
			sorted = append(sorted, c)
		} else {
			sorted = append(sorted, new(category)) // make space
			copy(sorted[i+1:], sorted[i:])

			sorted[i] = c
		}
	}

	return sorted
}

// channelSearchFunc returns a search func that determines the position of the
// passed category, as the client would sort it.
// Only works for guild text, guild news and guild voice channels.
func channelSearchFunc(c discord.Channel, targets []discord.Channel) func(int) bool {
	return func(i int) bool {
		target := targets[i]

		// text based channels should appear before voice channels
		if (c.Type == discord.GuildText || c.Type == discord.GuildNews) &&
			target.Type == discord.GuildVoice {
			return true
		} else if c.Type == discord.GuildVoice &&
			(target.Type == discord.GuildText || target.Type == discord.GuildNews) {
			return false
		}

		// sort by category id on equal position
		return target.Position > c.Position || (target.Position == c.Position && target.ID > c.ID)
	}
}
