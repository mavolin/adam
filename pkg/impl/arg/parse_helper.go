package arg

import (
	"reflect"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/i18nutil"
)

// parseHelper is a helper struct that aids in parsing plugin.ArgConfigs.
// It assumes flags always start with a single minus ('-').
type parseHelper struct {
	rargData []requiredArg
	oargData []optionalArg
	flagData []flag
	variadic bool

	state *state.State
	ctx   *plugin.Context

	args plugin.Args
	// variadicSlice is the reflect.Value of the slice containing the variadic
	// argument, if there is one
	variadicSlice reflect.Value
	argIndex      int

	flags      plugin.Flags
	multiFlags map[string]reflect.Value
}

type (
	requiredArg struct {
		name *i18nutil.Text
		typ  Type
	}

	optionalArg struct {
		name   *i18nutil.Text
		typ    Type
		dfault interface{}
	}

	flag struct {
		name    string
		aliases []string
		typ     Type
		dfault  interface{}
		multi   bool
	}
)

func newParseHelper( //nolint:dupl
	rargs []RequiredArg, oargs []OptionalArg, flags []Flag, variadic bool, s *state.State, ctx *plugin.Context,
) *parseHelper {
	p := &parseHelper{
		variadic: variadic,
		state:    s,
		ctx:      ctx,
		args:     make(plugin.Args, 0, len(rargs)+len(oargs)),
		flags:    make(plugin.Flags, len(flags)),
	}

	if len(rargs) > 0 {
		p.rargData = make([]requiredArg, len(rargs))

		for i, arg := range rargs {
			p.rargData[i] = requiredArg{
				name: i18nutil.NewText(arg.Name),
				typ:  arg.Type,
			}
		}
	}

	if len(oargs) > 0 {
		p.oargData = make([]optionalArg, len(oargs))

		for i, arg := range oargs {
			p.oargData[i] = optionalArg{
				name:   i18nutil.NewText(arg.Name),
				typ:    arg.Type,
				dfault: arg.Default,
			}
		}
	}

	if len(flags) > 0 {
		p.flagData = make([]flag, len(flags))

		for i, f := range flags {
			p.flagData[i] = flag{
				name:    f.Name,
				aliases: f.Aliases,
				typ:     f.Type,
				dfault:  f.Default,
				multi:   f.Multi,
			}
		}
	}

	var numMultiFlags int

	for _, f := range flags {
		if f.Multi {
			numMultiFlags++
		}
	}

	if numMultiFlags > 0 {
		p.multiFlags = make(map[string]reflect.Value, numMultiFlags)
	}

	return p
}

func newParseHelperl( //nolint:dupl
	rargs []LocalizedRequiredArg, oargs []LocalizedOptionalArg, flags []LocalizedFlag, variadic bool,
	s *state.State, ctx *plugin.Context,
) *parseHelper {
	p := &parseHelper{
		variadic: variadic,
		state:    s,
		ctx:      ctx,
		args:     make(plugin.Args, 0, len(rargs)+len(oargs)),
		flags:    make(plugin.Flags, len(flags)),
	}

	if len(rargs) > 0 {
		p.rargData = make([]requiredArg, len(rargs))

		for i, arg := range rargs {
			p.rargData[i] = requiredArg{
				name: i18nutil.NewTextl(arg.Name),
				typ:  arg.Type,
			}
		}
	}

	if len(oargs) > 0 {
		p.oargData = make([]optionalArg, len(oargs))

		for i, arg := range oargs {
			p.oargData[i] = optionalArg{
				name:   i18nutil.NewTextl(arg.Name),
				typ:    arg.Type,
				dfault: arg.Default,
			}
		}
	}

	if len(flags) > 0 {
		p.flagData = make([]flag, len(flags))

		for i, f := range flags {
			p.flagData[i] = flag{
				name:    f.Name,
				aliases: f.Aliases,
				typ:     f.Type,
				dfault:  f.Default,
				multi:   f.Multi,
			}
		}
	}

	var numMultiFlags int

	for _, f := range flags {
		if f.Multi {
			numMultiFlags++
		}
	}

	if numMultiFlags > 0 {
		p.multiFlags = make(map[string]reflect.Value, numMultiFlags)
	}

	return p
}

func (h *parseHelper) get() (plugin.Args, plugin.Flags, error) {
	if h.variadicSlice.IsValid() {
		h.args = append(h.args, h.variadicSlice.Interface())
	}

	if len(h.args) < len(h.rargData) {
		return nil, nil, errors.NewArgumentErrorl(notEnoughArgsError)
	}

	h.mergeFlags()
	h.fillFlagDefaults()
	h.fillArgDefaults()

	return h.args, h.flags, nil
}

func (h *parseHelper) mergeFlags() {
	for name, flag := range h.multiFlags {
		h.flags[name] = flag.Interface()
	}
}

func (h *parseHelper) fillFlagDefaults() {
	for _, f := range h.flagData {
		if _, ok := h.flags[f.name]; !ok {
			val := f.dfault
			if val == nil {
				val = f.typ.Default()
			}

			if f.multi {
				rval := reflect.ValueOf(val)
				var t reflect.Type

				if val == nil {
					t = interfaceType
				} else {
					t = rval.Type()
				}

				sliceType := reflect.SliceOf(t)

				slice := reflect.MakeSlice(sliceType, 1, 1)
				slice.Index(0).Set(rval)

				val = slice.Interface()
			}

			h.flags[f.name] = val
		}
	}
}

