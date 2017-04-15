package validation

import (
	"reflect"

	ozzo "github.com/go-ozzo/ozzo-validation"
	"github.com/lara-go/larago/support/utils"
)

// OzzoErrorsConverter converts ozzo-validation errors into valid meta.
type OzzoErrorsConverter struct{}

// ConvertValidationErrors to field - message format.
func (c *OzzoErrorsConverter) ConvertValidationErrors(err error, validator SelfValidator) *FieldsErrors {
	fails := err.(ozzo.Errors)

	// Populate ValidationErrors.
	if validator != nil {
		return c.formatStruct(validator, fails)
	}

	return c.formatRaw(fails)
}

// Format validator struct fields.
func (c *OzzoErrorsConverter) formatStruct(validator SelfValidator, fails ozzo.Errors) *FieldsErrors {
	validationErrors := &FieldsErrors{}

	reflection := reflect.ValueOf(validator).Elem().Type()
	for field, message := range fails {
		validationErrors.Errors = append(validationErrors.Errors, FieldError{
			Field:   c.normalizeFieldName(reflection, field),
			Message: c.normalizeMessage(message),
		})
	}

	return validationErrors
}

// Format raw fails.
func (c *OzzoErrorsConverter) formatRaw(fails ozzo.Errors) *FieldsErrors {
	validationErrors := &FieldsErrors{}

	for field, message := range fails {
		validationErrors.Errors = append(validationErrors.Errors, FieldError{
			Field:   field,
			Message: c.normalizeMessage(message),
		})
	}

	return validationErrors
}

// Convert request attribute name to normal.
func (c *OzzoErrorsConverter) normalizeFieldName(reflection reflect.Type, fieldName string) string {
	field, ok := reflection.FieldByName(fieldName)
	if ok {
		// Try json tag.
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			return jsonTag
		}

		// Try schema tag.
		shemaTag := field.Tag.Get("schema")
		if shemaTag != "" {
			return shemaTag
		}
	}

	// Return field name.r
	return fieldName
}

// Upper case first letter of the message.
func (c *OzzoErrorsConverter) normalizeMessage(message error) string {
	return utils.UcFirst(message.Error())
}
