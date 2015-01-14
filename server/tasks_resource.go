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

type Task struct {
	Version string `json:"version"`
	Command string `json:"command"`
}

func (tr *TasksResource) Create(u *url.URL, h http.Header, req *TaskRequest) (int, http.Header, interface{}, error) {
	serviceName := u.Query().Get("name")
	taskName := fmt.Sprintf("%s-%s-task", serviceName, req.Task.Version)
	imageName := fmt.Sprintf("%s/%s:%s", tr.DockerHubUsername, serviceName, req.Task.Version)

	container, err := tr.Docker.CreateContainer(docker.CreateContainerOptions{
		Name: taskName,
		Config: &docker.Config{
			Image:        imageName,
			Cmd:          []string{req.Task.Command},
			AttachStdout: true,
			AttachStderr: true,
		},
	})
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	log.Printf("%#v", container)
	tr.Docker.StartContainer(container.ID, &docker.HostConfig{})

	var buf bytes.Buffer
	err = tr.Docker.AttachToContainer(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: &buf,
		Logs:         true,
		Stdout:       true,
		Stderr:       true,
		Stream:       true,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	log.Println(buf.String())

	err = tr.Docker.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	return http.StatusCreated, nil, nil, nil
}
