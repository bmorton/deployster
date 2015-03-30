package poller

import fleet "github.com/coreos/fleet/schema"

type Handler interface {
	Handle(*fleet.UnitState)
}
