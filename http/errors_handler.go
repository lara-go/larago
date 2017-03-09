package http

import (
	"github.com/lara-go/larago/http/errors"
	"github.com/lara-go/larago/http/responses"
	"github.com/lara-go/larago/logger"
	"github.com/lara-go/larago/validation"
)

// ErrorsHandler to handle errors during http calls.
type ErrorsHandler struct {
	Logger                    *logger.Logger
	Debug                     bool `di:"Config.App.Debug"`
	ValidationErrorsConverter validation.ErrorsConverter
}

// Report error to logger.
func (h *ErrorsHandler) Report(err error) {
	// Convert error to HTTPError
	httpErr := h.makeHTTPError(err)

	if httpErr.WantsToBeReported() {
		log := h.Logger.From("lara-go/larago").WithContext(httpErr.Context)

		if httpErr.WantsToShowTrace() {
			log = log.WithTrace()
		}

		log.Warning("HTTP Error: %s", httpErr.Body.Message)
	}
}

// Render error to the client
func (h *ErrorsHandler) Render(request *Request, err error) responses.Response {
	var response responses.Response

	// Convert error to HTTPError
	httpErr := h.makeHTTPError(err)

	// Convert response due to what client wants: JSON / HTML / plain text.
	switch true {
	case request.WantsJSON():
		response = responses.NewJSON(httpErr.HTTPStatus, httpErr)
	case request.WantsHTML():
		response = responses.NewHTML(httpErr.HTTPStatus, httpErr.Body.Message)
	default:
		response = responses.NewText(httpErr.HTTPStatus, httpErr.Body.Message)
	}

	return response
}

// Make HTTPError instance from custom error.
func (h *ErrorsHandler) makeHTTPError(err error) *errors.HTTPError {
	if httpErr, ok := err.(*errors.HTTPError); ok {
		return httpErr
	}

	if vError, ok := err.(*validation.Error); ok {
		return h.makeValidationError(vError).WithContext(err)
	}

	return errors.UnknownError(err)
}

// Make http validation error.
func (h *ErrorsHandler) makeValidationError(err *validation.Error) *errors.HTTPError {
	httpError := errors.ValidationFailedHTTPError()

	if h.ValidationErrorsConverter != nil {
		httpError.WithMeta(h.ValidationErrorsConverter.ConvertValidationErrors(
			err.GetError(),
			err.GetValidator(),
		))
	} else {
		httpError.WithMeta(map[string]interface{}{"errors": err})
	}

	return httpError
}
