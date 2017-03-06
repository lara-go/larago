package console

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/lara-go/larago/logger"
	"github.com/lara-go/larago/support/stubs"
	"github.com/lara-go/larago/support/utils"

	"github.com/urfave/cli"
)

const (
	commandsPath = "./app/commands"
	commandName  = commandsPath + "/%s.go"
)

// CommandMakeCommand to apply DB changes.
type CommandMakeCommand struct {
	Logger *logger.Logger
	name   string
}

// GetCommand for the cli to register.
func (c *CommandMakeCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:      "make:command",
		Usage:     "Make new console command",
		UsageText: "Makes new Command file in ./app/commands directory.\n",
		Category:  "Code generators",
		ArgsUsage: "[CommandName]",
	}
}

// Handle command.
func (c *CommandMakeCommand) Handle(args cli.Args) error {
	c.name = args.Get(0)

	if c.name == "" {
		return errors.New("Name can not be blank")
	}

	fileName, err := c.makeFile()
	if err != nil {
		return fmt.Errorf("Can't make new command: %s", err)
	}
	c.Logger.Success("New command created at: %s", fileName)

	return nil
}

// Make commands file.
func (c *CommandMakeCommand) makeFile() (string, error) {
	// Open command file.
	fileName := c.getResultFileName()
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Parse and execute template into command file.
	t := template.Must(template.New("command").Parse(stubs.CommandStub))
	var vars = struct {
		Name    string
		Command string
	}{
		Name:    c.name,
		Command: utils.ToSnake(c.name),
	}
	err = t.Execute(f, vars)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// Make path to the result command file.
func (c *CommandMakeCommand) getResultFileName() string {
	return fmt.Sprintf(commandName, utils.ToSnake(c.name))
}
