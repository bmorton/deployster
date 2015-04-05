package clients

import "github.com/fsouza/go-dockerclient"

// DockerClient is the interface required for TasksResource to be able to
// create, start, attach, inspect, and remove Docker containers.
type Docker interface {
	CreateContainer(docker.CreateContainerOptions) (*docker.Container, error)
	StartContainer(string, *docker.HostConfig) error
	AttachToContainer(docker.AttachToContainerOptions) error
	InspectContainer(string) (*docker.Container, error)
	RemoveContainer(docker.RemoveContainerOptions) error
}
