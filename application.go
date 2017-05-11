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

	config       *ConfigRepository
	configLoader ConfigLoader

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

// SetConfig callback that loads config while bootstrapping.
func (app *Application) SetConfig(loader ConfigLoader) *Application {
	app.configLoader = loader

	return app
}

// ImportConfig of the application.
func (app *Application) ImportConfig() *Application {
	// Resolve config repository.
	app.config = &ConfigRepository{
		config: app.configLoader(),
	}

	// Save it to container.
	app.Instance(app.config, "config", (*Config)(nil))

	return app
}

// Config getter.
func (app *Application) Config() *ConfigRepository {
	return app.config
}

// Env checks if application works in this environment.
func (app *Application) Env(name string) bool {
	return app.config.Env() == name
}

// Register service.
func (app *Application) Register(providers ...ServiceProvider) {
	for _, provider := range providers {
		provider.Register(app)

		if method := reflect.ValueOf(provider).MethodByName("Boot"); method.IsValid() {
			if app.booted {
				app.bootProvider(method)
			} else {
				app.toBoot = append(app.toBoot, method)
			}
		}
	}
}

// Facade registers application facades.
func (app *Application) Facade(wrappers ...*Facade) {
	for _, wrapper := range wrappers {
		wrapper.Application = app
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
	app.catchSignals()

	defer app.exitHandler()

	// Get Kernel from container and handle request.
	var kernel Kernel
	app.Assign(&kernel)

	kernel.Handle()
}

func (app *Application) exitHandler() {
	var exitHandler ExitHandler
	app.Assign(&exitHandler)

	// Set defer panic handle.
	exitHandler.Defer()
}

func (app *Application) catchSignals() {
	var signalsHandler SignalsHandler
	app.Assign(&signalsHandler)

	signalsHandler.CatchInterrupt()
}
