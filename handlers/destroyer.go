package handlers

import (
	"log"

	"github.com/bmorton/deployster/clients"
	"github.com/bmorton/deployster/poller"
	"github.com/bmorton/deployster/schema"
)

type Destroyer struct {
	PreviousVersion *schema.Deploy
	Client          clients.Fleet
}

func (d *Destroyer) Handle(event *poller.Event) {
	marked := d.PreviousVersion.ServiceInstance(event.ServiceInstance.Instance)
	log.Printf("Destroying %s due to launched instance replacement.\n", marked.FleetUnitName())
	err := d.Client.DestroyUnit(marked.FleetUnitName())
	if err != nil {
		log.Println(err)
	}
	return
}