func (h *parseHelper) fillArgDefaults() {
	argIndex := h.argIndex - len(h.rargData)

	if argIndex >= len(h.oargData) {
		return
	}

	for i := argIndex; i < len(h.oargData)-1; i++ {
		arg := h.oargData[i]

		val := arg.dfault
		if val == nil {
			val = arg.typ.Default()
		}

		h.args = append(h.args, val)
	}

	last := h.oargData[len(h.oargData)-1]

	val := last.dfault
	if val == nil {
		val = last.typ.Default()
	}

	if h.variadic {
		rval := reflect.ValueOf(val)
		var t reflect.Type

		if val == nil {
			t = interfaceType
		} else {
			t = rval.Type()
		}

		sliceType := reflect.SliceOf(t)

		slice := reflect.MakeSlice(sliceType, 1, 1)
		slice.Index(0).Set(rval)

		val = slice.Interface()
	}

	h.args = append(h.args, val)
}

func (h *parseHelper) flag(name string) *flag {
	for _, flag := range h.flagData {
		if flag.name == name {
			return &flag
		}

		for _, alias := range flag.aliases {
			if alias == name {
				return &flag
			}
		}
	}

	return nil
}

func (h *parseHelper) addFlag(flag *flag, usedName, content string) (err error) {
	var val interface{}

	if flag.typ == Switch {
		val = true
	} else {
		ctx := &Context{
			Context:  h.ctx,
			Raw:      content,
			Name:     "-" + flag.name,
			UsedName: "-" + usedName,
			Kind:     KindFlag,
		}

		val, err = flag.typ.Parse(h.state, ctx)
		if err != nil {
			return err
		}
	}

	if !flag.multi {
		return h.setSingleFlag(flag.name, usedName, val)
	}

	h.setMultiFlag(flag.name, val)
	return nil
}

func (h *parseHelper) setSingleFlag(name, usedName string, val interface{}) error {
	if _, ok := h.flags[name]; ok {
		return errors.NewArgumentErrorl(flagUsedMultipleTimesError.
			WithPlaceholders(flagUsedMultipleTimesErrorPlaceholders{
				Name: usedName,
			}))
	}

	h.flags[name] = val

	return nil
}

func (h *parseHelper) setMultiFlag(name string, val interface{}) {
	rval := reflect.ValueOf(val)

	if flags, ok := h.multiFlags[name]; ok {
		h.multiFlags[name] = reflect.Append(flags, rval)
		return
	}

	var t reflect.Type

	if val == nil {
		t = interfaceType
	} else {
		t = rval.Type()
	}

	sliceType := reflect.SliceOf(t)

	flags := reflect.MakeSlice(sliceType, 1, 1)
	flags.Index(0).Set(rval)

	h.multiFlags[name] = flags
}

// nextArg returns meta information about the next argument.
func (h *parseHelper) nextArg() (name string, typ Type, variadic bool, err error) {
	totalArgs := len(h.rargData) + len(h.oargData)
	if totalArgs == 0 {
		return "", nil, false, errors.NewArgumentErrorl(tooManyArgsError)
	}

	if h.argIndex >= totalArgs {
		if !h.variadic {
			return "", nil, false, errors.NewArgumentErrorl(tooManyArgsError)
		}

		if len(h.oargData) > 0 {
			arg := h.oargData[len(h.oargData)-1]

			name, err = arg.name.Get(h.ctx.Localizer)
			return name, arg.typ, true, err
		}

		arg := h.rargData[len(h.rargData)-1]

		name, err = arg.name.Get(h.ctx.Localizer)
		return name, arg.typ, true, err
	}

	variadic = h.argIndex == totalArgs-1 && h.variadic

	if h.argIndex < len(h.rargData) {
		arg := h.rargData[h.argIndex]

		name, err = arg.name.Get(h.ctx.Localizer)
		return name, arg.typ, variadic, err
	}

	arg := h.oargData[h.argIndex-len(h.rargData)]
	name, err = arg.name.Get(h.ctx.Localizer)
	return name, arg.typ, variadic, err
}

func (h *parseHelper) addArg(content string) error {
	name, typ, variadic, err := h.nextArg()
	if err != nil {
		return err
	}

	ctx := &Context{
		Context:  h.ctx,
		Raw:      content,
		Name:     name,
		UsedName: name,
		Index:    h.argIndex,
		Kind:     KindArg,
	}

	val, err := typ.Parse(h.state, ctx)
	if err != nil {
		return err
	}

	if !variadic {
		h.args = append(h.args, val)
	} else {
		rval := reflect.ValueOf(val)

		if h.variadicSlice.IsValid() {
			h.variadicSlice = reflect.Append(h.variadicSlice, rval)
		} else {
			var t reflect.Type

			if val == nil {
				t = interfaceType
			} else {
				t = rval.Type()
			}

			sliceType := reflect.SliceOf(t)

			h.variadicSlice = reflect.MakeSlice(sliceType, 1, 1)
			h.variadicSlice.Index(0).Set(rval)
		}
	}

	h.argIndex++

	return nil
}
