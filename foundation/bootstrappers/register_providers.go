package bootstrappers

import (
	"github.com/lara-go/larago"
)

// RegisterProviders registers all service providers.
func RegisterProviders(application *larago.Application) error {
	application.Register(application.ApplicationServiceProvider())

	return nil
}
