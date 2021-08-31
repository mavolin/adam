package msgbuilder

import (
	"fmt"
	"reflect"
	"strconv"
	"sync/atomic"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"

	"github.com/mavolin/adam/internal/errorutil"
	"github.com/mavolin/adam/pkg/i18n"
)

type (
	// ComponentBuilder is the abstraction of any component builder.
	ComponentBuilder interface {
		Build(l *i18n.Localizer) (discord.Component, error)
		disable()
	}

	// TopLevelComponentBuilder is the abstraction of a builder that produces
	// top-level components.
	TopLevelComponentBuilder interface {
		ComponentBuilder
		// handle handles the passed *gateway.InteractionData.
		handle(*gateway.InteractionData) (bool, error)
	}

	// ActionRowComponentBuilder is the abstraction of a builder that produces
	// components that can be put into an ActionRowBuilder.
	ActionRowComponentBuilder interface {
		ComponentBuilder
		is(data *gateway.InteractionData) bool
		value() interface{}
	}
)

var id uint64

// nextID is used to generate custom id's for components.
func nextID() string {
	return strconv.FormatUint(atomic.AddUint64(&id, 1), 10)
}

// =============================================================================
// ActionRow
// =====================================================================================

// ActionRowBuilder is a builder used to build a *discord.ActionRowComponent.
// It must not be used to wrap SelectMenuBuilders, as they automatically wrap
// themselves in an ActionRow.
type ActionRowBuilder struct {
	components []ActionRowComponentBuilder

	resultVar interface{}
}

var _ TopLevelComponentBuilder = new(ActionRowBuilder)

// NewActionRow creates a new *ActionRowBuilder that stores the value of its
// components in the passed resultVar.
// resultVar must be a pointer to a variable.
func NewActionRow(resultVar interface{}) *ActionRowBuilder {
	return &ActionRowBuilder{
		// 5 is the max components an action row can hold
		components: make([]ActionRowComponentBuilder, 0, 5),
		resultVar:  resultVar,
	}
}

// With adds the passed ComponentBuilder to the ActionRowBuilder.
func (b *ActionRowBuilder) With(c ActionRowComponentBuilder) *ActionRowBuilder {
	b.components = append(b.components, c)
	return b
}

func (b *ActionRowBuilder) disable() {
	for _, c := range b.components {
		c.disable()
	}
}

func (b *ActionRowBuilder) handle(data *gateway.InteractionData) (bool, error) {
	for _, c := range b.components {
		if c.is(data) {
			result := reflect.ValueOf(c.value())
			reflect.ValueOf(b.resultVar).Elem().Set(result)

			return true, nil
		}
	}

	return false, nil
}

// Build builds the ActionRowBuilder.
// Errors returned by Build will be of type *ActionRowError.
func (b *ActionRowBuilder) Build(l *i18n.Localizer) (discord.Component, error) {
	r := &discord.ActionRowComponent{
		Components: make([]discord.Component, len(b.components)),
	}

	for i, cb := range b.components {
		c, err := cb.Build(l)
		if err != nil {
			return nil, NewActionRowError(i, reflect.TypeOf(cb).String(), err)
		}

		r.Components[i] = c
	}

	return r, nil
}

// =============================================================================
// ButtonBuilder
// =====================================================================================

type ButtonBuilder struct {
	label *i18n.Config
	style discord.ButtonStyle
	emoji *discord.ButtonEmoji

	url      discord.URL
	disabled bool

	id  string
	val interface{}
}

var _ ActionRowComponentBuilder = new(ButtonBuilder)

// NewButton creates a new *ButtonBuilder with the given label and the
// corresponding go value.
// val must be the element type of the ButtonBuilder's parent ActionRowBuilder.
func NewButton(style discord.ButtonStyle, label string, val interface{}) *ButtonBuilder {
	return NewButtonl(style, i18n.NewStaticConfig(label), val)
}

// NewButtonlt creates a new *ButtonBuilder with the given label and the
// corresponding go value.
// val must be the element type of the ButtonBuilder's parent ActionRowBuilder.
func NewButtonlt(style discord.ButtonStyle, label i18n.Term, val interface{}) *ButtonBuilder {
	return NewButtonl(style, label.AsConfig(), val)
}

// NewButtonl creates a new *ButtonBuilder with the given label and the
// corresponding go value.
// val must be the element type of the ButtonBuilder's parent ActionRowBuilder.
func NewButtonl(style discord.ButtonStyle, label *i18n.Config, val interface{}) *ButtonBuilder {
	return &ButtonBuilder{style: style, label: label, id: nextID(), val: val}
}

// WithEmoji assigns the passed emoji to the button.
func (b *ButtonBuilder) WithEmoji(emoji discord.ButtonEmoji) *ButtonBuilder {
	b.emoji = &emoji
	return b
}

