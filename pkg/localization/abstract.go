package localization

type (
	// Func is the func used to retrieve a LangFunc, used to create
	// localized messages.
	// If there is no data available for the passed language or the passed
	// language is invalid, it should return a fallback language.
	//
	// If LangFunc is nil, the Localizer will use the fallback
	// translation.
	Func func(lang string) LangFunc
	// LangFunc is a language specific function for localizing messages.
	//
	// The first value is the unique id of the translation.
	//
	// The second parameter is a map with the filled placeholders, or, if there
	// are no placeholders a nil map.
	//
	// The third parameter is a number or a string of such defining the plural.
	// If it is nil, there is no pluralization.
	LangFunc func(term string, placeholders map[string]interface{}, plural interface{}) (string, error)
)
