package arg

import (
	"strconv"

	"github.com/mavolin/disstate/v4/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type mockType struct {
	name string
	desc string

	parseFunc func(s *state.State, ctx *plugin.ParseContext) (interface{}, error)

	Default interface{}
}

func (m mockType) GetName(*i18n.Localizer) string        { return m.name }
func (m mockType) GetDescription(*i18n.Localizer) string { return m.desc }

func (m mockType) Parse(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	return m.parseFunc(s, ctx)
}

func (m mockType) GetDefault() interface{} { return m.Default }

var (
	mockTypeInt = mockType{
		parseFunc: func(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
			return strconv.Atoi(ctx.Raw)
		},
		Default: 0,
	}
	mockTypeString = mockType{
		parseFunc: func(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
			return ctx.Raw, nil
		},
		Default: "",
	}
)
