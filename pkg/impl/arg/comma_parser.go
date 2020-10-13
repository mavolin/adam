package arg

import (
	"reflect"
	"strings"

	"github.com/mavolin/disstate/v2/pkg/state"

	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
)

var commaConfigEscaper = strings.NewReplacer(",,", ",", "--", "-")

var interfaceType = reflect.TypeOf(func(interface{}) {}).In(0)

type commaParser struct {
	config CommaConfig
	state  *state.State
	ctx    *plugin.Context

	lexer *commaLexer

	args plugin.Args
	// variadicSlice is the reflect.Value of the slice containing the variadic
	// argument, if there is one
	variadicSlice reflect.Value
	argIndex      int

	flags      plugin.Flags
	multiFlags map[string]reflect.Value
}

func newCommaParser(args string, cfg CommaConfig, s *state.State, ctx *plugin.Context) *commaParser {
	p := &commaParser{
		config: cfg,
		state:  s,
		ctx:    ctx,
		lexer:  newCommaLexer(args, len(cfg.RequiredArgs), len(cfg.Flags) > 0),
		args:   make(plugin.Args, 0, len(cfg.RequiredArgs)+len(cfg.OptionalArgs)),
		flags:  make(plugin.Flags, len(cfg.Flags)),
	}

	var numMultiFlags int

	for _, f := range cfg.Flags {
		if f.Multi {
			numMultiFlags++
		}
	}

	if numMultiFlags > 0 {
		p.multiFlags = make(map[string]reflect.Value, numMultiFlags)
	}

	return p
}

func (p *commaParser) parse() (plugin.Args, plugin.Flags, error) {
	err := p.startParse()
	if err != nil {
		return nil, nil, err
	}

	if p.variadicSlice.IsValid() {
		p.args = append(p.args, p.variadicSlice.Interface())
	}

	if len(p.args) < len(p.config.RequiredArgs) {
		return nil, nil, errors.NewArgumentParsingErrorl(notEnoughArgsError)
	}

	p.mergeFlags()
	p.fillDefaults()

	return p.args, p.flags, nil
}

func (p *commaParser) mergeFlags() {
	for name, flag := range p.multiFlags {
		p.flags[name] = flag.Interface()
	}
}

func (p *commaParser) fillDefaults() {
	p.fillFlagDefaults()
	p.fillArgDefaults()
}

func (p *commaParser) fillFlagDefaults() {
	for _, f := range p.config.Flags {
		if _, ok := p.flags[f.Name]; !ok {
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

			p.flags[f.Name] = val
		}
	}
}

