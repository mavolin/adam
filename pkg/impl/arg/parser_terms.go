package arg

import "github.com/mavolin/adam/pkg/i18n"

var (
	notEnoughArgsError = i18n.NewFallbackConfig(
		"args.parser.errors.not_enough_args", "There are not enough arguments to execute the command.")

	tooManyArgsError = i18n.NewFallbackConfig(
		"args.parser.errors.too_many_args", "Hold it chief! Those are too many arguments.")

	noArgsError = i18n.NewFallbackConfig(
		"args.parser.errors.no_args", "This command has no arguments and flags.")

	unknownFlagError = i18n.NewFallbackConfig(
		"args.parser.errors.unknown_flag", "I don't know a flag by the name of `-{{.name}}`.")

	flagUsedMultipleTimesError = i18n.NewFallbackConfig(
		"args.parser.errors.flag_used_multiple_times", "You can't use the `-{{.name}}`-flag multiple times.")

	emptyFlagError = i18n.NewFallbackConfig(
		"args.parser.errors.empty_flag", "You can't leave the `-{{.name}}`-flag empty.")

	emptyArgError = i18n.NewFallbackConfig(
		"args.parser.errors.empty_arg", "The argument at position {{.position}} may not be empty.")

	groupNotClosedError = i18n.NewFallbackConfig(
		"args.parser.errors.group_not_closed", "You need to close the {{.quote}}.")
)

type (
	unknownFlagErrorPlaceholders struct {
		Name string
	}

	flagUsedMultipleTimesErrorPlaceholders struct {
		Name string
	}

	emptyFlagErrorPlaceholders struct {
		Name string
	}

	emptyArgErrorPlaceholders struct {
		Position int
	}

	groupNotClosedErrorPlaceholders struct {
		Quote string
	}
)
