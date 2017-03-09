package validation

// SelfValidator interface for json requests.
type SelfValidator interface {
	Validate() error
}

// ErrorsConverter interface for requests validators.
type ErrorsConverter interface {
	// ConvertValidationErrors to field - message format.
	ConvertValidationErrors(err error, validator SelfValidator) *FieldsErrors
}
