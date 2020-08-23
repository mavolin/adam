package plugin

import "github.com/mavolin/adam/pkg/localization"

var (
	guildTextType     = localization.NewFallbackConfig("channel_types.guild_text", "text channel")
	guildNewsType     = localization.NewFallbackConfig("channel_types.guild_news", "announcement channel")
	directMessageType = localization.NewFallbackConfig("channel_types.direct_message", "direct message")
)
