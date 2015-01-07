package main

import (
	"flag"
	"github.com/bmorton/deployster/server"
	"log"
)

// A version string that can be set at compile time with:
//  -ldflags "-X main.AppVersion VERSION"
var AppVersion string

// Options that are configurable from flags
var listen string
var dockerHubUsername string
var username string
var password string

func init() {
	flag.StringVar(&listen, "listen", "0.0.0.0:3000", "Specifies the IP and port that the HTTP server will listen on")
	flag.StringVar(&dockerHubUsername, "docker-hub-username", "", "The username of the Docker Hub account that all deployable images are hosted under")
	flag.StringVar(&username, "username", "deployster", "Username that will be used to authenticate with Deployster via HTTP basic auth")
	flag.StringVar(&password, "password", "mmmhm", "Password that will be used to authenticate with Deployster via HTTP basic auth")
	flag.Parse()
}

func main() {
	log.Printf("Starting deployster on %s...\n", listen)
	service := server.NewDeploysterService(listen, AppVersion, username, password, dockerHubUsername)
	service.ListenAndServe()
}
