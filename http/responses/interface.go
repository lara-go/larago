package responses

import net_http "net/http"

// Response interface.
type Response interface {
	// Headers get headers to send.
	Headers() map[string]string

	// WithHeader attaches header to response.
	WithHeader(name, value string) Response

	// WithCookies attaches cookies to response.
	WithCookies(cookie ...*net_http.Cookie) Response

	// Cookies returns set of cookies to send.
	Cookies() []*net_http.Cookie

	// SetStatus sets HTTP status.
	SetStatus(int) Response

	// Status returns HTTP status.
	Status() int

	// ContentType returns Content-Type header.
	ContentType() string

	// Body returns content.
	Body() []byte
}
