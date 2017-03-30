package errors

// PanicHandlerInterface interface for panic handlers.
type PanicHandlerInterface interface {
	Defer()
}
