package responses

import "fmt"

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
