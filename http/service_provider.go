package http

import (
	"github.com/lara-go/larago"
	"github.com/lara-go/larago/http/errors"
)

// ServiceProvider struct.
type ServiceProvider struct{}

// Register service.
func (p *ServiceProvider) Register(application *larago.Application) {
	application.Commands(
		&CommandDown{},
		&CommandUp{},
		&CommandServe{},
	)

	// Register server itself.
	p.registerRouter(application)
	p.registerErrorsHandler(application)
	p.registerRequestsValidator(application)
}

func (p *ServiceProvider) registerRouter(application *larago.Application) {
	application.Bind(NewRouter(), "router")
}

func (p *ServiceProvider) registerErrorsHandler(application *larago.Application) {
	application.Bind(&ErrorsHandler{}, (*ErrorsHandlerContract)(nil))
}

func (p *ServiceProvider) registerRequestsValidator(application *larago.Application) {
	application.Bind(&RequestsValidator{})
	application.Bind(&errors.OzzoValidationErrorsConverter{}, (*ValidationErrorsConverter)(nil))
}
