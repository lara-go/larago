package cache_test

import (
	"testing"

	"github.com/lara-go/larago/cache"
)

func TestMemoryStore_Expire(t *testing.T) {
	testExpiration(t, cache.NewInMemoryStore())
}

func TestMemoryStore_Forget(t *testing.T) {
	testForget(t, cache.NewInMemoryStore())
}

func TestMemoryStore_Clear(t *testing.T) {
	testClear(t, cache.NewInMemoryStore())
}
