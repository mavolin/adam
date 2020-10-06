package locutil

import "github.com/mavolin/adam/pkg/i18n"

// Text is an optionally localizable text.
// It can either be static, or be defined using a Config.
type Text struct {
	string string
	config i18n.Config
}

// NewStaticText returns a new static unlocalized Text.
func NewStaticText(src string) Text {
	return Text{
		string: src,
	}
}

// NewLocalizedText returns a localized Text.
func NewLocalizedText(src i18n.Config) Text {
	return Text{
		config: src,
	}
}

// IsEmpty checks if the Text has no content.
func (t Text) IsEmpty() bool {
	return len(t.string) == 0 && !t.config.IsValid()
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
