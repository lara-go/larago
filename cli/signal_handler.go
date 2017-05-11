package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/asaskevich/EventBus"
	"github.com/lara-go/larago"
)

// SignalsHandler handles cli signals.
type SignalsHandler struct {
	Events      *EventBus.EventBus
	ExitHandler larago.ExitHandler
}

// CatchInterrupt handles sigterm.
func (h *SignalsHandler) CatchInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		if h.Events != nil {
			h.Events.Publish("sigterm")
		}

		h.ExitHandler.Exit("\nGracefull exit...", 0)
	}()
}
