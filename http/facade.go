package http

import "github.com/lara-go/larago"

// FacadeWrapper for facade.
var FacadeWrapper = &larago.Facade{}

// Facade for router.
func Facade() *Router {
	return FacadeWrapper.Resolve("router").(*Router)
}
