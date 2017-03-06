package container_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lara-go/larago/container"
)

type Test struct {
	value string
}

func (t *Test) Do() {
	fmt.Println("Tester does.")
}

type Tester interface {
	Do()
}

func TestImplentsInterface(t *testing.T) {
	c := container.New()

	assert.Implements(t, (*container.Interface)(nil), c)
}

func TestConcreteBind(t *testing.T) {
	var a *Test
	var b *Test

	c := container.New()
	s := &Test{}

	c.Bind(s)

	assert.True(t, c.Bound((*Test)(nil)))

	assert.NotPanics(t, func() {
		a = c.Get(a).(*Test)
		assert.Exactly(t, s, a)
	})

	// check for re-resolve
	assert.NotPanics(t, func() {
		b = c.Get(b).(*Test)
		assert.Exactly(t, s, b)
	})
	assert.Exactly(t, a, b)
}

func TestInterfaceBind(t *testing.T) {
	c := container.New()
	b := &Test{}

	c.Bind(b, (*Tester)(nil))

	assert.True(t, c.Bound((*Tester)(nil)))

	assert.NotPanics(t, func() {
		r := c.Get((*Tester)(nil))
		assert.Exactly(t, b, r)
	})
}

func TestAliasBind(t *testing.T) {
	c := container.New()
	b := &Test{}

	c.Bind(b, "alias1", "alias2")

	assert.True(t, c.Bound("alias1"))
	assert.True(t, c.Bound("alias2"))
	assert.False(t, c.Bound("alias3"))

	assert.NotPanics(t, func() {
		r1 := c.Get("alias1")
		assert.Exactly(t, b, r1)

		r2 := c.Get("alias2")
		assert.Exactly(t, b, r2)
	})

	assert.Panics(t, func() {
		r1 := c.Get("alias3")
		assert.Exactly(t, b, r1)
	})
}

func TestBindWithCustomResolver(t *testing.T) {
	c := container.New()

	c.Bind(func() (*Test, error) {
		return &Test{}, nil
	})

	assert.True(t, c.Bound(&Test{}))

	assert.NotPanics(t, func() {
		r := c.Get(&Test{})
		assert.NotNil(t, r)
	})
}

func TestBindWithCustomResolverWithInterfaceReturn(t *testing.T) {
	c := container.New()

	c.Bind(func() (Tester, error) {
		return &Test{}, nil
	})

	assert.True(t, c.Bound((*Tester)(nil)))

	assert.NotPanics(t, func() {
		r := c.Get((*Tester)(nil))
		assert.NotNil(t, r)
		assert.Implements(t, (*Tester)(nil), r)
	})
}

func TestBindWithCustomResolverWithInterfaceReturnAndStringAlias(t *testing.T) {
	c := container.New()

	c.Bind(func() (Tester, error) {
		return &Test{}, nil
	}, "alias")

	assert.True(t, c.Bound((*Tester)(nil)))
	assert.True(t, c.Bound("alias"))

	assert.NotPanics(t, func() {
		r1 := c.Get((*Tester)(nil))
		assert.NotNil(t, r1)

		r2 := c.Get("alias")
		assert.NotNil(t, r2)

		assert.Exactly(t, r1, r2)
	})
}

func TestBindWithCustomResolverWithError(t *testing.T) {
	c := container.New()

	c.Bind(func() (Tester, error) {
		return nil, errors.New("Error")
	})

	assert.Panics(t, func() {
		r1 := c.Get((*Tester)(nil))
		assert.NotNil(t, r1)
	})
}

func TestBindWithBadCustomResolver(t *testing.T) {
	c := container.New()

	assert.Panics(t, func() {
		c.Bind(func() {

		})
	})
}

type Stub1 struct {
	Test *Test
}

func TestMake(t *testing.T) {
	c := container.New()
	b := &Test{}

	c.Bind(b)

	var stub Stub1
	c.Make(&stub)
	assert.NotNil(t, stub.Test)
}

type Stub2 struct {
	Test Tester
}

