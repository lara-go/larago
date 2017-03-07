package responses

import (
	"fmt"
	net_http "net/http"
)

// HTML response.
type HTML struct {
	AbstractResponse

	html string
}

// NewHTML send plain text response.
func NewHTML(status int, html string, a ...interface{}) *HTML {
	response := &HTML{
		html: fmt.Sprintf(html, a...),
	}
	response.SetStatus(status)

	return response
}

// WithStatus sets HTTP status.
func (r *HTML) WithStatus(status int) Response {
	r.SetStatus(status)

	return r
}

// WithHeader attaches header to response.
func (r *HTML) WithHeader(name, value string) Response {
	r.SetHeader(name, value)

	return r
}

// WithCookies attaches cookies to response.
func (r *HTML) WithCookies(cookie ...*net_http.Cookie) Response {
	r.SetCookies(cookie)

	return r
}

// ContentType returns Content-Type header.
func (r *HTML) ContentType() string {
	return "text/html"
}

// Body returns content.
func (r *HTML) Body() []byte {
	return []byte(r.html)
}

// String returns response body as string.
func (r *HTML) String() string {
	return string(r.Body())
}
