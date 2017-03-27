package cache

import (
	"time"

	"github.com/asaskevich/EventBus"
)

// Repository .
type Repository struct {
	Events *EventBus.EventBus
	store  Store
}

// NewRepository constructor.
func NewRepository(store Store) *Repository {
	return &Repository{
		store: store,
	}
}

// Has checks if there is such item.
func (r *Repository) Has(key string) bool {
	return r.Get(key) != nil
}

// Put value in cache by key.
func (r *Repository) Put(key string, value interface{}, duration time.Duration) {
	r.store.Put(key, value, duration)

	r.event("cache.write", key)
}

// Forever remembers the value for ever.
func (r *Repository) Forever(key string, value interface{}) {
	r.store.Forever(key, value)
}

// Remember returns value if it was saved, or saves the value from the callback.
func (r *Repository) Remember(key string, duration time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	return r.doRemember(key, callback, func(value interface{}) {
		if duration != 0 {
			r.store.Put(key, value, duration)
		}
	})
}

// RememberForever returns value if it was saved, or saves the value from the callback for ever.
func (r *Repository) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return r.doRemember(key, callback, func(value interface{}) {
		r.store.Forever(key, value)
	})
}

// RememberForever returns value if it was saved, or saves the value from the callback for ever.
func (r *Repository) doRemember(
	key string,
	callback func() (interface{}, error),
	saveCallback func(value interface{}),
) (interface{}, error) {

	// Try to get existing item.
	if value := r.Get(key); value != nil {
		r.event("cache.hit", key)

		return value, nil
	}

	r.event("cache.miss", key)

	// Retrieve and save new item from the callback.
	value, err := callback()
	if err != nil {
		return nil, err
	}

	saveCallback(value)

	return value, nil
}

// Get saved value by the key.
func (r *Repository) Get(key string, defaultValue ...interface{}) interface{} {
	item := r.store.Get(key)
	if item != nil {
		r.event("cache.hit", key)

		return item
	}

	r.event("cache.miss", key)

	if len(defaultValue) == 0 {
		defaultValue = append(defaultValue, nil)
	}

	return defaultValue[0]
}

// Pull saved value by the key.
func (r *Repository) Pull(key string, defaultValue ...interface{}) interface{} {
	item := r.Get(key)
	if item != nil {
		r.Forget(key)
		r.event("cache.hit", key)

		return item
	}

	r.event("cache.miss", key)

	if len(defaultValue) == 0 {
		defaultValue = append(defaultValue, nil)
	}

	return defaultValue[0]
}

// Forget the value.
func (r *Repository) Forget(key string) {
	r.store.Forget(key)

	r.event("cache.delete", key)
}

// Clear storage.
func (r *Repository) Clear() {
	r.store.Clear()

	r.event("cache.clear")
}

// Fire event.
func (r *Repository) event(event string, payload ...interface{}) {
	if r.Events == nil {
		return
	}

	r.Events.Publish(event, payload...)
}
