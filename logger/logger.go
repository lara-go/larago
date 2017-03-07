package logger

import (
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"time"

	ansii "github.com/lara-go/larago/logger/format"
	"github.com/lara-go/larago/support/utils"
)

const defaultSource = "larago/logger"

// Logger struct.
type Logger struct {
	DateTimeFormat string
	DebugMode      bool
	Logger         *log.Logger

	raw     bool
	from    []string
	trace   bool
	context interface{}
}

// Raw record. Without date, file:line, etc.
func (l Logger) Raw() *Logger {
	l.raw = true

	return &l
}

// From sets the actual source of the log record.
func (l Logger) From(from ...string) *Logger {
	l.from = from

	return &l
}

// WithTrace adds trace to the record.
func (l Logger) WithTrace() *Logger {
	l.trace = true

	return &l
}

// WithContext adds additional context to the record.
func (l Logger) WithContext(context interface{}) *Logger {
	l.context = context

	return &l
}

// SetOutput changes default output.
func (l *Logger) SetOutput(w io.Writer) *Logger {
	l.Logger.SetOutput(w)

	return l
}

// Error text.
func (l *Logger) Error(err error) {
	message := err.Error()

	if IsTTY() {
		message = ansii.Red(message).Format()
	}

	l.Println(l.prepare(message))
}

// Warning text.
func (l *Logger) Warning(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)

	if IsTTY() {
		message = ansii.Yellow(message).Format()
	}

	l.Println(l.prepare(message))
}

// Success text.
func (l *Logger) Success(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)

	if IsTTY() {
		message = ansii.Green(message).Format()
	}

	l.Println(l.prepare(message))
}

// Info text.
func (l *Logger) Info(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)

	if IsTTY() {
		message = ansii.Cyan(message).Format()
	}

	l.Println(l.prepare(message))
}

// Debug text.
func (l *Logger) Debug(format string, a ...interface{}) {
	if !l.DebugMode {
		return
	}

	message := fmt.Sprintf(format, a...)

	if IsTTY() {
		message = ansii.Blue(message).Format()
	}

	l.Println(l.prepare(message))
}

// Println message.
func (l *Logger) Println(message string) {
	l.Logger.Println(message)

	if l.context != nil {
		l.Logger.Printf("Context: %#v", l.context)
	}

	if l.trace {
		l.Logger.Printf("Trace: %s", string(debug.Stack()))
	}
}

// Writeln message to output.
func (l *Logger) prepare(message string) string {
	// If raw record requested, return only message.
	// Set raw to false, so next records will not be by default.
	if l.raw {
		return message
	}

	dateTime := fmt.Sprintf("[%s]", time.Now().Format(l.DateTimeFormat))

	if IsTTY() {
		dateTime = ansii.Yellow(dateTime).Format()
	}

	// If debug mode on, add source file:line to every record.
	if l.DebugMode {
		source := append(l.from, defaultSource)
		if file, line := utils.Thrower(source...); file != "" {
			linefile := fmt.Sprintf("(%s:%d)", file, line)

			if IsTTY() {
				linefile = ansii.Gray(linefile).Format()
			}

			return fmt.Sprintf("%s\n%s  %s",
				linefile,
				dateTime, message,
			)
		}
	}

	return fmt.Sprintf("%s  %s", dateTime, message)
}
