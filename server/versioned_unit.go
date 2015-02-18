package server

import (
	"github.com/bmorton/deployster/fleet"
)

// VersionedUnit is our representation of a Fleet unit.  Largely, the difference
// is that Fleet units don't care about versioning, so this lets us bridge the
// gap by exploding data that is encoded in the Fleet unit name into the proper
// fields appropriate for deployster.
type VersionedUnit struct {
	Service      string `json:"service"`
	Instance     string `json:"instance"`
	Version      string `json:"version"`
	CurrentState string `json:"current_state"`
	DesiredState string `json:"desired_state"`
	MachineID    string `json:"machine_id"`
}

// FindServiceUnits parses an array of units returned from fleet and looks for
// only the units that match the given service name, which is a subset of the
// Fleet unit name.  It collects all those units and returns an array of
// VersionedUnit structs that have had their additional deployster-specific
// fields populated from the Fleet unit name.
func FindServiceUnits(serviceName string, units []fleet.Unit) []VersionedUnit {
	versionedUnits := []VersionedUnit{}

	for _, u := range units {
		extractable := ExtractableUnit(u)
		if extractable.ExtractBaseName() == serviceName {
			i := VersionedUnit{
				Service:      serviceName,
				Instance:     extractable.ExtractInstance(),
				Version:      extractable.ExtractVersion(),
				CurrentState: extractable.CurrentState,
				DesiredState: extractable.DesiredState,
				MachineID:    extractable.MachineID,
			}
			versionedUnits = append(versionedUnits, i)
		}
	}

	return versionedUnits
}

// FindServiceUnits parses an array of units returned from fleet and collects
// all the versions present for the given service.  This doesn't currently look
// at the state of the units, it simply looks for matching service names and
// returns an array of all the unique versions found.
func FindServiceVersions(serviceName string, units []fleet.Unit) []string {
	uniqueVersions := make(map[string]bool)

	for _, u := range units {
		extractable := ExtractableUnit(u)
		if extractable.ExtractBaseName() == serviceName {
			uniqueVersions[extractable.ExtractVersion()] = true
		}
	}

	versions := make([]string, 0, len(uniqueVersions))
	for k := range uniqueVersions {
		versions = append(versions, k)
	}

	return versions
}
