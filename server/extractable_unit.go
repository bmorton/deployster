package server

import (
	"strings"

	"github.com/coreos/fleet/schema"
)

// ExtractableUnit is the local struct for a fleet.Unit with added functions
// for extracting the name, version, and instance that deployster encodes into
// the Fleet unit name.
type ExtractableUnit schema.Unit

// ExtractBaseName returns the name of the service from the Fleet unit name.
// Given "railsapp:cf2e8ac_2013-06-05T14:10:43Z@1.service" this returns "railsapp"
func (eu *ExtractableUnit) ExtractBaseName() string {
	s := strings.Index(eu.Name, ":")
	return eu.Name[0:s]
}

// ExtractVersion returns the version of the service from the Fleet unit name.
// Given "railsapp:cf2e8ac_2013-06-05T14:10:43Z@1.service" this returns "cf2e8ac"
func (eu *ExtractableUnit) ExtractVersion() string {
	start := strings.Index(eu.Name, ":")
	end := strings.LastIndex(eu.Name, "_")
	return eu.Name[start+1 : end]
}

// ExtractInstance returns the instance of the service from the Fleet unit name.
// Given "railsapp:cf2e8ac_2013-06-05T14:10:43Z@1.service" this returns "1"
func (eu *ExtractableUnit) ExtractInstance() string {
	start := strings.LastIndex(eu.Name, "@")
	end := strings.LastIndex(eu.Name, ".")
	return eu.Name[start+1 : end]
}

// ExtractTimestamp returns the deploy timestamp that is appended to the Fleet
// unit name.
// Given "railsapp:cf2e8ac_2013-06-05T14:10:43Z@1.service" this returns
// "2013-06-05T14:10:43Z"
func (eu *ExtractableUnit) ExtractTimestamp() string {
	start := strings.LastIndex(eu.Name, "_")
	end := strings.LastIndex(eu.Name, "@")
	return eu.Name[start+1 : end]
}
