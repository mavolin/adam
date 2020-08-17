package errors

import (
	"github.com/mavolin/adam/internal/constant"
	"github.com/mavolin/adam/internal/errorutil"
	"github.com/mavolin/adam/pkg/localization"
	"github.com/mavolin/adam/pkg/utils/discordutil"
)

// stackTrace attempts to extract the stacktrace from the error.
// If that does not succeed iw till generate a stack trace.
func stackTrace(err error, skip int) (stack []uintptr) {
	if s, ok := err.(stackTracer); ok {
		stack = s.StackTrace()
	} else {
		stack = errorutil.GenerateStackTrace(1 + skip)
	}

	return
}

func newErrorEmbedBuilder(l *localization.Localizer) *discordutil.EmbedBuilder {
	// the error can be ignored, because there is fallback
	title, _ := l.Localize(errorTitle)

	return discordutil.NewEmbedBuilder().
		WithSimpleTitle(title).
		WithColor(constant.ErrorColor)
}

func newInfoEmbedBuild(l *localization.Localizer) *discordutil.EmbedBuilder {
	// the error can be ignored, because there is fallback
	title, _ := l.Localize(infoTitle)

	return discordutil.NewEmbedBuilder().
		WithSimpleTitle(title).
		WithColor(constant.InfoColor)
}
