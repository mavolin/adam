// Package channelutil provides utilities for interacting with channels.
package channelutil

import (
	"sort"

	"github.com/diamondburned/arikawa/discord"
)

type channel struct {
	parent   discord.Channel
	children []discord.Channel
}

// ResolvePositions resolves the position of the channels, as displayed in the
// client.
func ResolvePositions(c []discord.Channel) []discord.Channel {
	if len(c) == 0 {
		return nil
	}

	tree := make(map[discord.ChannelID]*channel)

	for _, c := range c {
		switch {
		case c.Type == discord.GuildCategory:
			if v, ok := tree[c.ID]; ok {
				v.parent = c
			} else {
				tree[c.ID] = &channel{parent: c}
			}
		case c.CategoryID.IsValid():
			if v, ok := tree[c.CategoryID]; ok {
				i := sort.Search(len(v.children), func(i int) bool {
					target := v.children[i]

					// text based channels should appear before voice channels
					if (c.Type == discord.GuildText || c.Type == discord.GuildNews) &&
						target.Type == discord.GuildVoice {
						return true
					} else if c.Type == discord.GuildVoice &&
						(target.Type == discord.GuildText || target.Type == discord.GuildNews) {
						return false
					}

					// sort by channel id on equal position
					return target.Position > c.Position || (target.Position == c.Position && target.ID > c.ID)
				})

				if i == len(v.children) {
					v.children = append(v.children, c)
				} else {
					v.children = append(v.children, discord.Channel{}) // make space
					copy(v.children[i+1:], v.children[i:])

					v.children[i] = c
				}
			} else {
				tree[c.CategoryID] = &channel{children: []discord.Channel{c}}
			}
		default:
			tree[c.ID] = &channel{parent: c}
		}
	}

	sorted := sortTree(tree)

	resolved := make([]discord.Channel, 0, len(c))

	for _, c := range sorted {
		resolved = append(resolved, c.parent)
		resolved = append(resolved, c.children...)
	}

	return resolved
}

// ResolveCategories extracts the individual categories and returns them.
// ResolveCategories(c)[0] will always contain all those channels, that are not
// in a category.
// For all remaining values the following applies:
//
// ResolveCategories(c)[n][0] will be the n-th category channel.
// ResolveCategories(c)[n][1:] will contain the channels in n-th category.
func ResolveCategories(c []discord.Channel) [][]discord.Channel {
	if len(c) == 0 {
		return nil
	}

	tree := make(map[discord.ChannelID]*channel)

	var nonCatChannels int // channels without a category

	for _, c := range c {
		switch {
		case c.Type == discord.GuildCategory:
			if v, ok := tree[c.ID]; ok {
				v.parent = c
			} else {
				tree[c.ID] = &channel{parent: c}
			}
		case c.CategoryID.IsValid():
			if v, ok := tree[c.CategoryID]; ok {
				i := sort.Search(len(v.children), func(i int) bool {
					target := v.children[i]

					// text based channels should appear before voice channels
					if (c.Type == discord.GuildText || c.Type == discord.GuildNews) &&
						target.Type == discord.GuildVoice {
						return true
					} else if c.Type == discord.GuildVoice &&
						(target.Type == discord.GuildText || target.Type == discord.GuildNews) {
						return false
					}

					// sort by channel id on equal position
					return target.Position > c.Position || (target.Position == c.Position && target.ID > c.ID)
				})

				if i == len(v.children) {
					v.children = append(v.children, c)
				} else {
					v.children = append(v.children, discord.Channel{}) // make space
					copy(v.children[i+1:], v.children[i:])

					v.children[i] = c
				}
			} else {
				tree[c.CategoryID] = &channel{children: []discord.Channel{c}}
			}
		default:
			tree[c.ID] = &channel{parent: c}
			nonCatChannels++
		}
	}

	sorted := sortTree(tree)

	resolved := make([][]discord.Channel, 0, len(sorted))

	if nonCatChannels == 0 {
		resolved = append(resolved, nil)
	} else {
		resolved = append(resolved, make([]discord.Channel, nonCatChannels))
	}

	for i, c := range sorted[:nonCatChannels] {
		resolved[0][i] = c.parent
	}

	for _, c := range sorted[nonCatChannels:] {
		cat := append([]discord.Channel{c.parent}, c.children...)
		resolved = append(resolved, cat)
	}

	return resolved
}

func sortTree(tree map[discord.ChannelID]*channel) []*channel {
	sorted := make([]*channel, 0, len(tree))

	for _, c := range tree {
		i := sort.Search(len(sorted), func(i int) bool {
			target := sorted[i]

			// sort in the following order: text channels, voice channels, categories
			switch {
			case (c.parent.Type == discord.GuildText || c.parent.Type == discord.GuildNews) &&
				(target.parent.Type != discord.GuildText && target.parent.Type != discord.GuildNews):
				return true
			case c.parent.Type == discord.GuildVoice:
				if target.parent.Type == discord.GuildText || target.parent.Type == discord.GuildNews {
					return false
				} else if target.parent.Type == discord.GuildCategory {
					return true
				}
			case c.parent.Type == discord.GuildCategory && target.parent.Type != discord.GuildCategory:
				return false
			}

			// sort by channel id on equal position
			return target.parent.Position > c.parent.Position ||
				(target.parent.Position == c.parent.Position && target.parent.ID > c.parent.ID)
		})

		if i >= len(sorted) {
			sorted = append(sorted, c)
		} else {
			sorted = append(sorted, new(channel)) // make space
			copy(sorted[i+1:], sorted[i:])

			sorted[i] = c
		}
	}
	return sorted
}
