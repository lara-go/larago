package http

import (
	base_http "net/http"
	"strings"
)

// Request handles http request.
type Request struct {
	*base_http.Request
	Route *Route
}

// Param returns route param.
func (r *Request) Param(name string) string {
	return r.Route.Params.ByName(name)
}

// WantsJSON checks if client wants JSON answer.
func (r *Request) WantsJSON() bool {
	return strings.Contains(r.Header.Get("accept"), "application/json")
}

// WantsHTML checks if client wants HTML answer.
func (r *Request) WantsHTML() bool {
	return strings.Contains(r.Header.Get("accept"), "text/html")
}

// WantsPlainText checks if client wants plain text answer.
func (r *Request) WantsPlainText() bool {
	return strings.Contains(r.Header.Get("accept"), "text/plain")
}
