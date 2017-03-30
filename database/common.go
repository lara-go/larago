package database

import (
	"fmt"
	"reflect"
)

// ModelNotFoundError .
type ModelNotFoundError struct {
	Model interface{}
	IDs   []interface{}
}

// Error returns error message.
func (e *ModelNotFoundError) Error() string {
	if e.Model != nil {
		t := reflect.TypeOf(e.Model)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		return fmt.Sprintf("Model [%s] not found with ids: %v", t, e.IDs)
	}

	return fmt.Sprintf("Model not found with ids: %v", e.IDs)
}

// ModelNotFound error constructor.
func ModelNotFound(model interface{}, ids []interface{}) *ModelNotFoundError {
	return &ModelNotFoundError{
		Model: model,
		IDs:   ids,
	}
}
