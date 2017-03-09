package container

import (
	"fmt"
	"reflect"
)

// Normalize alias to internal form.
func normalizeAbstract(abstract interface{}) string {
	var t reflect.Type
	var ok bool

	// If already is reflect.Type, use it.
	// Or get type.
	if t, ok = abstract.(reflect.Type); !ok {
		t = reflect.TypeOf(abstract)
	}

	// For common string use it at once.
	if t.Kind() == reflect.String {
		return makeAlias(abstract)
	}

	// Always save by only element type.
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return makeAlias(t)
}

// Make alias string from abstract value.
func makeAlias(abstract interface{}) string {
	return fmt.Sprintf("%s", abstract)
}
