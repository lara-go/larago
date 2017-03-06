package middleware

import (
	"github.com/lara-go/larago/container"
	"github.com/lara-go/larago/http"
	"github.com/lara-go/larago/http/responses"
)

// FormRequests middleware.
type FormRequests struct {
	Container *container.Interface
}

// Handle request.
func (m *FormRequests) Handle(request *http.Request, next http.Handler) responses.Response {
	return next(request)
}
