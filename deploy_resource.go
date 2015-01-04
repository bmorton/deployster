package main

import (
	"bytes"
	"fmt"
	"github.com/bmorton/deployster/fleet"
	"github.com/coreos/fleet/schema"
	"github.com/coreos/fleet/unit"
	"net/http"
	"net/url"
	"text/template"
)

type DeployResource struct {
	Fleet fleet.Client
}

type Deploy struct {
	Version string `json:"version"`
}

type DeployRequest struct {
	Deploy Deploy `json:"deploy"`
}

type UnitTemplate struct {
	Name    string
	Version string
}

func (self *DeployResource) Create(u *url.URL, h http.Header, req *DeployRequest) (int, http.Header, interface{}, error) {
	serviceName := u.Query().Get("name")
	unitContents := buildUnitFile(serviceName, req.Deploy.Version)
	unitFile, _ := unit.NewUnitFile(unitContents)
	options := schema.MapUnitFileToSchemaUnitOptions(unitFile)
	serviceWithVersion := fmt.Sprintf("%s-%s", serviceName, req.Deploy.Version)

	resp, err := self.Fleet.StartUnit(serviceWithVersion, schemaToLocalUnit(options))
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}
	fmt.Printf("%#v\n", resp)

	return http.StatusCreated, nil, "", nil
}

func buildUnitFile(name string, version string) string {
	var unitFile bytes.Buffer
	t, _ := template.New("test").Parse(DOCKER_UNIT_TEMPLATE)
	t.Execute(&unitFile, UnitTemplate{name, version})

	return unitFile.String()
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
