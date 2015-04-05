package poller

import (
	"github.com/bmorton/deployster/schema"
	fleet "github.com/coreos/fleet/schema"
)

type Event struct {
	ServiceInstance    *schema.ServiceInstance
	SystemdActiveState string
	SystemdLoadState   string
	SystemdSubState    string
}

func NewEvent(instance *schema.ServiceInstance, unitState *fleet.UnitState) *Event {
	return &Event{
		ServiceInstance:    instance,
		SystemdActiveState: unitState.SystemdActiveState,
		SystemdLoadState:   unitState.SystemdLoadState,
		SystemdSubState:    unitState.SystemdSubState,
	}
}
