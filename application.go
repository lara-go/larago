package larago

import (
	"reflect"
	"strings"

	"github.com/lara-go/larago/container"
)

// Application struct.
type Application struct {
	*container.Container

	Version        string
	Name           string
	Description    string
	HomeDirectory  string
	DateTimeFormat string

	bootstrapped bool

	booted bool
	toBoot []reflect.Value

	providers []ServiceProvider

	config       Config
	configLoader func() Config

	commands []ConsoleCommand
}

// New constructor of the application.
func New() *Application {
	instance := &Application{
		Container: container.New(),

		HomeDirectory:  HomeDirectory,
		DateTimeFormat: DateTimeFormat,
	}

	instance.Container.Instance(instance)
	instance.Container.SetTagsResolver(instance)

	return instance
}

// SetConfigLoader callback that loads config while bootstrapping.
func (app *Application) SetConfigLoader(loader func() Config) *Application {
	app.configLoader = loader

	return app
}

// ConfigLoader returns callback that loads config while bootstrapping.
func (app *Application) ConfigLoader() func() Config {
	return app.configLoader
}

// SetConfig of the application.
func (app *Application) SetConfig(config Config) *Application {
	app.config = config
	app.Instance(config, "config", (*Config)(nil))

	return app
}

// Config getter.
func (app *Application) Config() Config {
	return app.config
}

// Env checks if application works in this environment.
func (app *Application) Env(name string) bool {
	return app.config.Env() == name
}

// Register service.
func (app *Application) Register(providers ...ServiceProvider) {
	var method reflect.Value

	for _, provider := range providers {
		provider.Register(app)

		method = reflect.ValueOf(provider).MethodByName("Boot")
		if method.IsValid() {
			if app.booted {
				app.bootProvider(method)
			} else {
				app.toBoot = append(app.toBoot, method)
			}
		}
	}
}

// Commands registers commands.
func (app *Application) Commands(commands ...ConsoleCommand) {
	for _, command := range commands {
		app.commands = append(app.commands, command)
	}
}

// GetCommands retrieve registered console commands.
func (app *Application) GetCommands() []ConsoleCommand {
	return app.commands
}

// Boot Application.
func (app *Application) Boot() error {
	if app.booted {
		return nil
	}

	for _, bootable := range app.toBoot {
		if err := app.bootProvider(bootable); err != nil {
			return err
		}
	}

	app.booted = true

	return nil
}

// Boot provider.
func (app *Application) bootProvider(boot reflect.Value) error {
	_, err := app.Call(boot)

	return err
}

// BootstrapWith runs bootstrappers one by one to populate and prepare application.
func (app *Application) BootstrapWith(boostrappers ...Bootstrapper) error {
	if app.bootstrapped {
		return nil
	}

	for _, bootstrapper := range boostrappers {
		if err := bootstrapper(app); err != nil {
			return err
		}
	}

	app.bootstrapped = true

	return nil
}

// ResolveTag resolves custom tags.
func (app *Application) ResolveTag(tag string, container *container.Container) interface{} {
	var resolved interface{}

	if strings.HasPrefix(tag, "Config.") {
		resolved = app.config.Get(strings.TrimPrefix(tag, "Config."))
	}

	return resolved
}

// Run application.
func (app *Application) Run() {
	defer app.panicHandler()

	var kernel Kernel
	app.Assign(&kernel)

	kernel.Handle()
}

// Handle panics.
func (app *Application) panicHandler() {
	var panicHandler PanicHandler
	app.Assign(&panicHandler)

	panicHandler.Defer()
}
