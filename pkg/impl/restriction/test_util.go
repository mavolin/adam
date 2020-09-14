package restriction

import (
	"github.com/mavolin/disstate/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

// plugin.RestrictionFuncs useful in testing

var (
	errorFuncReturn1 = errors.NewRestrictionError("abc")
	errorFuncReturn2 = errors.NewRestrictionError("def")
	errorFuncReturn3 = errors.NewRestrictionError("ghi")
	errorFuncReturn4 = errors.NewRestrictionError("jkl")

	fatalErrorFuncReturn = errors.NewFatalRestrictionError("mno")

	embeddableErrorFuncReturn = &EmbeddableError{
		EmbeddableVersion: errors.NewRestrictionError("pqr"),
		DefaultVersion:    errors.New("stu"),
	}

	fatalEmbeddableErrorFuncReturn = &EmbeddableError{
		EmbeddableVersion: errors.NewFatalRestrictionError("vwx"),
		DefaultVersion:    errors.New("yza"),
	}

	unexpectedErrorFuncReturn = errors.New("bcd")
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

func unexpectedErrorFunc(*state.State, *plugin.Context) error { return unexpectedErrorFuncReturn }

func defaultRestrictionErrorFunc(*state.State, *plugin.Context) error {
	return errors.DefaultRestrictionError
}

func defaultFatalRestrictionErrorFunc(*state.State, *plugin.Context) error {
	return errors.DefaultFatalRestrictionError
}

func nilFunc(*state.State, *plugin.Context) error { return nil }