func TestMakeWithInterfaceAlias(t *testing.T) {
	c := container.New()

	c.Bind(&Test{}, (*Tester)(nil))

	var stub Stub2
	c.Make(&stub)
	assert.NotNil(t, stub.Test)
}

type Stub3 struct {
	Test interface{} `di:"test"`
	Foo  interface{} `di:"-"`
}

func TestMakeWithTag(t *testing.T) {
	c := container.New()
	b := &Test{}

	c.Bind(b, "test")

	var stub Stub3
	c.Make(&stub)
	assert.NotNil(t, stub.Test)
}

func TestFunctionResolving(t *testing.T) {
	c := container.New()

	b := &Test{}
	c.Bind(b)

	assert.NotPanics(t, func() {
		result, err := c.Call(func(foo string, test *Test) (string, error) {
			assert.NotNil(t, test)

			return foo, nil
		}, "foo")

		assert.Nil(t, err)
		assert.Equal(t, "foo", result)
	})
}

func TestFunctionResolvingWithCustomResolver(t *testing.T) {
	c := container.New()

	c.Bind(func() (*Test, error) {
		return &Test{}, nil
	})

	assert.NotPanics(t, func() {
		result, err := c.Call(func(test *Test) (string, error) {
			assert.NotNil(t, test)

			return "foo", nil
		})

		assert.Nil(t, err)
		assert.Equal(t, "foo", result)
	})
}

func TestFunctionResolvingWithCustomResolverAndInterfaceReturn(t *testing.T) {
	c := container.New()

	c.Bind(func() (Tester, error) {
		return &Test{}, nil
	})

	assert.NotPanics(t, func() {
		result, err := c.Call(func(test Tester) (string, error) {
			assert.NotNil(t, test)

			return "foo", nil
		})

		assert.Nil(t, err)
		assert.Equal(t, "foo", result)
	})
}

func TestFunctionResolvingWithInterfaceAlias(t *testing.T) {
	c := container.New()

	c.Bind(&Test{}, (*Tester)(nil))

	assert.NotPanics(t, func() {
		result, err := c.Call(func(test Tester) (string, error) {
			assert.NotNil(t, test)

			return "foo", nil
		})

		assert.Nil(t, err)
		assert.Equal(t, "foo", result)
	})
}

func TestFunctionCallWithoutReturn(t *testing.T) {
	c := container.New()

	c.Bind(&Test{}, (*Tester)(nil))

	assert.NotPanics(t, func() {
		c.Call(func(test Tester) {
			assert.NotNil(t, test)
		})
	})
}

func TestFunctionCallWithErrorReturn(t *testing.T) {
	c := container.New()

	c.Bind(&Test{}, (*Tester)(nil))

	assert.NotPanics(t, func() {
		result, err := c.Call(func(test Tester) error {
			assert.NotNil(t, test)

			return errors.New("error")
		})
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.NotNil(t, err)
	})
}

func TestFunctionCallWithCommonReturn(t *testing.T) {
	c := container.New()

	c.Bind(&Test{}, (*Tester)(nil))

	assert.NotPanics(t, func() {
		result, err := c.Call(func(test Tester) string {
			assert.NotNil(t, test)

			return "foo"
		})
		assert.Nil(t, err)
		assert.Equal(t, "foo", result)
	})
}

func TestLazyFunctionResolving(t *testing.T) {
	c := container.New()
	c.Bind(&Test{})

	lazy := c.Wrap(func(foo string, test *Test) (string, error) {
		assert.NotNil(t, test)

		return foo, nil
	})

	assert.NotPanics(t, func() {
		result, err := lazy("foo", "bar")

		assert.Nil(t, err)
		assert.Equal(t, "foo", result)
	})
}

func TestLazyStructResolving(t *testing.T) {
	c := container.New()
	c.Bind(&Test{})

	var stub Stub1
	lazy := c.Factory(&stub)
	assert.Nil(t, stub.Test)

	assert.NotPanics(t, func() {
		lazy()
		assert.NotNil(t, stub.Test)
	})
}
