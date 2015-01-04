package fleet

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
)

type Client struct {
	http *http.Client
}

type Unit struct {
	CurrentState string `json:"currentState"`
	DesiredState string `json:"desiredState"`
	MachineID    string `json:"machineID"`
	Name         string `json:"name"`
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

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var parsedResponse UnitsResponse
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return nil, err
	}

	return parsedResponse.Units, nil
}
