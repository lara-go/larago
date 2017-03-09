package http

import "github.com/lara-go/larago/http/errors"

const (
	jsonTag   = "json"
	schemaTag = "schema"
)

// RequestsInjector populates and validates
type RequestsInjector struct {
	ErrorsConverter ValidationErrorsConverter
}

// Inject custom params to for the action.
func (rv *RequestsInjector) Inject(params []interface{}, request *Request) ([]interface{}, error) {
	// Validate request.
	for _, validator := range request.Route.ToValidate {
		if err := rv.validateRequest(request, validator); err != nil {
			return nil, err
		}

		params = append(params, validator)
	}

	return params, nil
}

// Validates form request.
func (rv *RequestsInjector) validateRequest(request *Request, validator SelfValidator) *errors.HTTPError {
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
