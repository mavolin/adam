package arg

import (
	"regexp"
	"strings"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

// Code is the type used for code enclosed in a markdown code block.
// Single, double, and triple backticks are permitted.
//
// Go type: *CodeBlock
var Code plugin.ArgType = new(code)

type code struct{}

// CodeBlock is the type returned by Code.
type CodeBlock struct {
	// Language is the language the user specified, if any.
	Language string
	// Code is the code itself.
	// Language and backticks are removed.
	Code string
	// QtyBackticks is the number of backticks the user used.
	// Guaranteed to be 1, 2, or 3 in error != nil scenarios.
	QtyBackticks int
}

func (c code) GetName(l *i18n.Localizer) string {
	name, _ := l.Localize(codeName) // we have a fallback
	return name
}

func (c code) GetDescription(l *i18n.Localizer) string {
	desc, _ := l.Localize(codeDescription) // we have a fallback
	return desc
}

var (
	// these regexps aren't perfect (e.g. they allow ```a``), but they should cover
	// most cases
	singleBacktickRegexp = regexp.MustCompile(`^\x60(?P<code>[^\x60]+)\x60$`)
	doubleBacktickRegexp = regexp.MustCompile(`^\x60\x60(?P<code>(?:\x60[^\x60]|[^\x60])+)\x60\x60$`)
	tripleBacktickRegexp = regexp.MustCompile(
		`^\x60\x60\x60(?:(?P<lang>\S+)\n)?(?P<code>(?:\x60\x60[^\x60]|\x60[^\x60]|[^\x60])+)\x60\x60\x60$`)
)

func (c code) Parse(_ *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	if matches := singleBacktickRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		return &CodeBlock{
			Code:         strings.Trim(matches[1], "\n"),
			QtyBackticks: 1,
		}, nil
	} else if matches := doubleBacktickRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		return &CodeBlock{
			Code:         strings.Trim(matches[1], "\n"),
			QtyBackticks: 2,
		}, nil
	} else if matches := tripleBacktickRegexp.FindStringSubmatch(ctx.Raw); len(matches) >= 2 {
		if len(matches) >= 3 { // with language
			return &CodeBlock{
				Language:     matches[1],
				Code:         strings.Trim(matches[2], "\n"),
				QtyBackticks: 3,
			}, nil
		}

		return &CodeBlock{
			Code:         strings.Trim(matches[1], "\n"),
			QtyBackticks: 3,
		}, nil
	}

	return nil, newArgumentError2(codeInvalidErrorArg, codeInvalidErrorFlag, ctx, nil)
}

func (c code) GetDefault() interface{} {
	return (*CodeBlock)(nil)
}
