package server

import (
	"github.com/bmorton/deployster/fleet"
)

type VersionedUnit struct {
	Service      string `json:"service"`
	Instance     string `json:"instance"`
	Version      string `json:"version"`
	CurrentState string `json:"current_state"`
	DesiredState string `json:"desired_state"`
	MachineID    string `json:"machine_id"`
}

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
