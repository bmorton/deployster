package main

import (
	"github.com/bmorton/deployster/fleet"
	"strings"
)

type VersionedUnit struct {
	Service      string `json:"service"`
	Instance     string `json:"instance"`
	Version      string `json:"version"`
	CurrentState string `json:"current_state"`
	DesiredState string `json:"desired_state"`
	MachineID    string `json:"machine_id"`
}

type ExtractableUnit fleet.Unit

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

func (self *ExtractableUnit) ExtractBaseName() string {
	s := strings.Split(self.Name, "-")
	return s[0]
}

func (self *ExtractableUnit) ExtractVersion() string {
	s := strings.Split(self.Name, "-")
	end := strings.Index(s[1], "@")
	return s[1][:end]
}

func (self *ExtractableUnit) ExtractInstance() string {
	s := strings.Split(self.Name, "@")
	end := strings.Index(s[1], ".")
	return s[1][:end]
}
