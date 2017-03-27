package cache_test

import (
	"testing"
	"time"

	"github.com/lara-go/larago/cache"
	"github.com/stretchr/testify/assert"
)

func testStore(t *testing.T, store cache.Store) {
	assert.Nil(t, store.Get("test1"))

	store.Put("test1", 1, time.Second)
	store.Put("test2", 2, time.Second*5)
	store.Forever("test3", 3)
	assert.Equal(t, 1, store.Get("test1").(int))

	// Test saving
	time.Sleep(time.Second * 2)
	assert.Nil(t, store.Get("test1"))
	assert.Equal(t, 2, store.Get("test2").(int))
	assert.Equal(t, 3, store.Get("test3").(int))

	store.Forget("test2")
	assert.Nil(t, store.Get("test2"))
	assert.Equal(t, 3, store.Get("test3").(int))

	store.Clear()
	assert.Nil(t, store.Get("test1"))
	assert.Nil(t, store.Get("test2"))
	assert.Nil(t, store.Get("test3"))
}
