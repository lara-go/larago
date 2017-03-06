package larago

import "github.com/urfave/cli"

// Bootstrapper function.
type Bootstrapper func(application *Application) error

// ServiceProvider interface.
type ServiceProvider interface {
	// Register service.
	Register(application *Application)
}

// ConsoleCommand interface.
type ConsoleCommand interface {
	// GetCommand for the cli to register.
	GetCommand() cli.Command

	// Handle command.
	Handle(args cli.Args) error
}

// Config interface.
type Config interface {
	// Env returns current environment name.
	Env() string

	// Debug returs debug mode state.
	Debug() bool

	// Get value from config using dot-notation.
	Get(key string) interface{}

	// Set value to config using dot-notation.
	Set(key string, value interface{})
}
