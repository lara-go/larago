package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Bindings struct.
type Bindings struct {
	lock  *sync.RWMutex
	items map[string]*reflect.Value
}

// NewBindings contrsuctor.
func NewBindings() *Bindings {
	return &Bindings{
		lock:  new(sync.RWMutex),
		items: make(map[string]*reflect.Value),
	}
}

// Set alias for binding.
func (b *Bindings) Set(abstract interface{}, resolver *reflect.Value) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.items[normalizeAbstract(abstract)] = resolver
}

// Has if there is such binding.
func (b *Bindings) Has(abstract interface{}) bool {
	b.lock.RLock()
	defer b.lock.RUnlock()

	_, exists := b.items[normalizeAbstract(abstract)]

	return exists
}

// Get binding by abstract value.
func (b *Bindings) Get(abstract interface{}) *reflect.Value {
	alias := normalizeAbstract(abstract)

	if !b.Has(alias) {
		panic(fmt.Errorf("Unknown service: %s", alias))
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.items[alias]
}

// Remove binding by abstract value.
func (b *Bindings) Remove(abstract interface{}) {
	alias := normalizeAbstract(abstract)

	if !b.Has(alias) {
		panic(fmt.Errorf("Unknown service: %s", alias))
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	delete(b.items, alias)
}

// Count bindings.
func (b *Bindings) Count() int {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return len(b.items)
}
