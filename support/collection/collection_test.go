package collection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lara-go/larago/support/collection"
)

func TestItCanBeCreatedWithMap(t *testing.T) {
	c := collection.NewWithItems(map[interface{}]interface{}{
		"key": "value",
	})

	assert.Equal(t, 1, c.Count())
}

func TestItSetsValues(t *testing.T) {
	c := collection.New()

	c.Set("key", "value")
	assert.Equal(t, 1, c.Count())
}

func TestItGetsValues(t *testing.T) {
	c := collection.New()

	c.Set("key", "value")
	assert.True(t, c.Has("key"))
	assert.Equal(t, "value", c.Get("key").(string))
}

func TestItDeletesValues(t *testing.T) {
	c := collection.New()

	c.Set("key", "value")
	assert.True(t, c.Has("key"))
	c.Delete("key")
	assert.False(t, c.Has("key"))
}

func TestItIteratesValues(t *testing.T) {
	c := collection.New()

	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Set("key3", "value3")

	for item := range c.Foreach() {
		assert.Contains(t, item.Key, "key")
		assert.Contains(t, item.Value, "value")
	}
}

func TestItMapsValuesToNewCollection(t *testing.T) {
	c := collection.New()

	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Set("key3", "value3")

	result := c.Map(func(item interface{}, key interface{}) interface{} {
		return "v" + item.(string)
	})

	for item := range result.Foreach() {
		assert.Contains(t, item.Key, "key")
		assert.Contains(t, item.Value, "vvalue")
	}
}

func TestItTransformsCollection(t *testing.T) {
	c := collection.New()

	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Set("key3", "value3")

	c.Transform(func(item interface{}, key interface{}) interface{} {
		return "v" + item.(string)
	})

	for item := range c.Foreach() {
		assert.Contains(t, item.Key, "key")
		assert.Contains(t, item.Value, "vvalue")
	}
}

func TestItReturnsAllValues(t *testing.T) {
	c := collection.New()

	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Set("key3", "value3")

	assert.Len(t, c.All(), 3)
	assert.Contains(t, c.All(), "value2")
}

func TestItReturnsAllKeys(t *testing.T) {
	c := collection.New()

	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Set("key3", "value3")

	assert.Len(t, c.Keys(), 3)
	assert.Contains(t, c.Keys(), "key2")
}
