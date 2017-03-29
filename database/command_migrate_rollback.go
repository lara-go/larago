package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lara-go/larago/logger"

	"github.com/urfave/cli"
)

// CommandMigrateRollback to apply DB changes.
type CommandMigrateRollback struct {
	DB       *gorm.DB
	Migrator *Migrator
	Logger   *logger.Logger
}

// GetCommand for the cli to register.
func (c *CommandMigrateRollback) GetCommand() cli.Command {
	return cli.Command{
		Name:     "migrate:rollback",
		Usage:    "Rollback last migration",
		Category: "Migrations",
	}
}

// Handle command.
func (c *CommandMigrateRollback) Handle(args cli.Args) error {
	// Roll back migrations.
	err := c.Migrator.Rollback(c.DB)
	if err != nil {
		return fmt.Errorf("Could not rollback: %v", err.Error())
	}

	c.Logger.Success("Migrations were rolled back.")

	return nil
}
