package larago_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lara-go/larago"
)

/**
 * Simple bind test.
 */

type Bind struct {
	value string
}

type BindServiceProvider struct{}

func (p *BindServiceProvider) Boot() {

}

func (p *BindServiceProvider) Register(application *larago.Application) {
	application.Bind(func() (*Bind, error) {
		return &Bind{}, nil
	}, "bind")
}

func TestSimpleRetrieve(t *testing.T) {
	application := larago.New()
	application.Register(&BindServiceProvider{})

	// Resolve Bind manually.
	bind := application.Get("bind")
	assert.NotNil(t, bind)
}

/**
 * Test custom bind errors.
 */

type BindWithErrorServiceProvider struct{}

func (p *BindWithErrorServiceProvider) Register(application *larago.Application) {
	application.Bind(func() (interface{}, error) {
		return nil, errors.New("Custom bind error")
	}, "bind1")
}

func TestErrorRetrieve(t *testing.T) {
	application := larago.New()
	application.Register(&BindWithErrorServiceProvider{})

	assert.Panics(t, func() {
		application.Get("bind1")
	})
}

/**
 * Test bad resolver.
 */

type BindWithBadResolverServiceProvider struct{}

func (p *BindWithBadResolverServiceProvider) Register(application *larago.Application) {
	application.Bind(func() error {
		return nil
	}, "bind2")
}

func TestBadResolverError(t *testing.T) {
	application := larago.New()

	assert.Panics(t, func() {
		application.Register(&BindWithBadResolverServiceProvider{})
	})
}

/**
 * Test dependencies injecting.
 * Test singletons.
 * Test interfaces and hidden fields.
 */

type Single struct {
	Bind  *Bind `di:"bind"`
	value string
}

type SingleServiceProvider struct{}

func (p *SingleServiceProvider) Boot(single *Single) {
	single.value = "booted"
}

func (p *SingleServiceProvider) Register(application *larago.Application) {
	application.Bind(&Single{})
}

type Interface interface{}
type Foo struct {
	// Skip embeded interface resilving
	Interface

	// Resolve singleton.
	Single *Single
}
type Bar struct {
	// Resolve the same singleton by tag.
	Singleee *Single

	// Skip hidden field.
	hiddenValue string

	// Skip field resolving.
	FooBar string `di:"-"`
}

func TestResolveDependencies(t *testing.T) {
	application := larago.New()
	application.Register(&BindServiceProvider{})
	application.Boot()
	application.Register(&SingleServiceProvider{})

	// Resolve Single by attribute name.
	// Check that interface wll not be "resolved".
	var foo Foo
	application.Make(&foo)
	assert.NotNil(t, foo.Single)
	assert.Equal(t, "booted", foo.Single.value)

	// Resolve Single by tag value.
	// Check that hidden value wll not be "resolved" and changed.
	bar := &Bar{hiddenValue: "foo bar baz"}
	application.Make(bar)

	assert.Equal(t, "foo bar baz", bar.hiddenValue)
	assert.NotNil(t, bar.Singleee)
	assert.Empty(t, bar.FooBar)

	assert.Exactly(t, foo.Single, bar.Singleee, "Singleton doesn't work. Foo was resolved again.")
}

func TestFunctionCall(t *testing.T) {
	application := larago.New()
	application.Register(&BindServiceProvider{})
	application.Boot()
	application.Register(&SingleServiceProvider{})

	// Resolve Single by attribute name.
	// Check that interface wll not be "resolved".
	var foo Foo
	application.Make(&foo)

	application.Call(func(str string, num int, single *Single) {
		assert.Equal(t, "string", str)
		assert.Equal(t, 123, num)
		assert.NotNil(t, single)
		assert.Equal(t, "booted", single.value)
	}, "string", 123)

	result, err := application.Call(func(single *Single) (string, error) {
		return "result", nil
	})
	assert.Equal(t, "result", result.(string))
	assert.Nil(t, err)
}

// type ConfMe struct {
// 	// Retrieve full config.
// 	Config *conf.Config
//
// 	// Retieve part of the config.
// 	Database *conf.Database `di:"Config.Database"`
//
// 	// Retieve particular config value.
// 	Env string `di:"Config.App.Env"`
// }
//
// func TestResolveConfigValues(t *testing.T) {
// 	application := larago.New()
// 	// application.SetConfig(conf.DefaultConfig())
//
// 	// Resolve config value.
// 	var foo ConfMe
// 	application.Make(&foo)
// 	assert.Equal(t, "production", foo.Env)
// 	assert.Equal(t, "sqlite3", foo.Database.Driver)
// 	assert.Equal(t, foo.Config.Database.Driver, foo.Database.Driver)
//
// 	assert.True(t, application.Env("production"))
// }
