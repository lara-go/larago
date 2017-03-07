package responses

import net_http "net/http"

// Response interface.
type Response interface {

	// WithStatus sets HTTP status.
	WithStatus(status int) Response

	// WithHeader attaches header to response.
	WithHeader(name, value string) Response

	// WithCookies attaches cookies to response.
	WithCookies(cookie ...*net_http.Cookie) Response

	// Status returns HTTP status.
	Status() int

	// ContentType returns Content-Type header.
	ContentType() string

	// Headers get headers to send.
	Headers() map[string]string

	// Cookies returns set of cookies to send.
	Cookies() []*net_http.Cookie

	// Body returns content.
	Body() []byte
}
