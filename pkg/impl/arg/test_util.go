package arg

import (
	"strconv"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type mockType struct {
	name string
	desc string

	parseFunc func(s *state.State, ctx *plugin.ParseContext) (interface{}, error)

	dfault interface{}
}

func (m mockType) GetName(*i18n.Localizer) string        { return m.name }
func (m mockType) GetDescription(*i18n.Localizer) string { return m.desc }

func (m mockType) Parse(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
	return m.parseFunc(s, ctx)
}

func (m mockType) GetDefault() interface{} { return m.dfault }

var (
	mockTypeInt = mockType{
		parseFunc: func(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
			return strconv.Atoi(ctx.Raw)
		},
		dfault: 0,
	}
	mockTypeString = mockType{
		parseFunc: func(s *state.State, ctx *plugin.ParseContext) (interface{}, error) {
			return ctx.Raw, nil
		},
		dfault: "",
	}
)

type mockModule struct {
	name      string
	shortDesc string
	longDesc  string
	commands  []plugin.Command
	modules   []plugin.Module
}

var _ plugin.Module = mockModule{}

func (m mockModule) GetName() string                            { return m.name }
func (m mockModule) GetShortDescription(*i18n.Localizer) string { return m.shortDesc }
func (m mockModule) GetLongDescription(*i18n.Localizer) string  { return m.longDesc }
func (m mockModule) Commands() []plugin.Command                 { return m.commands }
func (m mockModule) Modules() []plugin.Module                   { return m.modules }
