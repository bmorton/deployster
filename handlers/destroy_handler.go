package handlers

import (
	"log"

	"github.com/bmorton/deployster/poller"
	"github.com/bmorton/deployster/schema"
)

type DestroyHandler struct {
	PreviousVersion *schema.Deploy
	Client          poller.FleetClient
}

func (d *DestroyHandler) Handle(event *poller.Event) {
	log.Println("handling destroy!")
	err := d.Client.DestroyUnit(d.PreviousVersion.ServiceInstance(event.ServiceInstance.Instance).FleetUnitName())
	if err != nil {
		log.Println(err)
	}
	return
}
