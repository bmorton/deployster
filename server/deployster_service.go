package server

import (
	"net/http"

	"github.com/bmorton/deployster/fleet"
	"github.com/fsouza/go-dockerclient"
	"github.com/rcrowley/go-tigertonic"
)

// DeploysterService is the HTTP server that ties together all the resources,
// configures routing and dependencies, and authenticates requests.  The server
// will listen for requests in the address:port that is passed as the listen
// string.  The provided username/password will be used for HTTP basic auth for
// all requests.  An image prefix can be either a private registry address:port
// or a username on the public registry (basically something that'll be appended
// to the service name, e.g. mmmmhm/servicename or my.registry:5000/servicename).
type DeploysterService struct {
	AppVersion  string
	Listen      string
	Username    string
	Password    string
	ImagePrefix string
	RootMux     *tigertonic.TrieServeMux
	Mux         *tigertonic.TrieServeMux
	Server      *tigertonic.Server
}

// NewDeploysterService returns a configured DeploysterService, ready to listen
// for HTTP requests via the provided listen string.
func NewDeploysterService(listen string, version string, username string, password string, imagePrefix string) *DeploysterService {
	service := DeploysterService{
		Listen:      listen,
		AppVersion:  version,
		Username:    username,
		Password:    password,
		ImagePrefix: imagePrefix,
	}
	service.RootMux = tigertonic.NewTrieServeMux()
	service.Mux = tigertonic.NewTrieServeMux()
	service.RootMux.HandleNamespace("/v1", service.Mux)
	service.Server = tigertonic.NewServer(service.Listen, tigertonic.ApacheLogged(service.RootMux))
	service.ConfigureRoutes()

	return &service
}

// ConfigureRoutes sets up resources and their dependencies so that we can
// configure all the HTTP routes that will be supported by the server.
func (ds *DeploysterService) ConfigureRoutes() {
	fleetClient := fleet.NewClient("/var/run/fleet.sock")
	dockerClient, _ := docker.NewClient("unix:///var/run/docker.sock")
	deploys := DeploysResource{&fleetClient, ds.ImagePrefix}
	units := UnitsResource{&fleetClient}
	tasks := TasksResource{dockerClient, ds.ImagePrefix}

	ds.Mux.Handle("GET", "/version", ds.authenticated(tigertonic.Version(ds.AppVersion)))
	ds.Mux.Handle("POST", "/services/{name}/deploys", ds.authenticated(tigertonic.Marshaled(deploys.Create)))
	ds.Mux.Handle("DELETE", "/services/{name}/deploys/{version}", ds.authenticated(tigertonic.Marshaled(deploys.Destroy)))
	ds.Mux.Handle("GET", "/services/{name}/units", ds.authenticated(tigertonic.Marshaled(units.Index)))
	ds.Mux.Handle("POST", "/services/{name}/tasks", ds.authenticated(http.HandlerFunc(tasks.Create)))
}

// ListenAndServe starts the HTTP server.
func (ds *DeploysterService) ListenAndServe() error {
	return ds.Server.ListenAndServe()
}

// ListenAndServe starts the HTTPS server with the given certificate and key.
func (ds *DeploysterService) ListenAndServeTLS(certPath string, keyPath string) error {
	return ds.Server.ListenAndServeTLS(certPath, keyPath)
}

// Close gracefully stops listening for new requests
func (ds *DeploysterService) Close() error {
	return ds.Server.Close()
}

// authenticated wraps an HTTP handler with HTTP basic authentication
func (ds *DeploysterService) authenticated(h http.Handler) tigertonic.FirstHandler {
	return tigertonic.HTTPBasicAuth(
		map[string]string{ds.Username: ds.Password},
		"Deployster",
		h)
}
