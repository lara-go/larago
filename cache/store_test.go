package cache_test

import (
	"testing"
	"time"

	"github.com/lara-go/larago/cache"
	"github.com/stretchr/testify/assert"
)

func testExpiration(t *testing.T, store cache.Store) {
	var test int

	assert.Equal(t, cache.ErrorMissed, store.Get("expire", &test))

	store.Put("expire", 1, time.Second)

	assert.True(t, store.Has("expire"))
	assert.Nil(t, store.Get("expire", &test))
	assert.Equal(t, 1, test)

	time.Sleep(time.Second * 2)
	assert.False(t, store.Has("expire"))
}

func testForget(t *testing.T, store cache.Store) {
	var forget int

	assert.Equal(t, cache.ErrorMissed, store.Get("forget", &forget))

	store.Forever("forget", 1)

	assert.True(t, store.Has("forget"))
	assert.Nil(t, store.Get("forget", &forget))
	assert.Equal(t, 1, forget)

	store.Forget("forget")

	assert.False(t, store.Has("forget"))
}

func testClear(t *testing.T, store cache.Store) {
	store.Put("test1", 1, time.Second)
	store.Put("test2", 2, time.Second*5)
	store.Forever("test3", 3)

	store.Clear()

	assert.False(t, store.Has("test1"))
	assert.False(t, store.Has("test2"))
	assert.False(t, store.Has("test3"))
}
