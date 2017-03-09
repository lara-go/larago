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
	b.bindings.Set(normalizeAbstract(abstract), resolver)
}

// Has if there is such binding .
func (b *Bindings) Has(abstract interface{}) bool {
	return b.bindings.Has(normalizeAbstract(abstract))
}

// Get binding by abstract value.
func (b *Bindings) Get(abstract interface{}) *reflect.Value {
	alias := normalizeAbstract(abstract)

	if !b.bindings.Has(alias) {
		panic(fmt.Errorf("Unknown service: %s", alias))
	}

	return b.bindings.Get(alias).(*reflect.Value)
}

// Count bindings.
func (b *Bindings) Count() int {
	return b.bindings.Count()
}
