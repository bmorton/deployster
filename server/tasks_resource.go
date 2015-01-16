package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"

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

func (tr *TasksResource) Create(u *url.URL, h http.Header, req *TaskRequest) (int, http.Header, *TaskResponse, error) {
	serviceName := u.Query().Get("name")
	taskName := fmt.Sprintf("%s-%s-task", serviceName, req.Task.Version)
	imageName := fmt.Sprintf("%s/%s:%s", tr.DockerHubUsername, serviceName, req.Task.Version)

	container, err := tr.runContainer(taskName, imageName, req.Task.Command)
	if err != nil {
		return http.StatusInternalServerError, nil, &TaskResponse{}, err
	}

	output, err := tr.readContainerOutput(container.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, &TaskResponse{}, err
	}

	err = tr.Docker.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
	if err != nil {
		log.Printf("WARNING: Container could not be cleaned up (%s)", err)
	}

	return http.StatusCreated, nil, &TaskResponse{output}, nil
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

func (tr *TasksResource) readContainerOutput(containerID string) (string, error) {
	var buf bytes.Buffer
	err := tr.Docker.AttachToContainer(docker.AttachToContainerOptions{
		Container:    containerID,
		OutputStream: &buf,
		Logs:         true,
		Stdout:       true,
		Stderr:       true,
		Stream:       true,
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
