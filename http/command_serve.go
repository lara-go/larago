package http

import (
	"github.com/lara-go/larago"
	"github.com/urfave/cli"
)

// CommandServe command.
type CommandServe struct {
	Router *Router
	Config larago.Config

	listen string
}

// GetCommand for the cli to register.
func (c *CommandServe) GetCommand() cli.Command {
	return cli.Command{
		Name:     "http:serve",
		Usage:    "Start server",
		Category: "HTTP server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "listen, l",
				Usage:       "address to listen to (ex. 0.0.0.0:8080)",
				Destination: &c.listen,
			},
		},
	}
}

// Handle command.
func (c *CommandServe) Handle(args cli.Args) error {
	if c.listen != "" {
		c.Config.Set("HTTP.Listen", c.listen)
	}

	return c.Router.
		Bootstrap().
		Listen(c.Config.Get("HTTP.Listen").(string))
}
