package server

import (
	"log"
	"net/http"
	"net/url"

	"github.com/bmorton/deployster/fleet"
)

// UnitsResource is the HTTP resource responsible for getting basic information
// on all units that exist for a given service.
type UnitsResource struct {
	Fleet fleet.Client
}

// UnitsResponse is the wrapper struct for the JSON payload returned by the
// Index action.
type UnitsResponse struct {
	Units []VersionedUnit `json:"units"`
}

// Index is the GET endpoint for listing all units that exist for a given
// service.
//
// This function assumes that it is nested inside `/services/{name}`
// and that Tigertonic is extracting the service name and providing it via query
// params.
func (ur *UnitsResource) Index(u *url.URL, h http.Header, req interface{}) (int, http.Header, *UnitsResponse, error) {
	statusCode := http.StatusOK
	response := &UnitsResponse{}

	units, err := ur.Fleet.Units()
	if err != nil {
		log.Printf("%#v\n", err)
		return http.StatusInternalServerError, nil, nil, err
	}
	response.Units = FindServiceUnits(u.Query().Get("name"), units)

	return statusCode, nil, response, nil
}
