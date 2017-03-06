package http

import (
	"github.com/lara-go/larago/container"
	"github.com/lara-go/larago/http/responses"
)

// Middleware interface.
type Middleware interface {
	Handle(request *Request, next Handler) responses.Response
}

// Handler function.
type Handler func(request *Request) responses.Response

// Pipeline struct.
type Pipeline struct {
	container  container.Interface
	request    *Request
	middleware []Middleware
}

// NewPipeline constructor.
func NewPipeline(container container.Interface) *Pipeline {
	return &Pipeline{
		container: container,
	}
}

// Send value througn every middleware.
func (p *Pipeline) Send(request *Request) *Pipeline {
	p.request = request

	return p
}

// Through what pipes run passable.
func (p *Pipeline) Through(middleware []Middleware) *Pipeline {
	p.middleware = middleware

	return p
}

// Then run pipeline with the last handler.
func (p *Pipeline) Then(last Handler) (response responses.Response) {
	out := p.getInitialHandler(last)
	for _, middleware := range p.reverseMiddleware(p.middleware) {
		out = p.getHander(middleware, out)
	}

	return out(p.request)
}

// Make handler from the middleware.
func (p *Pipeline) getHander(middleware Middleware, next Handler) Handler {
	return func(request *Request) responses.Response {
		p.container.Make(middleware)

		return middleware.Handle(request, next)
	}
}

// Make initial handler.
func (p *Pipeline) getInitialHandler(next Handler) Handler {
	return func(request *Request) responses.Response {
		return next(p.request)
	}
}

// Reverse middleware to run them from the last one to the first.
func (p *Pipeline) reverseMiddleware(middleware []Middleware) []Middleware {
	for i, j := 0, len(middleware)-1; i < j; i, j = i+1, j-1 {
		middleware[i], middleware[j] = middleware[j], middleware[i]
	}

	return middleware
}
