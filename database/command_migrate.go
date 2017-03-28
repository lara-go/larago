package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lara-go/larago/logger"

	"github.com/urfave/cli"
)

// CommandMigrate to apply DB changes.
type CommandMigrate struct {
	DB       *gorm.DB
	Migrator *Migrator
	Logger   *logger.Logger
}

// GetCommand for the cli to register.
func (c *CommandMigrate) GetCommand() cli.Command {
	return cli.Command{
		Name:     "migrate",
		Usage:    "Migrate database",
		Category: "Migrations",
	}
}

// Handle command.
func (c *CommandMigrate) Handle(args cli.Args) error {
	// Run migrations.
	err := c.Migrator.Migrate(c.DB)
	if err != nil {
		return fmt.Errorf("Could not migrate: %v", err.Error())
	}

	c.Logger.Success("Database was migrated.")

	return nil
}
