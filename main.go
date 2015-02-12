package main

import (
	"flag"
	"github.com/bmorton/deployster/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// A version string that can be set at compile time with:
//  -ldflags "-X main.AppVersion VERSION"
var AppVersion string

// Options that are configurable from flags
var listen string
var dockerHubUsername string
var registryURL string
var username string
var password string
var certPath string
var keyPath string

func init() {
	flag.StringVar(&listen, "listen", "0.0.0.0:3000", "Specifies the IP and port that the HTTP server will listen on")
	flag.StringVar(&dockerHubUsername, "docker-hub-username", "deployster", "The username of the Docker Hub account that all deployable images are hosted under")
	flag.StringVar(&registryURL, "registry-url", "", "If using a private registry, this is the address:port of that registry (if supplied, docker-hub-username will be ignored)")
	flag.StringVar(&username, "username", "deployster", "Username that will be used to authenticate with Deployster via HTTP basic auth")
	flag.StringVar(&password, "password", "mmmhm", "Password that will be used to authenticate with Deployster via HTTP basic auth")
	flag.StringVar(&certPath, "cert", "", "Path to certificate to be used for serving HTTPS")
	flag.StringVar(&keyPath, "key", "", "Path to private key to bse used for serving HTTPS")
	flag.Parse()
}

func main() {
	var imagePrefix string
	if registryURL != "" {
		log.Printf("Starting deployster on %s using private registry at %s...\n", listen, registryURL)
		imagePrefix = registryURL
	} else {
		log.Printf("Starting deployster on %s using the public Docker Hub registry with user %s...\n", listen, dockerHubUsername)
		imagePrefix = dockerHubUsername
	}
	service := server.NewDeploysterService(listen, AppVersion, username, password, imagePrefix)

	go func() {
		var err error
		if certPath != "" && keyPath != "" {
			log.Println("Certificate and private key provided, HTTPS enabled.")
			err = service.ListenAndServeTLS(certPath, keyPath)
		} else {
			err = service.ListenAndServe()
		}
		if err != nil {
			log.Println(err)
		}
	}()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	log.Println(<-ch)
	service.Close()
}
