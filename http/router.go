package http

import (
	"fmt"
	net_http "net/http"

	"github.com/asaskevich/EventBus"
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
	Events        *EventBus.EventBus
	ErrorsHandler ErrorsHandlerInterface
	Config        larago.Config

	// Basic httprouter.
	router *httprouter.Router

	// Global middleware to run.
	middleware []Middleware

	// Map of path:*Route.
	routes []*Route

	// Map of alias:*Route.
	aliases map[string]*Route

	// Stack of group routes.
	groupsStack []*GroupRoute

	// Map of bindings and their callbacks.
	bindings map[string]BindingCallback

	// Action handler's arguments injectors.
	argsInjectors []ArgsInjector
}

// NewRouter constructor.
func NewRouter() *Router {
	router := &Router{
		aliases:     make(map[string]*Route),
		bindings:    make(map[string]BindingCallback),
		groupsStack: make([]*GroupRoute, 0),
		argsInjectors: []ArgsInjector{
			&RouteParamsInjector{},
		},
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

// Group routes.
func (r *Router) Group(path string, callback func(), middleware ...Middleware) {
	group := NewGroupRoute(path).Middleware(middleware...)

	// Add group to the begining of the stack.
	r.groupsStack = append([]*GroupRoute{group}, r.groupsStack...)

	// Call group callback.
	callback()

	// Remove group from the stack.
	r.groupsStack = r.groupsStack[:r.groupsStackLen()-1]
}

func (r *Router) inGroup() bool {
	return r.groupsStackLen() != 0
}

func (r *Router) groupsStackLen() int {
	return len(r.groupsStack)
}

// Add route helper.
func (r *Router) addRoute(method, path string) *Route {
	route := NewRoute(method, path)

	for _, group := range r.groupsStack {
		route.extendWithGroup(group)
	}

	r.routes = append(r.routes, route)

	return route
}

// Middleware sets global middleware to run over every request.
func (r *Router) Middleware(middleware ...Middleware) *Router {
	r.middleware = append(r.middleware, middleware...)

	return r
}

// Bind route param to the specified callback return value.
func (r *Router) Bind(param string, callback BindingCallback) {
	r.bindings[param] = callback
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
		request.Params = ps

		// Handle panics during pipeline.
		defer r.panicHandler(w, request)

		// Save request to container.
		r.Container.Instance(request)

		// Run request through all middleware chaining one by one.
		// then dispatch action handler itself, obtain response
		// and lift it back.
		response := NewPipeline(r.Container).
			Send(request).
			Through(middleware).
			Then(r.dispatchRequest)

		r.send(response, request, w)
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

		r.send(r.formatErrorResponse(request, err), request, w)
	}
}

// Substitute bindings from request url params.
func (r *Router) substituteBindings(request *Request) {
	for _, param := range request.Params {
		if callback, ok := r.bindings[param.Key]; ok {
			result, err := callback(param.Value)
			if err != nil {
				panic(err)
			}

			request.Bindings = append(request.Bindings, result)
		}
	}
}

// Final pipeline callback.
// Dispatches route handler and returns Response.
func (r *Router) dispatchRequest(request *Request) responses.Response {
	action := r.Container.Wrap(request.Route.Handler)

	// Substitute bindings.
	r.substituteBindings(request)

	// Prepare params to pass to the action.
	params := request.Bindings
	for _, injector := range r.argsInjectors {
		var err error
		params, err = injector.Inject(params, request)

		if err != nil {
			return r.formatErrorResponse(request, err)
		}
	}

	// Dispatch route action.
	result, err := action(params...)
	if err != nil {
		return r.formatErrorResponse(request, err)
	}

	// Return appopriate response.
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
		return responses.NewText(200, fmt.Sprintf("%e", v))
	case bool:
		return responses.NewText(200, fmt.Sprintf("%t", v))
	case error:
		return r.formatErrorResponse(request, v)
	default:
		return responses.NewJSON(200, v)
	}
}

// Format error.
func (r *Router) formatErrorResponse(request *Request, err error) responses.Response {
	r.ErrorsHandler.Report(err)

	return r.ErrorsHandler.Render(request, err)
}

// Send abstract response to the client.
func (r *Router) send(response responses.Response, request *Request, w net_http.ResponseWriter) {
	// Send response.
	switch resp := response.(type) {
	case *responses.Redirect:
		r.sendRedirect(resp, request, w)
	default:
		r.sendResponse(resp, request, w)
	}
}

// Send redirect response via native http.Redirect.
func (r *Router) sendRedirect(redirect *responses.Redirect, request *Request, w net_http.ResponseWriter) {
	// Find alias if there is one and set its real path as a redirect location.
	alias := redirect.GetRoute()
	if route, ok := r.aliases[alias]; ok {
		redirect.To(route.Path)
	}

	// Send redirect.
	net_http.Redirect(w, request.BaseRequest(), redirect.GetLocation(), redirect.Status())
}

// Send common response.
func (r *Router) sendResponse(response responses.Response, request *Request, w net_http.ResponseWriter) {
	if r.Events != nil {
		r.Events.Publish("router:request-handled", request, response)
	}

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

// GetRoutes returns all registered routes.
func (r *Router) GetRoutes() []*Route {
	return r.routes
}

// GetMiddleware returns all glbally registered middleware.
func (r *Router) GetMiddleware() []Middleware {
	return r.middleware
}

// SetArgsInjectors sets additional components that can inject more custom arguments to route action handler.
func (r *Router) SetArgsInjectors(injectors ...ArgsInjector) {
	for _, injector := range injectors {
		r.Container.Make(injector)

		r.argsInjectors = append(r.argsInjectors, injector)
	}
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

	r.sendResponse(r.ErrorsHandler.Render(request, errors.NotFoundHTTPError()), request, w)
}

// Custom handler for MethodNotAllowed errors.
func (r *Router) handleMethodNotAllowed(w net_http.ResponseWriter, req *net_http.Request) {
	request := NewRequest(req)

	r.sendResponse(r.ErrorsHandler.Render(request, errors.MethodNotAllowedHTTPError()), request, w)
}
