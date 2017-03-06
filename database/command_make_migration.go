package database

import (
	"errors"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/lara-go/larago/logger"
	"github.com/lara-go/larago/support/stubs"

	"github.com/urfave/cli"
)

const (
	dateFormat     = "2006_01_02_150405"
	migrationsPath = "./app/database/migrations"
	migrationName  = migrationsPath + "/%s_%s.go"
)

// CommandMakeMigration to apply DB changes.
type CommandMakeMigration struct {
	Logger *logger.Logger
	name   string
}

// GetCommand for the cli to register.
func (c *CommandMakeMigration) GetCommand() cli.Command {
	return cli.Command{
		Name:      "make:migration",
		Usage:     "Make new migration",
		UsageText: "Makes new migration file in ./app/database/migrations directory.\n",
		Description: "Please follow this naming convension:\n" +
			"     - CreateUsersTable\n" +
			"     - AddNameColumnToUsersTable",
		Category:  "Code generators",
		ArgsUsage: "[MigrationName]",
	}
}

// Handle command.
func (c *CommandMakeMigration) Handle(args cli.Args) error {
	c.name = args.Get(0)

	if c.name == "" {
		return errors.New("Migration name can not be blank")
	}

	fileName, err := c.makeFile()
	if err != nil {
		return fmt.Errorf("Can't make new migration: %s", err)
	}
	c.Logger.Success("New migration created at: %s", fileName)

	return nil
}

// Make migrations file.
func (c *CommandMakeMigration) makeFile() (string, error) {
	// Open migration file.
	fileName := c.getResultFileName()
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Parse and execute template into migration file.
	t := template.Must(template.New("migration").Parse(stubs.MigrationStub))
	var vars = struct {
		Name string
	}{
		Name: c.name,
	}
	err = t.Execute(f, vars)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// Make path to the result migration file.
func (c *CommandMakeMigration) getResultFileName() string {
	return fmt.Sprintf(migrationName, time.Now().Format(dateFormat), c.name)
}
