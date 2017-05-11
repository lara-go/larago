package larago

import "errors"

// Facade struct.
type Facade struct {
	Application *Application
	resolved    interface{}
}

// Resolve instance.
func (f *Facade) Resolve(accessor interface{}) interface{} {
	if f.Application == nil {
		panic(errors.New("Events facade was't registered properly"))
	}

	if f.resolved != nil {
		return f.resolved
	}

	f.resolved = f.Application.Get(accessor)

	return f.resolved
}

// Clear resolved facade.
func (f *Facade) Clear() {
	f.resolved = nil
}
