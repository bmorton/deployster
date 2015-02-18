// fleet is a package for communicating with Fleet's HTTP API to list units and
// their states, launch new units, and destroy units.  It uses the API that is
// documented here:
//   https://github.com/coreos/fleet/blob/master/Documentation/api-v1.md
package fleet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// Client is the main interface for the Fleet client that is implemented in this
// package.
type Client interface {
	Units() ([]Unit, error)
	StartUnit(name string, options []UnitOption) (*http.Response, error)
	DestroyUnit(name string) (*http.Response, error)
	UnitState(name string) (UnitState, error)
}

// client is the Fleet client implemenation for the interface described above.
// Use NewClient in this package to get a new client with the HTTP client wired
// up to talk to fleet's socket.
type client struct {
	http *http.Client
}

var _ Client = &client{}

// Unit is the representation of a Fleet unit and is used to deserialze the JSON
// returned from the Fleet API.  This is also the struct that is returned when
// this package returns units via Units().
type Unit struct {
	CurrentState string       `json:"currentState,omitempty"`
	DesiredState string       `json:"desiredState"`
	MachineID    string       `json:"machineID,omitempty"`
	Name         string       `json:"name,omitempty"`
	Options      []UnitOption `json:"options"`
}

// UnitOption is the representation for each of a Fleet unit's individual
// options.  It is used to deserialze the JSON returned from the Fleet API as
// well as provide the information for starting a new unit.
type UnitOption struct {
	Section string `json:"section"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

// UnitState is the representation for detailed information about the current
// state of a Fleet unit.  It is used to deserialize the JSON returned from the
// Fleet API when requesting the state of a given unit.  This is also the struct
// that is returned when this package returns the state via UnitState().
type UnitState struct {
	Name        string `json:"name"`
	Hash        string `json:"hash"`
	MachineID   string `json:"machineID"`
	LoadState   string `json:"systemdLoadState"`
	ActiveState string `json:"systemdActiveState"`
	SubState    string `json:"systemdSubState"`
}

// UnitsResponse is the wrapper for a Fleet API response containing an array of
// units and is used to deserialize JSON.
type UnitsResponse struct {
	Units []Unit `json:"units"`
}

// StatesResponse is the wrapper for a Fleet API response containing an array of
// unit states and is used to deserialize JSON.
type StatesResponse struct {
	States []UnitState `json:"states"`
}

// NewClient construct an HTTP client using the given path as the socket to dial
// out to for making requests and returns a new Fleet client with the HTTP
// client wired up.  Currently this method *only* supports communicating with
// Fleet over a socket.
func NewClient(path string) client {
	dialFunc := func(string, string) (net.Conn, error) {
		return net.Dial("unix", path)
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			Dial: dialFunc,
		},
	}

	return client{&httpClient}
}

// Units queries the Fleet API for a list of all units and returns them as an
// array of Unit structs.
func (c *client) Units() ([]Unit, error) {
	response, err := c.http.Get("http://sock/fleet/v1/units")
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)
	var parsedResponse UnitsResponse
	err = decoder.Decode(&parsedResponse)
	if err != nil {
		return nil, err
	}

	return parsedResponse.Units, nil
}

// StartUnit uses the given name and unit options to PUT a request to launch
// a unit.  The UnitOption array is similar to a unit file that you'd use to
// start a unit via fleetctl.
func (c *client) StartUnit(name string, options []UnitOption) (*http.Response, error) {
	url := fmt.Sprintf("http://sock/fleet/v1/units/%s", name)
	unit := Unit{
		DesiredState: "launched",
		Options:      options,
	}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.Encode(unit)

	r, err := http.NewRequest("PUT", url, &b)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Content-Type", "application/json")

	return c.http.Do(r)
}

// DestroyUnit issues a DELETE request to the Fleet API to remove the unit for
// the name passed.
func (c *client) DestroyUnit(name string) (*http.Response, error) {
	url := fmt.Sprintf("http://sock/fleet/v1/units/%s", name)

	r, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return c.http.Do(r)
}

// UnitState gets detailed state information about the unit that is passed.  We
// assume that a full name is given and that only one match will be returned
// from the Fleet API for that name.
func (c *client) UnitState(name string) (UnitState, error) {
	url := fmt.Sprintf("http://sock/fleet/v1/state?unitName=%s", name)
	response, err := c.http.Get(url)
	if err != nil {
		return UnitState{}, err
	}

	decoder := json.NewDecoder(response.Body)
	var parsedResponse StatesResponse
	err = decoder.Decode(&parsedResponse)
	if err != nil {
		return UnitState{}, err
	}

	if len(parsedResponse.States) > 0 {
		return parsedResponse.States[0], nil
	} else {
		return UnitState{}, nil
	}
}
