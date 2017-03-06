package collection

import "sync"

// Item structure used while iterating collection.
type Item struct {
	Key   interface{}
	Value interface{}
}

// Callback is a function that will be accepted on iterate methods.
type Callback func(item interface{}, key interface{}) interface{}

// Collection struct.
type Collection struct {
	lock  *sync.RWMutex
	items map[interface{}]interface{}
}

// New collection constructor.
func New() *Collection {
	return &Collection{
		lock:  new(sync.RWMutex),
		items: make(map[interface{}]interface{}),
	}
}

// NewWithItems to init collection with the set of items.
func NewWithItems(items map[interface{}]interface{}) *Collection {
	return &Collection{
		lock:  new(sync.RWMutex),
		items: items,
	}
}

// Set value.
func (c *Collection) Set(key interface{}, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items[key] = value
}

// Get value.
func (c *Collection) Get(key interface{}) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.items[key]
}

// Has value.
func (c *Collection) Has(key interface{}) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, exists := c.items[key]

	return exists
}

// Delete value.
func (c *Collection) Delete(key interface{}) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	delete(c.items, key)
}

// Count values in the map.
func (c *Collection) Count() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return len(c.items)
}

// Foreach iterates through collection via a channel.
// Returns Item struct for every iteration.
func (c *Collection) Foreach() <-chan Item {
	ch := make(chan Item)

	go func() {
		c.lock.RLock()
		defer c.lock.RUnlock()

		for key, value := range c.items {
			ch <- Item{key, value}
		}

		close(ch)
	}()

	return ch
}

// Map iterates through collection, applying callback function on every item.
// Returns new instance.
func (c *Collection) Map(callback Callback) *Collection {
	result := New()

	for item := range c.Foreach() {
		result.Set(item.Key, callback(item.Value, item.Key))
	}

	return result
}

// Transform iterates through collection, applying callback function on every item.
// Rather then Map method, transform changes collection itself.
func (c *Collection) Transform(callback Callback) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for key, value := range c.items {
		c.items[key] = callback(value, key)
	}
}

// All returns collection values.
func (c *Collection) All() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	values := make([]interface{}, len(c.items))
	i := 0
	for _, v := range c.items {
		values[i] = v
		i++
	}

	return values
}

// Keys returns collection keys.
func (c *Collection) Keys() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	keys := make([]interface{}, len(c.items))
	i := 0
	for k := range c.items {
		keys[i] = k
		i++
	}

	return keys
}
