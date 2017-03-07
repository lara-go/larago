package responses

import (
	"encoding/json"
	"fmt"
	net_http "net/http"
)

// JSON response.
type JSON struct {
	AbstractResponse

	data interface{}
}

// NewJSON semd JSON formatted response.
func NewJSON(status int, data interface{}) *JSON {
	response := &JSON{
		data: data,
	}
	response.SetStatus(status)

	return response
}

// WithStatus sets HTTP status.
func (r *JSON) WithStatus(status int) Response {
	r.SetStatus(status)

	return r
}

// WithHeader attaches header to response.
func (r *JSON) WithHeader(name, value string) Response {
	r.SetHeader(name, value)

	return r
}

// WithCookies attaches cookies to response.
func (r *JSON) WithCookies(cookie ...*net_http.Cookie) Response {
	r.SetCookies(cookie)

	return r
}

// ContentType returns Content-Type header.
func (r *JSON) ContentType() string {
	return "application/json"
}

// GetData returns passed data.
func (r *JSON) GetData() interface{} {
	return &r.data
}

// Body returns content.
func (r *JSON) Body() []byte {
	body, err := json.Marshal(r.data)
	if err != nil {
		panic(fmt.Errorf("Invalid data to transform to JSON reponse: %s", body))
	}

	return body
}

// String returns response body as string.
func (r *JSON) String() string {
	return string(r.Body())
}
