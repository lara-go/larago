package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/lara-go/larago/validation"
)

// Route struct.
type Route struct {
	Method      string
	Path        string
	Name        string
	Middlewares []Middleware
	Handler     interface{}
	Params      httprouter.Params
	ToValidate  []validation.SelfValidator
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
	r.Middlewares = append(r.Middlewares, middleware...)

	return r
}

// Action to use.
func (r *Route) Action(handler interface{}) *Route {
	r.Handler = handler

	return r
}

// Validate request.
func (r *Route) Validate(requests ...validation.SelfValidator) *Route {
	r.ToValidate = requests

	return r
}

// Extend current route with data from group.
func (r *Route) extendWithGroup(group *GroupRoute) {
	if len(group.Middlewares) > 0 {
		// Merge group middleware with route ones.
		r.Middlewares = append(group.Middlewares, r.Middlewares...)
	}

	// Merge path only if route path not equal to /
	if r.Path != "/" {
		r.Path = group.Path + r.Path
	} else {
		// Set path to group path only if it is not /
		if group.Path != "/" {
			r.Path = group.Path
		}
	}
}
