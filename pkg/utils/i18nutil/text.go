package i18nutil

import "github.com/mavolin/adam/pkg/i18n"

// Text is an optionally localizable text.
// It can either be static, or be defined using a Config.
type Text struct {
	string string
	config i18n.Config
}

// NewText returns a new unlocalized Text.
func NewText(src string) Text {
	return Text{string: src}
}

// NewTextl returns a new localized Text using the passed i18n.Config.
func NewTextl(src i18n.Config) Text {
	return Text{config: src}
}

// NewTextl returns a new localized Text using the passed i18n.Term.
func NewTextlt(src i18n.Term) Text {
	return NewTextl(src.AsConfig())
}

// IsValid checks if the Text has no content.
func (t Text) IsValid() bool {
	return len(t.string) != 0 || t.config.IsValid()
}

// Get retrieves the value of the Text and localizes it, if possible.
func (t Text) Get(l *i18n.Localizer) (string, error) {
	if t.string != "" {
		return t.string, nil
	} else if t.config.IsValid() {
		return l.Localize(t.config)
	}

	return "", nil
}
