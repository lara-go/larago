package database

import (
	"github.com/lara-go/larago/logger"
	"github.com/urfave/cli"
)

// CommandDBSeed to apply DB changes.
type CommandDBSeed struct {
	Logger *logger.Logger
}

// GetCommand for the cli to register.
func (c *CommandDBSeed) GetCommand() cli.Command {
	return cli.Command{
		Name:     "db:seed",
		Usage:    "Seed database with initial info",
		Category: "Database",
	}
}

// Handle command.
func (c *CommandDBSeed) Handle(args cli.Args) error {
	c.Logger.Success("Database was seeded.")

	return nil
}
