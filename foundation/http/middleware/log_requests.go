package middleware

import (
	"time"

	"github.com/lara-go/larago/http"
	"github.com/lara-go/larago/http/responses"
	"github.com/lara-go/larago/logger"
)

// LogRequests for logger.
type LogRequests struct {
	NeedToLog bool `di:"Config.HTTP.LogRequests"`
	Logger    *logger.Logger
}

// Handle request.
func (m *LogRequests) Handle(request *http.Request, next http.Handler) responses.Response {
	path := request.RequestURI
	method := request.Method
	startTime := time.Now()

	// Retrieve response.
	response := next(request)

	if m.NeedToLog {
		// Write request info to the log.
		m.Logger.Info(
			"%d %4v %s %s %s",

			response.Status(),         // HTTP status code.
			time.Now().Sub(startTime), // Latency
			request.RemoteAddr,        // Remote IP
			method,                    // HTTP method
			path,                      // Request URI
		)
	}

	return response
}
