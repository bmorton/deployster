package server

import (
	"github.com/coreos/fleet/schema"
	"github.com/fsouza/go-dockerclient"
)

// DockerClient is the interface required for TasksResource to be able to
// create, start, attach, inspect, and remove Docker containers.
type DockerClient interface {
	CreateContainer(docker.CreateContainerOptions) (*docker.Container, error)
	StartContainer(string, *docker.HostConfig) error
	AttachToContainer(docker.AttachToContainerOptions) error
	InspectContainer(string) (*docker.Container, error)
	RemoveContainer(docker.RemoveContainerOptions) error
}

type FleetClient interface {
	Units() ([]*schema.Unit, error)
	CreateUnit(*schema.Unit) error
	DestroyUnit(string) error
	UnitStates() ([]*schema.UnitState, error)
}
