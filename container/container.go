package container

import (
	"errors"
	"fmt"
	"reflect"
)

// TagName is the name of field tag to resolve dependency via alias or custom resolver.
const TagName = "di"

// Container is a registry for all application objects.
// It can store/resolve/retrieve instances.
// Also it can resolve objects and functions dependencies.
type Container struct {
	Bindings  *Bindings
	Instances *Bindings

	tagsResolvers []TagsResolver
}

// New constructor.
func New() *Container {
	container := &Container{
		Bindings:  NewBindings(),
		Instances: NewBindings(),
	}

	container.Instance(container, (*Interface)(nil))

	return container
}

// Bind registers instance in container.
func (c *Container) Bind(concrete interface{}, aliases ...interface{}) {
	binding := c.makeBinding(concrete)
	resolver := reflect.ValueOf(concrete)

	c.Bindings.Set(binding, &resolver)

	for _, alias := range aliases {
		c.Bindings.Set(alias, &resolver)
	}
}

// Instance saves concrete as an already resolved instance.
func (c *Container) Instance(concrete interface{}, aliases ...interface{}) {
	binding := c.makeBinding(concrete)
	resolver := reflect.ValueOf(concrete)

	c.Instances.Set(binding, &resolver)

	for _, alias := range aliases {
		c.Instances.Set(alias, &resolver)
	}
}

// Unbind removes binding or resolved instance from container.
func (c *Container) Unbind(abstract interface{}) {
	// Mute all "not-found" panics.
	defer func() {}()

	c.Bindings.Remove(abstract)
	c.Instances.Remove(abstract)
}

// Make abstract value from concrete binding.
func (c *Container) makeBinding(concrete interface{}) reflect.Type {
	t := reflect.TypeOf(concrete)

	// If concrete resolver is not a function, return it.
	if t.Kind() != reflect.Func {
		return t
	}

	// If resolver is a function, use first returned value type as service abstract binding.
	if reflect.ValueOf(concrete).IsValid() && t.NumOut() != 2 {
		panic(errors.New("Custom resolver should be a valid function and has to return result and error. E.g. (*Foo, error)"))
	}

	return t.Out(0)
}

// Bound returns true if container has requested binding.
func (c *Container) Bound(abstract interface{}) bool {
	return c.Bindings.Has(abstract) || c.Instances.Has(abstract)
}

// Get binding from container.
func (c *Container) Get(abstract interface{}) interface{} {
	return c.resolveService(abstract).Interface()
}

//
func (c *Container) resolveService(abstract interface{}) *reflect.Value {
	resolved, err := c.resolveBinding(abstract)
	if err != nil {
		panic(fmt.Errorf("Can't resolve %s: %s", abstract, err.Error()))
	}

	return resolved
}

func (c *Container) resolveBinding(abstract interface{}) (*reflect.Value, error) {
	// If instance was already resolved, do not try to do it again.
	if c.Instances.Has(abstract) {
		return c.Instances.Get(abstract), nil
	}

	if !c.Bindings.Has(abstract) {
		return nil, fmt.Errorf("Unknown service %s", abstract)
	}

	concrete := c.Bindings.Get(abstract)
	resolved, err := c.resolve(concrete)
	if err != nil {
		return nil, err
	}

	c.Instances.Set(abstract, resolved)

	return resolved, nil
}

func (c *Container) resolve(concrete *reflect.Value) (*reflect.Value, error) {
	t := concrete.Type()

	// If concrete is a function, resolve it and continue immideately.
	if t.Kind() == reflect.Func {
		return c.callFunction(*concrete, nil)
	}

	c.fillStructDependencies(concrete)

	return concrete, nil
}

