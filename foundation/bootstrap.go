package foundation

import (
	"github.com/lara-go/larago"
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

	return application
}
