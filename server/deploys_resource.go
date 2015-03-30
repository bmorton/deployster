package server

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/bmorton/deployster/units"
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

// UnitTemplate is the view model that is passed to the template parser that
// renders a unit file.
type UnitTemplate struct {
	Name        string
	Version     string
	ImagePrefix string
	Timestamp   string
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
	req.Deploy.ServiceName = u.Query().Get("name")

	if req.Deploy.Timestamp == "" {
		req.Deploy.Timestamp = time.Now().UTC().Format("2006.01.02-15.04.05")
	}

	allUnits, err := dr.Fleet.Units()
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, nil, nil, err
	}

	previousVersions := units.FindTimestampedServiceVersions(req.Deploy.ServiceName, allUnits)
	previousUnits := units.FindServiceUnits(req.Deploy.ServiceName, "", allUnits)

	if req.Deploy.DestroyPrevious {
		if len(previousVersions) > 1 {
			return http.StatusBadRequest, nil, nil, errors.New("Too many versions are running.  Destroying previous units is not supported when more than one version is currently running.")
		}

		if req.Deploy.InstanceCount != 0 && req.Deploy.InstanceCount < len(previousUnits) {
			return http.StatusBadRequest, nil, nil, errors.New("A greater number of instances than what was specified is already running.  Make sure this number is less than or equal to the number already running or disable destroying previous units.")
		}

		for _, unit := range previousUnits {
			oldInstance := &ServiceInstance{Name: req.Deploy.ServiceName, Version: unit.Version, Timestamp: unit.Timestamp, Instance: unit.Instance}
			newInstance := req.Deploy.ServiceInstance(unit.Instance)
			log.Printf("Launching watcher for %s.\n", oldInstance.FleetUnitName())
			go dr.destroyPrevious(oldInstance.FleetUnitName(), newInstance.FleetUnitName(), destroyPreviousCheckDelay)
		}
	}

	req.Deploy.InstanceCount = determineNumberOfInstances(req.Deploy.InstanceCount, len(previousVersions), len(previousUnits))
	err = dr.startUnits(&req.Deploy)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	return http.StatusCreated, nil, nil, nil
}

// Destroy is the DELETE endpoint for destroying the units associated with
// the service name and version provided.  It will destroy all instances of a
// unit that exists within Fleet.  If a timestamp query parameter is provided,
// only units that match that timestamp will be destroyed.
//
// This function assumes that it is nested inside
// `/services/{name}/versions/{version}` and that Tigertonic is extracting the
// service name/version and providing it via query params.
func (dr *DeploysResource) Destroy(u *url.URL, h http.Header, req interface{}) (int, http.Header, interface{}, error) {
	deploy := &Deploy{
		ServiceName: u.Query().Get("name"),
		Version:     u.Query().Get("version"),
	}

	allUnits, err := dr.Fleet.Units()
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, nil, nil, err
	}
	serviceUnits := units.FindServiceUnits(deploy.ServiceName, deploy.Version, allUnits)

	for _, unit := range serviceUnits {
		if shouldDestroyUnit(u.Query().Get("timestamp"), unit.Timestamp) {
			instance := deploy.ServiceInstance(unit.Instance)
			instance.Timestamp = unit.Timestamp
			err := dr.Fleet.DestroyUnit(instance.FleetUnitName())
			if err != nil {
				return http.StatusInternalServerError, nil, nil, err
			}
		}
	}

	return http.StatusNoContent, nil, nil, nil
}

// startUnits is a helper function for ensuring that Fleet has all the units
// configured and for launching those units.
func (dr *DeploysResource) startUnits(deploy *Deploy) error {
	options := getUnitOptions(UnitTemplate{deploy.ServiceName, deploy.Version, dr.ImagePrefix, deploy.Timestamp})

	// Make sure all units exist before we start setting their target states
	for i := 1; i <= deploy.InstanceCount; i++ {
		instance := deploy.ServiceInstance(strconv.Itoa(i))
		log.Printf("Creating %s.\n", instance.FleetUnitName())
		err := dr.Fleet.CreateUnit(&schema.Unit{Name: instance.FleetUnitName(), Options: options})
		if err != nil {
			return err
		}
	}

	// Now that all the units exist, we can launch each of them
	for i := 1; i <= deploy.InstanceCount; i++ {
		instance := deploy.ServiceInstance(strconv.Itoa(i))
		log.Printf("Launching %s.\n", instance.FleetUnitName())
		err := dr.Fleet.SetUnitTargetState(instance.FleetUnitName(), "launched")
		if err != nil {
			return err
		}
	}

	return nil
}

// determineNumberOfInstances is a helper function to either return the number
// of instances specified or provide a default value based on the number of
// running versions and units.
func determineNumberOfInstances(instanceCount int, numberOfVersions int, numberOfUnits int) int {
	if instanceCount != 0 {
		return instanceCount
	}

	if numberOfVersions == 1 {
		return numberOfUnits
	} else {
		return 1
	}
}

// getUnitOptions renders the unit file and converts it to an array of
// UnitOption structs.
func getUnitOptions(unitViewTemplate UnitTemplate) []*schema.UnitOption {
	var unitTemplate bytes.Buffer
	t, _ := template.New("test").Parse(dockerUnitTemplate)
	t.Execute(&unitTemplate, unitViewTemplate)

	unitFile, _ := unit.NewUnitFile(unitTemplate.String())

	return schema.MapUnitFileToSchemaUnitOptions(unitFile)
}

// shouldDestroyUnit takes an optional timestamp (from the query string) and, if
// specified, ensures that it matches the unitTimestamp.  If the optional
// timestamp is left blank, we'll return true.  If the optional timestamp is
// present and it doesn't match the timestamp, we'll return false.
func shouldDestroyUnit(blankOrTimestampToMatch string, unitTimestamp string) bool {
	if blankOrTimestampToMatch == "" {
		return true
	} else if blankOrTimestampToMatch == unitTimestamp {
		return true
	}
	return false
}

// destroyPrevious is responsible for watching for a new version of a service to
// complete launching on the `destroyPreviousCheckDelay` interval so that it can
// fire off a request to destroy the previous version.  If this doesn't complete
// before the destroyPreviousCheckTimeout, the attempt will be abandoned.
func (dr *DeploysResource) destroyPrevious(previousFleetUnit string, currentFleetUnit string, checkDelay time.Duration) {
	timeoutChan := make(chan bool, 1)

	go func() {
		time.Sleep(destroyPreviousCheckTimeout)
		timeoutChan <- true
	}()

	for {
		startCheck := time.After(checkDelay)

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
							log.Println(err)
						}
						return
					}
					if state.SystemdSubState == "failed" {
						log.Printf("%s failed to launch, bailing out.\n", currentFleetUnit)
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
