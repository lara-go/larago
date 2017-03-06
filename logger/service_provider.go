package logger

import (
	"log"
	"os"

	"github.com/lara-go/larago"
)

// ServiceProvider struct.
type ServiceProvider struct{}

// Register service.
func (p *ServiceProvider) Register(application *larago.Application) {
	application.Bind(func() (*Logger, error) {
		return &Logger{
			dateTimeFormat: application.DateTimeFormat,
			debugMode:      !application.Env("production") || application.Config().Debug(),
			logger:         log.New(os.Stdout, "", 0),
		}, nil
	})
}
