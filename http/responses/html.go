package responses

// HTML response.
type HTML struct {
	AbstractResponse

	html string
}

// NewHTML send plain text response.
func NewHTML(status int, html string) *HTML {
	response := &HTML{
		html: html,
	}
	response.SetStatus(status)

	return response
}

// ContentType returns ContentType header.
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
