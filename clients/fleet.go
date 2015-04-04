package clients

import fleet "github.com/coreos/fleet/schema"

type Fleet interface {
	Units() ([]*fleet.Unit, error)
	CreateUnit(*fleet.Unit) error
	DestroyUnit(string) error
	UnitStates() ([]*fleet.UnitState, error)
	SetUnitTargetState(string, string) error
}
