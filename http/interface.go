package http

import "github.com/lara-go/larago/http/responses"

// ErrorsHandlerContract for every handler to resolve errors during http calls..
type ErrorsHandlerContract interface {
	// Report error to logger.
	Report(err error)

	// Render error to return to the client.
	Render(request *Request, err error) responses.Response
}
