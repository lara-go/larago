package foundation

import (
	"github.com/lara-go/larago"
	"github.com/lara-go/larago/cli"
	"github.com/lara-go/larago/errors"
	"github.com/lara-go/larago/events"
	"github.com/lara-go/larago/logger"
)

// MakeApplication start application and set default Service Providers.
func MakeApplication(name, version, description string) *larago.Application {
	application := larago.New()
	application.Name = name
	application.Version = version
	application.Description = description

	// Register basic services providers.
	application.Register(
		&events.ServiceProvider{},
		&logger.ServiceProvider{},
	)

	// Register panic handler.
	application.Bind(&errors.PanicHandler{}, (*larago.PanicHandler)(nil))

	// Register CLI Kernel.
	application.Bind(cli.NewKernel(), (*larago.Kernel)(nil))

	return application
}
