package errors

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/lara-go/larago/logger"
	"github.com/lara-go/larago/logger/format"
	"github.com/lara-go/larago/support/utils"
	"github.com/urfave/cli"
)

// PanicHandlerInterface interface for panic handlers.
type PanicHandlerInterface interface {
	Defer()
}

// PanicHandler handles all panics.
type PanicHandler struct {
}

// Defer to handle panics.
func (h *PanicHandler) Defer() {
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
func (h *PanicHandler) handleError(err error) {
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

	// Handle exit with formatted error message.
	cli.HandleExitCoder(
		cli.NewExitError(
			fmt.Sprintf("%s\n--> %s:%d\n%s", message, file, line, debug.Stack()),
			exitCode,
		),
	)
}
