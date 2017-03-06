package errors

import (
	"reflect"

	"git.acronis.com/ci/ci-2x-ipn/app/support/utils"
	validation "github.com/go-ozzo/ozzo-validation"
)

// OzzoValidationErrorsConverter converts ozzo-validation errors into valid meta.
type OzzoValidationErrorsConverter struct{}

// ConvertValidationErrors to field - message format.
func (c *OzzoValidationErrorsConverter) ConvertValidationErrors(err error, tagName string, validator interface{}) *ValidationErrors {
	fails := err.(validation.Errors)

	validationErrors := &ValidationErrors{}

	reflection := reflect.ValueOf(validator).Elem().Type()

	// Populate ValidationErrors.
	for field, message := range fails {
		validationErrors.Errors = append(validationErrors.Errors, ValidationError{
			Field:   c.normalizeFieldName(reflection, tagName, field),
			Message: c.normalizeMessage(message),
		})
	}

	return validationErrors
}

// Convert request attribute name to normal.
func (c *OzzoValidationErrorsConverter) normalizeFieldName(reflection reflect.Type, tagName string, fieldName string) string {
	field, ok := reflection.FieldByName(fieldName)
	if ok {
		return field.Tag.Get(tagName)
	}

	return fieldName
}

// Upper case first letter of the message.
func (c *OzzoValidationErrorsConverter) normalizeMessage(message error) string {
	return utils.UcFirst(message.Error())
}
