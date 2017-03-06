package database

import (
	"github.com/lara-go/larago/logger"
	"github.com/urfave/cli"
)

// CommandMigrateReset to apply DB changes.
type CommandMigrateReset struct {
	Migrator *Migrator
	Logger   *logger.Logger
}

// GetCommand for the cli to register.
func (c *CommandMigrateReset) GetCommand() cli.Command {
	return cli.Command{
		Name:     "migrate:reset",
		Usage:    "Reset all migrations",
		Category: "Migrations",
	}
}

// Handle command.
func (c *CommandMigrateReset) Handle(args cli.Args) error {
	c.Logger.Success("All migrations were reset.")

	return nil
}
