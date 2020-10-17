package arg

import (
	"reflect"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

type parseHelper struct {
	rargData []RequiredArg
	oargData []OptionalArg
	flagData []Flag
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

func newParseHelper(
	rargs []RequiredArg, oargs []OptionalArg, flags []Flag, variadic bool, s *state.State, ctx *plugin.Context,
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
		return nil, nil, errors.NewArgumentParsingErrorl(notEnoughArgsError)
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
		if _, ok := h.flags[f.Name]; !ok {
			val := f.Default
			if val == nil {
				val = f.Type.Default()
			}

			if f.Multi {
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

			h.flags[f.Name] = val
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

		val := arg.Default
		if val == nil {
			val = arg.Type.Default()
		}

		h.args = append(h.args, val)
	}

	last := h.oargData[len(h.oargData)-1]

	val := last.Default
	if val == nil {
		val = last.Type.Default()
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

func (h *parseHelper) flag(name string) *Flag {
	for _, flag := range h.flagData {
		if flag.Name == name {
			return &flag
		}

		for _, alias := range flag.Aliases {
			if alias == name {
				return &flag
			}
		}
	}

	return nil
}

func (h *parseHelper) addFlag(flag *Flag, usedName, content string) (err error) {
	var val interface{}

	if flag.Type == Switch {
		val = true
	} else {
		ctx := &Context{
			Context:  h.ctx,
			Raw:      content,
			Name:     flag.Name,
			UsedName: usedName,
			Kind:     KindFlag,
		}

		val, err = flag.Type.Parse(h.state, ctx)
		if err != nil {
			return err
		}
	}

	if !flag.Multi {
		return h.setSingleFlag(flag.Name, usedName, val)
	}

	h.setMultiFlag(flag.Name, val)
	return nil
}

func (h *parseHelper) setSingleFlag(name, usedName string, val interface{}) error {
	if _, ok := h.flags[name]; ok {
		return errors.NewArgumentParsingErrorl(flagUsedMultipleTimesError.
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
		return "", nil, false, errors.NewArgumentParsingErrorl(tooManyArgsError)
	}

	if h.argIndex >= totalArgs {
		if !h.variadic {
			return "", nil, false, errors.NewArgumentParsingErrorl(tooManyArgsError)
		}

		if len(h.oargData) > 0 {
			arg := h.oargData[len(h.oargData)-1]

			name, err = arg.Name.Get(h.ctx.Localizer)
			return name, arg.Type, true, err
		}

		arg := h.rargData[len(h.rargData)-1]

		name, err = arg.Name.Get(h.ctx.Localizer)
		return name, arg.Type, true, err
	}

	variadic = h.argIndex == totalArgs-1 && h.variadic

	if h.argIndex < len(h.rargData) {
		arg := h.rargData[h.argIndex]

		name, err = arg.Name.Get(h.ctx.Localizer)
		return name, arg.Type, variadic, err
	}

	arg := h.oargData[h.argIndex-len(h.rargData)]
	name, err = arg.Name.Get(h.ctx.Localizer)
	return name, arg.Type, variadic, err
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
		Kind:     KindArgument,
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
