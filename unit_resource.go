package main

import (
	"github.com/bmorton/deployster/fleet"
	"log"
	"net/http"
	"net/url"
)

type UnitResource struct {
	Fleet fleet.Client
}

type UnitResponse struct {
	Units []VersionedUnit `json:"units"`
}

func (self *UnitResource) Index(u *url.URL, h http.Header, req interface{}) (int, http.Header, *UnitResponse, error) {
	statusCode := http.StatusOK
	response := &UnitResponse{}

	units, err := self.Fleet.Units()
	if err != nil {
		log.Printf("%#v\n", err)
		return http.StatusInternalServerError, nil, nil, err
	}
	response.Units = FindServiceUnits(u.Query().Get("name"), units)

	return statusCode, nil, response, nil
}
