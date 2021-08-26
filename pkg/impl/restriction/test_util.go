package restriction

import (
	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

// plugin.RestrictionFuncs useful in testing

var (
	errorFunc1Description = "abc"
	errorFunc2Description = "def"
	errorFunc3Description = "ghi"
	errorFunc4Description = "jkl"

	fatalErrorFuncDescription = "mno"

	embeddableErrorFuncEmbeddableDescription = "pqr"
	embeddableErrorFuncDefaultDescription    = "stu"

	fatalEmbeddableErrorFuncEmbeddableDescription = "vwx"
	fatalEmbeddableErrorFuncDefaultDescription    = "yza"

	errUnexpectedErrorFuncReturn = errors.New("bcd")
)

func errorFunc1(*state.State, *plugin.Context) error {
	return plugin.NewRestrictionError(errorFunc1Description)
}
func errorFunc2(*state.State, *plugin.Context) error {
	return plugin.NewRestrictionError(errorFunc2Description)
}
func errorFunc3(*state.State, *plugin.Context) error {
	return plugin.NewRestrictionError(errorFunc3Description)
}
func errorFunc4(*state.State, *plugin.Context) error {
	return plugin.NewRestrictionError(errorFunc4Description)
}
func fatalErrorFunc(*state.State, *plugin.Context) error {
	return plugin.NewFatalRestrictionError(
		fatalErrorFuncDescription)
}
func embeddableErrorFunc(*state.State, *plugin.Context) error {
	return &EmbeddableError{
		EmbeddableVersion: plugin.NewRestrictionError(embeddableErrorFuncEmbeddableDescription),
		DefaultVersion:    plugin.NewRestrictionError(embeddableErrorFuncDefaultDescription),
	}
}

func fatalEmbeddableErrorFunc(*state.State, *plugin.Context) error {
	return &EmbeddableError{
		EmbeddableVersion: plugin.NewFatalRestrictionError(fatalEmbeddableErrorFuncEmbeddableDescription),
		DefaultVersion:    plugin.NewFatalRestrictionError(fatalEmbeddableErrorFuncDefaultDescription),
	}
}

func unexpectedErrorFunc(*state.State, *plugin.Context) error { return errUnexpectedErrorFuncReturn }

func defaultRestrictionErrorFunc(*state.State, *plugin.Context) error {
	return plugin.DefaultRestrictionError
}

func defaultFatalRestrictionErrorFunc(*state.State, *plugin.Context) error {
	return plugin.DefaultFatalRestrictionError
}

func nilFunc(*state.State, *plugin.Context) error { return nil }
