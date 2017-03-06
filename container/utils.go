package container

import "reflect"

// Check if value is an error.
func isError(err reflect.Value) bool {
	errType := reflect.TypeOf((*error)(nil)).Elem()

	return err.Type().Kind() == reflect.Interface && !err.IsNil() && err.Type().Implements(errType)
}
