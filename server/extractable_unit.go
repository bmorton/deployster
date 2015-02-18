package server

import (
	"strings"

	"github.com/bmorton/deployster/fleet"
)

// ExtractableUnit is the local struct for a fleet.Unit with added functions
// for extracting the name, version, and instance that deployster encodes into
// the Fleet unit name.
type ExtractableUnit fleet.Unit

// ExtractBaseName returns the name of the service from the Fleet unit name.
// Given "railsapp-cf2e8ac@1.service" this returns "railsapp"
func (eu *ExtractableUnit) ExtractBaseName() string {
	s := strings.LastIndex(eu.Name, "-")
	return eu.Name[0:s]
}

// ExtractVersion returns the version of the service from the Fleet unit name.
// Given "railsapp-cf2e8ac@1.service" this returns "cf2e8ac"
func (eu *ExtractableUnit) ExtractVersion() string {
	start := strings.LastIndex(eu.Name, "-")
	end := strings.LastIndex(eu.Name, "@")
	return eu.Name[start+1 : end]
}

// ExtractInstance returns the instance of the service from the Fleet unit name.
// Given "railsapp-cf2e8ac@1.service" this returns "1"
func (eu *ExtractableUnit) ExtractInstance() string {
	s := strings.Split(eu.Name, "@")
	end := strings.Index(s[1], ".")
	return s[1][:end]
}
