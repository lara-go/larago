package console

import (
	"github.com/lara-go/larago"
	"github.com/lara-go/larago/logger"

	"github.com/urfave/cli"
)

// CommandEnv to apply DB changes.
type CommandEnv struct {
	Config larago.Config
	Logger *logger.Logger
}

// GetCommand for the cli to register.
func (c *CommandEnv) GetCommand() cli.Command {
	return cli.Command{
		Name:  "env",
		Usage: "Show current environment mode",
	}
}

// Handle command.
func (c *CommandEnv) Handle(args cli.Args) error {
	var debug string

	if c.Config.Debug() {
		debug = "on"
	} else {
		debug = "off"
	}

	c.Logger.Success("Application is in %s environment with debug mode %s.", c.Config.Env(), debug)

	return nil
}
