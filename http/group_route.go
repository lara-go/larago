package http

// GroupRoute struct.
type GroupRoute struct {
	Path        string
	Middlewares []Middleware
}

// NewGroupRoute constructor.
func NewGroupRoute(path string) *GroupRoute {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	return &GroupRoute{
		Path: path,
	}
}

// Middleware sets middleware for GroupRoute.
func (r *GroupRoute) Middleware(middleware ...Middleware) *GroupRoute {
	r.Middlewares = middleware

	return r
}