// WithUnicodeEmoji assigns the passed unicode emoji to the button.
func (b *ButtonBuilder) WithUnicodeEmoji(emoji string) *ButtonBuilder {
	return b.WithEmoji(discord.ButtonEmoji{Name: emoji})
}

// WithURL attaches the given url to the button.
// This must be called, and only if, the button is a link-style button.
func (b *ButtonBuilder) WithURL(url discord.URL) *ButtonBuilder {
	b.url = url
	return b
}

// Disable disables the button.
func (b *ButtonBuilder) Disable() *ButtonBuilder {
	b.disabled = true
	return b
}

func (b *ButtonBuilder) disable() {
	b.Disable()
}

func (b *ButtonBuilder) is(data *gateway.InteractionData) bool {
	return data.CustomID == b.id
}

func (b *ButtonBuilder) value() interface{} {
	return b.val
}

func (b *ButtonBuilder) Build(l *i18n.Localizer) (c discord.Component, err error) {
	button := &discord.ButtonComponent{
		CustomID: b.id,
		Style:    b.style,
		Emoji:    b.emoji,
		URL:      b.url,
		Disabled: b.disabled,
	}

	button.Label, err = l.Localize(b.label)
	if err != nil {
		return nil, err
	}

	return button, nil
}

// =============================================================================
// SelectBuilder
// =====================================================================================

// SelectBuilder is a builder used to build a *discord.SelectComponent.
type SelectBuilder struct {
	id          string
	options     []*SelectOptionBuilder
	placeholder *i18n.Config
	minValues   option.Int
	maxValues   int
	disabled    bool

	resultVar interface{}
}

var _ TopLevelComponentBuilder = new(SelectBuilder)

// NewSelect creates a new *SelectBuilder that stores the value(s) of its
// components in the passed resultVar.
// If using the default bounds (1, 1), or (0, 1), resultVar must be a pointer.
// Otherwise, resultVar must be a pointer to a slice.
func NewSelect(resultVar interface{}) *SelectBuilder {
	return &SelectBuilder{
		id:        nextID(),
		options:   make([]*SelectOptionBuilder, 0, 25),
		resultVar: resultVar,
	}
}

// WithPlaceholder adds the passed placeholder to the select.
func (b *SelectBuilder) WithPlaceholder(placeholder string) *SelectBuilder {
	return b.WithPlaceholderl(i18n.NewStaticConfig(placeholder))
}

// WithPlaceholderlt adds the passed placeholder to the select.
func (b *SelectBuilder) WithPlaceholderlt(placeholder i18n.Term) *SelectBuilder {
	return b.WithPlaceholderl(placeholder.AsConfig())
}

// WithPlaceholderl adds the passed placeholder to the select.
func (b *SelectBuilder) WithPlaceholderl(placeholder *i18n.Config) *SelectBuilder {
	b.placeholder = placeholder
	return b
}

// WithBounds sets the passed bounds as min and max values.
func (b *SelectBuilder) WithBounds(min, max int) *SelectBuilder {
	b.minValues = option.NewInt(min)
	b.maxValues = max

	return b
}

// Disable disables the select.
func (b *SelectBuilder) Disable() *SelectBuilder {
	b.disabled = true
	return b
}

func (b *SelectBuilder) disable() {
	b.Disable()
}

// With adds the passed *SelectOptionBuilder to the SelectBuilder
func (b *SelectBuilder) With(option *SelectOptionBuilder) *SelectBuilder {
	b.options = append(b.options, option)
	return b
}

// WithDefault adds the passed *SelectOptionBuilder as the default option to the
// SelectBuilder.
func (b *SelectBuilder) WithDefault(option *SelectOptionBuilder) *SelectBuilder {
	option._default = true
	b.options = append(b.options, option)

	return b
}

func (b *SelectBuilder) handle(data *gateway.InteractionData) (bool, error) {
	if data.CustomID != b.id {
		return false, nil
	}

	if len(data.Options) == 0 {
		return true, nil
	} else if b.maxValues == 1 {
		var optionID string
		if err := data.Options[0].Value.UnmarshalTo(&optionID); err != nil {
			return false, errorutil.WithStack(err)
		}

		for _, optBuilder := range b.options {
			if optionID == optBuilder.id {
				result := reflect.ValueOf(optBuilder.val)
				reflect.ValueOf(b.resultVar).Elem().Set(result)

				return true, nil
			}
		}

		return false, errorutil.WithStack(fmt.Errorf("msgbuilder: SelectBuilder: found unknown option value %s",
			optionID))
	}

	optionIDs := make([]string, len(data.Options))
	for i, opt := range data.Options {
		if err := opt.Value.UnmarshalTo(&optionIDs[i]); err != nil {
			return false, errorutil.WithStack(err)
		}
	}

	resultV := reflect.ValueOf(b.resultVar)
	resultElem := resultV.Elem()

OptionBuilders:
	for _, optBuilder := range b.options {
		for i, optID := range optionIDs {
			if optBuilder.id == optID {
				resultElem = reflect.Append(resultElem, reflect.ValueOf(optBuilder.val))

				copy(optionIDs[i:], optionIDs[i+1:])
				optionIDs = optionIDs[:len(optionIDs)-1]
				continue OptionBuilders
			}
		}

		return false, errorutil.WithStack(fmt.Errorf(
			"msgbuilder: SelectBuilder: unable to find option value %s, only have %v", optBuilder.id, optionIDs))
	}

	if len(optionIDs) == 0 {
		resultV.Elem().Set(resultElem)
		return true, nil
	}

	return false, errorutil.WithStack(
		fmt.Errorf("msgbuilder: SelectBuilder: found unknown option values %v", optionIDs))
}

