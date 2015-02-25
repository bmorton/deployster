package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fsouza/go-dockerclient"
)

// TasksResource is the HTTP resource responsible for launching new tasks via
// the Docker API in an opinionated and conventional way.  Using the provided
// ImagePrefix and the payload passed to the Create endpoint, we can construct
// the image name to pull from the Docker Hub Registry so that the task can be
// launched.
type TasksResource struct {
	Docker      DockerClient
	ImagePrefix string
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

// defaultTaskTimeout is the amount of time that we allow for a task to run
// before it is forcefully killed.  This timeout is currently set to 10 minutes.
const defaultTaskTimeout time.Duration = 600 * time.Second

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
	imageName := fmt.Sprintf("%s/%s:%s", tr.ImagePrefix, serviceName, req.Task.Version)

	container, err := tr.runContainer(taskName, imageName, req.Task.Command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("ERROR: %s\n", err))
		return
	}
	w.WriteHeader(http.StatusOK)

	err = tr.streamContainerOutputWithTimeout(container.ID, w, defaultTaskTimeout)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("ERROR: %s\n", err))
	}

	err = tr.Docker.RemoveContainer(docker.RemoveContainerOptions{
		ID:    container.ID,
		Force: true,
	})
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

// streamContainerOutputWithTimeout will spawn two goroutines: one for
// fulfilling the streamContainerOutput request and one for managing the
// timeout.
func (tr *TasksResource) streamContainerOutputWithTimeout(containerID string, writer io.Writer, timeout time.Duration) error {
	timeoutChan := make(chan bool, 1)
	go func() {
		time.Sleep(timeout)
		timeoutChan <- true
	}()

	successChan := make(chan error, 1)
	go func() {
		successChan <- tr.streamContainerOutput(containerID, writer)
	}()

	select {
	case err := <-successChan:
		if err != nil {
			return err
		}
	case <-timeoutChan:
		io.WriteString(writer, fmt.Sprintf("\nExited (124) The task timed out after %s. Forcefully removing container.\n", timeout))
	}

	return nil
}
