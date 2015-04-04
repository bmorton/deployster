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

	"github.com/bmorton/deployster/clients"
	"github.com/bmorton/deployster/handlers"
	"github.com/bmorton/deployster/poller"
	"github.com/bmorton/deployster/schema"
	"github.com/bmorton/deployster/units"
	fleet "github.com/coreos/fleet/schema"
	"github.com/coreos/fleet/unit"
)

// DeploysResource is the HTTP resource responsible for creating and destroying
// deployments of services.
type DeploysResource struct {
	Fleet       clients.Fleet
	ImagePrefix string
}

// DeployRequest is the wrapper struct used to deserialize the JSON payload that
// is sent for creating a new deploy.
type DeployRequest struct {
	Deploy *schema.Deploy `json:"deploy"`
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

		if len(previousVersions) == 1 {
			req.Deploy.PreviousVersion = &schema.Deploy{
				ServiceName:   req.Deploy.ServiceName,
				Version:       previousUnits[0].Version,
				Timestamp:     previousUnits[0].Timestamp,
				InstanceCount: len(previousUnits),
			}
		} else {
			req.Deploy.DestroyPrevious = false
		}
	}

	req.Deploy.InstanceCount = determineNumberOfInstances(req.Deploy.InstanceCount, len(previousVersions), len(previousUnits))
	err = dr.startUnits(req.Deploy)
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
	deploy := &schema.Deploy{
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
func (dr *DeploysResource) startUnits(deploy *schema.Deploy) error {
	options := getUnitOptions(UnitTemplate{deploy.ServiceName, deploy.Version, dr.ImagePrefix, deploy.Timestamp})

	if deploy.DestroyPrevious {
		log.Printf("Polling %s:%s.\n", deploy.ServiceName, deploy.Version)
		poller := poller.New(deploy, dr.Fleet)
		poller.AddSuccessHandler(&handlers.Destroyer{PreviousVersion: deploy.PreviousVersion, Client: dr.Fleet})
		go poller.Watch()
	}

	for i := 1; i <= deploy.InstanceCount; i++ {
		instance := deploy.ServiceInstance(strconv.Itoa(i))
		log.Printf("Creating %s.\n", instance.FleetUnitName())
		err := dr.Fleet.CreateUnit(&fleet.Unit{Name: instance.FleetUnitName(), Options: options})
		if err != nil {
			return err
		}

		log.Printf("Launching %s.\n", instance.FleetUnitName())
		err = dr.Fleet.SetUnitTargetState(instance.FleetUnitName(), "launched")
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
func getUnitOptions(unitViewTemplate UnitTemplate) []*fleet.UnitOption {
	var unitTemplate bytes.Buffer
	t, _ := template.New("test").Parse(dockerUnitTemplate)
	t.Execute(&unitTemplate, unitViewTemplate)

	unitFile, _ := unit.NewUnitFile(unitTemplate.String())

	return fleet.MapUnitFileToSchemaUnitOptions(unitFile)
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
