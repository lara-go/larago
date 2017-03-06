package responses

// AbstractResponse handles response.
type AbstractResponse struct {
	Response

	status  int
	headers map[string]string
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

// WithHeader sets header to send.
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
