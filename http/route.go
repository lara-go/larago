package http

import "github.com/julienschmidt/httprouter"

// Route struct.
type Route struct {
	Method      string
	Path        string
	Name        string
	Middlewares []Middleware
	Handler     interface{}
	Params      httprouter.Params
}

// NewRoute constructor.
func NewRoute(method string, path string) *Route {
	return &Route{
		Method: method,
		Path:   path,
	}
}

// As name.
func (r *Route) As(name string) *Route {
	r.Name = name

	return r
}

// Middleware sets middleware for route.
func (r *Route) Middleware(middleware ...Middleware) *Route {
	r.Middlewares = middleware

	return r
}

// Action to use.
func (r *Route) Action(handler interface{}) *Route {
	r.Handler = handler

	return r
}
