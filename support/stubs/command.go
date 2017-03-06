package stubs

// CommandStub template.
const CommandStub = `
package commands

import (
	"github.com/lara-go/larago/logger"

	"github.com/urfave/cli"
)

// {{.Name}} .
type {{.Name}} struct {
	Logger *logger.Logger
}

// GetCommand for the cli to register.
func (c *{{.Name}}) GetCommand() cli.Command {
	return cli.Command{
		Name:  "{{.Command}}",
		Usage: "",
	}
}

// Handle command.
func (c *{{.Name}}) Handle(args cli.Args) error {
	c.Logger.Info("Command info")

	return nil
}
`
