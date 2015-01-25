package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmorton/deployster/server/mocks"
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TasksResourceTestSuite struct {
	suite.Suite
	Subject          TasksResource
	DockerClientMock *mocks.DockerClient
	Service          *DeploysterService
}

var validRequestBody []byte = []byte(`{"task":{"version":"abc123", "command":"bundle exec rake db:migrate"}}`)

func (suite *TasksResourceTestSuite) SetupSuite() {
	suite.Service = NewDeploysterService("0.0.0.0:3000", "v1.0", "username", "password", "mmmhm")
}

func (suite *TasksResourceTestSuite) SetupTest() {
	suite.DockerClientMock = new(mocks.DockerClient)
	suite.Subject = TasksResource{suite.DockerClientMock, "mmmhm"}
}

func (suite *TasksResourceTestSuite) TestCreateTellsDockerToCreateContainer() {
	suite.setupSuccessfulDockerMock()
	req, _ := http.NewRequest("POST", "http://example.com/services/carousel/tasks?name=carousel", bytes.NewBuffer(validRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Subject.Create(w, req)

	suite.DockerClientMock.Mock.AssertCalled(suite.T(), "CreateContainer", docker.CreateContainerOptions{
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

	suite.DockerClientMock.Mock.AssertCalled(suite.T(), "StartContainer", "c0c0c0c0c0", &docker.HostConfig{})
}

func (suite *TasksResourceTestSuite) TestCreateAttachesToContainer() {
	suite.setupSuccessfulDockerMock()
	req, _ := http.NewRequest("POST", "http://example.com/services/carousel/tasks?name=carousel", bytes.NewBuffer(validRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Subject.Create(w, req)
	fw := newFlushWriter(w)

	suite.DockerClientMock.Mock.AssertCalled(suite.T(), "AttachToContainer", docker.AttachToContainerOptions{
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

	suite.DockerClientMock.Mock.AssertCalled(suite.T(), "InspectContainer", "c0c0c0c0c0")
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
	suite.DockerClientMock.On("CreateContainer", mock.AnythingOfType("docker.CreateContainerOptions")).Return(&docker.Container{ID: "c123d123bb"}, nil)
	suite.DockerClientMock.On("StartContainer", "c123d123bb", &docker.HostConfig{}).Return(errors.New("failed"))

	req, _ := http.NewRequest("POST", "http://example.com/services/carousel/tasks?name=carousel", bytes.NewBuffer(validRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.Subject.Create(w, req)

	assert.Equal(suite.T(), 500, w.Code)
}

func TestTasksResourceTestSuite(t *testing.T) {
	suite.Run(t, new(TasksResourceTestSuite))
}

func (suite *TasksResourceTestSuite) setupSuccessfulDockerMock() {
	suite.DockerClientMock.On("CreateContainer", mock.AnythingOfType("docker.CreateContainerOptions")).Return(&docker.Container{ID: "c0c0c0c0c0"}, nil)
	suite.DockerClientMock.On("StartContainer", "c0c0c0c0c0", &docker.HostConfig{}).Return(nil)
	suite.DockerClientMock.On("AttachToContainer", mock.AnythingOfType("docker.AttachToContainerOptions")).Return(nil)
	suite.DockerClientMock.On("InspectContainer", "c0c0c0c0c0").Return(&docker.Container{}, nil)
	suite.DockerClientMock.On("RemoveContainer", mock.AnythingOfType("docker.RemoveContainerOptions")).Return(nil)
}
