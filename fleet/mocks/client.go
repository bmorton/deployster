package mocks

import (
	"github.com/bmorton/deployster/fleet"
	"github.com/stretchr/testify/mock"
	"net/http"
)

type Client struct {
	mock.Mock
}

func (m *Client) Units() ([]fleet.Unit, error) {
	ret := m.Called()

	r0 := ret.Get(0).([]fleet.Unit)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Client) StartUnit(name string, options []fleet.UnitOption) (*http.Response, error) {
	ret := m.Called(name, options)

	r0 := ret.Get(0).(*http.Response)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Client) DestroyUnit(name string) (*http.Response, error) {
	ret := m.Called(name)

	r0 := ret.Get(0).(*http.Response)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *Client) UnitState(name string) (fleet.UnitState, error) {
	ret := m.Called(name)

	r0 := ret.Get(0).(fleet.UnitState)
	r1 := ret.Error(1)

	return r0, r1
}
