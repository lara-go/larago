package stubs

// MiddlewareStub template.
const MiddlewareStub = `
package middleware

import (
	"github.com/lara-go/larago/http"
	"github.com/lara-go/larago/http/responses"
)

// {{.Name}} test.
type {{.Name}} struct{}

// Handle request.
func (m *{{.Name}}) Handle(request *http.Request, next http.Handler) responses.Response {
	return next(request)
}
`
