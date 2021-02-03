// Package i18n provides abstractions for field-based localization libraries.
package i18n

// Func is the function used to translate to a specific language.
//
// The first parameter is the unique id of the translation.
//
// The second parameter is a map with the filled placeholders, or, if there
// are no placeholders a nil map.
//
// The third parameter is the interface{} defining the plural.
// Valid plural types are number types or strings containing a number.
// If plural is nil, Other should be used.
//
// If the Func returns an error, the fallback translation will be used.
type Func func(term Term, placeholders map[string]interface{}, plural interface{}) (string, error)
