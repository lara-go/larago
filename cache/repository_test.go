package cache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/lara-go/larago/cache"
	"github.com/stretchr/testify/assert"
)

func repositoryFactory() *cache.Repository {
	store := cache.NewInMemoryStore()

	repo := cache.NewRepository(store)
	repo.Events = events()

	return repo
}

func events() *EventBus.EventBus {
	events := EventBus.New()
	events.Subscribe("cache.write", func(key string, duration time.Duration) {
		fmt.Printf("Put '%s' for %s\n", key, duration)
	})
	events.Subscribe("cache.hit", func(key string) {
		fmt.Printf("Hit '%s'\n", key)
	})
	events.Subscribe("cache.miss", func(key string) {
		fmt.Printf("Missed '%s'\n", key)
	})

	return events
}

func TestRemember(t *testing.T) {
	repo := repositoryFactory()

	var value int
	err := repo.Remember("remember", time.Second, func() (interface{}, error) {
		return 1, nil
	}, &value)

	assert.Nil(t, err)
	assert.Equal(t, 1, value)

	assert.True(t, repo.Has("remember"))

	time.Sleep(time.Second * 2)
	assert.False(t, repo.Has("remember"))
}

func TestRememberForever(t *testing.T) {
	repo := repositoryFactory()

	var value int
	err := repo.RememberForever("forever", func() (interface{}, error) {
		return 1, nil
	}, &value)

	assert.Nil(t, err)
	assert.Equal(t, 1, value)

	assert.True(t, repo.Has("forever"))

	time.Sleep(time.Second * 2)
	assert.True(t, repo.Has("forever"))
}

func TestPull(t *testing.T) {
	repo := repositoryFactory()
	var test int

	assert.Equal(t, cache.ErrorMissed, repo.Pull("pull", &test))

	repo.Put("pull", 1, time.Second)

	assert.Nil(t, repo.Pull("pull", &test))
	assert.Equal(t, 1, test)
	assert.False(t, repo.Has("pull"))
}
