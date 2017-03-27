package cache_test

import (
	"testing"
	"time"

	"github.com/lara-go/larago/cache"
	"github.com/stretchr/testify/assert"
)

func factory() *cache.Repository {
	store := cache.NewInMemoryStore()

	return cache.NewRepository(store)
}

func TestDefaultValues(t *testing.T) {
	repo := factory()

	assert.Equal(t, "default", repo.Get("test", "default"))
	assert.Equal(t, "default", repo.Pull("test", "default"))
}

func TestRemember(t *testing.T) {
	repo := factory()

	value, err := repo.Remember("test", time.Second, func() (interface{}, error) {
		return 1, nil
	})

	assert.Nil(t, err)
	assert.Equal(t, 1, value.(int))

	assert.True(t, repo.Has("test"))
	assert.Equal(t, 1, repo.Get("test").(int))

	time.Sleep(time.Second * 2)
	assert.False(t, repo.Has("test"))
}

func TestRememberForever(t *testing.T) {
	repo := factory()

	value, err := repo.RememberForever("test", func() (interface{}, error) {
		return 1, nil
	})

	assert.Nil(t, err)
	assert.Equal(t, 1, value.(int))

	assert.True(t, repo.Has("test"))
	assert.Equal(t, 1, repo.Get("test").(int))

	time.Sleep(time.Second * 2)
	assert.True(t, repo.Has("test"))
}

func TestPull(t *testing.T) {
	repo := factory()

	assert.Nil(t, repo.Pull("test"))
	repo.Put("test", 1, time.Minute)
	assert.Equal(t, 1, repo.Pull("test").(int))
	assert.False(t, repo.Has("test"))
}
