package responses

import net_http "net/http"

// AbstractResponse handles response.
type AbstractResponse struct {
	Response

	status  int
	headers map[string]string
	cookies []*net_http.Cookie
}

// Status returns HTTP status.
func (r *AbstractResponse) Status() int {
	return r.status
}

// SetStatus sets HTTP status.
func (r *AbstractResponse) SetStatus(status int) Response {
	r.status = status

	return r
}

// WithHeader attaches header to response.
func (r *AbstractResponse) WithHeader(name, value string) Response {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}

	r.headers[name] = value

	return r
}

// Headers get headers to send.
func (r *AbstractResponse) Headers() map[string]string {
	return r.headers
}

// WithCookies attaches cookies to response.
func (r *AbstractResponse) WithCookies(cookie ...*net_http.Cookie) Response {
	r.cookies = cookie

	return r
}

// Cookies returns set of cookies to send.
func (r *AbstractResponse) Cookies() []*net_http.Cookie {
	return r.cookies
}
