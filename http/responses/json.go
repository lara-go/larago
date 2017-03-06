package responses

import (
	"encoding/json"
	"fmt"
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

// ContentType returns ContentType header.
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
