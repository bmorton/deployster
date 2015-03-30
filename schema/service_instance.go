package schema

import "fmt"

// ServiceInstance represents a single unit of a possibly-many-unit deploy.
type ServiceInstance struct {
	Name      string
	Version   string
	Timestamp string
	Instance  string
}

// FleetUnitName generates a fleet unit name with the service name, version,
// and instance encoded within it.
func (s *ServiceInstance) FleetUnitName() string {
	return fmt.Sprintf("%s:%s:%s@%s.service", s.Name, s.Version, s.Timestamp, s.Instance)
}
