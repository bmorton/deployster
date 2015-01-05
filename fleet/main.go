package fleet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type Client struct {
	http *http.Client
}

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

type UnitsResponse struct {
	Units []Unit `json:"units"`
}

func NewClient(path string) Client {
	dialFunc := func(string, string) (net.Conn, error) {
		return net.Dial("unix", path)
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			Dial: dialFunc,
		},
	}

	return Client{&httpClient}
}

func (self *Client) Units() ([]Unit, error) {
	response, err := self.http.Get("http://sock/fleet/v1/units")
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

func (self *Client) StartUnit(name string, options []UnitOption) (*http.Response, error) {
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

	return self.http.Do(r)
}

func (self *Client) DestroyUnit(name string) (*http.Response, error) {
	url := fmt.Sprintf("http://sock/fleet/v1/units/%s", name)

	r, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return self.http.Do(r)
}
