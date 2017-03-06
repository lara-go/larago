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
	middewarePath  = "./app/middleware"
	middlewareName = middewarePath + "/%s.go"
)

// CommandMakeMiddleware to apply DB changes.
type CommandMakeMiddleware struct {
	Logger *logger.Logger
	name   string
}

// GetCommand for the cli to register.
func (c *CommandMakeMiddleware) GetCommand() cli.Command {
	return cli.Command{
		Name:      "make:middleware",
		Usage:     "Make new middleware",
		UsageText: "Makes new Middleware file in ./app/middleware directory.\n",
		Category:  "Code generators",
		ArgsUsage: "[MiddlewareName]",
	}
}

// Handle command.
func (c *CommandMakeMiddleware) Handle(args cli.Args) error {
	c.name = args.Get(0)

	if c.name == "" {
		return errors.New("Name can not be blank")
	}

	fileName, err := c.makeFile()
	if err != nil {
		return fmt.Errorf("Can't make new middleware: %s", err)
	}
	c.Logger.Success("New middleware created at: %s", fileName)

	return nil
}

// Make migrations file.
func (c *CommandMakeMiddleware) makeFile() (string, error) {
	// Open migration file.
	fileName := c.getResultFileName()
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Parse and execute template into migration file.
	t := template.Must(template.New("middleware").Parse(stubs.MiddlewareStub))
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
func (c *CommandMakeMiddleware) getResultFileName() string {
	return fmt.Sprintf(middlewareName, utils.ToSnake(c.name))
}
