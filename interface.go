package larago

import "github.com/urfave/cli"

// Bootstrapper function.
type Bootstrapper func(application *Application) error

// ServiceProvider interface.
type ServiceProvider interface {
	// Register service.
	Register(application *Application)
}

// PanicHandler provides interface to handle panic throws.
type PanicHandler interface {
	// Defer handles panics.
	Defer()
}

// Kernel interface.
type Kernel interface {
	// Handle console command.
	Handle()

	// SetBootstrappers sets Application bootstrappers.
	SetBootstrappers(bootstrappers ...Bootstrapper)
}

// ConsoleCommand interface.
type ConsoleCommand interface {
	// GetCommand for the cli to register.
	GetCommand() cli.Command

	// Handle command.
	Handle(args cli.Args) error
}

// ConfigLoader is a callback function for lazy loading application config.
type ConfigLoader func() Config

// Config interface.
type Config interface {
	// Env returns current environment name.
	Env() string

	// Debug returs debug mode state.
	Debug() bool
}
