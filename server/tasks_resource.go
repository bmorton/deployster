package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fsouza/go-dockerclient"
)

// TasksResource is the HTTP resource responsible for launching new tasks via
// the Docker API in an opinionated and conventional way.  Using the provided
// DockerHubUsername and the payload passed to the Create endpoint, we can
// construct the image name to pull from the Docker Hub Registry so that the
// task can be launched.
type TasksResource struct {
	Docker            DockerClient
	DockerHubUsername string
}

// DockerClient is the interface required for TasksResource to be able to
// create, start, attach, inspect, and remove Docker containers.
type DockerClient interface {
	CreateContainer(docker.CreateContainerOptions) (*docker.Container, error)
	StartContainer(string, *docker.HostConfig) error
	AttachToContainer(docker.AttachToContainerOptions) error
	InspectContainer(string) (*docker.Container, error)
	RemoveContainer(docker.RemoveContainerOptions) error
}

// TaskRequest is the top-level wrapper for the Task in the JSON payload sent by
// the client.
type TaskRequest struct {
	Task Task `json:"task"`
}

// Task is the JSON payload required to launch a new task.
type Task struct {
	Version string `json:"version"`
	Command string `json:"command"`
}

// Create handles launching new tasks and streaming the output back over the
// http.ResponseWriter.  It expects a JSON payload that can be decoded into a
// TaskRequest.
//
// If an error occurs decoding the JSON or creating/running the container, an
// Internal Server Error will be returned in the response.  However, if an error
// occurs after this point, we've already sent a 200 OK and started streaming
// the response body.  This means the task was successfully launched, but the
// task could have possibly errored out.  At the end of the task output, the
// exit code of the task will be printed so that it can be handled by the client
// if necessary.
func (tr *TasksResource) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	decoder := json.NewDecoder(r.Body)
	var req TaskRequest
	err := decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("ERROR: %s\n", err))
		return
	}
	serviceName := r.URL.Query().Get("name")
	taskName := fmt.Sprintf("%s-%s-task", serviceName, req.Task.Version)
	imageName := fmt.Sprintf("%s/%s:%s", tr.DockerHubUsername, serviceName, req.Task.Version)

	container, err := tr.runContainer(taskName, imageName, req.Task.Command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("ERROR: %s\n", err))
		return
	}
	w.WriteHeader(http.StatusOK)

	err = tr.streamContainerOutput(container.ID, w)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("ERROR: %s\n", err))
	}

	err = tr.Docker.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
	if err != nil {
		io.WriteString(w, fmt.Sprintf("WARNING: Container could not be cleaned up (%s)\n", err))
	}
}

// runContainer creates and starts a Docker container using the provided task
// name, image name, and command.
func (tr *TasksResource) runContainer(taskName string, imageName string, command string) (*docker.Container, error) {
	container, err := tr.Docker.CreateContainer(docker.CreateContainerOptions{
		Name: taskName,
		Config: &docker.Config{
			Image:        imageName,
			Cmd:          []string{command},
			AttachStdout: true,
			AttachStderr: true,
		},
	})
	if err != nil {
		return &docker.Container{}, err
	}

	err = tr.Docker.StartContainer(container.ID, &docker.HostConfig{})
	if err != nil {
		return &docker.Container{}, err
	}

	return container, nil
}

// streamContainerOutput attaches to the container ID's STDOUT/STDERR and
// streams the output to the provided io.Writer wrapped in a flushWriter so that
// we can continuously flush the output to the client as its provided from the
// Docker API.
func (tr *TasksResource) streamContainerOutput(containerID string, writer io.Writer) error {
	fw := newFlushWriter(writer)
	err := tr.Docker.AttachToContainer(docker.AttachToContainerOptions{
		Container:    containerID,
		OutputStream: &fw,
		ErrorStream:  &fw,
		Logs:         true,
		Stdout:       true,
		Stderr:       true,
		Stream:       true,
	})

	if err != nil {
		return err
	}

	container, err := tr.Docker.InspectContainer(containerID)
	if err != nil {
		return err
	}
	io.WriteString(&fw, fmt.Sprintf("\nExited (%d) %s\n", container.State.ExitCode, container.State.Error))

	return nil
}
