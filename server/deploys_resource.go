package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/bmorton/deployster/fleet"
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
	Fleet       fleet.Client
	ImagePrefix string
}

type DeployRequest struct {
	Deploy Deploy `json:"deploy"`
}

type Deploy struct {
	Version         string `json:"version"`
	DestroyPrevious bool   `json:"destroy_previous"`
}

type UnitTemplate struct {
	Name        string
	Version     string
	ImagePrefix string
}

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

	_, err := dr.Fleet.StartUnit(fleetName, options)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	return http.StatusCreated, nil, nil, nil
}

func (dr *DeploysResource) Destroy(u *url.URL, h http.Header, req interface{}) (int, http.Header, interface{}, error) {
	fleetName := fleetServiceName(u.Query().Get("name"), u.Query().Get("version"))

	_, err := dr.Fleet.DestroyUnit(fleetName)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	return http.StatusNoContent, nil, nil, nil
}

func getUnitOptions(name string, version string, imagePrefix string) []fleet.UnitOption {
	var unitTemplate bytes.Buffer
	t, _ := template.New("test").Parse(dockerUnitTemplate)
	t.Execute(&unitTemplate, UnitTemplate{name, version, imagePrefix})

	unitFile, _ := unit.NewUnitFile(unitTemplate.String())

	return schemaToLocalUnit(schema.MapUnitFileToSchemaUnitOptions(unitFile))
}

func schemaToLocalUnit(options []*schema.UnitOption) []fleet.UnitOption {
	convertedOptions := []fleet.UnitOption{}
	for _, o := range options {
		convertedOptions = append(convertedOptions, fleet.UnitOption{
			Section: o.Section,
			Name:    o.Name,
			Value:   o.Value,
		})
	}
	return convertedOptions
}

func fleetServiceName(name string, version string) string {
	return fmt.Sprintf("%s-%s@1.service", name, version)
}

func (dr *DeploysResource) destroyPrevious(name string, previousVersion string, currentVersion string) {
	timeoutChan := make(chan bool, 1)

	go func() {
		time.Sleep(destroyPreviousCheckTimeout)
		timeoutChan <- true
	}()

	for {
		startCheck := time.After(destroyPreviousCheckDelay)

		select {
		case <-startCheck:
			log.Printf("Checking if %s:%s has finished launching...\n", name, currentVersion)
			state, err := dr.Fleet.UnitState(fleetServiceName(name, currentVersion))
			if err != nil {
				log.Println(err)
				break
			}
			if state.SubState == "running" {
				log.Printf("%s:%s has launched, destroying %s.\n", name, currentVersion, previousVersion)
				_, err := dr.Fleet.DestroyUnit(fleetServiceName(name, previousVersion))
				if err != nil {
					log.Printf("%#v\n", err)
				}
				return
			}
			log.Printf("%s:%s isn't running (currently %s).  Trying again in %s.\n", name, currentVersion, state.SubState, destroyPreviousCheckDelay)
		case <-timeoutChan:
			log.Printf("Timed out trying to destroy %s:%s after %s.\n", name, previousVersion, destroyPreviousCheckTimeout)
			return
		}
	}
}
