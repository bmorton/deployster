package fleet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type client struct {
	http *http.Client
}

type Client interface {
	Units() ([]Unit, error)
	StartUnit(name string, options []UnitOption) (*http.Response, error)
	DestroyUnit(name string) (*http.Response, error)
	UnitState(name string) (UnitState, error)
}

var _ Client = &client{}

type Unit struct {
	CurrentState string       `json:"currentState,omitempty"`
	DesiredState string       `json:"desiredState"`
	MachineID    string       `json:"machineID,omitempty"`
	Name         string       `json:"name,omitempty"`
	Options      []UnitOption `json:"options"`
}

type UnitOption struct {
	Section string `json:"section"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

type UnitState struct {
	Name        string `json:"name"`
	Hash        string `json:"hash"`
	MachineID   string `json:"machineID"`
	LoadState   string `json:"systemdLoadState"`
	ActiveState string `json:"systemdActiveState"`
	SubState    string `json:"systemdSubState"`
}

type UnitsResponse struct {
	Units []Unit `json:"units"`
}

type StatesResponse struct {
	States []UnitState `json:"states"`
}

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

func (c *client) DestroyUnit(name string) (*http.Response, error) {
	url := fmt.Sprintf("http://sock/fleet/v1/units/%s", name)

	r, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return c.http.Do(r)
}

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
