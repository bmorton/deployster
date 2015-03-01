package server

import "github.com/coreos/fleet/schema"

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
	Timestamp    string `json:"deploy_timestamp"`
}

// FindServiceUnits parses an array of units returned from fleet and looks for
// only the units that match the given service name, which is a subset of the
// Fleet unit name.  It collects all those units and returns an array of
// VersionedUnit structs that have had their additional deployster-specific
// fields populated from the Fleet unit name.
//
// Optional filtering by version is available too.  If all versions are desired,
// set version to "".  If a specific version is desired, set version to that
// version.
func FindServiceUnits(serviceName string, version string, units []*schema.Unit) []VersionedUnit {
	versionedUnits := []VersionedUnit{}

	for _, u := range units {
		dereferencedUnit := *u
		extractable := ExtractableUnit(dereferencedUnit)
		if extractable.ExtractBaseName() == serviceName {
			i := VersionedUnit{
				Service:      serviceName,
				Instance:     extractable.ExtractInstance(),
				Version:      extractable.ExtractVersion(),
				Timestamp:    extractable.ExtractTimestamp(),
				CurrentState: extractable.CurrentState,
				DesiredState: extractable.DesiredState,
				MachineID:    extractable.MachineID,
			}
			if shouldIncludeVersion(version, i.Version) {
				versionedUnits = append(versionedUnits, i)
			}
		}
	}

	return versionedUnits
}

// FindServiceUnits parses an array of units returned from fleet and collects
// all the versions present for the given service.  This doesn't currently look
// at the state of the units, it simply looks for matching service names and
// returns an array of all the unique versions found.
func FindServiceVersions(serviceName string, units []*schema.Unit) []string {
	uniqueVersions := make(map[string]bool)

	for _, u := range units {
		dereferencedUnit := *u
		extractable := ExtractableUnit(dereferencedUnit)
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

// shouldIncludeVersion takes an optional version checker and, if specified,
// ensures that it matches the unitVersion.  If the optional version is left
// blank, we'll return true.  If the optional version is present and it doesn't
// match the unitVersion, we'll return false.
func shouldIncludeVersion(blankOrVersion string, unitVersion string) bool {
	if blankOrVersion == "" {
		return true
	} else if blankOrVersion == unitVersion {
		return true
	}
	return false
}
