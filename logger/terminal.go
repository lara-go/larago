package logger

import (
	"os"

	isatty "github.com/mattn/go-isatty"
)

var isTTY = isatty.IsTerminal(os.Stdout.Fd())

// IsTTY app runs in tty.
func IsTTY() bool {
	return isTTY
}
