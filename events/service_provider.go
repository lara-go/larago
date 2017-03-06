package events

import (
	"github.com/asaskevich/EventBus"
	"github.com/lara-go/larago"
)

// ServiceProvider for events service.
type ServiceProvider struct{}

// Register service.
func (p *ServiceProvider) Register(application *larago.Application) {
	application.Bind(EventBus.New(), "events")
}
