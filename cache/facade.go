package cache

import "github.com/lara-go/larago"

// FacadeWrapper for facade.
var FacadeWrapper = &larago.Facade{}

// Facade for cache.
func Facade() Cache {
	return FacadeWrapper.Resolve("cache").(Cache)
}
