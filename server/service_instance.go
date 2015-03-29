package server

import "fmt"

type ServiceInstance struct {
	Name      string
	Version   string
	Timestamp string
	Instance  string
}

// fleetServiceName generates a fleet unit name with the service name, version,
// and instance encoded within it.
func (s *ServiceInstance) FleetUnitName() string {
	return fmt.Sprintf("%s:%s:%s@%s.service", s.Name, s.Version, s.Timestamp, s.Instance)
}
