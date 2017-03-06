package http

import "github.com/lara-go/larago/http/errors"

const (
	jsonTag   = "json"
	schemaTag = "schema"
)

// RequestsValidator populates and validates
type RequestsValidator struct {
	ErrorsConverter ValidationErrorsConverter
}

// ValidateRequest validates form request.
func (rv *RequestsValidator) ValidateRequest(request *Request, validator SelfValidator) *errors.HTTPError {
	var tagName string
	var err error

	switch v := validator.(type) {
	case JSONRequestValidator:
		tagName = jsonTag

		request.ReadJSON(v)
		err = v.ValidateJSON()
	case FormRequestValidator:
		tagName = schemaTag

		request.ReadForm(v)
		err = v.ValidateForm()
	case QueryRequestValidator:
		tagName = schemaTag

		request.ReadQuery(v)
		err = v.ValidateQuery()
	case ParamsRequestValidator:
		tagName = schemaTag

		request.ReadParams(v)
		err = v.ValidateParams()
	}

	if err != nil {
		httpErr := errors.ValidationFailedHTTPError()
		httpErr.
			WithMeta(rv.ErrorsConverter.ConvertValidationErrors(err, tagName, validator)).
			WithContext(err)

		return httpErr
	}

	return nil
}
