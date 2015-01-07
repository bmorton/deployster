package server

import (
	"github.com/bmorton/deployster/fleet"
	"github.com/rcrowley/go-tigertonic"
	"net/http"
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

func (self *DeploysterService) ConfigureRoutes() {
	fleetClient := fleet.NewClient("/var/run/fleet.sock")
	deploys := DeploysResource{&fleetClient, self.DockerHubUsername}
	units := UnitsResource{&fleetClient}

	self.Mux.Handle("GET", "/version", self.authenticated(tigertonic.Version(self.AppVersion)))
	self.Mux.Handle("POST", "/services/{name}/deploys", self.authenticated(tigertonic.Marshaled(deploys.Create)))
	self.Mux.Handle("DELETE", "/services/{name}/deploys/{version}", self.authenticated(tigertonic.Marshaled(deploys.Destroy)))
	self.Mux.Handle("GET", "/services/{name}/units", self.authenticated(tigertonic.Marshaled(units.Index)))
}

func (self *DeploysterService) ListenAndServe() {
	self.Server.ListenAndServe()
}

func (self *DeploysterService) authenticated(h http.Handler) tigertonic.FirstHandler {
	return tigertonic.HTTPBasicAuth(
		map[string]string{self.Username: self.Password},
		"Deployster",
		h)
}
