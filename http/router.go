package http

import (
	"fmt"
	net_http "net/http"

	"github.com/lara-go/larago"
	"github.com/lara-go/larago/container"
	"github.com/lara-go/larago/http/errors"
	"github.com/lara-go/larago/http/responses"
	"github.com/lara-go/larago/logger"

	"github.com/julienschmidt/httprouter"
)

// Router struct.
type Router struct {
	// Dependencies
	Container     *container.Container
	Logger        *logger.Logger
	ErrorsHandler ErrorsHandlerContract
	Config        larago.Config

	// Basic httprouter.
	router *httprouter.Router

	// Global middleware to run.
	middleware []Middleware

	// Map of path:*Route.
	routes []*Route

	// Map of alias:*Route.
	aliases map[string]*Route
}

// NewRouter constructor.
func NewRouter() *Router {
	router := &Router{
		aliases: make(map[string]*Route),
	}

	router.router = httprouter.New()
	router.SetNotFoundHandler(router.handleNotFound)
	router.SetMethodNotAllowedHandler(router.handleMethodNotAllowed)

	return router
}

// GET request.
func (r *Router) GET(path string) *Route {
	return r.addRoute(net_http.MethodGet, path)
}

// POST request.
func (r *Router) POST(path string) *Route {
	return r.addRoute(net_http.MethodPost, path)
}

// PUT request.
func (r *Router) PUT(path string) *Route {
	return r.addRoute(net_http.MethodPut, path)
}

// PATCH request.
func (r *Router) PATCH(path string) *Route {
	return r.addRoute(net_http.MethodPatch, path)
}

// DELETE request.
func (r *Router) DELETE(path string) *Route {
	return r.addRoute(net_http.MethodDelete, path)
}

// Add route helper.
func (r *Router) addRoute(method, path string) *Route {
	route := NewRoute(method, path)
	r.routes = append(r.routes, route)

	return route
}

// Middleware sets global middleware to run over every request.
func (r *Router) Middleware(middleware ...Middleware) *Router {
	r.middleware = append(r.middleware, middleware...)

	return r
}

// Bootstrap router to be ready to handle requests.
func (r *Router) Bootstrap() *Router {
	for _, route := range r.routes {
		r.setHTTPRoute(route)
	}

	return r
}

// Listen to requests.
// Do not forget to run Bootstrap in order to prepare and set routes.
func (r *Router) Listen(listen string) error {
	// Set httprouter routes once before server starts listen.
	r.Logger.Info(
		"Serving app at %s with %s environment and debug mode %t.",
		listen,
		r.Config.Env(),
		r.Config.Debug(),
	)

	return net_http.ListenAndServe(listen, r.router)
}

// Set route to httprouter.
func (r *Router) setHTTPRoute(route *Route) {
	// Call httprouter.
	switch route.Method {
	case net_http.MethodGet:
		r.router.GET(route.Path, r.wrapHandlers(route))
	case net_http.MethodPost:
		r.router.POST(route.Path, r.wrapHandlers(route))
	case net_http.MethodPut:
		r.router.PUT(route.Path, r.wrapHandlers(route))
	case net_http.MethodPatch:
		r.router.PATCH(route.Path, r.wrapHandlers(route))
	case net_http.MethodDelete:
		r.router.DELETE(route.Path, r.wrapHandlers(route))
	}

	// Append alias if there is one.
	if route.Name != "" {
		r.aliases[route.Name] = route
	}
}

// Wrap handlers for httprouter.
func (r *Router) wrapHandlers(route *Route) httprouter.Handle {
	// Merge global middleware with route ones.
	middleware := append(r.middleware, route.Middlewares...)

	// Return httprouter handler.
	return func(w net_http.ResponseWriter, req *net_http.Request, ps httprouter.Params) {
		request := NewRequest(req)
		request.Route = route
		request.Route.Params = ps

		// Save request to container.
		r.Container.Instance(request)

		// Handle panics during pipeline.
		defer r.panicHandler(w, request)

		// Run request through all middleware chaining one by one.
		// then dispatch action handler itself, obtain response
		// and lift it back.
		result := NewPipeline(r.Container).
			Send(request).
			Through(middleware).
			Then(r.dispatchRequest)

		// Send response.
		switch response := result.(type) {
		case *responses.Redirect:
			r.sendRedirect(w, request, response)
		default:
			r.sendResponse(w, request, response)
		}
	}
}