func (p *commaParser) fillArgDefaults() {
	argIndex := p.argIndex - len(p.config.RequiredArgs)

	if argIndex >= len(p.config.OptionalArgs) {
		return
	}

	for i := argIndex; i < len(p.config.OptionalArgs)-1; i++ {
		arg := p.config.OptionalArgs[i]

		val := arg.Default
		if val == nil {
			val = arg.Type.Default()
		}

		p.args = append(p.args, val)
	}

	last := p.config.OptionalArgs[len(p.config.OptionalArgs)-1]

	val := last.Default
	if val == nil {
		val = last.Type.Default()
	}

	if p.config.Variadic {
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

	p.args = append(p.args, val)
}

func (p *commaParser) setFlag(name string, val interface{}, multi bool) {
	if !multi {
		p.flags[name] = val
		return
	}

	rval := reflect.ValueOf(val)

	if flags, ok := p.multiFlags[name]; ok {
		p.multiFlags[name] = reflect.Append(flags, rval)
	} else {
		var t reflect.Type

		if val == nil {
			t = interfaceType
		} else {
			t = rval.Type()
		}

		sliceType := reflect.SliceOf(t)

		flags := reflect.MakeSlice(sliceType, 1, 1)
		flags.Index(0).Set(rval)

		p.multiFlags[name] = flags
	}
}

// nextArg returns meta information about the next argument.
func (p *commaParser) nextArg() (name string, typ Type, variadic bool, err error) {
	totalArgs := len(p.config.RequiredArgs) + len(p.config.OptionalArgs)
	if totalArgs == 0 {
		return "", nil, false, errors.NewArgumentParsingErrorl(tooManyArgsError)
	}

	if p.argIndex >= totalArgs {
		if !p.config.Variadic {
			return "", nil, false, errors.NewArgumentParsingErrorl(tooManyArgsError)
		}

		if len(p.config.OptionalArgs) > 0 {
			arg := p.config.OptionalArgs[len(p.config.OptionalArgs)-1]

			name, err = arg.Name.Get(p.ctx.Localizer)
			return name, arg.Type, true, err
		}

		arg := p.config.RequiredArgs[len(p.config.RequiredArgs)-1]

		name, err = arg.Name.Get(p.ctx.Localizer)
		return name, arg.Type, true, err
	}

	variadic = p.argIndex == totalArgs-1 && p.config.Variadic

	if p.argIndex < len(p.config.RequiredArgs) {
		arg := p.config.RequiredArgs[p.argIndex]

		name, err = arg.Name.Get(p.ctx.Localizer)
		return name, arg.Type, variadic, err
	}

	arg := p.config.OptionalArgs[p.argIndex-len(p.config.RequiredArgs)]
	name, err = arg.Name.Get(p.ctx.Localizer)
	return name, arg.Type, variadic, err
}

func (p *commaParser) startParse() error {
	item, err := p.lexer.nextItem()
	if err != nil {
		return err
	}

	if len(p.config.RequiredArgs)+len(p.config.OptionalArgs) == 0 && len(p.config.Flags) == 0 {
		if item.typ != itemEOF {
			return errors.NewArgumentParsingErrorl(noArgsError)
		}
	}

	for ; err == nil && item.typ != itemEOF; item, err = p.lexer.nextItem() {
		// the lexer keeps track of the correct ordering, so we don't need to
		// worry about that
		switch item.typ {
		case itemFlagName:
			err = p.parseFlag(item)
		case itemArgContent:
			err = p.parseArg(item)
		default:
			return errors.NewWithStackf("arg: unexpected item during parsing: %s", item.typ)
		}

		if err != nil {
			return err
		}
	}

	return err
}

// parseFlag parses a flag.
// The passed commaItem must have the commaItemType itemFlagName.
func (p *commaParser) parseFlag(flagName commaItem) (err error) {
	f := p.config.flag(flagName.val)
	if f == nil {
		return errors.NewArgumentParsingErrorl(unknownFlagError.
			WithPlaceholders(unknownFlagErrorPlaceholders{
				Name: flagName.val,
			}))
	}

	if !f.Multi {
		_, ok := p.flags[f.Name]
		if ok {
			return errors.NewArgumentParsingErrorl(flagUsedMultipleTimesError.
				WithPlaceholders(flagUsedMultipleTimesErrorPlaceholders{
					Name: flagName.val,
				}))
		}
	}

	var val interface{}

	if f.Type == Switch {
		val = true
	} else {
		content, err := p.lexer.nextItem()
		if err != nil {
			return err
		} else if content.typ != itemFlagContent {
			return errors.NewArgumentParsingErrorl(emptyFlagError.
				WithPlaceholders(emptyFlagErrorPlaceholders{
					Name: flagName.val,
				}))
		}

		ctx := &Context{
			Context:  p.ctx,
			Raw:      commaConfigEscaper.Replace(content.val), // replace escapes
			Name:     f.Name,
			UsedName: flagName.val,
			Kind:     KindFlag,
		}

		val, err = f.Type.Parse(p.state, ctx)
		if err != nil {
			return err
		}
	}

	p.setFlag(f.Name, val, f.Multi)

	finalizer, err := p.lexer.nextItem()
	switch {
	case err != nil:
		return err
	case finalizer.typ == itemFlagContent && f.Type == Switch:
		return errors.NewArgumentParsingErrorl(switchWithContentError.
			WithPlaceholders(switchWithContentErrorPlaceholders{
				Name: flagName.val,
			}))
	case finalizer.typ != itemComma && finalizer.typ != itemEOF:
		return errors.NewWithStackf("arg: unexpected item during parsing: %s", finalizer.typ)
	default:
		return nil
	}
}

func (p *commaParser) parseArg(content commaItem) error {
	name, typ, variadic, err := p.nextArg()
	if err != nil {
		return err
	}

	ctx := &Context{
		Context:  p.ctx,
		Raw:      commaConfigEscaper.Replace(content.val), // replace escapes
		Name:     name,
		UsedName: name,
		Index:    p.argIndex,
		Kind:     KindArg,
	}

	val, err := typ.Parse(p.state, ctx)
	if err != nil {
		return err
	}

	if !variadic {
		p.args = append(p.args, val)
	} else {
		rval := reflect.ValueOf(val)

		if p.variadicSlice.IsValid() {
			p.variadicSlice = reflect.Append(p.variadicSlice, rval)
		} else {
			var t reflect.Type

			if val == nil {
				t = interfaceType
			} else {
				t = rval.Type()
			}

			sliceType := reflect.SliceOf(t)

			p.variadicSlice = reflect.MakeSlice(sliceType, 1, 1)
			p.variadicSlice.Index(0).Set(rval)
		}
	}

	finalizer, err := p.lexer.nextItem()
	if err != nil {
		return err
	} else if finalizer.typ != itemComma && finalizer.typ != itemEOF {
		return errors.NewWithStackf("arg: unexpected item during parsing: %s", finalizer.typ)
	}

	p.argIndex++

	return nil
}
