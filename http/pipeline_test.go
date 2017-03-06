package http_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lara-go/larago/container"
	"github.com/lara-go/larago/http"
	"github.com/lara-go/larago/http/responses"
)

type Foo struct{}

func (f *Foo) Bar() string {
	return "Foo"
}

type First struct {
	Foo *Foo
}

func (m *First) Handle(request *http.Request, next http.Handler) responses.Response {
	response := next(request)

	return responses.NewText(response.Status(), fmt.Sprintf("%s %s", response.Body(), m.Foo.Bar()))
}

type Second struct{}

func (m *Second) Handle(request *http.Request, next http.Handler) responses.Response {
	response := next(request)

	return responses.NewText(response.Status(), fmt.Sprintf("%s Second", response.Body()))
}

func TestMiddleware(t *testing.T) {
	request := &http.Request{}

	c := container.New()
	c.Bind(&Foo{})

	response := http.NewPipeline(c).
		Send(request).
		Through([]http.Middleware{
			&First{},
			&Second{},
		}).
		Then(func(request *http.Request) responses.Response {
			return responses.NewText(200, "Result")
		})

	assert.Equal(t, 200, response.Status())
	assert.Equal(t, "Result Second Foo", string(response.Body()))
}
