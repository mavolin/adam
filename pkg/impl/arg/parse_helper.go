package arg

import (
	"reflect"

	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/adam/pkg/plugin"
)

var interfaceType = reflect.TypeOf(func(interface{}) {}).In(0)

// parseHelper is a helper struct that aids in parsing plugin.ArgConfigs.
// It assumes flags always start with a single minus ('-').
type parseHelper struct {
	rargData []plugin.RequiredArg
	oargData []plugin.OptionalArg
	flagData []plugin.Flag
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

func newParseHelper( //nolint:dupl
	rargs []plugin.RequiredArg, oargs []plugin.OptionalArg, flags []plugin.Flag, variadic bool, s *state.State,
	ctx *plugin.Context,
) *parseHelper {
	p := &parseHelper{
		rargData: rargs,
		oargData: oargs,
		flagData: flags,
		variadic: variadic,
		state:    s,
		ctx:      ctx,
		args:     make(plugin.Args, 0, len(rargs)+len(oargs)),
		flags:    make(plugin.Flags, len(flags)),
	}

	var numMultiFlags int

	for _, f := range flags {
		if f.IsMulti() {
			numMultiFlags++
		}
	}

	if numMultiFlags > 0 {
		p.multiFlags = make(map[string]reflect.Value, numMultiFlags)
	}

	return p
}

// store stores the parsed arguments in the context.
func (h *parseHelper) store() error {
	if h.variadicSlice.IsValid() {
		h.args = append(h.args, h.variadicSlice.Interface())
	}

	if len(h.args) < len(h.rargData) {
		return plugin.NewArgumentErrorl(notEnoughArgsError)
	}

	h.mergeFlags()
	h.fillFlagDefaults()
	h.fillArgDefaults()

	h.ctx.Args = h.args
	h.ctx.Flags = h.flags

	return nil
}

func (h *parseHelper) mergeFlags() {
	for name, flag := range h.multiFlags {
		h.flags[name] = flag.Interface()
	}
}

func (h *parseHelper) fillFlagDefaults() {
	for _, f := range h.flagData {
		if _, ok := h.flags[f.GetName()]; !ok {
			var val interface{}

			if f.IsMulti() {
				if f.GetDefault() != nil {
					val = f.GetDefault()
				} else {
					var t reflect.Type

					if def := f.GetType().GetDefault(); def == nil {
						t = interfaceType
					} else {
						t = reflect.TypeOf(def)
					}

					t = reflect.SliceOf(t)
					val = reflect.Zero(t).Interface()
				}
			} else {
				val = f.GetDefault()
				if val == nil {
					val = f.GetType().GetDefault()
				}
			}

			h.flags[f.GetName()] = val
		}
	}
}

func (h *parseHelper) fillArgDefaults() {
	argIndex := h.argIndex - len(h.rargData)

	if argIndex >= len(h.oargData) {
		return
	}

	for _, arg := range h.oargData[argIndex : len(h.oargData)-1] {
		val := arg.GetDefault()
		if val == nil {
			val = arg.GetType().GetDefault()
		}

		h.args = append(h.args, val)
	}

	last := h.oargData[len(h.oargData)-1]

	var val interface{}

	if h.variadic {
		if last.GetDefault() != nil {
			val = last.GetDefault()
		} else {
			var t reflect.Type

			if def := last.GetType().GetDefault(); def == nil {
				t = interfaceType
			} else {
				t = reflect.TypeOf(def)
			}

			t = reflect.SliceOf(t)
			val = reflect.Zero(t).Interface()
		}
	} else {
		val = last.GetDefault()
		if val == nil {
			val = last.GetType().GetDefault()
		}
	}

	h.args = append(h.args, val)
}

func (h *parseHelper) flag(name string) plugin.Flag {
	for _, flag := range h.flagData {
		if flag.GetName() == name {
			return flag
		}

		for _, alias := range flag.GetAliases() {
			if alias == name {
				return flag
			}
		}
	}

	return nil
}

func (h *parseHelper) addFlag(flag plugin.Flag, usedName, content string) (err error) {
	var val interface{}

	if flag.GetType() == Switch {
		val = true
	} else {
		ctx := &plugin.ParseContext{
			Context:  h.ctx,
			Raw:      content,
			Name:     "-" + flag.GetName(),
			UsedName: "-" + usedName,
			Kind:     plugin.KindFlag,
		}

		val, err = flag.GetType().Parse(h.state, ctx)
		if err != nil {
			return err
		}
	}

	if !flag.IsMulti() {
		return h.setSingleFlag(flag.GetName(), usedName, val)
	}

	h.setMultiFlag(flag.GetName(), val)
	return nil
}

func (h *parseHelper) setSingleFlag(name, usedName string, val interface{}) error {
	if _, ok := h.flags[name]; ok {
		return plugin.NewArgumentErrorl(flagUsedMultipleTimesError.
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
func (h *parseHelper) nextArg() (name string, typ plugin.ArgType, variadic bool, err error) {
	totalArgs := len(h.rargData) + len(h.oargData)
	if totalArgs == 0 {
		return "", nil, false, plugin.NewArgumentErrorl(tooManyArgsError)
	}

	if h.argIndex >= totalArgs {
		if !h.variadic {
			return "", nil, false, plugin.NewArgumentErrorl(tooManyArgsError)
		}

		if len(h.oargData) > 0 {
			arg := h.oargData[len(h.oargData)-1]

			name = arg.GetName(h.ctx.Localizer)
			return name, arg.GetType(), true, nil
		}

		arg := h.rargData[len(h.rargData)-1]

		name = arg.GetName(h.ctx.Localizer)
		return name, arg.GetType(), true, nil
	}

	variadic = h.argIndex == totalArgs-1 && h.variadic

	if h.argIndex < len(h.rargData) {
		arg := h.rargData[h.argIndex]

		name = arg.GetName(h.ctx.Localizer)
		return name, arg.GetType(), variadic, nil
	}

	arg := h.oargData[h.argIndex-len(h.rargData)]
	name = arg.GetName(h.ctx.Localizer)
	return name, arg.GetType(), variadic, nil
}

func (h *parseHelper) addArg(content string) error {
	name, typ, variadic, err := h.nextArg()
	if err != nil {
		return err
	}

	ctx := &plugin.ParseContext{
		Context:  h.ctx,
		Raw:      content,
		Name:     name,
		UsedName: name,
		Index:    h.argIndex,
		Kind:     plugin.KindArg,
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
