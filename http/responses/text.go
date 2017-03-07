package responses

import (
	"fmt"
	net_http "net/http"
)

// Text response.
type Text struct {
	AbstractResponse

	text string
}

// NewText send plain text response.
func NewText(status int, text string, a ...interface{}) *Text {
	response := &Text{
		text: fmt.Sprintf(text, a...),
	}
	response.SetStatus(status)

	return response
}

// WithStatus sets HTTP status.
func (r *Text) WithStatus(status int) Response {
	r.SetStatus(status)

	return r
}

// WithHeader attaches header to response.
func (r *Text) WithHeader(name, value string) Response {
	r.SetHeader(name, value)

	return r
}

// WithCookies attaches cookies to response.
func (r *Text) WithCookies(cookie ...*net_http.Cookie) Response {
	r.SetCookies(cookie)

	return r
}

// ContentType returns Content-Type header.
func (r *Text) ContentType() string {
	return "text/plain"
}

// Body returns content.
func (r *Text) Body() []byte {
	return []byte(r.text)
}

// String returns response body as string.
func (r *Text) String() string {
	return string(r.Body())
}
