package larago

import (
	"os"
	"path"

	"github.com/urfave/cli"
)

// Kernel for cli commands.
type Kernel struct {
	application *Application

	bootstrappers   []Bootstrapper
	commands        []ConsoleCommand
	onReadyCallback func() error
}

// UseBootstrappers sets application bootstrappers.
func (k *Kernel) UseBootstrappers(bootstrappers ...Bootstrapper) *Kernel {
	k.bootstrappers = bootstrappers

	return k
}

// WithApplication sets application instance.
func (k *Kernel) WithApplication(application *Application) *Kernel {
	k.application = application

	return k
}

// Handle console commands.
func (k *Kernel) Handle() {
	k.application.BootstrapWith(k.bootstrappers...)

	app := cli.NewApp()

	app.Version = k.application.Version
	app.Name = k.application.Name
	app.Usage = k.application.Description

	app.Flags = k.GetGlobalFlags()
	app.Commands = k.makeCommands(k.application.GetCommands())

	app.Run(os.Args)
}

// GetGlobalFlags registers global flags.
func (k *Kernel) GetGlobalFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: path.Join(k.application.HomeDirectory, ".env"),
			Usage: "path to .env config file",
		},
		cli.StringFlag{
			Name:  "home, r",
			Value: k.application.HomeDirectory,
			Usage: "path to home directory",
		},
	}
}

// Cli commands factory.
func (k *Kernel) makeCommands(commands []ConsoleCommand) []cli.Command {
	var cliCommands []cli.Command

	for _, command := range commands {
		cliCommands = append(cliCommands, k.makeCommand(command))
	}

	return cliCommands
}

// Make command for the cli package.
func (k *Kernel) makeCommand(command ConsoleCommand) cli.Command {
	cliCommand := command.GetCommand()

	// Cli command handler.
	cliCommand.Action = func(c *cli.Context) error {
		// Resolve command's dependencies.
		k.application.Make(command)

		// Run Handler.
		if err := command.Handle(c.Args()); err != nil {
			panic(err)
		}

		return nil
	}

	return cliCommand
}

// OnReady sets onReady callback.
func (k *Kernel) OnReady(callback func() error) *Kernel {
	k.onReadyCallback = callback

	return k
}

func (k Kernel) defaultOnReadyCallback() error {
	return nil
}
