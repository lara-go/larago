package errors

import "net/http"

// Body of the error.
// Contains basic error info.
type Body struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// HTTPError custom API error.
type HTTPError struct {
	Body         Body        `json:"error"`
	Meta         interface{} `json:"meta,omitempty"`
	HTTPStatus   int         `json:"-"`
	Context      interface{} `json:"-"`
	ShouldReport bool        `json:"-"`
	HasTrace     bool        `json:"-"`
}

// UnknownError generate HTTP error from unknown error.
func UnknownError(err error) *HTTPError {
	return &HTTPError{
		Body: Body{
			ID:      "internal_server_error",
			Message: err.Error(),
		},
		HTTPStatus:   http.StatusInternalServerError,
		ShouldReport: true,
		HasTrace:     true,
	}
}

// Error returns error message.
func (e *HTTPError) Error() string {
	return e.Body.Message
}

// WithContext to the error.
// Only logged.
func (e *HTTPError) WithContext(context interface{}) *HTTPError {
	e.Context = context

	return e
}

// WithMeta to the error.
// Will be returned to the user.
func (e *HTTPError) WithMeta(meta interface{}) *HTTPError {
	e.Meta = meta

	return e
}

// WithTrace enables stacktrace reporting in logs.
func (e *HTTPError) WithTrace() *HTTPError {
	e.HasTrace = true

	return e
}

// Report error to the logs.
func (e *HTTPError) Report() *HTTPError {
	e.ShouldReport = true

	return e
}

// WantsToBeReported in logs or not.
func (e *HTTPError) WantsToBeReported() bool {
	return e.ShouldReport
}

// WantsToShowTrace in logs or not.
func (e *HTTPError) WantsToShowTrace() bool {
	return e.HasTrace
}
