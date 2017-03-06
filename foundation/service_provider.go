package foundation

import (
	"github.com/lara-go/larago"
	"github.com/lara-go/larago/foundation/console"
)

// ServiceProvider .
type ServiceProvider struct{}

// Register service.
func (p *ServiceProvider) Register(application *larago.Application) {
	application.Commands(
		&console.CommandEnv{},
		&console.CommandMakeCommand{},
		&console.CommandMakeMiddleware{},
		&console.CommandMakeModel{},
		&console.CommandMakeProvider{},
	)
}