func (b *SelectBuilder) Build(l *i18n.Localizer) (c discord.Component, err error) {
	sel := &discord.SelectComponent{
		CustomID:  b.id,
		Options:   make([]discord.SelectComponentOption, len(b.options)),
		MinValues: b.minValues,
		MaxValues: b.maxValues,
		Disabled:  b.disabled,
	}

	if b.placeholder != nil {
		sel.Placeholder, err = l.Localize(b.placeholder)
		if err != nil {
			return nil, err
		}
	}

	for i, optBuilder := range b.options {
		opt, err := optBuilder.Build(l)
		if err != nil {
			return nil, NewSelectError(i, err)
		}

		sel.Options[i] = opt
	}

	return sel, nil
}

// =============================================================================
// SelectOptionBuilder
// =====================================================================================

type SelectOptionBuilder struct {
	label       *i18n.Config
	description *i18n.Config
	emoji       *discord.ButtonEmoji
	_default    bool

	id  string
	val interface{}
}

// NewSelectOption creates a new *SelectOptionBuilder with the given label and
// the corresponding go value.
//
// If the parent SelectBuilder uses the bounds (0, 1) or (1, 1), val must be of
// the elem type of the SelectBuilder's resultVar.
// Otherwise, val must be of the element type of the SelectBuilder's slice
// type.
func NewSelectOption(label string, val interface{}) *SelectOptionBuilder {
	return NewSelectOptionl(i18n.NewStaticConfig(label), val)
}

// NewSelectOptionlt creates a new *SelectOptionBuilder with the given label
// and the corresponding go value.
//
// If the parent SelectBuilder uses the bounds (0, 1) or (1, 1), val must be of
// the elem type of the SelectBuilder's resultVar.
// Otherwise, val must be of the element type of the SelectBuilder's slice
// type.
func NewSelectOptionlt(label i18n.Term, val interface{}) *SelectOptionBuilder {
	return NewSelectOptionl(label.AsConfig(), val)
}

// NewSelectOptionl creates a new *SelectOptionBuilder with the given label and
// the corresponding go value.
//
// If the parent SelectBuilder uses the bounds (0, 1) or (1, 1), val must be of
// the elem type of the SelectBuilder's resultVar.
// Otherwise, val must be of the element type of the SelectBuilder's slice
// type.
func NewSelectOptionl(label *i18n.Config, val interface{}) *SelectOptionBuilder {
	return &SelectOptionBuilder{label: label, id: nextID(), val: val}
}

// WithDescription adds the passed description to the SelectOptionBuilder.
func (b *SelectOptionBuilder) WithDescription(description string) *SelectOptionBuilder {
	return b.WithDescriptionl(i18n.NewStaticConfig(description))
}

// WithDescriptionlt adds the passed description to the SelectOptionBuilder.
func (b *SelectOptionBuilder) WithDescriptionlt(description i18n.Term) *SelectOptionBuilder {
	return b.WithDescriptionl(description.AsConfig())
}

// WithDescriptionl adds the passed description to the SelectOptionBuilder.
func (b *SelectOptionBuilder) WithDescriptionl(description *i18n.Config) *SelectOptionBuilder {
	b.description = description
	return b
}

// WithEmoji assigns the passed emoji to the SelectOptionBuilder.
func (b *SelectOptionBuilder) WithEmoji(emoji discord.ButtonEmoji) *SelectOptionBuilder {
	b.emoji = &emoji
	return b
}

// WithUnicodeEmoji assigns the passed unicode emoji to the select option.
func (b *SelectOptionBuilder) WithUnicodeEmoji(emoji string) *SelectOptionBuilder {
	return b.WithEmoji(discord.ButtonEmoji{Name: emoji})
}

func (b *SelectOptionBuilder) Build(l *i18n.Localizer) (selectOption discord.SelectComponentOption, err error) {
	selectOption = discord.SelectComponentOption{
		Value: b.id,
		Emoji: b.emoji,
	}

	selectOption.Label, err = l.Localize(b.label)
	if err != nil {
		return discord.SelectComponentOption{}, err
	}

	if b.description != nil {
		selectOption.Description, err = l.Localize(b.description)
	}

	return selectOption, err
}
