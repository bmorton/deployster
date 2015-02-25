package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/coreos/fleet/schema"
	"github.com/coreos/fleet/unit"
)

const (
	// destroyPreviousCheckTimeout is the amount of time to attempt checking for
	// the new version of the service to boot.  If checking exceeds this time, the
	// previous version will not be destroyed.
	destroyPreviousCheckTimeout time.Duration = 5 * time.Minute

	// destroyPreviousCheckDelay is the amount of to wait between checks for the
	// boot completion of the new version.
	destroyPreviousCheckDelay time.Duration = 1 * time.Second
)

// DeploysResource is the HTTP resource responsible for creating and destroying
// deployments of services.
type DeploysResource struct {
	Fleet       FleetClient
	ImagePrefix string
}

// DeployRequest is the wrapper struct used to deserialize the JSON payload that
// is sent for creating a new deploy.
type DeployRequest struct {
	Deploy Deploy `json:"deploy"`
}

// Deploy is the struct that defines all the options for creating a new deploy
// and is wrapped by DeployRequest and deserialized in the Create function.
type Deploy struct {
	Version         string `json:"version"`
	DestroyPrevious bool   `json:"destroy_previous"`
}

// UnitTemplate is the view model that is passed to the template parser that
// renders a unit file.
type UnitTemplate struct {
	Name        string
	Version     string
	ImagePrefix string
}

// Create is the POST endpoint for kicking off a new deployment of the service
// and version provided.  It uses these parameters to spin up tasks that will
// asyncronously start new units via Fleet and optionally wait for units to
// complete launching so that it can destroy old versions of the service that
// are no longer desired.
//
// This function assumes that it is nested inside `/services/{name}`
// and that Tigertonic is extracting the service name and providing it via query
// params.
func (dr *DeploysResource) Create(u *url.URL, h http.Header, req *DeployRequest) (int, http.Header, interface{}, error) {
	serviceName := u.Query().Get("name")
	options := getUnitOptions(serviceName, req.Deploy.Version, dr.ImagePrefix)
	fleetName := fleetServiceName(serviceName, req.Deploy.Version)

	if req.Deploy.DestroyPrevious {
		units, err := dr.Fleet.Units()
		if err != nil {
			log.Printf("%#v\n", err)
			return http.StatusInternalServerError, nil, nil, err
		}
		versions := FindServiceVersions(u.Query().Get("name"), units)

		if len(versions) != 1 {
			log.Printf("Can't destroy previous versions (%d previous versions), disabling destroy.", len(versions))
			req.Deploy.DestroyPrevious = false
		} else {
			go dr.destroyPrevious(serviceName, versions[0], req.Deploy.Version)
		}
	}

	unit := schema.Unit{
		Name:         fleetName,
		DesiredState: "launched",
		Options:      options,
	}
	err := dr.Fleet.CreateUnit(&unit)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	return http.StatusCreated, nil, nil, nil
}

// Destroy is the DELETE endpoint for destroying the units associated with
// the service name and version provided.
//
// This function assumes that it is nested inside
// `/services/{name}/versions/{version}` and that Tigertonic is extracting the
// service name/version and providing it via query params.
func (dr *DeploysResource) Destroy(u *url.URL, h http.Header, req interface{}) (int, http.Header, interface{}, error) {
	fleetName := fleetServiceName(u.Query().Get("name"), u.Query().Get("version"))

	err := dr.Fleet.DestroyUnit(fleetName)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	return http.StatusNoContent, nil, nil, nil
}

// getUnitOptions renders the unit file and converts it to an array of
// UnitOption structs.
func getUnitOptions(name string, version string, imagePrefix string) []*schema.UnitOption {
	var unitTemplate bytes.Buffer
	t, _ := template.New("test").Parse(dockerUnitTemplate)
	t.Execute(&unitTemplate, UnitTemplate{name, version, imagePrefix})

	unitFile, _ := unit.NewUnitFile(unitTemplate.String())

	return schema.MapUnitFileToSchemaUnitOptions(unitFile)
}

// fleetServiceName generates a fleet unit name with the service name, version,
// and instance encoded within it.
func fleetServiceName(name string, version string) string {
	return fmt.Sprintf("%s-%s@1.service", name, version)
}

// destroyPrevious is responsible for watching for a new version of a service to
// complete launching on the `destroyPreviousCheckDelay` interval so that it can
// fire off a request to destroy the previous version.  If this doesn't complete
// before the destroyPreviousCheckTimeout, the attempt will be abandoned.
func (dr *DeploysResource) destroyPrevious(name string, previousVersion string, currentVersion string) {
	timeoutChan := make(chan bool, 1)

	go func() {
		time.Sleep(destroyPreviousCheckTimeout)
		timeoutChan <- true
	}()

	currentFleetUnit := fleetServiceName(name, currentVersion)
	previousFleetUnit := fleetServiceName(name, previousVersion)

	for {
		startCheck := time.After(destroyPreviousCheckDelay)

		select {
		case <-startCheck:
			log.Printf("Checking if %s has finished launching...\n", currentFleetUnit)
			states, err := dr.Fleet.UnitStates()
			if err != nil {
				log.Println(err)
				break
			}
			for _, state := range states {
				if state.Name == currentFleetUnit {
					if state.SystemdSubState == "running" {
						log.Printf("%s has launched, destroying %s.\n", currentFleetUnit, previousFleetUnit)
						err := dr.Fleet.DestroyUnit(previousFleetUnit)
						if err != nil {
							log.Printf("%#v\n", err)
						}
						return
					}
					log.Printf("%s isn't running (currently %s).  Trying again in %s.\n", currentFleetUnit, state.SystemdSubState, destroyPreviousCheckDelay)
				}
			}
		case <-timeoutChan:
			log.Printf("Timed out trying to destroy %s after %s.\n", previousFleetUnit, destroyPreviousCheckTimeout)
			return
		}
	}
}
