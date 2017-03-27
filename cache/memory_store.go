package cache

import (
	"time"

	"github.com/lara-go/larago/support/collection"
	"github.com/uniplaces/carbon"
)

// memoryItem to store.
type memoryItem struct {
	value      interface{}
	expiration *carbon.Carbon
}

// InMemoryStore .
type InMemoryStore struct {
	store *collection.Collection
}

// NewInMemoryStore .
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: collection.New(),
	}
}

// Put value in cache by key.
func (s *InMemoryStore) Put(key string, value interface{}, duration time.Duration) {
	s.store.Set(key, &memoryItem{
		value:      value,
		expiration: carbon.NewCarbon(time.Now().Add(duration)),
	})
}

// Get saved value by the key.
func (s *InMemoryStore) Get(key string) interface{} {
	// Check if is still alive and if not, forget it.
	hit := s.store.Get(key)
	if hit == nil {
		return nil
	}

	item := hit.(*memoryItem)
	if item.expiration.IsPast() {
		s.Forget(key)

		return nil
	}

	return item.value
}

// Forever put value in store by key forever.
func (s *InMemoryStore) Forever(key string, value interface{}) {
	s.Put(key, value, time.Minute*5256000)
}

// Forget the value.
func (s *InMemoryStore) Forget(key string) {
	s.store.Delete(key)
}

// Clear storage.
func (s *InMemoryStore) Clear() {
	s.store = collection.New()
}
