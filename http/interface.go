package http

import "github.com/lara-go/larago/http/responses"

// ErrorsHandlerContract for every handler to resolve errors during http calls..
type ErrorsHandlerContract interface {
	Report(err error)
	Render(request *Request, err error) responses.Response
}
