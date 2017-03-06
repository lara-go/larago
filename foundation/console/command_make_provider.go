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
	providersPath = "./app/providers"
	providersName = providersPath + "/%s.go"
)

// CommandMakeProvider to apply DB changes.
type CommandMakeProvider struct {
	Logger *logger.Logger
	name   string
}

// GetCommand for the cli to register.
func (c *CommandMakeProvider) GetCommand() cli.Command {
	return cli.Command{
		Name:      "make:provider",
		Usage:     "Make new service provider",
		UsageText: "Makes new ServiceProvider file in ./app/providers directory.\n",
		Category:  "Code generators",
		ArgsUsage: "[ServiceProviderName]",
	}
}

// Handle command.
func (c *CommandMakeProvider) Handle(args cli.Args) error {
	c.name = args.Get(0)

	if c.name == "" {
		return errors.New("Name can not be blank")
	}

	fileName, err := c.makeFile()
	if err != nil {
		return fmt.Errorf("Can't make new service provider: %s", err)
	}
	c.Logger.Success("New service provider created at: %s", fileName)

	return nil
}

// Make migrations file.
func (c *CommandMakeProvider) makeFile() (string, error) {
	// Open migration file.
	fileName := c.getResultFileName()
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Parse and execute template into migration file.
	t := template.Must(template.New("provider").Parse(stubs.ServiceProviderStub))
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
func (c *CommandMakeProvider) getResultFileName() string {
	return fmt.Sprintf(providersName, utils.ToSnake(c.name))
}
