package cache

import "time"

// Cache interface
type Cache interface {
	Has(key string) bool

	// Put value in cache by key.
	Put(key string, value interface{}, duration time.Duration)

	// Forever remembers the value for ever.
	Forever(key string, value interface{})

	// Remember returns value if it was saved, or saves the value from the callback.
	Remember(key string, duration time.Duration, callback func() (interface{}, error)) (interface{}, error)

	// RememberForever returns value if it was saved, or saves the value from the callback for ever.
	RememberForever(key string, callback func() (interface{}, error)) (interface{}, error)

	// Get saved value by the key.
	Get(key string, defaultValue ...interface{}) interface{}

	// Pull saved value by the key.
	Pull(key string, defaultValue ...interface{}) interface{}

	// Forget the value.
	Forget(key string)

	// Clear storage.
	Clear()
}

// Store interface.
type Store interface {
	// Put value in store by key.
	Put(key string, value interface{}, duration time.Duration)

	// Forever put value in store by key forever.
	Forever(key string, value interface{})

	// Get saved value by the key.
	Get(key string) interface{}

	// Forget the value.
	Forget(key string)

	// Clear storage.
	Clear()
}
