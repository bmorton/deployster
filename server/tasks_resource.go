package server

import (
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

	return http.StatusCreated, nil, nil, nil
}
