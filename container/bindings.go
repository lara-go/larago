package container

import (
	"fmt"
	"reflect"

	"github.com/lara-go/larago/support/collection"
)

// Bindings struct.
type Bindings struct {
	bindings *collection.Collection
}

// NewBindings contrsuctor.
func NewBindings() *Bindings {
	return &Bindings{collection.New()}
}

// Set alias for binding.
func (b *Bindings) Set(abstract interface{}, resolver *reflect.Value) {
	b.bindings.Set(b.normalizeAbstract(abstract), resolver)
}

// Has if there is such binding .
func (b *Bindings) Has(abstract interface{}) bool {
	return b.bindings.Has(b.normalizeAbstract(abstract))
}

// Get binding by abstract value.
func (b *Bindings) Get(abstract interface{}) *reflect.Value {
	alias := b.normalizeAbstract(abstract)

	if !b.bindings.Has(alias) {
		panic(fmt.Errorf("Unknown service: %s", alias))
	}

	return b.bindings.Get(alias).(*reflect.Value)
}

// Count bindings.
func (b *Bindings) Count() int {
	return b.bindings.Count()
}

// Normalize alias to internal form.
func (b *Bindings) normalizeAbstract(abstract interface{}) string {
	var t reflect.Type
	var ok bool

	// If already is reflect.Type, use it.
	// Or get type.
	if t, ok = abstract.(reflect.Type); !ok {
		t = reflect.TypeOf(abstract)
	}

	// For common string use it at once.
	if t.Kind() == reflect.String {
		return b.makeAlias(abstract)
	}

	// Always save by only element type.
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return b.makeAlias(t)
}

// Make alias string from abstract value.
func (b *Bindings) makeAlias(abstract interface{}) string {
	return fmt.Sprintf("%s", abstract)
}
