package handlers

import (
	"log"

	"github.com/bmorton/deployster/poller"
	"github.com/bmorton/deployster/schema"
	"github.com/bmorton/deployster/units"
	fleet "github.com/coreos/fleet/schema"
)

type DestroyHandler struct {
	PreviousVersion *schema.Deploy
	Client          poller.FleetClient
}

func (d *DestroyHandler) Handle(unit *fleet.UnitState) {
	log.Println("handling destroy!")
	launched := &units.ExtractableUnit{Name: unit.Name}
	err := d.Client.DestroyUnit(d.PreviousVersion.ServiceInstance(launched.ExtractInstance()).FleetUnitName())
	if err != nil {
		log.Println(err)
	}
	return
}
