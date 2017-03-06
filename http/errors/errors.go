package errors

import "net/http"

// InternalServerErrorHTTPError error.
func InternalServerErrorHTTPError() *HTTPError {
	return &HTTPError{
		Body: Body{
			ID:      "internal_server_error",
			Message: "The server encountered an internal error or misconfiguration and was unable to complete your request.",
		},
		HTTPStatus:   http.StatusInternalServerError,
		ShouldReport: true,
		HasTrace:     true,
	}
}

// NotFoundHTTPError error.
func NotFoundHTTPError() *HTTPError {
	return &HTTPError{
		Body: Body{
			ID:      "not_found",
			Message: "Requested object not found.",
		},
		HTTPStatus: http.StatusNotFound,
	}
}

// ForbiddenHTTPError error.
func ForbiddenHTTPError() *HTTPError {
	return &HTTPError{
		Body: Body{
			ID:      "forbidden",
			Message: "You don't have permissions to perform this request.",
		},
		HTTPStatus: http.StatusForbidden,
	}
}

// UnauthorizedHTTPError error.
func UnauthorizedHTTPError() *HTTPError {
	return &HTTPError{
		Body: Body{
			ID:      "invalid_credentials",
			Message: "Sent credentials are invalid.",
		},
		HTTPStatus: http.StatusUnauthorized,
	}
}

// BadRequestHTTPError error.
func BadRequestHTTPError() *HTTPError {
	return &HTTPError{
		Body: Body{
			ID:      "bad_request",
			Message: "The server cannot process the request due to its malformed syntax.",
		},
		HTTPStatus: http.StatusBadRequest,
	}
}

// MethodNotAllowedHTTPError error.
func MethodNotAllowedHTTPError() *HTTPError {
	return &HTTPError{
		Body: Body{
			ID:      "method_not_allowed",
			Message: "Method you requesting is not allowed for this endpoint.",
		},
		HTTPStatus: http.StatusMethodNotAllowed,
	}
}

// ValidationFailedHTTPError error.
func ValidationFailedHTTPError() *HTTPError {
	return &HTTPError{
		Body: Body{
			ID:      "validation_failed",
			Message: "Validation failed.",
		},
		HTTPStatus: http.StatusUnprocessableEntity,
	}
}

// ValidationError for single field.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a set of messages for every invalid field.
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}
