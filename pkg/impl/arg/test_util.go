package arg

import (
	"strconv"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
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
