package format

import (
	"errors"
	"fmt"
)

const (
	off     = "\x1b[0m"
	red     = "\x1b[31m"
	green   = "\x1b[32m"
	yellow  = "\x1b[33m"
	blue    = "\x1b[34m"
	magenta = "\x1b[35m"
	cyan    = "\x1b[36m"
	gray    = "\x1b[90m"
)

// Text to send.
type Text struct {
	Text  string
	Color string
}

// ToError conversion.
func (m *Text) ToError() error {
	return errors.New(m.Format())
}

// Format message.
func (m *Text) String() string {
	return m.Format()
}

// Format message.
func (m *Text) Format() string {
	if m.Color != "" {
		return fmt.Sprintf("%s%s%s", m.Color, m.Text, off)
	}

	return m.Text
}

// NewText constructor.
func NewText(message string, color string) *Text {
	return &Text{
		Text:  message,
		Color: color,
	}
}

// Red text.
func Red(message string) *Text {
	return NewText(message, red)
}

// Yellow text.
func Yellow(message string) *Text {
	return NewText(message, yellow)
}

// Green text.
func Green(message string) *Text {
	return NewText(message, green)
}

// Blue text.
func Blue(message string) *Text {
	return NewText(message, blue)
}

// Magenta text.
func Magenta(message string) *Text {
	return NewText(message, magenta)
}

// Cyan text.
func Cyan(message string) *Text {
	return NewText(message, cyan)
}

// Gray text.
func Gray(message string) *Text {
	return NewText(message, gray)
}
