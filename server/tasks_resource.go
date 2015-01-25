package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fsouza/go-dockerclient"
)

type TasksResource struct {
	Docker            *docker.Client
	DockerHubUsername string
}

type TaskRequest struct {
	Task Task `json:"task"`
}

type TaskResponse struct {
	Output string `json:"output"`
}

type Task struct {
	Version string `json:"version"`
	Command string `json:"command"`
}

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

func (tr *TasksResource) streamContainerOutput(containerID string, writer io.Writer) error {
	fw := flushWriter{w: writer}
	if f, ok := writer.(http.Flusher); ok {
		fw.f = f
	}
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
