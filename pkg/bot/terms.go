package bot

import "github.com/mavolin/adam/pkg/i18n"

var unknownCommandErrorDescription = i18n.NewFallbackConfig(
	"bot.errors.unknown_command.description",
	"I don't know a command with that name")
