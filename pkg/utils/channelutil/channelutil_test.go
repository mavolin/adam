package channelutil

import (
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
)

func TestResolvePositions(t *testing.T) {
	channels := []discord.Channel{
		{ID: 11, Position: 0, Type: discord.GuildCategory},
		{ID: 10, Position: 0, Type: discord.GuildCategory},
		{ID: 20, Position: 0, Type: discord.GuildText, CategoryID: 11},
		{ID: 30, Position: 1, Type: discord.GuildVoice, CategoryID: 10},
		{ID: 21, Position: 3, Type: discord.GuildText, CategoryID: 12},
		{ID: 31, Position: 0, Type: discord.GuildVoice, CategoryID: 11},
		{ID: 22, Position: 1, Type: discord.GuildText, CategoryID: 10},
		{ID: 23, Position: 2, Type: discord.GuildText},
		{ID: 12, Position: 1, Type: discord.GuildCategory},
		{ID: 26, Position: 4, Type: discord.GuildText, CategoryID: 12},
		{ID: 25, Position: 3, Type: discord.GuildText, CategoryID: 12},
		{ID: 13, Position: 2, Type: discord.GuildCategory},
		{ID: 24, Position: 4, Type: discord.GuildText},
		{ID: 32, Position: 2, Type: discord.GuildVoice},
	}

	expect := []discord.Channel{
		{ID: 23, Position: 2, Type: discord.GuildText},
		{ID: 24, Position: 4, Type: discord.GuildText},
		{ID: 32, Position: 2, Type: discord.GuildVoice},
		{ID: 10, Position: 0, Type: discord.GuildCategory},
		{ID: 22, Position: 1, Type: discord.GuildText, CategoryID: 10},
		{ID: 30, Position: 1, Type: discord.GuildVoice, CategoryID: 10},
		{ID: 11, Position: 0, Type: discord.GuildCategory},
		{ID: 20, Position: 0, Type: discord.GuildText, CategoryID: 11},
		{ID: 31, Position: 0, Type: discord.GuildVoice, CategoryID: 11},
		{ID: 12, Position: 1, Type: discord.GuildCategory},
		{ID: 21, Position: 3, Type: discord.GuildText, CategoryID: 12},
		{ID: 25, Position: 3, Type: discord.GuildText, CategoryID: 12},
		{ID: 26, Position: 4, Type: discord.GuildText, CategoryID: 12},
		{ID: 13, Position: 2, Type: discord.GuildCategory},
	}

	actual := ResolvePositions(channels)
	assert.Equal(t, expect, actual)
}

func TestResolveCategories(t *testing.T) {
	t.Run("with category-less channels", func(t *testing.T) {
		channels := []discord.Channel{
			{ID: 11, Position: 0, Type: discord.GuildCategory},
			{ID: 10, Position: 0, Type: discord.GuildCategory},
			{ID: 20, Position: 0, Type: discord.GuildText, CategoryID: 11},
			{ID: 30, Position: 1, Type: discord.GuildVoice, CategoryID: 10},
			{ID: 21, Position: 3, Type: discord.GuildText, CategoryID: 12},
			{ID: 31, Position: 0, Type: discord.GuildVoice, CategoryID: 11},
			{ID: 22, Position: 1, Type: discord.GuildText, CategoryID: 10},
			{ID: 23, Position: 2, Type: discord.GuildText},
			{ID: 12, Position: 1, Type: discord.GuildCategory},
			{ID: 26, Position: 4, Type: discord.GuildText, CategoryID: 12},
			{ID: 25, Position: 3, Type: discord.GuildText, CategoryID: 12},
			{ID: 13, Position: 2, Type: discord.GuildCategory},
			{ID: 24, Position: 4, Type: discord.GuildText},
			{ID: 32, Position: 2, Type: discord.GuildVoice},
		}

		expect := [][]discord.Channel{
			{
				{ID: 23, Position: 2, Type: discord.GuildText},
				{ID: 24, Position: 4, Type: discord.GuildText},
				{ID: 32, Position: 2, Type: discord.GuildVoice},
			},
			{
				{ID: 10, Position: 0, Type: discord.GuildCategory},
				{ID: 22, Position: 1, Type: discord.GuildText, CategoryID: 10},
				{ID: 30, Position: 1, Type: discord.GuildVoice, CategoryID: 10},
			},
			{
				{ID: 11, Position: 0, Type: discord.GuildCategory},
				{ID: 20, Position: 0, Type: discord.GuildText, CategoryID: 11},
				{ID: 31, Position: 0, Type: discord.GuildVoice, CategoryID: 11},
			},
			{
				{ID: 12, Position: 1, Type: discord.GuildCategory},
				{ID: 21, Position: 3, Type: discord.GuildText, CategoryID: 12},
				{ID: 25, Position: 3, Type: discord.GuildText, CategoryID: 12},
				{ID: 26, Position: 4, Type: discord.GuildText, CategoryID: 12},
			},
			{{ID: 13, Position: 2, Type: discord.GuildCategory}},
		}

		actual := ResolveCategories(channels)
		assert.Equal(t, expect, actual)
	})

	t.Run("only categories", func(t *testing.T) {
		channels := []discord.Channel{
			{ID: 11, Position: 0, Type: discord.GuildCategory},
			{ID: 10, Position: 0, Type: discord.GuildCategory},
			{ID: 20, Position: 0, Type: discord.GuildText, CategoryID: 11},
			{ID: 30, Position: 1, Type: discord.GuildVoice, CategoryID: 10},
			{ID: 21, Position: 3, Type: discord.GuildText, CategoryID: 12},
			{ID: 31, Position: 0, Type: discord.GuildVoice, CategoryID: 11},
			{ID: 22, Position: 1, Type: discord.GuildText, CategoryID: 10},
			{ID: 12, Position: 1, Type: discord.GuildCategory},
			{ID: 26, Position: 4, Type: discord.GuildText, CategoryID: 12},
			{ID: 25, Position: 3, Type: discord.GuildText, CategoryID: 12},
			{ID: 13, Position: 2, Type: discord.GuildCategory},
		}

		expect := [][]discord.Channel{
			nil,
			{
				{ID: 10, Position: 0, Type: discord.GuildCategory},
				{ID: 22, Position: 1, Type: discord.GuildText, CategoryID: 10},
				{ID: 30, Position: 1, Type: discord.GuildVoice, CategoryID: 10},
			},
			{
				{ID: 11, Position: 0, Type: discord.GuildCategory},
				{ID: 20, Position: 0, Type: discord.GuildText, CategoryID: 11},
				{ID: 31, Position: 0, Type: discord.GuildVoice, CategoryID: 11},
			},
			{
				{ID: 12, Position: 1, Type: discord.GuildCategory},
				{ID: 21, Position: 3, Type: discord.GuildText, CategoryID: 12},
				{ID: 25, Position: 3, Type: discord.GuildText, CategoryID: 12},
				{ID: 26, Position: 4, Type: discord.GuildText, CategoryID: 12},
			},
			{{ID: 13, Position: 2, Type: discord.GuildCategory}},
		}

		actual := ResolveCategories(channels)
		assert.Equal(t, expect, actual)
	})
}
