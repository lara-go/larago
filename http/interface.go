package http

import "github.com/lara-go/larago/http/responses"

// Middleware interface.
type Middleware interface {
	// Handle request.
	Handle(request *Request, next Handler) responses.Response
}

// Handler function.
type Handler func(request *Request) responses.Response

// ErrorsHandlerContract for every handler to resolve errors during http calls..
type ErrorsHandlerContract interface {
	// Report error to logger.
	Report(err error)

	// Render error to return to the client.
	Render(request *Request, err error) responses.Response
}

// BindingCallback is a function to resolve binded param value.
type BindingCallback func(param string) (interface{}, error)

// ArgsInjector uses as an external source of arguments that can be injected in the route action handler.
type ArgsInjector interface {
	// Injects custom params to for the action.
	Inject(params []interface{}, request *Request) ([]interface{}, error)
}

// JSONRequestValidator interface for json requests.
type JSONRequestValidator interface {
	ValidateJSON() error
}

// FormRequestValidator interface for form requests.
type FormRequestValidator interface {
	ValidateForm() error
}

// QueryRequestValidator interface for query.
type QueryRequestValidator interface {
	ValidateQuery() error
}

// ParamsRequestValidator interface for url params.
type ParamsRequestValidator interface {
	ValidateParams() error
}
