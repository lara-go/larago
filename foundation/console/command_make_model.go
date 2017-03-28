package console

import (
	"errors"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/lara-go/larago/logger"
	"github.com/lara-go/larago/support/stubs"
	"github.com/lara-go/larago/support/utils"

	"github.com/urfave/cli"
)

var (
	modelsPath = path.Join(".", "app", "models")
	modelsName = path.Join(modelsPath, "%s.go")
)

// CommandMakeModel to apply DB changes.
type CommandMakeModel struct {
	Logger *logger.Logger
	name   string
}

// GetCommand for the cli to register.
func (c *CommandMakeModel) GetCommand() cli.Command {
	return cli.Command{
		Name:      "make:model",
		Usage:     "Make new model",
		UsageText: "Makes new model file in ./app/models directory.\n",
		Category:  "Code generators",
		ArgsUsage: "[ModelName]",
	}
}

// Handle command.
func (c *CommandMakeModel) Handle(args cli.Args) error {
	c.name = args.Get(0)

	if c.name == "" {
		return errors.New("Model name can not be blank")
	}

	fileName, err := c.makeFile()
	if err != nil {
		return fmt.Errorf("Can't make new model: %s", err)
	}
	c.Logger.Success("New model created at: %s", fileName)

	return nil
}

// Make migrations file.
func (c *CommandMakeModel) makeFile() (string, error) {
	// Open migration file.
	fileName := c.getResultFileName()
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Parse and execute template into migration file.
	t := template.Must(template.New("model").Parse(stubs.ModelStub))
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
func (c *CommandMakeModel) getResultFileName() string {
	return fmt.Sprintf(modelsName, utils.ToSnake(c.name))
}
