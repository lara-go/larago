package foundation

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/asaskevich/EventBus"
	"github.com/lara-go/larago/logger"
	"github.com/lara-go/larago/logger/format"
	"github.com/lara-go/larago/support/utils"
	"github.com/urfave/cli"
)

// ExitHandler handles all panics.
type ExitHandler struct {
	Events *EventBus.EventBus
}

// Exit terminates application.
func (h *ExitHandler) Exit(message string, exitCode int) {
	cli.HandleExitCoder(
		cli.NewExitError(message, exitCode),
	)
}

// Defer to handle panics.
func (h *ExitHandler) Defer() {
	var err error
	var ok bool

	// Try to recover.
	if r := recover(); r != nil {
		if err, ok = r.(error); !ok {
			err = errors.New(r.(string))
		}

		h.handleError(err)
	}
}

// HandleError writes error to the output and terminates.
func (h *ExitHandler) handleError(err error) {
	message := err.Error()
	exitCode := 2

	// Check if error is already an instance of cli.ExitCoder
	// and update exit code
	if exitErr, ok := err.(cli.ExitCoder); ok {
		exitCode = exitErr.ExitCode()
	}

	// Color output.
	if logger.IsTTY() {
		message = format.Red(message).Format()
	}

	// Find trower line.
	file, line := utils.Thrower("larago/errors")

	// Publish event.
	if h.Events != nil {
		h.Events.Publish("panic", message, file, line, debug.Stack(), exitCode)
	}

	// Handle exit with formatted error message.
	h.Exit(
		fmt.Sprintf("%s\n--> %s:%d\n%s", message, file, line, debug.Stack()),
		exitCode,
	)
}
