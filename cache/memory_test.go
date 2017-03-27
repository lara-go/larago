package cache_test

import (
	"testing"

	"github.com/lara-go/larago/cache"
)

func TestMemoryStore(t *testing.T) {
	testStore(t, cache.NewInMemoryStore())
}
