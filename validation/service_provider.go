package validation

import "github.com/lara-go/larago"

// ServiceProvider struct.
type ServiceProvider struct{}

// Register service.
func (p *ServiceProvider) Register(application *larago.Application) {
	application.Bind(&OzzoErrorsConverter{}, (*ErrorsConverter)(nil))
}
