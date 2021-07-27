// Package arg provides implementations for the argument abstractions found in
// Package plugins.
package arg

import (
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/plugin"
)

type (
	// Config is an unlocalized plugin.ArgConfig.
	Config struct {
		// RequiredArgs are the required arguments.
		RequiredArgs []RequiredArg
		// OptionalArgs are the optional arguments.
		OptionalArgs []OptionalArg
		// Variadic specifies whether the last argument is variadic, i.e. it
		// may be specified more than once.
		Variadic bool

		// Flags are the flags.
		Flags []Flag

		iRequiredArgs []plugin.RequiredArg
		iOptionalArgs []plugin.OptionalArg
		iFlags        []plugin.Flag
	}

	// RequiredArg is an unlocalized required argument.
	RequiredArg struct {
		// Name is the name of the argument.
		Name string
		// Type is the type of the argument.
		Type plugin.ArgType
		// Description is an optional short description of the argument.
		Description string
	}

	// OptionalArg is an unlocalized optional argument.
	OptionalArg struct {
		// Name is the name of the argument.
		Name string
		// Type is the type of the argument.
		Type plugin.ArgType
		// Default is the default value of the argument.
		//
		// If Default is (interface{})(nil), Type.GetDefault() will be used.
		Default interface{}
		// Description is an optional short description of the argument.
		Description string
	}

	// Flag is an unlocalized flag.
	Flag struct {
		// Name is the name of the flag.
		Name string
		// Aliases contains the optional aliases of the flag.
		Aliases []string
		// Type is the type of the flag.
		Type plugin.ArgType
		// Default is the default value of the flag, and is used if the flag
		// isn't set.
		//
		// If Default is (interface{})(nil), Type.GetDefault() will be used.
		Default interface{}
		// Description is an optional short description of the flag.
		Description string
		// Multi specifies whether this flag can be used multiple times.
		Multi bool
	}
)

func (r RequiredArg) GetName(*i18n.Localizer) string        { return r.Name }
func (r RequiredArg) GetType() plugin.ArgType               { return r.Type }
func (r RequiredArg) GetDescription(*i18n.Localizer) string { return r.Description }

func (o OptionalArg) GetName(*i18n.Localizer) string        { return o.Name }
func (o OptionalArg) GetType() plugin.ArgType               { return o.Type }
func (o OptionalArg) GetDefault() interface{}               { return o.Default }
func (o OptionalArg) GetDescription(*i18n.Localizer) string { return o.Description }

func (f Flag) GetName() string                       { return f.Name }
func (f Flag) GetAliases() []string                  { return f.Aliases }
func (f Flag) GetType() plugin.ArgType               { return f.Type }
func (f Flag) GetDefault() interface{}               { return f.Default }
func (f Flag) GetDescription(*i18n.Localizer) string { return f.Description }
func (f Flag) IsMulti() bool                         { return f.Multi }

var _ plugin.ArgConfig = new(Config)

func (c *Config) GetRequiredArgs() []plugin.RequiredArg {
	c.setInterfaces()
	return c.iRequiredArgs
}

func (c *Config) GetOptionalArgs() []plugin.OptionalArg {
	c.setInterfaces()
	return c.iOptionalArgs
}

func (c *Config) IsVariadic() bool {
	return c.Variadic
}

func (c *Config) GetFlags() []plugin.Flag {
	c.setInterfaces()
	return c.iFlags
}

func (c *Config) setInterfaces() {
	if c.iRequiredArgs == nil && len(c.RequiredArgs) > 0 {
		c.iRequiredArgs = make([]plugin.RequiredArg, len(c.RequiredArgs))

		for i, rarg := range c.RequiredArgs {
			c.iRequiredArgs[i] = rarg
		}
	}

	if c.iFlags == nil && len(c.OptionalArgs) > 0 {
		c.iOptionalArgs = make([]plugin.OptionalArg, len(c.OptionalArgs))

		for i, oarg := range c.OptionalArgs {
			c.iOptionalArgs[i] = oarg
		}
	}

	if c.iFlags == nil && len(c.Flags) > 0 {
		c.iFlags = make([]plugin.Flag, len(c.Flags))

		for i, flag := range c.Flags {
			c.iFlags[i] = flag
		}
	}
}

