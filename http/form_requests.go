package http

import (
	"reflect"

	"git.acronis.com/ci/ci-2x-ipn/app/support/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/lara-go/larago/http/errors"
)

const (
	jsonTag   = "json"
	schemaTag = "schema"
)

// ValidateRequest validates form request.
func ValidateRequest(request *Request, validator SelfValidator) *errors.HTTPError {
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
		httpErr.WithMeta(convertValidationErrors(err, tagName, validator)).WithContext(err)

		return httpErr
	}

	return nil
}

// Convert errors to field - message format.
func convertValidationErrors(err error, tagName string, validator SelfValidator) *errors.ValidationErrors {
	fails := err.(validation.Errors)
	validationErrors := &errors.ValidationErrors{}

	reflection := reflect.ValueOf(validator).Elem().Type()

	for field, message := range fails {
		validationErrors.Errors = append(validationErrors.Errors, errors.ValidationError{
			Field:   normalizeFieldName(reflection, tagName, field),
			Message: normalizeMessage(message),
		})
	}

	return validationErrors
}

// Convert request attribute name to normal.
func normalizeFieldName(reflection reflect.Type, tagName string, fieldName string) string {
	field, ok := reflection.FieldByName(fieldName)
	if ok {
		return field.Tag.Get(tagName)
	}

	return fieldName
}

// Upper case first letter of the message.
func normalizeMessage(message error) string {
	return utils.UcFirst(message.Error())
}
