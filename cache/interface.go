package cache

import "time"

// Cache interface
type Cache interface {
	// Has checks if there is such item.
	Has(key string) bool

	// Put value in cache by key.
	Put(key string, value interface{}, duration time.Duration)

	// Forever remembers the value for ever.
	Forever(key string, value interface{})

	// Remember returns value if it was saved, or saves the value from the callback.
	Remember(key string, duration time.Duration, callback func() (interface{}, error), target interface{}) error

	// RememberForever returns value if it was saved, or saves the value from the callback for ever.
	RememberForever(key string, callback func() (interface{}, error), target interface{}) error

	// Get saved value by the key.
	Get(key string, target interface{}) error

	// Pull saved value by the key.
	Pull(key string, target interface{}) error

	// Forget the value.
	Forget(key string)

	// Clear storage.
	Clear()
}

// Store interface.
type Store interface {
	// Has checks if there is such item.
	Has(key string) bool

	// Put value in store by key.
	Put(key string, value interface{}, duration time.Duration) error

	// Forever put value in store by key forever.
	Forever(key string, value interface{}) error

	// Get saved value by the key.
	Get(key string, target interface{}) error

	// Forget the value.
	Forget(key string)

	// Clear storage.
	Clear()
}