type (
	LocalizedConfig struct {
		// RequiredArgs are the required arguments.
		RequiredArgs []LocalizedRequiredArg
		// OptionalArgs are the optional arguments.
		OptionalArgs []LocalizedOptionalArg
		// Variadic specifies whether the last argument is variadic, i.e. it
		// may be specified more than once.
		Variadic bool

		// Flags are the flags.
		Flags []LocalizedFlag

		iRequiredArgs []plugin.RequiredArg
		iOptionalArgs []plugin.OptionalArg
		iFlags        []plugin.Flag
	}

	// LocalizedRequiredArg is a localized required argument.
	LocalizedRequiredArg struct {
		// Name is the name of the argument.
		Name *i18n.Config
		// Type is the type of the argument.
		Type plugin.ArgType
		// Description is an optional short description of the argument.
		Description *i18n.Config
	}

	// LocalizedOptionalArg is an localized optional argument.
	LocalizedOptionalArg struct {
		// Name is the name of the argument.
		Name *i18n.Config
		// Type is the type of the argument.
		Type plugin.ArgType
		// Default is the default value of the argument.
		//
		// If Default is (interface{})(nil), Type.GetDefault() will be used.
		Default interface{}
		// Description is an optional short description of the argument.
		Description *i18n.Config
	}

	// LocalizedFlag is a localized flag.
	LocalizedFlag struct {
		// Name is the name of the flag.
		Name string
		// Aliases contains the optional aliases of the flag.
		Aliases []string
		// Type is the type of the flag.
		Type plugin.ArgType
		// Default is the default value of the flag, and is used if the flag
		// isn't set.
		//
		// If Default is (interface{})(nil), Type.GetDefault() will be used.
		Default interface{}
		// Description is an optional short description of the flag.
		Description *i18n.Config
		// Multi specifies whether this flag can be used multiple times.
		Multi bool
	}
)

func (r LocalizedRequiredArg) GetName(l *i18n.Localizer) string {
	if name, err := l.Localize(r.Name); err == nil {
		return name
	}

	return ""
}

func (r LocalizedRequiredArg) GetType() plugin.ArgType { return r.Type }

func (r LocalizedRequiredArg) GetDescription(l *i18n.Localizer) string {
	if desc, err := l.Localize(r.Description); err == nil {
		return desc
	}

	return ""
}

func (o LocalizedOptionalArg) GetName(l *i18n.Localizer) string {
	if name, err := l.Localize(o.Name); err == nil {
		return name
	}

	return ""
}

func (o LocalizedOptionalArg) GetType() plugin.ArgType { return o.Type }
func (o LocalizedOptionalArg) GetDefault() interface{} { return o.Default }

func (o LocalizedOptionalArg) GetDescription(l *i18n.Localizer) string {
	if desc, err := l.Localize(o.Description); err == nil {
		return desc
	}

	return ""
}

func (f LocalizedFlag) GetName() string         { return f.Name }
func (f LocalizedFlag) GetAliases() []string    { return f.Aliases }
func (f LocalizedFlag) GetType() plugin.ArgType { return f.Type }
func (f LocalizedFlag) GetDefault() interface{} { return f.Default }

func (f LocalizedFlag) GetDescription(l *i18n.Localizer) string {
	if desc, err := l.Localize(f.Description); err == nil {
		return desc
	}

	return ""
}

func (f LocalizedFlag) IsMulti() bool { return f.Multi }

func (c *LocalizedConfig) GetRequiredArgs() []plugin.RequiredArg {
	c.setInterfaces()
	return c.iRequiredArgs
}

func (c *LocalizedConfig) GetOptionalArgs() []plugin.OptionalArg {
	c.setInterfaces()
	return c.iOptionalArgs
}

func (c *LocalizedConfig) IsVariadic() bool {
	return c.Variadic
}

func (c *LocalizedConfig) GetFlags() []plugin.Flag {
	c.setInterfaces()
	return c.iFlags
}

func (c *LocalizedConfig) setInterfaces() {
	if c.iRequiredArgs == nil && len(c.RequiredArgs) > 0 {
		c.iRequiredArgs = make([]plugin.RequiredArg, len(c.RequiredArgs))

		for i, rarg := range c.RequiredArgs {
			c.iRequiredArgs[i] = rarg
		}
	}

	if c.iFlags == nil && len(c.OptionalArgs) > 0 {
		c.iOptionalArgs = make([]plugin.OptionalArg, len(c.OptionalArgs))

		for i, oarg := range c.OptionalArgs {
			c.iOptionalArgs[i] = oarg
		}
	}

	if c.iFlags == nil && len(c.Flags) > 0 {
		c.iFlags = make([]plugin.Flag, len(c.Flags))

		for i, flag := range c.Flags {
			c.iFlags[i] = flag
		}
	}
}
