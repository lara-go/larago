package container_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lara-go/larago/container"
)

type B struct{}
type I interface{}

var v = reflect.ValueOf(&B{})

func TestItSavesBinding(t *testing.T) {
	b := container.NewBindings()
	b.Set("binding", &v)

	assert.Equal(t, 1, b.Count())
	assert.True(t, b.Has("binding"))
	assert.NotPanics(t, func() {
		r := b.Get("binding")
		assert.Equal(t, &v, r)
	})
}

func TestItSavesBindingByStructPtr(t *testing.T) {
	b := container.NewBindings()
	b.Set(&B{}, &v)

	assert.True(t, b.Has((*B)(nil)))
	assert.NotPanics(t, func() {
		r := b.Get((*B)(nil))
		assert.Equal(t, &v, r)
	})
}

func TestItSavesBindingByInterfacePtr(t *testing.T) {
	b := container.NewBindings()
	b.Set((*I)(nil), &v)

	assert.True(t, b.Has((*I)(nil)))
	assert.NotPanics(t, func() {
		r := b.Get((*I)(nil))
		assert.Equal(t, &v, r)
	})
}

func TestItPanicsOnUnknownService(t *testing.T) {
	b := container.NewBindings()
	b.Set((*I)(nil), &v)

	assert.False(t, b.Has("foo"))
	assert.Panics(t, func() {
		b.Get("foo")
	})
}
