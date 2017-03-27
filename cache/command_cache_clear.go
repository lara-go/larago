package cache

import "github.com/urfave/cli"
import "github.com/lara-go/larago/logger"

// CommandCacheClear for the app.
type CommandCacheClear struct {
	Repository *Repository
	Logger     *logger.Logger
}

// GetCommand for the cli to register.
func (c *CommandCacheClear) GetCommand() cli.Command {
	return cli.Command{
		Name:     "cache:clear",
		Usage:    "Clear application cache",
		Category: "Cache",
	}
}

// Handle command.
func (c *CommandCacheClear) Handle(args cli.Args) error {
	c.Repository.Clear()

	c.Logger.Success("Cache was cleared.")

	return nil
}
