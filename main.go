package main

import (
	"github.com/bmorton/deployster/fleet"
	"github.com/rcrowley/go-tigertonic"
	"log"
)

// A version string that can be set at compile time with:
//  -ldflags "-X main.AppVersion VERSION"
var AppVersion string

func main() {
	listen := "0.0.0.0:3000"
	log.Printf("Starting deployster on %s...\n", listen)
	service := NewDeploysterService(listen, AppVersion)
	service.ListenAndServe()
}

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
	deploys := DeployResource{fleetClient}
	units := UnitResource{fleetClient}

	self.Mux.Handle("GET", "/version", tigertonic.Version(self.AppVersion))
	self.Mux.Handle("POST", "/services/{name}/deploys", tigertonic.Marshaled(deploys.Create))
	self.Mux.Handle("GET", "/services/{name}/units", tigertonic.Marshaled(units.Index))
}

func (self *DeploysterService) ListenAndServe() {
	self.Server.ListenAndServe()
}
