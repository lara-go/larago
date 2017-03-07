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
func (r *AbstractResponse) SetStatus(status int) {
	r.status = status
}

// SetHeader attaches header to response.
func (r *AbstractResponse) SetHeader(name, value string) {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}

	r.headers[name] = value
}

// Headers returns set of additional headers to send.
func (r *AbstractResponse) Headers() map[string]string {
	return r.headers
}

// SetCookies attaches cookies to response.
func (r *AbstractResponse) SetCookies(cookie []*net_http.Cookie) {
	r.cookies = cookie
}

// Cookies returns set of cookies to send.
func (r *AbstractResponse) Cookies() []*net_http.Cookie {
	return r.cookies
}

// ContentType returns Content-Type header.
func (r *AbstractResponse) ContentType() string {
	return "text/plain"
}

// Body returns content.
func (r *AbstractResponse) Body() []byte {
	return []byte("")
}

// String returns response body as string.
func (r *AbstractResponse) String() string {
	return ""
}
