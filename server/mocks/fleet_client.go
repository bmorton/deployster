package mocks

import "github.com/stretchr/testify/mock"

import "github.com/coreos/fleet/schema"

type FleetClient struct {
	mock.Mock
}

func (m *FleetClient) Units() ([]*schema.Unit, error) {
	ret := m.Called()

	r0 := ret.Get(0).([]*schema.Unit)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *FleetClient) CreateUnit(_a0 *schema.Unit) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *FleetClient) DestroyUnit(_a0 string) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
func (m *FleetClient) UnitStates() ([]*schema.UnitState, error) {
	ret := m.Called()

	r0 := ret.Get(0).([]*schema.UnitState)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *FleetClient) SetUnitTargetState(_a0 string, _a1 string) error {
	ret := m.Called(_a0, _a1)

	r0 := ret.Error(0)

	return r0
}
