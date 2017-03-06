package http

import (
	"os"
	"path"

	"github.com/lara-go/larago"
	"github.com/lara-go/larago/logger"

	"github.com/urfave/cli"
)

// CommandUp to apply DB changes.
type CommandUp struct {
	Application *larago.Application
	Logger      *logger.Logger
}

// GetCommand for the cli to register.
func (c *CommandUp) GetCommand() cli.Command {
	return cli.Command{
		Name:     "http:up",
		Usage:    "Bring server out of maintenance mode",
		Category: "HTTP server",
	}
}

// Handle command.
func (c *CommandUp) Handle(args cli.Args) error {
	os.Remove(path.Join(c.Application.HomeDirectory, downFile))

	c.Logger.Success("Server is now live.")

	return nil
}
