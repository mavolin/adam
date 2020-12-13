package i18n

type (
	// Func is the func used to retrieve a LangFunc, used to create
	// localized messages.
	// If there is no data available for the passed language or the passed
	// language is invalid, it should return a fallback language.
	//
	// If LangFunc is nil, the Localizer will use the fallback
	// translation.
	Func func(lang string) LangFunc

	// LangFunc is a function used to translate to a specific language.
	//
	// The first value is the unique id of the translation.
	//
	// The second parameter is a map with the filled placeholders, or, if there
	// are no placeholders a nil map.
	//
	// The third parameter is a number or a string of such defining the plural.
	// Valid plural data are number types or a string containing a number.
	// If plural is nil, Other should be used.
	//
	// If the LangFunc returns an error, the fallback translation will be used.
	LangFunc func(term Term, placeholders map[string]interface{}, plural interface{}) (string, error)
)
