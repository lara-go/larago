package foundation

import (
	"github.com/lara-go/larago"
	"github.com/lara-go/larago/errors"
	"github.com/lara-go/larago/events"
	"github.com/lara-go/larago/foundation/bootstrappers"
	"github.com/lara-go/larago/logger"
)

// Handle cli request.
func Handle(application *larago.Application) {
	// Handle all panics.
	defer makePanicHandler(application).Defer()

	// Handle request.
	makeKernel(application).
		UseBootstrappers(
			bootstrappers.DetectEnv,
			bootstrappers.RegisterProviders,
			bootstrappers.BootProviders,
		).
		WithApplication(application).
		Handle()
}

// MakeApplication start application and set default Service Providers.
func MakeApplication(name, version, description string) *larago.Application {
	application := larago.New()
	application.Name = name
	application.Version = version
	application.Description = description

	// Register panic handler.
	application.Bind(func() (errors.PanicHandlerInterface, error) {
		return &errors.PanicHandler{}, nil
	})

	// Register basic services providers.
	application.Register(
		&events.ServiceProvider{},
		&logger.ServiceProvider{},
	)

	return application
}

// Make kernel instance.
func makeKernel(application *larago.Application) *larago.Kernel {
	var kernel larago.Kernel
	application.Make(&kernel)

	return &kernel
}

// Make panic handler instance.
func makePanicHandler(application *larago.Application) errors.PanicHandlerInterface {
	var panicHander errors.PanicHandler
	application.Make(&panicHander)

	return &panicHander
}
