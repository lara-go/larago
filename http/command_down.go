package http

import (
	"os"
	"path"

	"github.com/lara-go/larago"
	"github.com/lara-go/larago/logger"

	"github.com/urfave/cli"
)

const (
	downFile = "down"
)

// CommandDown to apply DB changes.
type CommandDown struct {
	Application *larago.Application
	Logger      *logger.Logger
}

// GetCommand for the cli to register.
func (c *CommandDown) GetCommand() cli.Command {
	return cli.Command{
		Name:     "http:down",
		Usage:    "Put server in maintenance mode",
		Category: "HTTP server",
	}
}

// Handle command.
func (c *CommandDown) Handle(args cli.Args) error {
	os.OpenFile(path.Join(c.Application.HomeDirectory, downFile), os.O_RDONLY|os.O_CREATE, 0666)

	c.Logger.Success("Server is now in maintenance mode.")

	return nil
}
