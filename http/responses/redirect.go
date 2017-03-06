package responses

import (
	"net/url"
	"strings"
)

// Redirect response.
type Redirect struct {
	AbstractResponse

	location string
	query    url.Values

	route  string
	params map[string]string
}

// NewRedirect send redirect to route or custom url.
func NewRedirect(status int) *Redirect {
	response := &Redirect{
		query: make(url.Values),
	}
	response.SetStatus(status)

	return response
}

// To sets real url to redirect.
func (r *Redirect) To(location string) *Redirect {
	r.location = location

	return r
}

// Route sets route name to redirect to + params if it has any.
func (r *Redirect) Route(name string, params map[string]string) *Redirect {
	r.route = name
	r.params = params

	return r
}

// WithQuery adds key-value to url.
// Appends values if there is already one.
func (r *Redirect) WithQuery(key, value string) *Redirect {
	r.query.Add(key, value)

	return r
}

// GetRoute returns route name to redirect to.
func (r *Redirect) GetRoute() string {
	return r.route
}

// GetLocation formats location to redirect to and returns it.
func (r *Redirect) GetLocation() string {
	if !strings.Contains(r.location, ":") {
		return r.location
	}

	// Replace :key with value from params.
	result := r.location
	for key, value := range r.params {
		result = strings.Replace(result, ":"+key, value, -1)
	}

	return result
}
