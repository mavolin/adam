package arg

import (
	"fmt"
	"strconv"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type mockType struct {
	name string
	desc string

	parseFunc func(s *state.State, ctx *Context) (interface{}, error)

	dfault interface{}
}

func (m mockType) Name(*i18n.Localizer) string        { return m.name }
func (m mockType) Description(*i18n.Localizer) string { return m.desc }

func (m mockType) Parse(s *state.State, ctx *Context) (interface{}, error) {
	return m.parseFunc(s, ctx)
}

func (m mockType) Default() interface{} { return m.dfault }

var (
	mockTypeInt = mockType{
		parseFunc: func(s *state.State, ctx *Context) (interface{}, error) {
			return strconv.Atoi(ctx.Raw)
		},
		dfault: 0,
	}
	mockTypeString = mockType{
		parseFunc: func(s *state.State, ctx *Context) (interface{}, error) {
			return ctx.Raw, nil
		},
		dfault: "",
	}
)

var testArgFormatter plugin.ArgFormatter = func(i plugin.ArgInfo, optional, variadic bool) string {
	if optional {
		if variadic {
			return fmt.Sprintf("[%s:%s+]", i.Name, i.Type.Name)
		}

		return fmt.Sprintf("[%s:%s]", i.Name, i.Type.Name)
	}

	if variadic {
		return fmt.Sprintf("<%s:%s+>", i.Name, i.Type.Name)
	}

	return fmt.Sprintf("<%s:%s>", i.Name, i.Type.Name)
}
