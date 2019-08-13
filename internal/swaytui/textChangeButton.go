package swaytui

import (
  tui "github.com/marcusolsson/tui-go"
	"image"
	"strings"
        runewidth "github.com/mattn/go-runewidth"
)

var _ tui.Widget = &TextChangeButton{}

// Button is a widget that can be activated to perform some action, or to
// answer a question.
type TextChangeButton struct {
	tui.WidgetBase

	text string

	onActivated func(*TextChangeButton)
}

// NewButton returns a new Button with the given text as the label.
func NewTextChangeButton(text string) *TextChangeButton {
	return &TextChangeButton{
		text: text,
	}
}
// SetText sets the text of the button
func (b *TextChangeButton) SetText(text string) {
  b.text = text
}
// SetText gets the text of the button
func (b *TextChangeButton) GetText() string {
  return b.text
}

// Draw draws the button.
func (b *TextChangeButton) Draw(p *tui.Painter) {
	style := "button"
	if b.IsFocused() {
		style += ".focused"
	}
	p.WithStyle(style, func(p *tui.Painter) {
		lines := strings.Split(b.text, "\n")
		for i, line := range lines {
			p.FillRect(0, i, b.Size().X, 1)
			p.DrawText(0, i, line)
		}
	})
}

// SizeHint returns the recommended size hint for the button.
func (b *TextChangeButton) SizeHint() image.Point {
	if len(b.text) == 0 {
		return b.MinSizeHint()
	}

	var size image.Point
	lines := strings.Split(b.text, "\n")
	for _, line := range lines {
		if w := runewidth.StringWidth(line); w > size.X {
			size.X = w
		}
	}
	size.Y = len(lines)

	return size
}

// OnKeyEvent handles keys events.
func (b *TextChangeButton) OnKeyEvent(ev tui.KeyEvent) {
	if !b.IsFocused() {
		return
	}
	if ev.Key == tui.KeyEnter && b.onActivated != nil {
		b.onActivated(b)
	}
}


// OnActivated allows a custom function to be run whenever the button is activated.
func (b *TextChangeButton) OnActivated(fn func(b *TextChangeButton)) {
	b.onActivated = fn
}
