package server

// Deploy is the struct that defines all the options for creating a new deploy
// and is wrapped by DeployRequest and deserialized in the Create function.
type Deploy struct {
	ServiceName     string `json:"service_name,omitempty"`
	Version         string `json:"version"`
	DestroyPrevious bool   `json:"destroy_previous"`
	Timestamp       string `json:"timestamp,omitempty"`
	InstanceCount   int    `json:"instance_count,omitempty"`
}

func (d *Deploy) ServiceInstance(num string) *ServiceInstance {
	return &ServiceInstance{
		Name:      d.ServiceName,
		Version:   d.Version,
		Timestamp: d.Timestamp,
		Instance:  num,
	}
}
