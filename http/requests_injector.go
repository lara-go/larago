package http

import "github.com/lara-go/larago/validation"

// RequestsInjector populates and validates
type RequestsInjector struct{}

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
func (rv *RequestsInjector) validateRequest(request *Request, validator validation.SelfValidator) *validation.Error {
	var err error

	switch v := validator.(type) {
	case JSONRequestValidator:
		request.ReadJSON(v)
		err = v.ValidateJSON()
	case FormRequestValidator:
		request.ReadForm(v)
		err = v.ValidateForm()
	case QueryRequestValidator:
		request.ReadQuery(v)
		err = v.ValidateQuery()
	case ParamsRequestValidator:
		request.ReadParams(v)
		err = v.ValidateParams()
	}

	if err != nil {
		return validation.NewError(err, validator)
	}

	return nil
}