// Panic handler.
func (r *Router) panicHandler(w net_http.ResponseWriter, request *Request) {
	if re := recover(); re != nil {
		var err error
		var ok bool

		if err, ok = re.(error); !ok {
			err = fmt.Errorf("%s", re)
		}

		r.sendResponse(w, request, r.formatError(request, err))
	}
}

// Final pipeline callback.
// Dispatches route handler and returns Response.
func (r *Router) dispatchRequest(request *Request) responses.Response {
	action := r.Container.Wrap(request.Route.Handler)

	params := make([]interface{}, 0)
	for _, param := range request.Route.Params {
		params = append(params, param.Value)
	}

	result, err := action(params...)
	if err != nil {
		return r.formatError(request, err)
	}

	if response, ok := result.(responses.Response); ok {
		return response
	}

	return r.formatResponse(request, result)
}

// Formats every type into suitable Response.
func (r *Router) formatResponse(request *Request, result interface{}) responses.Response {
	switch v := result.(type) {
	case string:
		return responses.NewText(200, fmt.Sprintf("%s", v))
	case int:
		return responses.NewText(200, fmt.Sprintf("%d", v))
	case float64:
		return responses.NewText(200, fmt.Sprintf("%f", v))
	case bool:
		return responses.NewText(200, fmt.Sprintf("%t", v))
	case error:
		return r.formatError(request, v)
	default:
		return responses.NewJSON(200, v)
	}
}

// Format error.
func (r *Router) formatError(request *Request, err error) responses.Response {
	r.ErrorsHandler.Report(err)

	return r.ErrorsHandler.Render(request, err)
}

// Send redirect response via native http.Redirect.
func (r *Router) sendRedirect(w net_http.ResponseWriter, request *Request, redirect *responses.Redirect) {
	// Find alias if there is one.
	alias := redirect.GetRoute()
	if r, ok := r.aliases[alias]; ok {
		redirect.To(r.Path)
	}

	// Send redirect.
	net_http.Redirect(w, request.BaseRequest(), redirect.GetLocation(), redirect.Status())
}

// Send common Response to client.
func (r *Router) sendResponse(w net_http.ResponseWriter, request *Request, response responses.Response) {
	// Send content type.
	w.Header().Set("content-type", response.ContentType()+"; charset=utf-8")

	// Send additional headers.
	for name, value := range response.Headers() {
		w.Header().Set(name, value)
	}

	// Send cookies.
	for _, cookie := range response.Cookies() {
		net_http.SetCookie(w, cookie)
	}

	// Send status.
	w.WriteHeader(response.Status())

	// Send body.
	w.Write(response.Body())
}

// GetHTTPRouter returns httprouter instance.
func (r *Router) GetHTTPRouter() *httprouter.Router {
	return r.router
}

// SetNotFoundHandler allowes to set custom NotFound handler.
func (r *Router) SetNotFoundHandler(handler net_http.HandlerFunc) *Router {
	r.router.NotFound = handler

	return r
}

// SetMethodNotAllowedHandler allowes to set custom MethodNotAllowed handler.
func (r *Router) SetMethodNotAllowedHandler(handler net_http.HandlerFunc) *Router {
	r.router.MethodNotAllowed = handler

	return r
}

// Custom handler for NotFound errors.
func (r *Router) handleNotFound(w net_http.ResponseWriter, req *net_http.Request) {
	request := NewRequest(req)

	r.sendResponse(w, request, r.ErrorsHandler.Render(request, errors.NotFoundHTTPError()))
}

// Custom handler for MethodNotAllowed errors.
func (r *Router) handleMethodNotAllowed(w net_http.ResponseWriter, req *net_http.Request) {
	request := NewRequest(req)

	r.sendResponse(w, request, r.ErrorsHandler.Render(request, errors.MethodNotAllowedHTTPError()))
}
