package responses

import "fmt"

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
