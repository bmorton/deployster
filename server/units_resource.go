package server

import (
	"github.com/bmorton/deployster/fleet"
	"log"
	"net/http"
	"net/url"
)

type UnitsResource struct {
	Fleet fleet.Client
}

type UnitsResponse struct {
	Units []VersionedUnit `json:"units"`
}

func (self *UnitsResource) Index(u *url.URL, h http.Header, req interface{}) (int, http.Header, *UnitsResponse, error) {
	statusCode := http.StatusOK
	response := &UnitsResponse{}

	units, err := self.Fleet.Units()
	if err != nil {
		log.Printf("%#v\n", err)
		return http.StatusInternalServerError, nil, nil, err
	}
	response.Units = FindServiceUnits(u.Query().Get("name"), units)

	return statusCode, nil, response, nil
}
