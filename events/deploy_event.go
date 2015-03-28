package events

const deployEventType = "deploy"

type DeployEvent struct {
	Type          string `json:"type"`
	ServiceName   string `json:"service_name"`
	Version       string `json:"version"`
	Timestamp     string `json:"timestamp"`
	InstanceCount int    `json:"instance_count"`
}

func (d *DeployEvent) EventType() string {
	return deployEventType
}
