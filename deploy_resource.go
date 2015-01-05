package main

import (
	"bytes"
	"fmt"
	"github.com/bmorton/deployster/fleet"
	"github.com/coreos/fleet/schema"
	"github.com/coreos/fleet/unit"
	"log"
	"net/http"
	"net/url"
	"text/template"
)

type DeployResource struct {
	Fleet fleet.Client
}

type Deploy struct {
	Version         string `json:"version"`
	DestroyPrevious bool   `json:"destroy_previous"`
}

type DeployRequest struct {
	Deploy Deploy `json:"deploy"`
}

type UnitTemplate struct {
	Name    string
	Version string
}

func (self *DeployResource) Create(u *url.URL, h http.Header, req *DeployRequest) (int, http.Header, interface{}, error) {
	var previousVersion string
	serviceName := u.Query().Get("name")
	options := getUnitOptions(serviceName, req.Deploy.Version)
	fleetName := fleetServiceName(serviceName, req.Deploy.Version)

	if req.Deploy.DestroyPrevious {
		units, err := self.Fleet.Units()
		if err != nil {
			log.Printf("%#v\n", err)
			return http.StatusInternalServerError, nil, nil, err
		}
		versions := FindServiceVersions(u.Query().Get("name"), units)

		if len(versions) != 1 {
			log.Printf("Can't destroy previous versions (%d previous versions), disabling destroy.")
			req.Deploy.DestroyPrevious = false
		} else {
			previousVersion = versions[0]
		}
	}

	resp, err := self.Fleet.StartUnit(fleetName, options)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	fmt.Printf("%#v\n", resp)

	if req.Deploy.DestroyPrevious {
		self.destroyPrevious(serviceName, previousVersion)
	}

	return http.StatusCreated, nil, nil, nil
}

func (self *DeployResource) Destroy(u *url.URL, h http.Header, req interface{}) (int, http.Header, interface{}, error) {
	fleetName := fleetServiceName(u.Query().Get("name"), u.Query().Get("version"))

	resp, err := self.Fleet.DestroyUnit(fleetName)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	fmt.Printf("%#v\n", resp)

	return http.StatusNoContent, nil, nil, nil
}

func getUnitOptions(name string, version string) []fleet.UnitOption {
	var unitTemplate bytes.Buffer
	t, _ := template.New("test").Parse(DOCKER_UNIT_TEMPLATE)
	t.Execute(&unitTemplate, UnitTemplate{name, version})

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

func (self *DeployResource) destroyPrevious(name string, version string) {
	resp, err := self.Fleet.DestroyUnit(fleetServiceName(name, version))
	if err != nil {
		log.Printf("%#v\n", err)
	} else {
		fmt.Printf("%#v\n", resp)
	}
}
