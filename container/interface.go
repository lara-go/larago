package container

// Interface for container instance.
type Interface interface {
	// Bind registers instance in container.
	Bind(concrete interface{}, alias ...interface{})

	// Bound checks if container has requested binding.
	Bound(abstract interface{}) bool

	// Get binding from container.
	// Panics if there is no such binding.
	Get(abstract interface{}) interface{}

	// Make target with all its resolvable dependencies.
	// Panics if can't resolve missing field.
	Make(target interface{})

	// Factory wraps struct's dependency resolving to the lazy factory.
	Factory(target interface{}) func()

	// Call function resolving its params.
	// First it injects arguments,
	// then tries to resolve missing ones as dependencies.
	Call(function interface{}, args ...interface{}) (interface{}, error)

	// Wrap function call by the lazy callback.
	Wrap(function interface{}) func(args ...interface{}) (interface{}, error)
}

// TagsResolver interface to provide ability to resolve custom tags.
type TagsResolver interface {
	// ResolveTag resolves custom tags.
	ResolveTag(tag string, container *Container) interface{}
}
