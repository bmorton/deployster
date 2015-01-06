package main

import (
	"github.com/bmorton/deployster/fleet"
	"github.com/rcrowley/go-tigertonic"
	"net/http"
)

type DeploysterService struct {
	AppVersion string
	Listen     string
	RootMux    *tigertonic.TrieServeMux
	Mux        *tigertonic.TrieServeMux
	Server     *tigertonic.Server
}

func NewDeploysterService(listen string, version string) *DeploysterService {
	service := DeploysterService{Listen: listen, AppVersion: version}
	service.RootMux = tigertonic.NewTrieServeMux()
	service.Mux = tigertonic.NewTrieServeMux()
	service.RootMux.HandleNamespace("/v1", service.Mux)
	service.Server = tigertonic.NewServer(service.Listen, tigertonic.ApacheLogged(service.RootMux))
	service.ConfigureRoutes()

	return &service
}

func (self *DeploysterService) ConfigureRoutes() {
	fleetClient := fleet.NewClient("/var/run/fleet.sock")
	deploys := DeploysResource{&fleetClient}
	units := UnitsResource{&fleetClient}

	self.Mux.Handle("GET", "/version", authenticated(tigertonic.Version(self.AppVersion)))
	self.Mux.Handle("POST", "/services/{name}/deploys", authenticated(tigertonic.Marshaled(deploys.Create)))
	self.Mux.Handle("DELETE", "/services/{name}/deploys/{version}", authenticated(tigertonic.Marshaled(deploys.Destroy)))
	self.Mux.Handle("GET", "/services/{name}/units", authenticated(tigertonic.Marshaled(units.Index)))
}

func (self *DeploysterService) ListenAndServe() {
	self.Server.ListenAndServe()
}

func authenticated(h http.Handler) tigertonic.FirstHandler {
	return tigertonic.HTTPBasicAuth(
		map[string]string{username: password},
		"Deployster",
		h)
}
