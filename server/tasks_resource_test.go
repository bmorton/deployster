package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmorton/deployster/clients/mocks"
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TasksResourceTestSuite struct {
	suite.Suite
	Subject    TasksResource
	DockerMock *mocks.Docker
	Service    *DeploysterService
}

var validRequestBody []byte = []byte(`{"task":{"version":"abc123", "command":"bundle exec rake db:migrate"}}`)

func (suite *TasksResourceTestSuite) SetupSuite() {
	suite.Service = NewDeploysterService("0.0.0.0:3000", "v1.0", "username", "password", "mmmhm")
}

func (suite *TasksResourceTestSuite) SetupTest() {
	suite.DockerMock = new(mocks.Docker)
	suite.Subject = TasksResource{suite.DockerMock, "mmmhm"}
}

func (suite *TasksResourceTestSuite) TestCreateTellsDockerToCreateContainer() {
	suite.setupSuccessfulDockerMock()
	req, _ := http.NewRequest("POST", "http://example.com/services/carousel/tasks?name=carousel", bytes.NewBuffer(validRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Subject.Create(w, req)

	suite.DockerMock.Mock.AssertCalled(suite.T(), "CreateContainer", docker.CreateContainerOptions{
		Name: "carousel-abc123-task",
		Config: &docker.Config{
			Image:        "mmmhm/carousel:abc123",
			Cmd:          []string{"bundle exec rake db:migrate"},
			AttachStdout: true,
			AttachStderr: true,
		},
	})
}

func (suite *TasksResourceTestSuite) TestCreateTellsDockerToStartContainer() {
	suite.setupSuccessfulDockerMock()
	req, _ := http.NewRequest("POST", "http://example.com/services/carousel/tasks?name=carousel", bytes.NewBuffer(validRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Subject.Create(w, req)

	suite.DockerMock.Mock.AssertCalled(suite.T(), "StartContainer", "c0c0c0c0c0", &docker.HostConfig{})
}

func (suite *TasksResourceTestSuite) TestCreateAttachesToContainer() {
	suite.setupSuccessfulDockerMock()
	req, _ := http.NewRequest("POST", "http://example.com/services/carousel/tasks?name=carousel", bytes.NewBuffer(validRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Subject.Create(w, req)
	fw := newFlushWriter(w)

	suite.DockerMock.Mock.AssertCalled(suite.T(), "AttachToContainer", docker.AttachToContainerOptions{
		Container:    "c0c0c0c0c0",
		OutputStream: &fw,
		ErrorStream:  &fw,
		Logs:         true,
		Stdout:       true,
		Stderr:       true,
		Stream:       true,
	})
}

func (suite *TasksResourceTestSuite) TestCreateInspectsContainerForExitStatus() {
	suite.setupSuccessfulDockerMock()
	req, _ := http.NewRequest("POST", "http://example.com/services/carousel/tasks?name=carousel", bytes.NewBuffer(validRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Subject.Create(w, req)

	suite.DockerMock.Mock.AssertCalled(suite.T(), "InspectContainer", "c0c0c0c0c0")
}

func (suite *TasksResourceTestSuite) TestCreateIsSuccessfulWhenContainerStarted() {
	suite.setupSuccessfulDockerMock()
	req, _ := http.NewRequest("POST", "http://example.com/services/carousel/tasks?name=carousel", bytes.NewBuffer(validRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Subject.Create(w, req)

	assert.Equal(suite.T(), 200, w.Code)
}

func (suite *TasksResourceTestSuite) TestCreateReturnsServerErrorWhenContainerFailsToStart() {
	suite.DockerMock.On("CreateContainer", mock.AnythingOfType("docker.CreateContainerOptions")).Return(&docker.Container{ID: "c123d123bb"}, nil)
	suite.DockerMock.On("StartContainer", "c123d123bb", &docker.HostConfig{}).Return(errors.New("failed"))

	req, _ := http.NewRequest("POST", "http://example.com/services/carousel/tasks?name=carousel", bytes.NewBuffer(validRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Subject.Create(w, req)

	assert.Equal(suite.T(), 500, w.Code)
}

func (suite *TasksResourceTestSuite) TestCreateReturnsExitCodeOnSuccess() {
	suite.DockerMock.On("CreateContainer", mock.AnythingOfType("docker.CreateContainerOptions")).Return(&docker.Container{ID: "c0c0c0c0c0"}, nil)
	suite.DockerMock.On("StartContainer", "c0c0c0c0c0", &docker.HostConfig{}).Return(nil)
	suite.DockerMock.On("AttachToContainer", mock.AnythingOfType("docker.AttachToContainerOptions")).Return(nil)
	suite.DockerMock.On("InspectContainer", "c0c0c0c0c0").Return(&docker.Container{
		State: docker.State{
			ExitCode: 0,
		},
	}, nil)
	suite.DockerMock.On("RemoveContainer", mock.AnythingOfType("docker.RemoveContainerOptions")).Return(nil)

	req, _ := http.NewRequest("POST", "http://example.com/services/carousel/tasks?name=carousel", bytes.NewBuffer(validRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Subject.Create(w, req)

	assert.Equal(suite.T(), "\nExited (0) \n", w.Body.String())
}

func (suite *TasksResourceTestSuite) TestCreateReturnsExitCodeOnFailure() {
	suite.DockerMock.On("CreateContainer", mock.AnythingOfType("docker.CreateContainerOptions")).Return(&docker.Container{ID: "c0c0c0c0c0"}, nil)
	suite.DockerMock.On("StartContainer", "c0c0c0c0c0", &docker.HostConfig{}).Return(nil)
	suite.DockerMock.On("AttachToContainer", mock.AnythingOfType("docker.AttachToContainerOptions")).Return(nil)
	suite.DockerMock.On("InspectContainer", "c0c0c0c0c0").Return(&docker.Container{
		State: docker.State{
			ExitCode: 127,
			Error:    "Something went wrong",
		},
	}, nil)
	suite.DockerMock.On("RemoveContainer", mock.AnythingOfType("docker.RemoveContainerOptions")).Return(nil)

	req, _ := http.NewRequest("POST", "http://example.com/services/carousel/tasks?name=carousel", bytes.NewBuffer(validRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Subject.Create(w, req)

	assert.Equal(suite.T(), "\nExited (127) Something went wrong\n", w.Body.String())
}

func TestTasksResourceTestSuite(t *testing.T) {
	suite.Run(t, new(TasksResourceTestSuite))
}

func (suite *TasksResourceTestSuite) setupSuccessfulDockerMock() {
	suite.DockerMock.On("CreateContainer", mock.AnythingOfType("docker.CreateContainerOptions")).Return(&docker.Container{ID: "c0c0c0c0c0"}, nil)
	suite.DockerMock.On("StartContainer", "c0c0c0c0c0", &docker.HostConfig{}).Return(nil)
	suite.DockerMock.On("AttachToContainer", mock.AnythingOfType("docker.AttachToContainerOptions")).Return(nil)
	suite.DockerMock.On("InspectContainer", "c0c0c0c0c0").Return(&docker.Container{}, nil)
	suite.DockerMock.On("RemoveContainer", mock.AnythingOfType("docker.RemoveContainerOptions")).Return(nil)
}
