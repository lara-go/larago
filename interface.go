package larago

import "github.com/urfave/cli"

// Bootstrapper function.
type Bootstrapper func(application *Application) error

// ServiceProvider interface.
type ServiceProvider interface {
	Register(application *Application)
}

// ConsoleCommand interface.
type ConsoleCommand interface {
	GetCommand() cli.Command
	Handle(args cli.Args) error
}

// Config interface.
type Config interface {
	Env() string
	Debug() bool
	Get(key string) interface{}
	Set(key string, value interface{})
}
