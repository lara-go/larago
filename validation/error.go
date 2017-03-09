package validation

// Error .
type Error struct {
	err       error
	validator SelfValidator
}

// NewError constructor.
func NewError(err error, validator SelfValidator) *Error {
	return &Error{
		err:       err,
		validator: validator,
	}
}

// GetError returns validation error.
func (e *Error) GetError() error {
	return e.err
}

// GetValidator returns validator.
func (e *Error) GetValidator() SelfValidator {
	return e.validator
}

// Error returns error message.
func (e *Error) Error() string {
	return e.err.Error()
}

// FieldError for single field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FieldsErrors is a set of messages for every invalid field.
type FieldsErrors struct {
	Errors []FieldError `json:"errors"`
}
