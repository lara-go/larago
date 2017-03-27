package cache

import (
	"errors"

	larago "github.com/lara-go/larago"
)

// ServiceProvider struct.
type ServiceProvider struct{}

// Register service.
func (p *ServiceProvider) Register(application *larago.Application) {
	application.Commands(&CommandCacheClear{})

	application.Bind(func() (Cache, error) {
		store := application.Get("cache.store")

		if store == nil {
			return nil, errors.New("Cannot resolve cache store. 'cache.store' is empty")
		}

		repository := NewRepository(store.(Store))
		application.Make(repository)

		return repository, nil
	}, "cache")
}
