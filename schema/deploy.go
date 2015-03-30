package schema

// Deploy is the struct that defines all the options for creating a new deploy.
// It is further populated after the initial request payload to contain all the
// information needed to be passed around to various collaborators.
type Deploy struct {
	ServiceName     string  `json:"service_name,omitempty"`
	Version         string  `json:"version"`
	DestroyPrevious bool    `json:"destroy_previous"`
	Timestamp       string  `json:"timestamp,omitempty"`
	InstanceCount   int     `json:"instance_count,omitempty"`
	PreviousVersion *Deploy `json:"previous_version,omitempty"`
}

// ServiceInstance returns a single unit of a possibly-many-unit deploy given
// the instance number.
func (d *Deploy) ServiceInstance(num string) *ServiceInstance {
	return &ServiceInstance{
		Name:      d.ServiceName,
		Version:   d.Version,
		Timestamp: d.Timestamp,
		Instance:  num,
	}
}
