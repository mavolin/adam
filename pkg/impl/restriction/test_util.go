package restriction

import (
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

// plugin.RestrictionFuncs useful in testing

var (
	errorFuncReturn1 = plugin.NewRestrictionError("abc")
	errorFuncReturn2 = plugin.NewRestrictionError("def")
	errorFuncReturn3 = plugin.NewRestrictionError("ghi")
	errorFuncReturn4 = plugin.NewRestrictionError("jkl")

	fatalErrorFuncReturn = plugin.NewFatalRestrictionError("mno")

	embeddableErrorFuncReturn = &EmbeddableError{
		EmbeddableVersion: plugin.NewRestrictionError("pqr"),
		DefaultVersion:    errors.New("stu"),
	}

	fatalEmbeddableErrorFuncReturn = &EmbeddableError{
		EmbeddableVersion: plugin.NewFatalRestrictionError("vwx"),
		DefaultVersion:    errors.New("yza"),
	}

	errUnexpectedErrorFuncReturn = errors.New("bcd")
)

func errorFunc1(*state.State, *plugin.Context) error          { return errorFuncReturn1 }
func errorFunc2(*state.State, *plugin.Context) error          { return errorFuncReturn2 }
func errorFunc3(*state.State, *plugin.Context) error          { return errorFuncReturn3 }
func errorFunc4(*state.State, *plugin.Context) error          { return errorFuncReturn4 }
func fatalErrorFunc(*state.State, *plugin.Context) error      { return fatalErrorFuncReturn }
func embeddableErrorFunc(*state.State, *plugin.Context) error { return embeddableErrorFuncReturn }

func fatalEmbeddableErrorFunc(*state.State, *plugin.Context) error {
	return fatalEmbeddableErrorFuncReturn
}

func unexpectedErrorFunc(*state.State, *plugin.Context) error { return errUnexpectedErrorFuncReturn }

func defaultRestrictionErrorFunc(*state.State, *plugin.Context) error {
	return plugin.DefaultRestrictionError
}

func defaultFatalRestrictionErrorFunc(*state.State, *plugin.Context) error {
	return plugin.DefaultFatalRestrictionError
}

func nilFunc(*state.State, *plugin.Context) error { return nil }
