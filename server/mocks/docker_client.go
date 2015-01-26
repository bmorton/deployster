package mocks

import "github.com/stretchr/testify/mock"
import "github.com/fsouza/go-dockerclient"

type DockerClient struct {
	mock.Mock
}

func (m *DockerClient) RemoveContainer(_a0 docker.RemoveContainerOptions) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *DockerClient) CreateContainer(_a0 docker.CreateContainerOptions) (*docker.Container, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(*docker.Container)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DockerClient) StartContainer(_a0 string, _a1 *docker.HostConfig) error {
	ret := m.Called(_a0, _a1)

	r0 := ret.Error(0)

	return r0
}
func (m *DockerClient) AttachToContainer(_a0 docker.AttachToContainerOptions) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *DockerClient) InspectContainer(_a0 string) (*docker.Container, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(*docker.Container)
	r1 := ret.Error(1)

	return r0, r1
}