// Call function and return it's results.
func (c *Container) callFunction(fn reflect.Value, args []reflect.Value) (*reflect.Value, error) {
	if !fn.IsValid() {
		panic(fmt.Errorf("Invalid function to call %s", fn.Type()))
	}

	out := fn.Call(args)

	switch len(out) {
	case 2:
		if isError(out[1]) {
			return nil, out[1].Interface().(error)
		}

		return &out[0], nil
	case 1:
		if isError(out[0]) {
			return nil, out[0].Interface().(error)
		}

		return &out[0], nil
	case 0:
		return nil, nil
	}

	panic(fmt.Errorf("Unknown return type in function %s", fn.Type()))
}

// Build target with all its resolvable dependencies.
func (c *Container) Build(target interface{}) {
	if c.Bound(target) {
		c.fillStructDependencies(c.resolveService(target))
	} else {
		c.Make(target)
	}
}

// Make target with all its resolvable dependencies.
func (c *Container) Make(target interface{}) {
	v := reflect.ValueOf(target)
	c.fillStructDependencies(&v)
}

// Iterate through structure fields and try to resolve them.
func (c *Container) fillStructDependencies(target *reflect.Value) {
	v := target.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		c.resolveField(&fieldValue, &fieldType)
	}
}

// Resolve field by either its tag (if there is one) or its type.
func (c Container) resolveField(fieldValue *reflect.Value, fieldType *reflect.StructField) {
	// Check if we can set field's value
	if fieldType.Anonymous || !fieldValue.IsValid() || !fieldValue.CanSet() {
		return
	}

	tag := fieldType.Tag.Get(TagName)

	if tag != "" {
		if tag != "-" {
			fieldValue.Set(c.resolveTag(tag))
		}
	} else {
		fieldValue.Set(*c.resolveService(fieldValue.Type()))
	}
}

// Resolve service by field's tag value.
func (c *Container) resolveTag(tag string) reflect.Value {
	if resolved := c.tryTagsResolver(tag); resolved != nil {
		return reflect.ValueOf(resolved)
	}

	return *c.resolveService(tag)
}

// Try to resolve value via custom tags resolver.
func (c *Container) tryTagsResolver(tag string) interface{} {
	for _, tagsResolver := range c.tagsResolvers {
		if resolved := tagsResolver.ResolveTag(tag, c); resolved != nil {
			return resolved
		}
	}

	return nil
}

// SetTagsResolver setter for tags resolvers.
func (c *Container) SetTagsResolver(resolver TagsResolver) {
	c.tagsResolvers = append(c.tagsResolvers, resolver)
}

// Call function resolving its params.
func (c *Container) Call(function interface{}, args ...interface{}) (interface{}, error) {
	var result interface{}
	var fn reflect.Value
	var ok bool

	if fn, ok = function.(reflect.Value); !ok {
		fn = reflect.ValueOf(function)
	}

	ins := c.resolveFunctionArgs(fn, args)
	value, err := c.callFunction(fn, ins)
	if value != nil && value.IsValid() && value.CanInterface() {
		result = value.Interface()
	}

	return result, err
}

// Resolve function arguments.
func (c *Container) resolveFunctionArgs(function reflect.Value, args []interface{}) []reflect.Value {
	var f, i, j int

	t := function.Type()
	insLen := t.NumIn()
	ins := make([]reflect.Value, insLen)

	if insLen > 0 {
		// First inject passed in args that fit in ins.
		argsLen := len(args)
		for i = 0; i < argsLen && i < insLen; i++ {
			v := reflect.ValueOf(args[i])

			// If argument's type equals to ins type, use it. Otherwise skip.
			if normalizeAbstract(v.Type()) == normalizeAbstract(t.In(i)) {
				ins[i] = v
				f++
			}
		}

		// Then try to resolve left ones as dependencies.
		for j = f; j < insLen; j++ {
			ins[j] = *c.resolveService(t.In(j))
		}
	}

	return ins
}

// Factory wraps dependency resolving to the lazy factory.
func (c *Container) Factory(target interface{}) func() {
	return func() {
		c.Make(target)
	}
}

// Wrap function by the lazy callback.
func (c *Container) Wrap(function interface{}) func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		return c.Call(function, args...)
	}
}
