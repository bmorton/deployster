package server

import (
	"net/http"

	"github.com/bmorton/deployster/fleet"
	"github.com/fsouza/go-dockerclient"
	"github.com/rcrowley/go-tigertonic"
)

type DeploysterService struct {
	AppVersion        string
	Listen            string
	Username          string
	Password          string
	DockerHubUsername string
	RootMux           *tigertonic.TrieServeMux
	Mux               *tigertonic.TrieServeMux
	Server            *tigertonic.Server
}

func NewDeploysterService(listen string, version string, username string, password string, dockerHubUsername string) *DeploysterService {
	service := DeploysterService{
		Listen:            listen,
		AppVersion:        version,
		Username:          username,
		Password:          password,
		DockerHubUsername: dockerHubUsername,
	}
	service.RootMux = tigertonic.NewTrieServeMux()
	service.Mux = tigertonic.NewTrieServeMux()
	service.RootMux.HandleNamespace("/v1", service.Mux)
	service.Server = tigertonic.NewServer(service.Listen, tigertonic.ApacheLogged(service.RootMux))
	service.ConfigureRoutes()

	return &service
}

func (ds *DeploysterService) ConfigureRoutes() {
	fleetClient := fleet.NewClient("/var/run/fleet.sock")
	dockerClient, _ := docker.NewClient("unix:///var/run/docker.sock")
	deploys := DeploysResource{&fleetClient, ds.DockerHubUsername}
	units := UnitsResource{&fleetClient}
	tasks := TasksResource{dockerClient, ds.DockerHubUsername}

	ds.Mux.Handle("GET", "/version", ds.authenticated(tigertonic.Version(ds.AppVersion)))
	ds.Mux.Handle("POST", "/services/{name}/deploys", ds.authenticated(tigertonic.Marshaled(deploys.Create)))
	ds.Mux.Handle("DELETE", "/services/{name}/deploys/{version}", ds.authenticated(tigertonic.Marshaled(deploys.Destroy)))
	ds.Mux.Handle("GET", "/services/{name}/units", ds.authenticated(tigertonic.Marshaled(units.Index)))
	ds.Mux.Handle("POST", "/services/{name}/tasks", ds.authenticated(tigertonic.Marshaled(tasks.Create)))
}

func (ds *DeploysterService) ListenAndServe() error {
	return ds.Server.ListenAndServe()
}

func (ds *DeploysterService) ListenAndServeTLS(certPath string, keyPath string) error {
	return ds.Server.ListenAndServeTLS(certPath, keyPath)
}

func (ds *DeploysterService) Close() error {
	return ds.Server.Close()
}

func (ds *DeploysterService) authenticated(h http.Handler) tigertonic.FirstHandler {
	return tigertonic.HTTPBasicAuth(
		map[string]string{ds.Username: ds.Password},
		"Deployster",
		h)
}
