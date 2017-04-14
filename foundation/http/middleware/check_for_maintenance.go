package middleware

import (
	"os"
	"path"

	"github.com/lara-go/larago"
	"github.com/lara-go/larago/http"
	"github.com/lara-go/larago/http/errors"
	"github.com/lara-go/larago/http/responses"
)

// CheckForMaintenance middleware.
type CheckForMaintenance struct {
	Application *larago.Application
}

// Handle request.
func (m *CheckForMaintenance) Handle(request *http.Request, next http.Handler) responses.Response {
	if _, err := os.Stat(path.Join(m.Application.HomeDirectory, "down")); err == nil {
		panic(errors.ServiceUnavailableHTTPError())
	}

	return next(request)
}
