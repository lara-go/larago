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
	return r.store.Has(key)
}

// Put value in cache by key.
func (r *Repository) Put(key string, value interface{}, duration time.Duration) {
	if duration.Nanoseconds() > 0 {
		r.store.Put(key, value, duration)

		r.event("cache.write", key, duration)
	}
}

// Forever remembers the value for ever.
func (r *Repository) Forever(key string, value interface{}) {
	r.store.Forever(key, value)

	var duration time.Duration
	r.event("cache.write", key, duration)
}

// Remember returns value if it was saved, or saves the value from the callback.
func (r *Repository) Remember(key string, duration time.Duration, callback func() (interface{}, error), target interface{}) error {
	return r.doRemember(key, callback, target, func(value interface{}) {
		r.Put(key, value, duration)
	})
}

// RememberForever returns value if it was saved, or saves the value from the callback for ever.
func (r *Repository) RememberForever(key string, callback func() (interface{}, error), target interface{}) error {
	return r.doRemember(key, callback, target, func(value interface{}) {
		r.Forever(key, value)
	})
}

// RememberForever returns value if it was saved, or saves the value from the callback for ever.
func (r *Repository) doRemember(
	key string,
	callback func() (interface{}, error),
	target interface{},
	save func(interface{}),
) error {

	// Try to get existing item.
	err := r.Get(key, target)
	if err == nil {
		r.event("cache.hit", key)

		return nil
	}

	if err != ErrorMissed {
		return err
	}

	r.event("cache.miss", key)

	// Retrieve and save new item from the callback.
	value, err := callback()
	if err != nil {
		return err
	}

	save(value)

	return setValue(target, value)
}

// Get saved value by the key.
func (r *Repository) Get(key string, target interface{}) error {
	err := r.store.Get(key, target)
	if err != nil {
		r.event("cache.miss", key)

		return err
	}

	r.event("cache.hit", key)

	return nil
}

// Pull saved value by the key.
func (r *Repository) Pull(key string, target interface{}) error {
	err := r.store.Get(key, target)
	if err != nil {
		r.event("cache.miss", key)

		return err
	}

	r.event("cache.hit", key)
	r.Forget(key)

	return nil
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
	if r.Events != nil {
		r.Events.Publish(event, payload...)
	}
}
