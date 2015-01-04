package main

import (
	"github.com/bmorton/deployster/fleet"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type UnitResource struct {
	Fleet fleet.Client
}

type VersionedUnit struct {
	Service      string `json:"service"`
	Instance     string `json:"instance"`
	Version      string `json:"version"`
	CurrentState string `json:"current_state"`
	DesiredState string `json:"desired_state"`
	MachineID    string `json:"machine_id"`
}

type UnitResponse struct {
	Units []VersionedUnit `json:"units"`
}

type ExtractableUnit fleet.Unit

func (self *UnitResource) Index(u *url.URL, h http.Header, req interface{}) (int, http.Header, *UnitResponse, error) {
	statusCode := http.StatusOK
	response := &UnitResponse{}

	units, err := self.Fleet.Units()
	if err != nil {
		log.Printf("%#v\n", err)
		return http.StatusInternalServerError, nil, nil, err
	}
	response.Units = findServiceUnits(u.Query().Get("name"), units)

	return statusCode, nil, response, nil
}

func findServiceUnits(serviceName string, units []fleet.Unit) []VersionedUnit {
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
