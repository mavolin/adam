package locutil

import "github.com/mavolin/adam/pkg/localization"

// Text is an optionally localizable text.
// It can either be static, or be defined using a Config.
type Text struct {
	string string
	config localization.Config
}

// NewStaticText returns a new static unlocalized Text.
func NewStaticText(src string) Text {
	return Text{
		string: src,
	}
}

// NewLocalizedText returns a localized Text.
func NewLocalizedText(src localization.Config) Text {
	return Text{
		config: src,
	}
}

// IsEmpty checks if the Text has no content.
func (t Text) IsEmpty() bool {
	return len(t.string) == 0 && !t.config.IsValid()
}

// Get retrieves the value of the Text and localizes it, if possible.
func (t Text) Get(l *localization.Localizer) (string, error) {
	if t.string != "" {
		return t.string, nil
	} else if t.config.IsValid() {
		return l.Localize(t.config)
	}

	return "", nil
}
