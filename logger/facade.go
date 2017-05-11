package logger

import "github.com/lara-go/larago"

// FacadeWrapper for facade.
var FacadeWrapper = &larago.Facade{}

// Facade for logger.
func Facade() *Logger {
	return FacadeWrapper.Resolve("logger").(*Logger)
}
