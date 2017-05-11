package events

import (
	"github.com/asaskevich/EventBus"
	"github.com/lara-go/larago"
)

// FacadeWrapper for facade.
var FacadeWrapper = &larago.Facade{}

// Facade for events.
func Facade() *EventBus.EventBus {
	return FacadeWrapper.Resolve("events").(*EventBus.EventBus)
}
