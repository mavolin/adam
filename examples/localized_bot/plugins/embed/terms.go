package embed

import "github.com/mavolin/adam/pkg/i18n"

var (
	shortDescription = i18n.NewFallbackConfig("plugin.embed.meta.short_description", "Create an embed.")
	longDescription  = i18n.NewFallbackConfig(
		"plugin.embed.meta.long_description",
		"Create a custom embed from the input you give me.")

	titleQuestion = i18n.NewFallbackConfig(
		"plugin.embed.input.title.question",
		"What should the title of the embed be?")
	descriptionQuestion = i18n.NewFallbackConfig(
		"plugin.embed.input.description.question",
		"What should the description of the embed be?")
)
