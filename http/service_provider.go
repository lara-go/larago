package http

import "github.com/lara-go/larago"

// ServiceProvider struct.
type ServiceProvider struct{}

// Register service.
func (p *ServiceProvider) Register(application *larago.Application) {
	application.Commands(
		&CommandDown{},
		&CommandUp{},
		&CommandServe{},
		&CommandRoutes{},
	)

	// Register server itself.
	p.registerRouter(application)
	p.registerErrorsHandler(application)
}

func (p *ServiceProvider) registerRouter(application *larago.Application) {
	application.Bind(NewRouter(), "router")
}

func (p *ServiceProvider) registerErrorsHandler(application *larago.Application) {
	application.Bind(&ErrorsHandler{}, (*ErrorsHandlerContract)(nil))
}
